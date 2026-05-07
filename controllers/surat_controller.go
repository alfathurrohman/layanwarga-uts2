package controllers

import (
	"fmt"
	"net/http"
	"time"

	"layanwarga/config"
	"layanwarga/models"
	"layanwarga/utils"

	"github.com/gin-gonic/gin"
)

func BuatPengajuan(c *gin.Context) {
	// 1. Ambil ID Warga dari tiket JWT (yang dikirim via Middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, ID tidak ditemukan"})
		return
	}

	// 2. Ambil data teks dari Form-Data
	jenisSurat := c.PostForm("jenis_surat")
	keperluan := c.PostForm("keperluan")

	if jenisSurat == "" || keperluan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis surat dan keperluan wajib diisi"})
		return
	}

	// 3. Ambil file KTP dan KK dari Form-Data
	fileKTP, errKTP := c.FormFile("file_ktp")
	fileKK, errKK := c.FormFile("file_kk")

	if errKTP != nil || errKK != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File KTP dan KK wajib diupload"})
		return
	}

	// 4. Proses Upload ke S3 (via folder utils)
	// Membuat folder unik berdasarkan ID user dan waktu agar nama file tidak bentrok
	folderPath := fmt.Sprintf("uploads/user_%v/%d", userID, time.Now().Unix())
	
	urlKTP, err := utils.UploadToS3(fileKTP, folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload KTP ke S3"})
		return
	}

	urlKK, err := utils.UploadToS3(fileKK, folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload KK ke S3"})
		return
	}

	// 5. Simpan data pengajuan ke Database
	pengajuan := models.PengajuanSurat{
		UserID:     uint(userID.(float64)), // Konversi tipe data bawaan JWT
		JenisSurat: jenisSurat,
		Keperluan:  keperluan,
		FileKTPUrl: urlKTP, // URL ini sudah berformat CloudFront
		FileKKUrl:  urlKK,
		Status:     "Pending",
	}

	if err := config.DB.Create(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan ke database"})
		return
	}

	// 6. Kembalikan respons sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengajuan surat berhasil dibuat",
		"data":    pengajuan,
	})
}

// Fungsi untuk melihat daftar pengajuan
func GetPengajuan(c *gin.Context) {
	// Ambil data user dari tiket JWT
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var pengajuans []models.PengajuanSurat

	// Logika: Jika yang login adalah admin, tampilkan SEMUA pengajuan.
	// Jika yang login adalah warga, tampilkan HANYA pengajuan miliknya sendiri.
	if role == "admin" {
		// Mengambil semua data dan menggabungkan (Preload) dengan data User
		config.DB.Preload("User").Find(&pengajuans)
	} else {
		config.DB.Where("user_id = ?", userID).Preload("User").Find(&pengajuans)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data pengajuan",
		"data":    pengajuans,
	})
}

// Struct untuk menangkap request update dari admin
type UpdateStatusInput struct {
	Status     string `json:"status" binding:"required"`
	Keterangan string `json:"keterangan"`
}

// Fungsi khusus Admin untuk mengubah status
func UpdateStatusPengajuan(c *gin.Context) {
	// Pengecekan keamanan ganda: Pastikan hanya Admin yang bisa akses
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengubah status"})
		return
	}

	// Ambil ID Pengajuan dari URL (misal: /api/pengajuan/1)
	pengajuanID := c.Param("id")

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pengajuan models.PengajuanSurat
	// Cari data pengajuan di database berdasarkan ID
	if err := config.DB.First(&pengajuan, pengajuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pengajuan tidak ditemukan"})
		return
	}

	// Update status dan keterangan
	pengajuan.Status = input.Status
	pengajuan.Keterangan = input.Keterangan

	if err := config.DB.Save(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pengajuan berhasil diperbarui",
		"data":    pengajuan,
	})
}

func HapusPengajuan(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")

	var pengajuan models.PengajuanSurat
	// Cari data pengajuan berdasarkan ID
	if err := config.DB.First(&pengajuan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pengajuan tidak ditemukan"})
		return
	}

	// LOGIKA KEAMANAN:
	// 1. Jika Admin, bebas hapus siapapun.
	// 2. Jika Warga, hanya bisa hapus jika UserID di database cocok dengan UserID di token JWT-nya.
	if role != "admin" && uint(userID.(float64)) != pengajuan.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki izin untuk menghapus pengajuan ini"})
		return
	}

	// Eksekusi penghapusan (GORM akan melakukan Soft Delete secara default karena kita pakai gorm.Model)
	if err := config.DB.Delete(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengajuan berhasil dihapus/dibatalkan"})
}