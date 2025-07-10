package Controller

import (
	helper "kriptografi-zaidaan/Helper"
	"kriptografi-zaidaan/Model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(context *gin.Context) {
	var input Model.AuthenticationInput
	if err := context.ShouldBind(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := Model.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		Name:     input.Name,
		Role:     input.Role,
		Alamat:   input.Alamat,
		Telp:     input.Telp,
	}

	_, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User created successfully! "})
}

func GetUserByUsername(context *gin.Context) {
	username := context.Param("username")

	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	user, err := Model.FindUserByUsername(username)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	context.JSON(http.StatusOK, user)
}

func UpdateUser(context *gin.Context) {
	username := context.Param("username")

	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var input Model.AuthenticationInput
	if err := context.ShouldBind(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Model.FindUserByUsername(username)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Username = input.Username
	user.Password = input.Password
	user.Email = input.Email
	user.Name = input.Name
	user.Role = input.Role
	user.Alamat = input.Alamat
	user.Telp = input.Telp

	updatedUser, err := user.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	context.JSON(http.StatusOK, updatedUser)
}

func GetAllUsers(context *gin.Context) {
	users, err := Model.GetAllUsers()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	context.JSON(http.StatusOK, users)
}

func DeleteUser(context *gin.Context) {
	username := context.Param("username")

	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	err := Model.DeleteUserByUsername(username)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func Login(context *gin.Context) {
	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Model.FindUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login success", "token": jwt})
}
