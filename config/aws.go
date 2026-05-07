package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client adalah variabel global yang akan dipakai untuk upload file
var S3Client *s3.Client

func ConnectAWS() {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Membangun konfigurasi AWS menggunakan kredensial dari .env
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi AWS: %v", err)
	}

	// Inisialisasi klien S3
	S3Client = s3.NewFromConfig(cfg)
	log.Println("Koneksi ke AWS S3 berhasil diinisialisasi!")
}