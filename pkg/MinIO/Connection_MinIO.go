package minio

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	client *minio.Client
	Cf     *config.Config
}

// ContentType map now supports a wider range of image formats
var ContentType = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".tiff": "image/tiff",
}

func MinIOConnect(cf *config.Config) (*MinIO, error) {
	endpoint := cf.MinioUrl
	accessKeyID := cf.MinioUser
	secretAccessKey := cf.MinIOSecredKey
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bucketName := cf.MinIOBucketName

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists != nil && exists {
			log.Println(err)
			return nil, err
		}
	}

	policy := fmt.Sprintf(`{
      "Version": "2012-10-17",
      "Statement": [
          {
              "Effect": "Allow",
              "Principal": "*",
              "Action": ["s3:GetObject"],
              "Resource": ["arn:aws:s3:::%s/*"]
          }
      ]
  }`, bucketName)

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		log.Println("error while setting bucket policy : ", err)
		return nil, err
	}

	return &MinIO{
		client: minioClient,
		Cf:     cf,
	}, err
}

func (m *MinIO) Upload(fileName, filePath string) (string, error) {
	// Fayl kengaytmasini olish
	ext := filepath.Ext(fileName)
	contentType, ok := ContentType[ext]
	if !ok {
		// Agar kengaytma mos kelmasa, default ContentType qo'yamiz
		contentType = "application/octet-stream"
	}

	cfg, _ := config.NewConfig()

	// Faylni MinIOga yuklash
	_, err := m.client.FPutObject(context.Background(), cfg.MinIOBucketName, fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		slog.Error("Error while uploading file", "fileName", fileName, "bucket", cfg.MinIOBucketName, "error", err)
		return "", err
	}

	// MinIO URLni yaratish
	minioURL := fmt.Sprintf("http://localhost:9000/%s/%s", cfg.MinIOBucketName, fileName)

	return minioURL, nil
}
