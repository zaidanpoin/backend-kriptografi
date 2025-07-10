package helper

import (
	"fmt"
	"kriptografi-zaidaan/Model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user Model.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,

		"nama": user.Name,
		"exp":  time.Now().Add(time.Minute * time.Duration(tokenTTL)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("Token expired in:", tokenTTL, "minutes")

	return token.SignedString(privateKey)
}
