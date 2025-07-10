package Router

import (
	"fmt"
	"kriptografi-zaidaan/Controller"
	"kriptografi-zaidaan/Middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ServeApps() {

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Static("/file-plainteks", "./file-plainteks")
	router.Static("/file-enkrip", "./file-enkrip")
	router.Static("/file-dekrip", "./file-dekrip")
	router.Static("/uploads", "./uploads")

	authRoutes := router.Group("/auth")
	{
		AuthRoutes(authRoutes)
	}

	inventoryGroup := router.Group("/inventory")
	{
		InventoryRoutes(inventoryGroup)
	}

	fileRoutes := router.Group("/file")

	{
		fileRoutes.POST("/upload", Controller.UploadFile)
		fileRoutes.GET("/", Controller.GetFile)
		fileRoutes.GET("/detail/:id", Controller.GetFileByID)
		fileRoutes.DELETE("/:id", Controller.DeleteFile)

	}

	// 	c.JSON(200, gin.H{
	// 		"message": "Hello World",
	// 	})
	// })

	router.POST("/enkrip/:id", Controller.EncryptHandler)
	router.GET("/enkrip/data", Controller.GetAllDataEncrypt)
	router.GET("/enkrip/:id", Controller.GetDataEncryptByID)
	router.DELETE("/enkrip/:id", Controller.DeleteDataEncrypt)

	router.POST("/dekrip/:id", Controller.DecryptHandler)
	router.GET("/dekrip/data", Controller.GetAllDataDecrypt)
	router.DELETE("/dekrip/:id", Controller.DeleteDataDecrypt)

	router.Run(":8080")
	fmt.Println("Server is running on port 8080")
}

func InventoryRoutes(router *gin.RouterGroup) {
	router.GET("/barang", Controller.GetAllBarang)
	router.GET("/barang/:kode_barang", Controller.GetBarangByKode)
	router.POST("/barang", Controller.CreateBarang)
	router.PUT("/barang/:kode_barang", Controller.UpdateBarang)
	router.DELETE("/barang/:kode_barang", Controller.DeleteBarang)
	router.GET("/barang/kategori", Controller.GetKategori)
	router.GET("/cetak-barang-masuk/:asal/:tanggal_awal/:tanggal_akhir", Controller.CetakBarang)
	router.GET("/barang-masuk", Controller.GetAllBarangMasuk)
	router.POST("/barang-masuk", Controller.InsertBarangMasuk)
	router.DELETE("/barang-masuk/:kode_barang_masuk", Controller.DeleteBarangMasuk)

	router.GET("/barang-keluar", Controller.GetAllBarangKeluar)
	router.POST("/barang-keluar", Controller.InsertBarangKeluar)
	router.DELETE("/barang-keluar/:kode_barang_keluar", Controller.DeleteBarangKeluar)

}

func AuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", Controller.Register)
	router.POST("/login", Controller.Login)
	router.GET("/data", Controller.GetAllUsers)
	router.GET("/data/:username", Controller.GetUserByUsername)
	router.PUT("/update/:username", Controller.UpdateUser)
	router.DELETE("/delete/:username", Controller.DeleteUser)
	router.GET("/verify-token", Middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Token valid"})
	})
}
