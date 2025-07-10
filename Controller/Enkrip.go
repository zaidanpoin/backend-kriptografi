package Controller

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"kriptografi-zaidaan/Model"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Nb        = 4
	Nk        = 8
	Nr        = 14
	BlockSize = 16
)

var sBox = [256]byte{0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76, 0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0, 0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15, 0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75, 0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84, 0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf, 0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8, 0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2, 0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73, 0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb, 0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79, 0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08, 0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a, 0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e, 0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf, 0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16}

var rcon = [255]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1B, 0x36}

func subWord(word []byte) []byte {
	out := make([]byte, len(word))
	for i := 0; i < len(word); i++ {
		out[i] = sBox[word[i]]
	}
	return out
}

func rotWord(word []byte) []byte { return append(word[1:], word[0]) }

func keyExpansion(key []byte) []byte {
	expanded := make([]byte, Nb*4*(Nr+1))
	copy(expanded, key[:Nk*4])
	temp := make([]byte, 4)

	for i := Nk; i < Nb*(Nr+1); i++ {
		copy(temp, expanded[(i-1)*4:(i-1)*4+4])
		if i%Nk == 0 {
			temp = subWord(rotWord(temp))
			temp[0] ^= rcon[i/Nk-1]
		} else if Nk > 6 && i%Nk == 4 {
			temp = subWord(temp)
		}
		for j := 0; j < 4; j++ {
			expanded[i*4+j] = expanded[(i-Nk)*4+j] ^ temp[j]
		}
	}
	return expanded
}

