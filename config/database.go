package config

import (
	"fmt"
	"log"
	"os"

	"layanwarga/models" // Import package models kita

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan sistem")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database MySQL:", err)
	}

	DB = database
	fmt.Println("Koneksi ke database MySQL berhasil!")

	// Menjalankan migrasi otomatis untuk membuat tabel
	err = DB.AutoMigrate(&models.User{}, &models.PengajuanSurat{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}
	fmt.Println("Migrasi tabel database berhasil!")
}