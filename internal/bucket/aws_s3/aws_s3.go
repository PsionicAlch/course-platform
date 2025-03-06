package awss3

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/PsionicAlch/psionicalch-home/internal/bucket"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Bucket struct {
	utils.Loggers
	Client     *s3.Client
	Region     string
	BucketName string
}

func SetupS3Bucket(region, accessKeyID, secretAccessKey, bucketName string) (*S3Bucket, error) {
	loggers := utils.CreateLoggers("S3 BUCKET")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")))
	if err != nil {
		loggers.ErrorLog.Printf("Failed to load AWS default config: %s\n", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	s3Bucket := &S3Bucket{
		Loggers:    loggers,
		Client:     client,
		Region:     region,
		BucketName: bucketName,
	}

	bucketExists, err := s3Bucket.BucketExists(bucketName)
	if err != nil {
		return nil, err
	}

	if !bucketExists {
		return nil, fmt.Errorf("\"%s\" bucket does not exists", bucketName)
	}

	return s3Bucket, nil
}

func (b *S3Bucket) BucketExists(bucketName string) (bool, error) {
	_, err := b.Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return false, nil
			default:
				b.ErrorLog.Printf("Failed to check if \"%s\" bucket is available: %s\n", bucketName, err)
				return false, err
			}
		} else {
			b.ErrorLog.Printf("Failed to check if \"%s\" bucket is available: %s\n", bucketName, err)
			return false, err
		}
	}

	return true, nil
}

func (b *S3Bucket) GetFileMetadata(fileName string) (map[string]string, error) {
	resp, err := b.Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		b.ErrorLog.Printf("Failed to retrieve metadata for \"%s\": %s\n", fileName, err)
		return nil, err
	}

	return resp.Metadata, nil
}

func (b *S3Bucket) GetAllFiles() ([]*bucket.File, error) {
	var files []*bucket.File

	paginator := s3.NewListObjectsV2Paginator(b.Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(b.BucketName),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			b.ErrorLog.Printf("Failed to get a list of all files in \"%s\" bucket: %s\n", b.BucketName, err)
			return nil, err
		}

		for _, object := range output.Contents {
			objMetadata, err := b.GetFileMetadata(*object.Key)
			if err != nil {
				return nil, err
			}

			if _, has := objMetadata["checksum"]; !has {
				return nil, fmt.Errorf("\"%s\" does not contain checksum metadata", *object.Key)
			}

			file := bucket.File{
				Name:     *object.Key,
				Checksum: objMetadata["checksum"],
			}

			files = append(files, &file)
		}
	}

	return files, nil
}

func (b *S3Bucket) UploadFileFS(files embed.FS, fileName, checksum string) error {
	output, err := files.ReadFile(fileName)
	if err != nil {
		b.ErrorLog.Printf("Failed to read \"%s\": %s\n", fileName, err)
		return err
	}

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = b.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:            aws.String(b.BucketName),
		Key:               aws.String(fileName),
		Body:              strings.NewReader(string(output)),
		ContentType:       aws.String(contentType),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		Metadata: map[string]string{
			"checksum": checksum,
		},
	})
	if err != nil {
		b.ErrorLog.Printf("Failed to upload \"%s\" to \"%s\" bucket: %s\n", fileName, b.BucketName, err)
		return err
	}

	return nil
}

func (b *S3Bucket) DeleteFile(fileName string) error {
	_, err := b.Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		b.ErrorLog.Printf("Failed to delete \"%s\" from \"%s\" bucket: %s\n", fileName, b.BucketName, err)
		return err
	}

	return nil
}
