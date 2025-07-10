package Controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"kriptografi-zaidaan/Model"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {

	id := c.Query("id")
	var file Model.File
	files, err := file.GetDataFile(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve file"})
		return
	}

	if files == nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.JSON(200, files)
}

func GetFileByID(c *gin.Context) {
	id := c.Param("id")
	var file Model.File
	files, err := file.GetFileByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve file"})
		return
	}

	if files == nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.JSON(200, files)
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file from request"})
		return
	}

	// Validate file size (example: max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(400, gin.H{"error": "File size exceeds the limit of 10MB"})
		return
	}
	// Validate file type (example: only allow .txt files)
	allowedTypes := map[string]bool{
		"application/pdf":          true,
		"application/vnd.ms-excel": true, // .xls
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,

		"application/vnd.ms-excel.sheet.macroEnabled.12": true, // .xlsm
		"application/octet-stream":                       true,
		// sometimes .xlsx uploads as this

		"image/jpeg": true, // .jpg
		"image/png":  true, // .png
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Invalid file type. Only PDF, Excel, and Word files are allowed"})
		return
	}

	ext := filepath.Ext(file.Filename)

	// Hash the file name
	hash := sha256.Sum256([]byte(file.Filename + time.Now().String()))
	hashedFileName := hex.EncodeToString(hash[:]) + ext

	// Save the file to the server
	filePath := "./file-plainteks/" + hashedFileName
	url := "http://localhost:8080/file-plainteks/" + hashedFileName
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	sizeInBytes := file.Size
	var ukuranFile string
	if sizeInBytes >= 1024*1024 {
		sizeInMB := float64(sizeInBytes) / (1024 * 1024)
		ukuranFile = fmt.Sprintf("%.2f MB", sizeInMB)
	} else {
		sizeInKB := float64(sizeInBytes) / 1024
		ukuranFile = fmt.Sprintf("%.2f KB", sizeInKB)
	}

	upload := c.PostForm("file_uploaded_by")

	fileMetadata := Model.File{
		FileName:       file.Filename,
		FileSize:       ukuranFile,
		FileType:       ext,
		FilePath:       url,
		FileHash:       hashedFileName,
		FileDate:       time.Now().Format("2006-01-02"),
		FileUploadedBy: upload, // Replace with actual user information
		FileStatus:     "uploaded",
		UserID:         1,
	}

	_, err = fileMetadata.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file metadata to database"})
		return
	}

	c.JSON(200, gin.H{
		"message":        "File successfully uploaded and saved to database",
		"fileName":       fileMetadata.FileName,
		"fileSize":       fileMetadata.FileSize,
		"fileType":       fileMetadata.FileType,
		"filePath":       fileMetadata.FilePath,
		"fileDate":       fileMetadata.FileDate,
		"fileUploadedBy": fileMetadata.FileUploadedBy,
		"fileStatus":     fileMetadata.FileStatus,
	})

}

func DeleteFile(c *gin.Context) {
	id := c.Param("id")
	var file Model.File
	err := file.DeleteFile(id)

	if err == nil {
		fileData, _ := file.GetDataFile(id)
		if len(fileData) > 0 {
			filePath := "./file-plainteks/" + fileData[0].FileHash
			_ = os.Remove(filePath)
		}
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(200, gin.H{"message": "File successfully deleted"})
}
