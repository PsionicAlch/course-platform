package assets

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/PsionicAlch/psionicalch-home/internal/bucket"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

//go:embed */**
var assetFiles embed.FS

type Assets struct {
	utils.Loggers
	Bucket bucket.Bucket
}

func SetupAssets(b bucket.Bucket) *Assets {
	return &Assets{
		Loggers: utils.CreateLoggers("ASSETS"),
		Bucket:  b,
	}
}

func (a *Assets) SyncAssets() error {
	bucketFiles, err := a.Bucket.GetAllFiles()
	if err != nil {
		a.ErrorLog.Printf("Failed to get all files from bucket: %s\n", err)
		return err
	}

	bucketMap := make(map[string]string, len(bucketFiles))
	for _, bucketFile := range bucketFiles {
		bucketMap[bucketFile.Name] = bucketFile.Checksum
	}

	fileMap, err := CreateFileMap()
	if err != nil {
		a.ErrorLog.Printf("Failed to read assets directory: %v\n", err)
		return err
	}

	filesToAdd, filesToDelete := SortFiles(fileMap, bucketMap)

	for _, file := range filesToAdd {
		a.InfoLog.Printf("Adding \"%s\" to the bucket\n", file.Name)

		err := a.Bucket.UploadFileFS(assetFiles, file.Name, file.Checksum)
		if err != nil {
			a.ErrorLog.Printf("Failed to upload \"%s\" to bucket: %v\n", file, err)
			return err
		}
	}

	for _, file := range filesToDelete {
		a.InfoLog.Printf("Deleting \"%s\" from the bucket\n", file)

		err := a.Bucket.DeleteFile(file)
		if err != nil {
			a.ErrorLog.Printf("Failed to delete \"%s\" from bucket: %v\n", file, err)
			return err
		}
	}

	return nil
}

func CreateFileMap() (map[string]string, error) {
	fileMap := make(map[string]string)

	err := fs.WalkDir(assetFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing \"%s\": %w", path, err)
		}

		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			return nil
		}

		if !d.IsDir() {
			file, err := assetFiles.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open \"%s\": %w", path, err)
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("failed to read \"%s\": %w", path, err)
			}

			fileMap[path] = GenerateChecksum(content)
		}

		return nil
	})
	if err != nil {
		return fileMap, err
	}

	return fileMap, nil
}

func GenerateChecksum(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func SortFiles(fileMap map[string]string, bucketMap map[string]string) (filesToAdd []*bucket.File, filesToDelete []string) {
	// Determine which files to add.
	for filePath, checksum := range fileMap {
		file := &bucket.File{
			Name:     filePath,
			Checksum: checksum,
		}

		if bucketChecksum, has := bucketMap[filePath]; has {
			if checksum != bucketChecksum {
				fmt.Printf("Checksums for %s:\n\tBucket File's Checksum: %s\n\tLocal File's Checksum: %s\n", filePath, bucketChecksum, checksum)

				filesToAdd = append(filesToAdd, file)
			}
		} else {
			filesToAdd = append(filesToAdd, file)
		}
	}

	// Determine which files to delete.
	for filePath := range bucketMap {
		if _, has := fileMap[filePath]; !has {
			filesToDelete = append(filesToDelete, filePath)
		}
	}

	return
}
