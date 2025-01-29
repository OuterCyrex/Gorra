package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

func JWTAuthMiddleware(JwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Token 为空",
			})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(authHeader, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		})
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "Token 已过期",
				})
				c.Abort()
				return
			default:
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "无效 Token",
				})
				c.Abort()
				return
			}
		}
		c.Set("claims", token.Claims.(*CustomClaims))
		c.Set("userId", token.Claims.(*CustomClaims).ID)
		return
	}
}
