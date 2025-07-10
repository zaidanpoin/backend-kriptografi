package Model

import (
	"kriptografi-zaidaan/Database"

	"gorm.io/gorm"
)

type Dekrip struct {
	gorm.Model
	FileName        string `gorm:"type:varchar(255)" json:"file_name" binding:"required" form:"file_name"`
	FileSize        string `json:"file_size" binding:"required" form:"file_size"`
	FileType        string `gorm:"type:varchar(255)" json:"file_type" binding:"required" form:"file_type"`
	FilePath        string `gorm:"type:varchar(255)" json:"file_path" binding:"required" form:"file_path"`
	FileHash        string `gorm:"type:varchar(255)" json:"file_hash" binding:"required" form:"file_hash"`
	FileDate        string `gorm:"type:varchar(255)" json:"file_date" binding:"required" form:"file_date"`
	FileEncryptedBy string `gorm:"type:varchar(255)" json:"file_encrypted_by" binding:"required" form:"file_encrypted_by"`
	FileStatus      string `gorm:"type:varchar(255)" json:"file_status" binding:"required" form:"file_status"`
	FileKey         string `gorm:"type:varchar(255)" json:"file_key" binding:"required" form:"file_key"`
	Key             string `gorm:"type:varchar(255)" json:"custom_key" binding:"required" form:"custom_key"`
	ExecutionTime   string `gorm:"type:varchar(255)" json:"execution_time" binding:"required" form:"execution_time"`
	EnkripID        uint   `json:"enkrip_id" binding:"required" form:"enkrip_id"`
	Enkrip          Enkrip `gorm:"foreignKey:EnkripID"`
}

func (d *Dekrip) SaveDataDekrip() (*Dekrip, error) {
	err := Database.Database.Create(&d).Error
	if err != nil {
		return nil, err
	}
	return d, nil
}
func (d *Dekrip) GetAllDataDekrip() ([]Dekrip, error) {
	var dekrip []Dekrip
	err := Database.Database.Preload("Enkrip").Find(&dekrip).Error
	if err != nil {
		return nil, err
	}
	return dekrip, nil
}
func (d *Dekrip) GetDataDekripByID(id uint) (*Dekrip, error) {
	var dekrip Dekrip
	err := Database.Database.Where("id = ?", id).First(&dekrip).Error
	if err != nil {
		return nil, err
	}
	return &dekrip, nil
}
func (d *Dekrip) DeleteDataDekrip(id uint) error {
	var dekrip Dekrip
	err := Database.Database.Unscoped().Where("id = ?", id).First(&dekrip).Error
	if err != nil {
		return err
	}

	err = Database.Database.Unscoped().Delete(&dekrip).Error
	if err != nil {
		return err
	}

	return nil
}
