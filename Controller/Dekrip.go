package Controller

import (
	"fmt"
	"io/ioutil"
	"kriptografi-zaidaan/Model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteDataDecrypt(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)

	if id == "" {
		c.JSON(400, gin.H{
			"error": "ID file tidak boleh kosong.",
		})
		return
	}

	var dekrip Model.Dekrip
	err = dekrip.DeleteDataDekrip(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{
			"error": "File tidak ditemukan.",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "File berhasil dihapus.",
	})
}

func GetAllDataDecrypt(c *gin.Context) {
	var dekrip Model.Dekrip
	data, err := dekrip.GetAllDataDekrip()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal mengambil data dekripsi.",
		})
		return
	}

	if len(data) == 0 {
		c.JSON(404, gin.H{
			"message": "Tidak ada data dekripsi.",
		})
		return
	}

	c.JSON(200, data)
}

func GetDataDecryptByID(c *gin.Context) {

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)

	if id == "" {
		c.JSON(400, gin.H{
			"error": "ID file tidak boleh kosong.",
		})
		return
	}

	var dekrip Model.Dekrip
	data, err := dekrip.GetDataDekripByID(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{
			"error": "File tidak ditemukan.",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": data,
	})
}

func DecryptHandler(c *gin.Context) {

	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	customKey := c.PostForm("custom_key")

	if id == "" {
		c.JSON(400, gin.H{
			"error": "ID file tidak boleh kosong.",
		})
		return
	}

	var file Model.Enkrip

	data, err := file.GetDataEnkripByID(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{
			"error": "File tidak ditemukan.",
		})
		return
	}
	fmt.Println("Data file:", data)

	if customKey == "" {
		c.JSON(400, gin.H{
			"error": "Custom key tidak boleh kosong.",
		})
		return
	}

	if customKey != data.Key {
		c.JSON(400, gin.H{
			"error": "Custom key tidak valid.",
		})
		return
	}

	key := stringToKey(customKey)

	encryptedData, err := ioutil.ReadFile("./file-enkrip/" + data.FileHash)
	if err != nil {
		panic(err)
	}

	var files Model.File
	fileData, err := files.GetDataFile(strconv.Itoa(int(data.FileID)))
	start := time.Now()
	decryptedData := aesCBCDecrypt(encryptedData, key)
	elapsed := time.Since(start)
	if decryptedData == nil {
		fmt.Println("Error: Kunci dekripsi tidak valid!")
		return
	}
	sizeInBytes := len(decryptedData)
	var ukuranFile string
	if sizeInBytes >= 1024*1024 {
		sizeInMB := float64(sizeInBytes) / (1024 * 1024)
		ukuranFile = fmt.Sprintf("%.2f MB", sizeInMB)
	} else {
		sizeInKB := float64(sizeInBytes) / 1024
		ukuranFile = fmt.Sprintf("%.2f KB", sizeInKB)
	}

	_ = ioutil.WriteFile("./file-dekrip/dekrip_"+fileData[0].FileName, decryptedData, 0644)
	oldEnkrip, err := file.GetDataEnkripByID(uint(idUint))

	if err != nil {
		c.JSON(404, gin.H{
			"error": "File enkripsi tidak ditemukan.",
		})
		return
	}

	oldEnkrip.FileStatus = "decrypted"
	_, err = oldEnkrip.UpdateDataEnkrip(uint(idUint))

	if err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal memperbarui status file enkripsi.",
		})
		return
	}

	var waktu string
	if elapsed.Seconds() < 1 {
		waktu = fmt.Sprintf("%d ms", elapsed.Milliseconds())
	} else {
		waktu = fmt.Sprintf("%.2f detik", elapsed.Seconds())
	}

	// Simpan data dekripsi ke database
	var dekrip Model.Dekrip
	dekrip.FileName = fileData[0].FileName
	dekrip.FileSize = ukuranFile
	dekrip.FileType = fileData[0].FileType
	dekrip.FilePath = "http://localhost:8080/file-dekrip/dekrip_" + fileData[0].FileName
	dekrip.FileHash = fileData[0].FileHash
	dekrip.FileDate = fileData[0].FileDate
	dekrip.FileEncryptedBy = "example-user"
	dekrip.FileStatus = "dekrip"
	dekrip.ExecutionTime = waktu
	dekrip.Key = customKey
	dekrip.EnkripID = uint(idUint)
	_, err = dekrip.SaveDataDekrip()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal menyimpan data dekripsi.",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "File berhasil didekripsi.",
		"waktu":   waktu,
	})
}
