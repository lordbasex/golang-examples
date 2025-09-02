package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)


const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
	AWS_REGION            = "us-east-1"
	URL_EXPIRE_TIME       = 24
	BUCKET                = "TU_BUCKET_NAME"
)

func uploadAndPresign(ctx context.Context, bucket, key, localPath string, expire time.Duration) (string, error) {
	// Crear configuración con credenciales explícitas
	cfg := aws.Config{
		Region: AWS_REGION,
		Credentials: aws.NewCredentialsCache(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     AWS_ACCESS_KEY_ID,
					SecretAccessKey: AWS_SECRET_ACCESS_KEY,
					Source:          "HardcodedConfig",
				},
			}),
	}

	s3c := s3.NewFromConfig(cfg)

	// Abrir archivo local
	f, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Subir a S3
	uploader := manager.NewUploader(s3c)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		Body:               f,
		ACL:                types.ObjectCannedACLPrivate,
		ContentType:        aws.String("application/gzip"),
		ContentDisposition: aws.String(fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(localPath))),
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	// Generar URL prefirmada (expira en 10 minutos)
	ps := s3.NewPresignClient(s3c)
	resp, err := ps.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(localPath))),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}

	return resp.URL, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	const (
		bucket    = BUCKET
		localPath = "./reporte.csv.gz"
		key       = "backup/reporte.csv.gz"
	)

	url, err := uploadAndPresign(ctx, bucket, key, localPath, time.Duration(URL_EXPIRE_TIME)*time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("URL temporal ( %dhs ): %s\n", URL_EXPIRE_TIME, url)
}
