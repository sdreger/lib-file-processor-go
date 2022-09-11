package filestore

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"path/filepath"
	"testing"
)

const (
	testBucketName = "test-bucket"
	testFileName   = "1.png"
)

type minioContainer struct {
	testcontainers.Container
	BlobStore MinioStore
}

func TestMinioStore_CreateBucket(t *testing.T) {
	t.Log("Given the need to test Minio bucket creation.")

	ctx := context.Background()
	minioContainer := createMinioContainer(ctx, t)
	defer minioContainer.Terminate(ctx)

	bucketExists, err := minioContainer.BlobStore.client.BucketExists(ctx, testBucketName)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to check if Minio bucket exists: %s", failed, err)
	}
	if bucketExists {
		t.Fatalf("\t\t%s\tThe bucket %q should not exist", failed, testBucketName)
	}

	requestCtx := context.Background()
	err = minioContainer.BlobStore.CreateBucket(requestCtx, testBucketName)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to Minio bucket: %s", failed, err)
	}

	bucketExists, err = minioContainer.BlobStore.client.BucketExists(ctx, testBucketName)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to check if Minio bucket exists: %s", failed, err)
	}
	if !bucketExists {
		t.Fatalf("\t\t%s\tThe bucket %q should exist", failed, testBucketName)
	}

	t.Logf("\t\t%s\tShould be able to create Minio bucket", succeed)
}

func TestMinioStore_StoreObject(t *testing.T) {
	t.Log("Given the need to test Minio object store.")

	ctx := context.Background()
	minioContainer := createMinioContainer(ctx, t)
	defer minioContainer.Terminate(ctx)

	requestCtx := context.Background()
	err := minioContainer.BlobStore.CreateBucket(requestCtx, testBucketName)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to Minio bucket: %s", failed, err)
	}

	etag, err := minioContainer.BlobStore.
		StoreObject(ctx, testBucketName, testFileName, filepath.Join("testdata", testFileName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to store Minio object: %s", failed, err)
	}
	if etag == "" {
		t.Fatalf("\t\t%s\tShould be able to get Minio object Etag", failed)
	}

	var fileExists bool
	objects := minioContainer.BlobStore.client.ListObjects(ctx, testBucketName, minio.ListObjectsOptions{})
	for file := range objects {
		if file.Size == 0 || file.Key != testFileName {
			t.Fatalf("\t\t%s\tThe bucket should contain %q file", failed, testFileName)
		} else {
			fileExists = true
		}
	}
	if !fileExists {
		t.Fatalf("\t\t%s\tThe bucket should contain %q file", failed, testFileName)
	}

	t.Logf("\t\t%s\tShould be able to store Minio object", succeed)
}

func createMinioContainer(ctx context.Context, t *testing.T) *minioContainer {
	appConfig := config.GetAppConfig()
	req := testcontainers.ContainerRequest{
		Image: "quay.io/minio/minio:RELEASE.2022-06-03T01-40-53Z",
		Env: map[string]string{
			"MINIO_ROOT_USER":     appConfig.MinioAccessKeyID,
			"MINIO_ROOT_PASSWORD": appConfig.MinioSecretAccessKey,
		},
		ExposedPorts: []string{"9000"},
		Cmd:          []string{"server", "/data"},
		WaitingFor:   wait.ForLog("Status:         1 Online, 0 Offline."),
		Name:         "minio_test_container",
	}
	minioC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create Minio container: %s", failed, err)
	}

	ip, err := minioC.Host(ctx)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to get Minio container IP: %s", failed, err)
	}

	port, err := minioC.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to get Minio container port: %s", failed, err)
	}

	minioStore, err := NewMinioStore(fmt.Sprintf("%s:%s", ip, port.Port()), appConfig.MinioAccessKeyID,
		appConfig.MinioSecretAccessKey, appConfig.MinioUseSSL, log.Default())
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a Minio store: %s", failed, err)
	}

	return &minioContainer{
		Container: minioC,
		BlobStore: minioStore,
	}
}
