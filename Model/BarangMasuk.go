package Model

import (
	"kriptografi-zaidaan/Database"
	"time"

	"gorm.io/gorm"
)

type BarangMasuk struct {
	KodeBarangMasuk string `json:"kode_barang_masuk" gorm:"primaryKey"`
	KodeBarang      string `json:"kode_barang"`
	Barang          Barang `json:"barang" gorm:"foreignKey:KodeBarang;references:KodeBarang"`
	NamaBarang      string `json:"nama_barang"`
	JumlahMasuk     int    `json:"jumlah_masuk"`
	TanggalMasuk    string `json:"tanggal_masuk"`
	Deskripsi       string `json:"deskripsi"`
	Asal            string `json:"asal"`
	Admin           string `json:"admin"`
}

func (b *BarangMasuk) CetakBarangMasuk(asal, tanggalAwal, tanggalAkhir string) ([]BarangMasuk, error) {
	var barangMasuk []BarangMasuk
	if err := Database.Database.Where("asal = ? AND tanggal_masuk BETWEEN ? AND ?", asal, tanggalAwal, tanggalAkhir).Find(&barangMasuk).Error; err != nil {
		return nil, err
	}
	return barangMasuk, nil
}

func (b *BarangMasuk) BeforeCreate(tx *gorm.DB) (err error) {
	if b.KodeBarangMasuk == "" {
		// Generate a unique KodeBarangMasuk using timestamp and KodeBarang
		b.KodeBarangMasuk = "BM-" + b.KodeBarang + "-" + time.Now().Format("20060102150405.000")
	}
	if b.TanggalMasuk == "" {
		b.TanggalMasuk = time.Now().Format("2006-01-02")
	}
	return nil
}

func (b *BarangMasuk) Save() (*BarangMasuk, error) {
	if err := Database.Database.Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}
func (b *BarangMasuk) GetAllBarangMasuk() ([]BarangMasuk, error) {
	var barangMasuk []BarangMasuk
	if err := Database.Database.Find(&barangMasuk).Error; err != nil {
		return nil, err
	}
	return barangMasuk, nil
}

func (b *BarangMasuk) GetBarangMasukByKode(kodeBarangMasuk string) (*BarangMasuk, error) {
	var barangMasuk BarangMasuk
	if err := Database.Database.Where("kode_barang_masuk = ?", kodeBarangMasuk).First(&barangMasuk).Error; err != nil {
		return nil, err
	}
	return &barangMasuk, nil
}

func (b *BarangMasuk) DeleteBarangMasuk(kodeBarangMasuk string) error {
	if err := Database.Database.Delete(&BarangMasuk{}, "kode_barang_masuk = ?", kodeBarangMasuk).Error; err != nil {
		return err
	}
	return nil
}
