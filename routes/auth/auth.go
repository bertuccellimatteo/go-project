package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devmatteo/go-project/claims"
	"github.com/devmatteo/go-project/inputs"
	"github.com/devmatteo/go-project/middlewares"
	"github.com/devmatteo/go-project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Routes(router *gin.Engine, db *gorm.DB) {
	authRouter := router.Group("/auth")

	authRouter.POST("/register", func(ctx *gin.Context) {
		var input inputs.RegisterInput

		if err := ctx.BindJSON(&input); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		user := models.User{
			Username: input.Username,
			Password: string(password),
		}

		result := db.Create(&user)

		if result.Error != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.AbortWithStatus(http.StatusNoContent)
	})

	authRouter.POST("/login", func(ctx *gin.Context) {
		var input inputs.LoginInput

		if err := ctx.BindJSON(&input); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var user models.User

		result := db.Where("username = ?", input.Username).Find(&user)

		if result.Error != nil {
			switch {
			case errors.Is(result.Error, gorm.ErrRecordNotFound):
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "User not found",
				})
			default:
				ctx.AbortWithError(http.StatusBadRequest, result.Error)
			}

			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Wrong password",
			})
			return
		}

		jwtSigner := jwt.NewWithClaims(jwt.SigningMethodHS256,
			claims.JwtAccessTokenClaim{
				UserID: user.ID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(time.Hour * 24))),
				},
			})

		signed, err := jwtSigner.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))

		if err != nil {
			fmt.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"accessToken": signed})
	})

	authRouter.GET("/user/me", middlewares.IsAuthenticated(), func(ctx *gin.Context) {
		uid := ctx.GetUint("UserID")

		var user models.User

		err := db.Where("id = ?", uid).Find(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
			} else {
				ctx.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}

		ctx.JSON(http.StatusOK, user)
	})
}
