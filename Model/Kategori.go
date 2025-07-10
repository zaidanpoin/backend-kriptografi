package Model

import (
	"kriptografi-zaidaan/Database"

	"gorm.io/gorm"
)

type Kategori struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"primaryKey"`
	NamaKategori string `json:"nama_kategori" gorm:"unique;not null"`
}

func (k *Kategori) Save(db *gorm.DB) (*Kategori, error) {
	if err := db.Create(k).Error; err != nil {
		return nil, err
	}
	return k, nil
}

func FindKategoriByID(db *gorm.DB, id uint) (*Kategori, error) {
	var kategori Kategori
	if err := db.First(&kategori, id).Error; err != nil {
		return nil, err
	}
	return &kategori, nil
}

func UpdateKategori(db *gorm.DB, id uint, updatedKategori Kategori) (*Kategori, error) {

	var kategori Kategori
	if err := db.First(&kategori, id).Error; err != nil {
		return nil, err
	}

	kategori.NamaKategori = updatedKategori.NamaKategori

	if err := db.Save(&kategori).Error; err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (k *Kategori) GetAllKategori() ([]Kategori, error) {
	var kategoris []Kategori
	if err := Database.Database.Find(&kategoris).Error; err != nil {
		return nil, err
	}
	return kategoris, nil
}

func DeleteKategori(db *gorm.DB, id uint) error {
	var kategori Kategori
	if err := db.First(&kategori, id).Error; err != nil {
		return err
	}
	if err := db.Delete(&kategori).Error; err != nil {
		return err
	}
	return nil
}
