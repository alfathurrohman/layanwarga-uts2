package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"layanwarga/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// UploadToS3 mengirim file fisik dari memori server langsung ke AWS S3 Bucket
func UploadToS3(fileHeader *multipart.FileHeader, folderPath string) (string, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	cloudFrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := fmt.Sprintf("%s/%s", folderPath, fileHeader.Filename)

	// 1. Ambil tipe file asli (contoh: image/jpeg, application/pdf)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Default jika tidak terdeteksi
	}

	// 2. Kirim ke S3 dengan menyertakan ContentType
	_, err = config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType), // Baris baru untuk memberitahu browser agar tidak mendownload
	})

	if err != nil {
		return "", err
	}

	fileUrl := fmt.Sprintf("https://%s/%s", cloudFrontDomain, filename)
	
	return fileUrl, nil
}