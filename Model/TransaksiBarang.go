package Model

import (
	"time"

	"gorm.io/gorm"
)

type Transaksi struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	KodeBarang     uint           `gorm:"not null" json:"barang_id"`
	Barang         Barang         `gorm:"foreignKey:KodeBarang;references:KodeBarang" json:"barang"`
	Jumlah         int            `gorm:"not null" json:"jumlah"`
	JenisTransaksi string         `gorm:"type:enum('masuk','keluar');not null" json:"jenis_transaksi"`
	Tanggal        time.Time      `gorm:"not null" json:"tanggal"`
	Keterangan     string         `gorm:"type:text" json:"keterangan,omitempty"`
	PetugasID      *uint          `json:"petugas_id,omitempty"` // opsional
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
