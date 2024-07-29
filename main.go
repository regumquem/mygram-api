package main

import (
	"abid-prakerja-uji-kemampuan/db"
	"abid-prakerja-uji-kemampuan/dto"
	"abid-prakerja-uji-kemampuan/entity"
	"abid-prakerja-uji-kemampuan/middleware"
	"abid-prakerja-uji-kemampuan/pkg/internal_jwt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserLogin(ctx *gin.Context) {
	requestBody := dto.LoginRequestDto{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg := db.GetDB()
	user := entity.User{}

	if err := pg.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid email/password",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid email/password",
		})
		return
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	}

	token, err := internal_jwt.GenerateToken(claims)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	result := dto.LoginResponseDto{
		Status: "success",
		Data:   token,
	}

	ctx.JSON(http.StatusOK, result)
}

func UserRegister(ctx *gin.Context) {
	requestBody := dto.LoginRequestDto{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	if len([]rune(requestBody.Password)) < 6 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "password must be at least 6 characters",
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
		Password: string(hashedPassword),
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
		Data:   nil,
	}

	ctx.JSON(http.StatusCreated, result)
}

func CreatePhoto(ctx *gin.Context) {
	userId, ok := ctx.MustGet("userId").(float64)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	requestBody := dto.CreatePhotoRequestDto{}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	newPhoto := entity.Photo{
		Title:    requestBody.Title,
		Caption:  requestBody.Caption,
		PhotoURL: requestBody.PhotoURL,
		UserID:   uint(userId),
	}

	pg := db.GetDB()

	if err := pg.Create(&newPhoto).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

func GetPhotos(ctx *gin.Context) {
	userId, ok := ctx.MustGet("userId").(float64)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	var photos []entity.Photo

	pg := db.GetDB()
	if err := pg.Where("user_id = ?", uint(userId)).Find(&photos).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

func UpdatePhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid photo ID",
		})
		return
	}

	requestBody := dto.UpdatePhotoRequestDto{}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	pg := db.GetDB()
	photo := entity.Photo{}

	if err := pg.Where("id = ?", photoId).First(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "photo not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	photo.Title = requestBody.Title
	photo.Caption = requestBody.Caption
	photo.PhotoURL = requestBody.PhotoURL

	if err := pg.Save(&photo).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "photo updated successfully",
	})
}

func DeletePhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid photo ID",
		})
		return
	}

	pg := db.GetDB()
	photo := entity.Photo{}

	if err := pg.Where("id = ?", photoId).First(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "photo not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := pg.Delete(&photo).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "photo deleted successfully",
	})
}

func main() {
	db.InitializeDB()

	pg := db.GetDB()
	pg.AutoMigrate(entity.User{}, entity.Photo{})

	r := gin.Default()

	userRoute := r.Group("/users")
	{
		userRoute.POST("/register", UserRegister)
		userRoute.POST("/login", UserLogin)
	}

	photoRoute := r.Group("/photos")
	{
		photoRoute.Use(middleware.Authentication())
		photoRoute.GET("", GetPhotos)
		photoRoute.POST("", CreatePhoto)
		photoRoute.PUT("/:id", UpdatePhoto)
		photoRoute.DELETE("/:id", DeletePhoto)
	}

	r.Run(":8080")
}
