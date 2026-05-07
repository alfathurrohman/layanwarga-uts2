package controllers

import (
	"net/http"
	"os"
	"time"

	"layanwarga/config"
	"layanwarga/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struct untuk menerima data dari user (request body)
type RegisterInput struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	NIK         string `json:"nik" binding:"required,min=16,max=16"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Fungsi Registrasi
func Register(c *gin.Context) {
	var input RegisterInput
	// Mengecek apakah data yang dikirim user sesuai dengan struct RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hashing Password menggunakan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengamankan password"})
		return
	}

	// Menyusun data user baru untuk dimasukkan ke database
	newUser := models.User{
		NamaLengkap: input.NamaLengkap,
		NIK:         input.NIK,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Role:        "warga", // Default otomatis menjadi warga
	}

	// Menyimpan ke Database
	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email atau NIK sudah terdaftar!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil, silakan login!"})
}

// Fungsi Login
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Mencari user berdasarkan Email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Membandingkan password yang diinput dengan password hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Membuat JWT Token jika login berhasil (berlaku 24 jam)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Menandatangani token menggunakan JWT_SECRET dari file .env
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
	})
}