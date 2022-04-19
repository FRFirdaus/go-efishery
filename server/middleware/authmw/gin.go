package authmw

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/efishery/go-efishery/libs/token"
	libtoken "bitbucket.org/efishery/go-efishery/libs/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GinAuthMiddleware is middleware to validate token is valid from server
func GinAuthMiddleware(signatureKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.Contains(authHeader, bearerSchema) {
			tokenString = authHeader[len(bearerSchema):]
		}

		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}
		tokenString = strings.TrimSpace(tokenString)
		token, err := token.NewService(signatureKey).Validate(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			log.Println(err)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			// currently using marshal unmarshal
			// readable but need tunning performace
			byteClaims, err := json.Marshal(claims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": http.StatusText(http.StatusUnauthorized),
				})
				return
			}

			myClaim := libtoken.AuthServiceJWT{}
			err = json.Unmarshal(byteClaims, &myClaim)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": http.StatusText(http.StatusUnauthorized),
				})
				return
			}

			c.Set(contextClaimKey, myClaim)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": http.StatusText(http.StatusUnauthorized),
			})
		}
	}
}

// GetGinClaim is get token claim from user
func GetGinClaim(c *gin.Context) libtoken.AuthServiceJWT {
	return GetClaim(c.Request)
}
