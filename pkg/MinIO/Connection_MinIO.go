package minio

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	client *minio.Client
	Cf     *config.Config
}

var ContentType = map[string]string{
	".png": "image/png",
	".pdf": "application/pdf",
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

func (m *MinIO) Upload(fileName string, contentType string) (*string, error) {

	uploadPath := fileName
	c_type := ContentType[contentType]
	_, err := m.client.FPutObject(context.Background(), m.Cf.MinIOBucketName, fileName, uploadPath, minio.PutObjectOptions{
		ContentType: c_type,
	})

	if err != nil {
		return nil, fmt.Errorf("error while uploading to minio: %v", err)
	}

	// Delete the media in uploadPath after uploading to minio
	err = os.Remove(uploadPath)
	if err != nil {
		return nil, fmt.Errorf("error while deleting the file: %v", err)
	}

	minioURL := fmt.Sprintf("http://%s/%s/%s", m.Cf.MinioUrl, m.Cf.MinIOBucketName, fileName)
	
	return &minioURL, nil
}