func padPKCS7(data []byte) []byte {
	padding := BlockSize - len(data)%BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func xorBlocks(a, b []byte) []byte {
	result := make([]byte, BlockSize)
	for i := 0; i < BlockSize; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func subBytes(state []byte) {
	for i := range state {
		state[i] = sBox[state[i]]
	}
}

func shiftRows(state []byte) {
	tmp := make([]byte, 16)
	copy(tmp, state)
	state[1], state[5], state[9], state[13] = tmp[5], tmp[9], tmp[13], tmp[1]
	state[2], state[6], state[10], state[14] = tmp[10], tmp[14], tmp[2], tmp[6]
	state[3], state[7], state[11], state[15] = tmp[15], tmp[3], tmp[7], tmp[11]
}

func xtime(a byte) byte {
	if a&0x80 != 0 {
		return (a << 1) ^ 0x1b
	}
	return a << 1
}

func mixColumns(state []byte) {
	for i := 0; i < 4; i++ {
		a := state[i*4:]
		t := a[0] ^ a[1] ^ a[2] ^ a[3]
		u := a[0]
		a[0] ^= t ^ xtime(a[0]^a[1])
		a[1] ^= t ^ xtime(a[1]^a[2])
		a[2] ^= t ^ xtime(a[2]^a[3])
		a[3] ^= t ^ xtime(a[3]^u)
	}
}

func addRoundKey(state []byte, roundKey []byte) {
	for i := 0; i < BlockSize; i++ {
		state[i] ^= roundKey[i]
	}
}

func encryptBlock(input, keySchedule []byte) []byte {
	state := make([]byte, BlockSize)
	copy(state, input)

	addRoundKey(state, keySchedule[:BlockSize])

	for round := 1; round < Nr; round++ {
		subBytes(state)
		shiftRows(state)
		mixColumns(state)
		addRoundKey(state, keySchedule[round*BlockSize:(round+1)*BlockSize])
	}

	subBytes(state)
	shiftRows(state)
	addRoundKey(state, keySchedule[Nr*BlockSize:(Nr+1)*BlockSize])

	return state

}

func aesCBCEncrypt(data, key []byte) []byte {
	data = padPKCS7(data)
	iv := make([]byte, BlockSize)
	rand.Read(iv)
	ciphertext := make([]byte, BlockSize)
	copy(ciphertext, iv)

	keySchedule := keyExpansion(key)
	prev := iv
	for i := 0; i < len(data); i += BlockSize {
		block := xorBlocks(data[i:i+BlockSize], prev)
		encrypted := encryptBlock(block, keySchedule)
		ciphertext = append(ciphertext, encrypted...)
		prev = encrypted
	}
	return ciphertext

}

var invSBox = [256]byte{0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb, 0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb, 0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e, 0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25, 0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92, 0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84, 0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06, 0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b, 0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73, 0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e, 0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b, 0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4, 0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f, 0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef, 0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61, 0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d}

func invSubBytes(state []byte) {
	for i := range state {
		state[i] = invSBox[state[i]]
	}
}

func invShiftRows(state []byte) {
	tmp := make([]byte, 16)
	copy(tmp, state)
	state[1], state[5], state[9], state[13] = tmp[13], tmp[1], tmp[5], tmp[9]
	state[2], state[6], state[10], state[14] = tmp[10], tmp[14], tmp[2], tmp[6]
	state[3], state[7], state[11], state[15] = tmp[7], tmp[11], tmp[15], tmp[3]
}

func invMixColumns(state []byte) {
	for i := 0; i < 4; i++ {
		a := state[i*4:]
		u := xtime(xtime(a[0] ^ a[2]))
		v := xtime(xtime(a[1] ^ a[3]))
		a[0] ^= u
		a[1] ^= v
		a[2] ^= u
		a[3] ^= v
	}
	mixColumns(state)
}

func decryptBlock(input, keySchedule []byte) []byte {
	state := make([]byte, BlockSize)
	copy(state, input)

	addRoundKey(state, keySchedule[Nr*BlockSize:(Nr+1)*BlockSize])

	for round := Nr - 1; round > 0; round-- {
		invShiftRows(state)
		invSubBytes(state)
		addRoundKey(state, keySchedule[round*BlockSize:(round+1)*BlockSize])
		invMixColumns(state)
	}

	invShiftRows(state)
	invSubBytes(state)
	addRoundKey(state, keySchedule[:BlockSize])

	return state
}

func unpadPKCS7(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	padding := int(data[len(data)-1])
	if padding > BlockSize || padding == 0 {
		return data
	}
	return data[:len(data)-padding]
}

func aesCBCDecrypt(data, key []byte) []byte {
	if len(data) < BlockSize {
		return nil
	}

	iv := data[:BlockSize]
	ciphertext := data[BlockSize:]
	plaintext := make([]byte, 0, len(ciphertext))

	keySchedule := keyExpansion(key)
	prev := iv

	for i := 0; i < len(ciphertext); i += BlockSize {
		block := decryptBlock(ciphertext[i:i+BlockSize], keySchedule)
		decrypted := xorBlocks(block, prev)
		plaintext = append(plaintext, decrypted...)
		prev = ciphertext[i : i+BlockSize]
	}

	return unpadPKCS7(plaintext)
}

// Fungsi untuk mengkonversi string menjadi kunci AES
func stringToKey(keyString string) []byte {
	// Buat slice byte dengan panjang 32 (untuk AES-256)
	key := make([]byte, 32)

	// Konversi string ke byte dan salin ke key
	keyBytes := []byte(keyString)

	// Jika string lebih pendek dari 32 byte, isi sisa dengan 0
	// Jika lebih panjang, ambil 32 byte pertama
	for i := 0; i < 32; i++ {
		if i < len(keyBytes) {
			key[i] = keyBytes[i]
		} else {
			key[i] = 0
		}
	}

	return key
}

func EncryptHandler(c *gin.Context) {

	id := c.Param("id")
	username := c.PostForm("username")
	fmt.Println("Username:", username)
	customKey := c.PostForm("custom_key")
	fmt.Println(customKey)
	fmt.Println("adasda")

	var fileName Model.File
	data, err := fileName.GetDataFile(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error retrieving file data",
		})
		return
	}

	key := stringToKey(customKey)

	inputFile := "./file-plainteks/" + data[0].FileHash

	plaintext, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)

	}

	start := time.Now() // Mulai pengukuran waktu
	keyHex := hex.EncodeToString(key)

	hash := sha256.Sum256([]byte(data[0].FileName + time.Now().String()))
	hashedFileName := hex.EncodeToString(hash[:])

	ciphertext := aesCBCEncrypt(plaintext, key)
	elapsed := time.Since(start)
	_ = ioutil.WriteFile("./file-enkrip/"+hashedFileName, ciphertext, 0644)

	// prosess database

	file := Model.File{
		FileStatus: "encrypted",
	}

	_, err2 := file.UpdateFile(id)
	if err2 != nil {
		c.JSON(500, gin.H{
			"message": "Error updating file status",
		})
		return
	}

	var waktu string
	if elapsed.Seconds() < 1 {
		waktu = fmt.Sprintf("%d ms", elapsed.Milliseconds())
	} else {
		waktu = fmt.Sprintf("%.2f detik", elapsed.Seconds())
	}

	sizeInBytes := len(ciphertext)
	var ukuranFile string
	if sizeInBytes >= 1024*1024 {
		sizeInMB := float64(sizeInBytes) / (1024 * 1024)
		ukuranFile = fmt.Sprintf("%.2f MB", sizeInMB)
	} else {
		sizeInKB := float64(sizeInBytes) / 1024
		ukuranFile = fmt.Sprintf("%.2f KB", sizeInKB)
	}
	enkrip := Model.Enkrip{
		FileName: data[0].FileName,

		FileSize: ukuranFile,

		FilePath:        "http://localhost:8080/file-enkrip/" + hashedFileName,
		FileHash:        hashedFileName,
		Key:             customKey,
		FileKey:         keyHex,
		FileType:        data[0].FileType,
		FileDate:        time.Now().Format("2006-01-02"),
		FileEncryptedBy: username, // Ganti dengan informasi pengguna yang sebenarnya
		FileStatus:      "encrypted",
		FileID:          uint(data[0].ID),
		Excecution_time: waktu,
	}

	_, err3 := enkrip.SaveDataEnkrip()
	if err3 != nil {
		c.JSON(500, gin.H{
			"message": err3.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "File berhasil dienkripsi.",
		"key":     keyHex,
		"waktu":   waktu,
	})
}

