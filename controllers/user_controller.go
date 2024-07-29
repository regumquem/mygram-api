package controllers

import (
	"abid-prakerja-uji-kemampuan/db"
	"abid-prakerja-uji-kemampuan/dto"
	"abid-prakerja-uji-kemampuan/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func UserRegister(ctx *gin.Context) {
	requestBody := dto.RegisterRequestDto{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Validate request
	err := validate.Struct(requestBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	newUser := entity.User{
		Email:    requestBody.Email,
		Username: requestBody.Username,
		Password: string(hashedPassword),
		Age:      requestBody.Age,
	}

	pg := db.GetDB()

	if err := pg.Create(&newUser).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := dto.RegisterResponseDto{
		Status: "success",
		Data: dto.UserData{
			ID:       newUser.ID,
			Email:    newUser.Email,
			Username: newUser.Username,
		},
	}

	ctx.JSON(http.StatusCreated, result)
}
