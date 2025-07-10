package Model

import (
	"kriptografi-zaidaan/Database"

	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	KodeBarang string   `form:"kode_barang" json:"kode_barang" gorm:"primaryKey"`
	NamaBarang string   `form:"nama_barang" json:"nama_barang"`
	Stok       int      `form:"stok" json:"stok"`
	Gambar     string   `form:"gambar" json:"gambar"`
	Url        string   `form:"url" json:"url"`
	Deskripsi  string   `form:"deskripsi" json:"deskripsi"`
	KategoriID *uint    `json:"kategori_id"`
	Kategori   Kategori `json:"kategori" gorm:"foreignKey:KategoriID"`
}

// BeforeCreate GORM hook to set default KodeBarang if empty
func (b *Barang) BeforeCreate(tx *gorm.DB) (err error) {
	if b.KodeBarang == "" {
		rand.Seed(time.Now().UnixNano())
		b.KodeBarang = fmt.Sprintf("KB-%06d", rand.Intn(1000000))
	}
	return nil
}

func (b *Barang) Save() (*Barang, error) {
	if err := Database.Database.Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}
func FindBarangByKode(db *gorm.DB, kodeBarang string) (*Barang, error) {
	var barang Barang
	if err := db.First(&barang, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (b *Barang) UpdateBarang(kodeBarang string) (*Barang, error) {
	var barang Barang
	if err := Database.Database.First(&barang, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}
	barang.NamaBarang = b.NamaBarang
	barang.Deskripsi = b.Deskripsi
	barang.Gambar = b.Gambar
	barang.Url = b.Url
	barang.KategoriID = b.KategoriID
	if err := Database.Database.Save(&barang).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (b *Barang) UpdateStok(kodeBarang string, stok int) (*Barang, error) {
	var barang Barang
	if err := Database.Database.First(&barang, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}
	barang.Stok = stok
	if err := Database.Database.Save(&barang).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (b *Barang) DeleteBarang(kodeBarang string) error {
	var barang Barang
	if err := Database.Database.First(&barang, "kode_barang = ?", kodeBarang).Error; err != nil {
		return err
	}
	if err := Database.Database.Delete(&barang).Error; err != nil {
		return err
	}
	return nil
}

func (b *Barang) GetBarangByKode(kodeBarang string) (*Barang, error) {
	var barang Barang
	fmt.Println("Fetching barang with kode:", kodeBarang)
	if err := Database.Database.Preload("Kategori").First(&barang, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}
	return &barang, nil
}

func (b *Barang) GetAllBarang() ([]Barang, error) {
	var barangs []Barang
	if err := Database.Database.Preload("Kategori").Find(&barangs).Error; err != nil {
		return barangs, err
	}

	return barangs, nil
}
