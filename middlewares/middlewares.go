package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/devmatteo/go-project/claims"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")

		if authorization == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization error missing",
			})
			return
		}

		splittedHeader := strings.Split(authorization, " ")

		tokenType, token := splittedHeader[0], splittedHeader[1]

		if tokenType != "Bearer" || token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is not valid",
			})
			return
		}

		parsedToken, err := jwt.ParseWithClaims(token, &claims.JwtAccessTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
		})

		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Expired access token",
			})
			return
		} else if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if claims, ok := parsedToken.Claims.(*claims.JwtAccessTokenClaim); ok && parsedToken.Valid {
			ctx.Set("UserID", claims.UserID)
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("err"))
		}
	}
}
