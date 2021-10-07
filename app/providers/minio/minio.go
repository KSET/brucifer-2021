package minio

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minio_ *minio.Client

type minioProvider struct{}

func MinioProvider() minioProvider {
	return minioProvider{}
}

func (p minioProvider) log(format string, v ...interface{}) {
	log.Printf("|MINIO> "+format, v...)
}

func (p minioProvider) Client() *minio.Client {
	return minio_
}

func (p minioProvider) BucketName() string {
	return "brucosijada-2021-uploads"
}

func (p minioProvider) Register() (err error) {
	endpoint := os.Getenv("MINIO_HOST")
	accessKeyID := os.Getenv("MINIO_USER")
	secretAccessKey := os.Getenv("MINIO_PASS")
	useSSL := false

	p.log("Connecting to MinIO (H: `%s` U: `%s` P: `%s`)", endpoint, accessKeyID, secretAccessKey)
	minio_, err = minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		},
	)
	if err != nil {
		return err
	}
	p.log("Connected to MinIO")

	if err := p.createBucket(); err != nil {
		return err
	}

	return nil
}

func (p minioProvider) createBucket() (err error) {
	ctx := context.Background()

	bucketName := p.BucketName()
	location := "us-east-1"

	p.log("Creating bucket `%s` on `%s`", bucketName, location)

	err = p.Client().MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := p.Client().BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			p.log("We already own %s\n", bucketName)
		} else {
			p.log("Error %+v\n", errBucketExists)
			return err
		}
	} else {
		p.log("Successfully created %s\n", bucketName)
	}

	return nil
}
