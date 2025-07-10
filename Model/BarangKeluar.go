package Model

import (
	"kriptografi-zaidaan/Database"
	"time"

	"gorm.io/gorm"
)

type BarangKeluar struct {
	KodeBarangKeluar string `json:"kode_barang_keluar" gorm:"primaryKey"`
	KodeBarang       string `json:"kode_barang"`
	Barang           Barang `json:"barang" gorm:"foreignKey:KodeBarang;references:KodeBarang"`
	NamaBarang       string `json:"nama_barang"`
	JumlahKeluar     int    `json:"jumlah_keluar"`
	TanggalKeluar    string `json:"tanggal_keluar"`
	Deskripsi        string `json:"deskripsi"`
	Manufaktur       string `json:"manufaktur"`
}

func (b *BarangKeluar) BeforeCreate(tx *gorm.DB) (err error) {
	if b.KodeBarangKeluar == "" {

		b.KodeBarangKeluar = "BK-" + b.KodeBarang + "-" + time.Now().Format("20060102150405.000")
	}
	if b.TanggalKeluar == "" {
		b.TanggalKeluar = time.Now().Format("2006-01-02")
	}
	return nil
}

func (b *BarangKeluar) Save() (*BarangKeluar, error) {
	if err := Database.Database.Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}

func FindBarangKeluarByKode(e, kodeBarangKeluar string) (*BarangKeluar, error) {
	var barangKeluar BarangKeluar
	if err := Database.Database.First(&barangKeluar, "kode_barang_keluar = ?", kodeBarangKeluar).Error; err != nil {
		return nil, err
	}
	return &barangKeluar, nil
}

func (b *BarangKeluar) UpdateBarangKeluar(kodeBarangKeluar string, updatedBarangKeluar BarangKeluar) (*BarangKeluar, error) {
	var barangKeluar BarangKeluar
	if err := Database.Database.First(&barangKeluar, "kode_barang_keluar = ?", kodeBarangKeluar).Error; err != nil {
		return nil, err
	}

	barangKeluar.KodeBarang = updatedBarangKeluar.KodeBarang
	barangKeluar.NamaBarang = updatedBarangKeluar.NamaBarang
	barangKeluar.JumlahKeluar = updatedBarangKeluar.JumlahKeluar
	barangKeluar.TanggalKeluar = updatedBarangKeluar.TanggalKeluar
	barangKeluar.Deskripsi = updatedBarangKeluar.Deskripsi

	if err := Database.Database.Save(&barangKeluar).Error; err != nil {
		return nil, err
	}
	return &barangKeluar, nil
}

func (b *BarangKeluar) GetAllBarangKeluar() ([]BarangKeluar, error) {
	var barangKeluar []BarangKeluar
	if err := Database.Database.Find(&barangKeluar).Error; err != nil {
		return nil, err
	}
	return barangKeluar, nil
}

func (b *BarangKeluar) DeleteBarangKeluar(kodeBarangKeluar string) error {
	var barangKeluar BarangKeluar
	if err := Database.Database.First(&barangKeluar, "kode_barang_keluar = ?", kodeBarangKeluar).Error; err != nil {
		return err
	}
	if err := Database.Database.Delete(&barangKeluar).Error; err != nil {
		return err
	}
	return nil
}
