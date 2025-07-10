package Controller

import (
	"kriptografi-zaidaan/Model"

	"github.com/gin-gonic/gin"
)

func GetAllBarangMasuk(c *gin.Context) {
	var barangMasuk Model.BarangMasuk

	data, err := barangMasuk.GetAllBarangMasuk()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang masuk"})
		return
	}
	c.JSON(200, gin.H{"message": "Berhasil mengambil data barang masuk", "data": data})
}

func InsertBarangMasuk(c *gin.Context) {
	var barangMasuk Model.BarangMasuk

	if err := c.ShouldBindJSON(&barangMasuk); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	if barangMasuk.KodeBarang == "" || barangMasuk.NamaBarang == "" || barangMasuk.JumlahMasuk <= 0 {
		c.JSON(400, gin.H{"error": "KodeBarang, NamaBarang, and JumlahMasuk are required"})
		return
	}

	savedBarangMasuk, err := barangMasuk.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save barang masuk"})
		return
	}

	// Tambah stok di tabel barang

	var barang Model.Barang
	barangResult, err := barang.GetBarangByKode(barangMasuk.KodeBarang)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to find barang for updating stock"})
		return
	}
	if barangResult == nil {
		c.JSON(404, gin.H{"error": "Barang not found"})
		return
	}
	total := barangResult.Stok + barangMasuk.JumlahMasuk
	if _, err := barang.UpdateStok(barangResult.KodeBarang, total); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update barang stock"})
		return
	}

	c.JSON(200, gin.H{"message": "Berhasil menyimpan data barang masuk dan update stok", "data": savedBarangMasuk})
}

func CetakBarang(c *gin.Context) {

	asal := c.Query("asal")
	tanggalAwal := c.Query("tanggal_awal")
	tanggalAkhir := c.Query("tanggal_akhir")

	var barangMasuk Model.BarangMasuk
	data, err := barangMasuk.CetakBarangMasuk(asal, tanggalAwal, tanggalAkhir)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve barang masuk for cetak"})
		return
	}

	c.JSON(200, gin.H{"message": "Berhasil mengambil data barang masuk untuk cetak", "data": data})

}

func DeleteBarangMasuk(c *gin.Context) {
	kodeBarangMasuk := c.Param("kode_barang_masuk")
	var barangMasuk Model.BarangMasuk

	if err := barangMasuk.DeleteBarangMasuk(kodeBarangMasuk); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete barang masuk"})
		return
	}

	c.JSON(200, gin.H{"message": "Berhasil menghapus data barang masuk"})
}
