package Controller

import (
	"kriptografi-zaidaan/Model"

	"github.com/gin-gonic/gin"
)

func GetAllBarangKeluar(c *gin.Context) {
	var barangKeluar Model.BarangKeluar

	data, err := barangKeluar.GetAllBarangKeluar()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang keluar"})
		return
	}
	c.JSON(200, gin.H{"message": "Berhasil mengambil data barang keluar", "data": data})
}

func InsertBarangKeluar(c *gin.Context) {
	var barangKeluar Model.BarangKeluar

	if err := c.ShouldBindJSON(&barangKeluar); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	if barangKeluar.KodeBarang == "" || barangKeluar.NamaBarang == "" || barangKeluar.JumlahKeluar <= 0 {
		c.JSON(400, gin.H{"error": "KodeBarang, NamaBarang, and JumlahKeluar are required"})
		return
	}

	savedBarangKeluar, err := barangKeluar.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save barang keluar"})
		return
	}

	// Kurangi stok di tabel barang
	var barang Model.Barang
	barangResult, err := barang.GetBarangByKode(barangKeluar.KodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to find barang for updating stock"})
		return
	}
	if barangResult == nil {
		c.JSON(404, gin.H{"error": "Barang not found"})
		return
	}

	if barangKeluar.JumlahKeluar > barangResult.Stok {
		c.JSON(400, gin.H{"error": "JumlahKeluar melebihi stok barang yang tersedia"})
		return
	}

	total := barangResult.Stok - barangKeluar.JumlahKeluar
	if _, err := barang.UpdateStok(barangResult.KodeBarang, total); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update barang stock"})
		return
	}

	c.JSON(200, gin.H{"message": "Berhasil menyimpan data barang keluar dan update stok", "data": savedBarangKeluar})
}

func DeleteBarangKeluar(c *gin.Context) {
	kodeBarangKeluar := c.Param("kode_barang_keluar")
	var barangKeluar Model.BarangKeluar

	if err := barangKeluar.DeleteBarangKeluar(kodeBarangKeluar); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete barang keluar"})
		return
	}
	c.JSON(200, gin.H{"message": "Berhasil menghapus data barang keluar"})
}
