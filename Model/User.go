package Model

import (
	"kriptografi-zaidaan/Database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `form:"username" json:"username" binding:"required" gorm:"type:varchar(255);unique"`
	Password string `form:"password" json:"password" binding:"required" gorm:"type:varchar(255)"`
	Email    string `form:"email" json:"email" binding:"required" gorm:"type:varchar(255);unique"`
	Name     string `form:"name" json:"name" binding:"required" gorm:"type:varchar(255)"`
	Role     string `form:"role" json:"role" binding:"required" gorm:"type:varchar(100)"`
	Alamat   string `form:"alamat" json:"alamat" binding:"required" gorm:"type:varchar(255)"`
	Telp     string `form:"telp" json:"telp" gorm:"type:varchar(50)"`
}

func (u *User) BeforeSave(*gorm.DB) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Update() (*User, error) {
	var err error

	// check if user exists
	err = Database.Database.Model(&u).Where("id = ?", u.ID).Updates(u).Error
	if err != nil {
		return &User{}, err
	}

	// return updated user
	return u, nil
}

func (u *User) Save() (*User, error) {
	var err error

	// check duplicate username
	err = Database.Database.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}
func GetAllUsers() ([]User, error) {
	var users []User
	err := Database.Database.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUserByUsername(username string) error {
	var user User
	err := Database.Database.Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}

	err = Database.Database.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := Database.Database.Where("username = ?", username).First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}
