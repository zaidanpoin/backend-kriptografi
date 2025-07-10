package main

import (
	"kriptografi-zaidaan/Database"
	"kriptografi-zaidaan/Model"
	"kriptografi-zaidaan/Router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	Router.ServeApps()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadDatabase() {

	Database.Connect()
	Database.Database.AutoMigrate(&Model.User{})
	Database.Database.AutoMigrate(&Model.File{})
	Database.Database.AutoMigrate(&Model.Enkrip{})
	Database.Database.AutoMigrate(&Model.Dekrip{})
	Database.Database.AutoMigrate(&Model.Barang{})
	Database.Database.AutoMigrate(&Model.Kategori{})
	Database.Database.AutoMigrate(&Model.BarangMasuk{})
	Database.Database.AutoMigrate(&Model.BarangKeluar{})

}
