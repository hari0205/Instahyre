package middleware

import (
	"net/http"
	"os"
	"time"

	ini "example.com/Instahyre/teleapi/init"
	model "example.com/Instahyre/teleapi/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authentication")
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Authentication Failed",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		//log.Fatalf("Failed to get cookie:%v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}

		return []byte(os.Getenv("SECRET_KEY")), nil //"69CE28610443221E486BB78C26C91ABDA5F2B057DC742C615D0C547A3B121B6F"
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	//fmt.Println(err)

	// Token validation
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user model.UserData
		ini.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
