package claims

import "github.com/golang-jwt/jwt/v5"

type JwtAccessTokenClaim struct {
	UserID uint
	jwt.RegisteredClaims
}
