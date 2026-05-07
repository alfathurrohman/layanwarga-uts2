package routes

import (
	"net/http"

	"layanwarga/controllers"
	"layanwarga/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// ==========================================
	// 1. RUTE HALAMAN WEB (FRONTEND / ANTARMUKA)
	// ==========================================
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})

	r.GET("/admin/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_dashboard.html", nil)
	})

	// ==========================================
	// 2. RUTE API (BACKEND / PENGOLAHAN DATA)
	// ==========================================
	
	// Grup rute Auth (Terbuka untuk umum)
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Grup rute khusus API yang membutuhkan Login (Dilindungi Middleware)
	api := r.Group("/api")
	api.Use(middlewares.RequireAuth)
	{
		api.POST("/pengajuan", controllers.BuatPengajuan)
		api.GET("/pengajuan", controllers.GetPengajuan)
		api.PUT("/pengajuan/:id/status", controllers.UpdateStatusPengajuan)
		api.DELETE("/pengajuan/:id", controllers.HapusPengajuan)
	}
}