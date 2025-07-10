package Controller

import (
	"crypto/md5"
	"fmt"
	"kriptografi-zaidaan/Model"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllBarang(c *gin.Context) {
	var barangs Model.Barang

	data, err := barangs.GetAllBarang()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang"})
		return
	}
	c.JSON(200, gin.H{
		"message":     "Berhasil mengambil data barang",
		"data":        data,
		"jumlah_data": len(data),
	})
}

func GetBarangByKode(c *gin.Context) {
	kodeBarang := c.Param("kode_barang")
	var barang Model.Barang

	data, err := barang.GetBarangByKode(kodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang"})
		return
	}

	if data == nil {
		c.JSON(404, gin.H{"error": "Barang not found"})
		return
	}

	c.JSON(200, data)
}

func GetKategori(c *gin.Context) {
	var kategori Model.Kategori
	data, err := kategori.GetAllKategori()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve kategori"})
		return
	}
	c.JSON(200, gin.H{"message": "Berhasil mengambil data kategori", "data": data})
}

func CreateBarang(c *gin.Context) {
	var barang Model.Barang
	file, err := c.FormFile("gambar")
	if err != nil {
		c.JSON(400, gin.H{"error": "Gambar is required"})
		return
	}

	allowedExtensions := map[string]bool{
		".webp": true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		c.JSON(400, gin.H{"error": "Invalid image format. Allowed formats: .webp, .jpg, .jpeg, .png"})
		return
	}

	md5Hash := md5.New()
	md5Hash.Write([]byte(file.Filename + fmt.Sprintf("%d", time.Now().UnixNano())))
	hashedFilename := fmt.Sprintf("%x", md5Hash.Sum(nil))
	file.Filename = hashedFilename + ext

	uploadPath := "./uploads/" + file.Filename
	url := "http://localhost:8080/uploads/" + file.Filename

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	barang.KodeBarang = c.PostForm("kode_barang")
	barang.NamaBarang = c.PostForm("nama_barang")
	barang.Deskripsi = c.PostForm("deskripsi")
	barang.Gambar = file.Filename
	barang.Url = url

	stokStr := c.PostForm("stok") // <- pastikan form input-nya juga pakai nama ini
	stokInt, err := strconv.Atoi(stokStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Stok harus berupa angka"})
		return
	}
	barang.Stok = stokInt

	kategoriID, _ := strconv.Atoi(c.PostForm("kategori_id"))
	kategoriIDUint := uint(kategoriID)
	barang.KategoriID = &kategoriIDUint

	_, err = barang.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create barang"})
		return
	}
	c.JSON(201, gin.H{"message": "Barang created successfully", "data": barang})
}

func DeleteBarang(c *gin.Context) {
	kodeBarang := c.Param("kode_barang")
	var barang Model.Barang

	data, err := barang.GetBarangByKode(kodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang"})
		return
	}
	if data == nil {
		c.JSON(404, gin.H{"error": "Barang not found"})
		return
	}

	// Hapus file gambar jika ada
	if data.Gambar != "" {
		filePath := "./uploads/" + data.Gambar
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			c.JSON(500, gin.H{"error": "Failed to delete image file"})
			return
		}
	}

	err1 := barang.DeleteBarang(kodeBarang)
	if err1 != nil {
		c.JSON(500, gin.H{"error": "Failed to delete barang"})
		return
	}

	c.JSON(200, gin.H{"message": "Barang deleted successfully"})
}

func UpdateBarang(c *gin.Context) {
	kodeBarang := c.Param("kode_barang")
	var barang Model.Barang

	// Cek apakah barang ada
	data, err := barang.GetBarangByKode(kodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang"})
		return
	}
	if data == nil {
		c.JSON(404, gin.H{"error": "Barang not found"})
		return
	}

	// Ambil data dari form
	namaBarang := c.PostForm("nama_barang")
	deskripsi := c.PostForm("deskripsi")
	stokStr := c.PostForm("stok")
	kategoriIDStr := c.PostForm("kategori_id")

	// Validasi stok
	stokInt, err := strconv.Atoi(stokStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Stok harus berupa angka"})
		return
	}

	kategoriID, err := strconv.Atoi(kategoriIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Kategori ID harus berupa angka"})
		return
	}
	kategoriIDUint := uint(kategoriID)

	// Update data barang
	data.NamaBarang = namaBarang
	data.Deskripsi = deskripsi
	data.Stok = stokInt
	data.KategoriID = &kategoriIDUint

	// Cek apakah ada file gambar baru
	file, err := c.FormFile("gambar")
	if err == nil {
		// Validasi ekstensi gambar
		allowedExtensions := map[string]bool{
			".webp": true,
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		ext := filepath.Ext(file.Filename)
		if !allowedExtensions[ext] {
			c.JSON(400, gin.H{"error": "Invalid image format. Allowed formats: .webp, .jpg, .jpeg, .png"})
			return
		}

		// Hapus gambar lama jika ada
		if data.Gambar != "" {
			oldFilePath := "./uploads/" + data.Gambar
			if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
				c.JSON(500, gin.H{"error": "Failed to delete old image file"})
				return
			}
		}

		// Simpan gambar baru
		md5Hash := md5.New()
		md5Hash.Write([]byte(file.Filename + fmt.Sprintf("%d", time.Now().UnixNano())))
		hashedFilename := fmt.Sprintf("%x", md5Hash.Sum(nil))
		file.Filename = hashedFilename + ext

		uploadPath := "./uploads/" + file.Filename
		url := "http://localhost:8080/uploads/" + file.Filename

		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save uploaded file"})
			return
		}

		data.Gambar = file.Filename
		data.Url = url
	}

	_, err = data.UpdateBarang(kodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update barang"})
		return
	}

	c.JSON(200, gin.H{"message": "Barang updated successfully", "data": data})
}
