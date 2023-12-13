package middleware

import (
	"github.com/Gonnekone/onlineStore/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckRoleAdmin(c *gin.Context) {
	user, _ := c.Get("user")

	if user != nil {
		if u, ok := user.(models.User); ok {
			if u.Role == "ADMIN" {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Pizdec": "Ты не админ"})
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Pizdec": "Как нахуй так, что ты вытащил по ключу user не user!?"})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Pizdec": "user = nil"})
	}
}
