package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth adalah fungsi untuk melindungi endpoint (hanya yang punya token yang bisa akses)
func RequireAuth(c *gin.Context) {
	// 1. Ambil tiket dari header request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak, token tidak ditemukan"})
		c.Abort() // Hentikan proses
		return
	}

	// 2. Bersihkan string tiket (biasanya formatnya "Bearer token_acak_panjang")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// 3. Validasi keaslian tiket menggunakan kunci rahasia di .env
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kadaluarsa"})
		c.Abort()
		return
	}

	// 4. Jika valid, ambil data user_id dari dalam tiket, dan simpan sementara di Context
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Simpan ID agar controller tahu siapa yang sedang mengajukan surat
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next() // Izinkan masuk ke controller
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca data token"})
		c.Abort()
		return
	}
}