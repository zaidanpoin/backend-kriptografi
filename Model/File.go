package Model

import (
	"kriptografi-zaidaan/Database"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileName       string `json:"file_name" binding:"required" form:"file_name" gorm:"type:varchar(255)"`
	FileSize       string `json:"file_size" binding:"required" form:"file_size"`
	FileType       string `json:"file_type" binding:"required" form:"file_type" gorm:"type:varchar(100)"`
	FilePath       string `json:"file_path" binding:"required" form:"file_path" gorm:"type:varchar(255)"`
	FileHash       string `json:"file_hash" binding:"required" form:"file_hash" gorm:"type:varchar(255)"`
	FileDate       string `json:"file_date" binding:"required" form:"file_date" gorm:"type:varchar(50)"`
	FileUploadedBy string `json:"file_uploaded_by" binding:"required" form:"file_uploaded_by" gorm:"type:varchar(100)"`
	FileStatus     string `json:"file_status" binding:"required" form:"file_status" gorm:"type:varchar(50)"`
	UserID         uint   `json:"user_id" binding:"required" form:"user_id"`
	User           User   `gorm:"foreignKey:UserID"`
}

func (f *File) GetDataFile(id string) ([]File, error) {
	if id == "" {
		var files []File
		err := Database.Database.Find(&files).Error
		if err != nil {
			return nil, err
		}
		return files, nil
	} else {
		var file File
		err := Database.Database.Where("id = ?", id).First(&file).Error
		if err != nil {
			return nil, err
		}
		return []File{file}, nil
	}

}

func (f *File) GetFileByID(id string) (*File, error) {
	var file File
	err := Database.Database.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *File) UpdateFile(id string) (*File, error) {
	var file File
	err := Database.Database.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}

	err = Database.Database.Model(&file).Updates(f).Error
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (f *File) DeleteFile(id string) error {
	var file File
	err := Database.Database.Where("id = ?", id).First(&file).Error
	if err != nil {
		return err
	}

	err = Database.Database.Unscoped().Delete(&file).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Save() (*File, error) {
	var err error

	// check duplicate file name
	err = Database.Database.Create(&f).Error
	if err != nil {
		return &File{}, err

	}

	return f, nil

}
