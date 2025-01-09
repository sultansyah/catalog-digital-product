package middleware

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"catalog-digital-product/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
			return
		}

		token, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helper.HandleErrorResponde(c, custom.ErrUnauthorized)
			c.Abort()
			return
		}

		userId := int(claims["user_id"].(float64))

		c.Set("userId", userId)
		c.Next()
	}
}
