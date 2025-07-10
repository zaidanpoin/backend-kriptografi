package Model

import (
	"kriptografi-zaidaan/Database"

	"gorm.io/gorm"
)

type Enkrip struct {
	gorm.Model
	FileName        string `json:"file_name" binding:"required" form:"file_name" gorm:"type:varchar(255)"`
	FileSize        string `json:"file_size" binding:"required" form:"file_size"`
	FileType        string `json:"file_type" binding:"required" form:"file_type" gorm:"type:varchar(255)"`
	FilePath        string `json:"file_path" binding:"required" form:"file_path" gorm:"type:varchar(255)"`
	FileHash        string `json:"file_hash" binding:"required" form:"file_hash" gorm:"type:varchar(255)"`
	FileDate        string `json:"file_date" binding:"required" form:"file_date" gorm:"type:varchar(255)"`
	FileEncryptedBy string `json:"file_encrypted_by" binding:"required" form:"file_encrypted_by" gorm:"type:varchar(255)"`
	FileStatus      string `json:"file_status" binding:"required" form:"file_status" gorm:"type:varchar(255)"`
	FileKey         string `json:"file_key" binding:"required" form:"file_key" gorm:"type:varchar(255)"`
	Key             string `json:"custom_key" binding:"required" form:"custom_key" gorm:"type:varchar(255)"`
	FileID          uint   `json:"file_id" binding:"required" form:"file_id"`
	File            File   `gorm:"foreignKey:FileID"`
	Excecution_time string `json:"execution_time" binding:"required" form:"execution_time" gorm:"type:varchar(255)"`
}

func (e *Enkrip) SaveDataEnkrip() (*Enkrip, error) {
	err := Database.Database.Create(&e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Enkrip) UpdateDataEnkrip(id uint) (*Enkrip, error) {
	var enkrip Enkrip
	err := Database.Database.Where("id = ?", id).First(&enkrip).Error
	if err != nil {
		return nil, err
	}

	err = Database.Database.Model(&enkrip).Updates(e).Error
	if err != nil {
		return nil, err
	}

	return &enkrip, nil
}

func (e *Enkrip) GetAllDataEnkrip() ([]Enkrip, error) {

	var enkrip []Enkrip
	err := Database.Database.Preload("File").Find(&enkrip).Error
	if err != nil {
		return nil, err
	}
	return enkrip, nil
}

func (e *Enkrip) GetDataEnkripByID(id uint) (*Enkrip, error) {

	var enkrip Enkrip
	err := Database.Database.Where("id = ?", id).First(&enkrip).Error
	if err != nil {
		return nil, err
	}
	return &enkrip, nil
}

func (e *Enkrip) DeleteDataEnkrip(id uint) error {
	var enkrip Enkrip
	err := Database.Database.Unscoped().Where("id = ?", id).First(&enkrip).Error
	if err != nil {
		return err
	}

	err = Database.Database.Unscoped().Delete(&enkrip).Error
	if err != nil {
		return err
	}

	return nil
}