func GetAllDataEncrypt(c *gin.Context) {

	var enkrip Model.Enkrip
	data, err := enkrip.GetAllDataEnkrip()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error retrieving encrypted data",
		})
		return
	}

	c.JSON(200, data)
}

func GetDataEncryptByID(c *gin.Context) {

	idString := c.Param("id")
	var enkrip Model.Enkrip
	idUint, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid ID format",
		})
		return
	}
	data, err := enkrip.GetDataEnkripByID(uint(idUint))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error retrieving encrypted data",
		})
		return
	}

	c.JSON(200, data)
}

func DeleteDataEncrypt(c *gin.Context) {
	idString := c.Param("id")
	var enkrip Model.Enkrip
	idUint, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid ID format",
		})
		return
	}
	data, err := enkrip.GetDataEnkripByID(uint(idUint))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error retrieving encrypted data",
		})
		return
	}

	fmt.Printf("DEBUG: Data to delete: %+v\n", data)

	err = enkrip.DeleteDataEnkrip(uint(idUint))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting encrypted data",
		})
		return
	}

	fmt.Printf("DEBUG: Deleted DB entry for ID %d\n", idUint)

	// Pastikan FilePath adalah path file lokal, bukan URL
	filePath := data.FilePath
	if len(filePath) > 0 && (filePath[:4] == "http" || filePath[:5] == "https") {
		// Ekstrak nama file dari URL
		lastSlash := -1
		for i := len(filePath) - 1; i >= 0; i-- {
			if filePath[i] == '/' {
				lastSlash = i
				break
			}
		}
		if lastSlash != -1 {
			filePath = "./file-enkrip/" + filePath[lastSlash+1:]
		}
	}

	fmt.Printf("DEBUG: File path to remove: %s\n", filePath)

	err = os.Remove(filePath)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting file",
			"debug":   err.Error(),
		})
		return
	}

	fmt.Println("DEBUG: File deleted successfully")

	c.JSON(200, gin.H{
		"message": "Data berhasil dihapus",
	})
}
