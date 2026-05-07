package main

import (
	"fmt"
	"layanwarga/config"
	"layanwarga/models" // <-- Ini baris yang tadi terlewat
	"layanwarga/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Memulai aplikasi LayanWarga...")
	
	// 1. Koneksi ke Database
	config.ConnectDB()

	// 2. Koneksi ke AWS S3
	config.ConnectAWS()

	// 3. Baris ajaib untuk menyulap akun Budi menjadi Admin secara otomatis
	config.DB.Model(&models.User{}).Where("email = ?", "budi@gmail.com").Update("role", "admin")

	// 4. Inisialisasi Web Server menggunakan Gin
	r := gin.Default()

	// 5. Membaca folder templates untuk merender HTML
	r.LoadHTMLGlob("templates/*")

	// 6. Daftarkan rute-rute
	routes.SetupRoutes(r)

	// 7. Jalankan server
	fmt.Println("Server web berjalan di http://localhost:8080")
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Gagal menjalankan server:", err)
	}
}