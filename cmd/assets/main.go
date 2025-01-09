package main

import (
	"time"

	awss3 "github.com/PsionicAlch/psionicalch-home/internal/bucket/aws_s3"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/envloader"
	"github.com/PsionicAlch/psionicalch-home/pkg/envloader/validators"
	"github.com/PsionicAlch/psionicalch-home/web/assets"
)

type AssetPipeline struct {
	utils.Loggers
}

func main() {
	pipeline := &AssetPipeline{
		Loggers: utils.CreateLoggers("ASSET PIPELINE"),
	}

	bucket := pipeline.CreateBucket()

	startTimer := time.Now()

	pipeline.InfoLog.Println("Starting asset sync..")

	asset := assets.SetupAssets(bucket)
	if err := asset.SyncAssets(); err != nil {
		pipeline.ErrorLog.Fatal(err)
	}

	endTimer := time.Since(startTimer)

	pipeline.InfoLog.Printf("Finished syncing assets in %s!\n", endTimer)
}

func (pipeline *AssetPipeline) CreateBucket() *awss3.S3Bucket {
	envloader.LoadEnvironment(map[string]validators.ValidationFunc{
		"REGION":            validators.NotEmpty,
		"ACCESS_KEY_ID":     validators.NotEmpty,
		"SECRET_ACCESS_KEY": validators.NotEmpty,
		"BUCKET_NAME":       validators.NotEmpty,
	})

	region, err := envloader.GetVariable[string]("REGION")
	if err != nil {
		pipeline.ErrorLog.Fatalf("Failed to get \"REGION\" from env: %s\n", err)
	}

	accessKeyID, err := envloader.GetVariable[string]("ACCESS_KEY_ID")
	if err != nil {
		pipeline.ErrorLog.Fatalf("Failed to get \"ACCESS_KEY_ID\" from env: %s\n", err)
	}

	secretAccessKey, err := envloader.GetVariable[string]("SECRET_ACCESS_KEY")
	if err != nil {
		pipeline.ErrorLog.Fatalf("Failed to get \"SECRET_ACCESS_KEY\" from env: %s\n", err)
	}

	bucketName, err := envloader.GetVariable[string]("BUCKET_NAME")
	if err != nil {
		pipeline.ErrorLog.Fatalf("Failed to get \"BUCKET_NAME\" from env: %s\n", err)
	}

	bucket, err := awss3.SetupS3Bucket(region, accessKeyID, secretAccessKey, bucketName)
	if err != nil {
		pipeline.ErrorLog.Fatal(err)
	}

	return bucket
}
