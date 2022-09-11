package filestore

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioStore struct {
	client *minio.Client
	logger *log.Logger
}

func NewMinioStore(endpoint, accessKeyID, secretAccessKey string, useSSL bool, logger *log.Logger) (MinioStore, error) {
	client, err := getMinioClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return MinioStore{}, err
	}

	return MinioStore{
		client: client,
		logger: logger,
	}, nil
}

// CreateBucket creates a new Minio bucket.
// If there is a bucket with the same name, creation will be skipped.
func (ms MinioStore) CreateBucket(ctx context.Context, bucketName string) error {
	bucketExists, err := ms.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if bucketExists {
		return nil
	}

	return ms.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

// StoreObject puts an object to a specified bucket.
// If the bucket doesn't exist, tries to create one.
func (ms MinioStore) StoreObject(ctx context.Context, bucketName string, fileName, filePath string) (string, error) {
	bucketExists, err := ms.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("can not check if a bucket exists: %w", err)
	}

	if !bucketExists {
		ms.logger.Printf("[INFO] - The bucket %q doesn't exist, creating...", bucketName)
		err := ms.CreateBucket(ctx, bucketName)
		if err != nil {
			return "", fmt.Errorf("can not create bucket: %w", err)
		}
	}

	object, err := ms.client.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("can not store file: %w", err)
	}

	return object.ETag, nil
}

// getMinioClient initializes a new Minio client.
func getMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	_, err = client.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}
