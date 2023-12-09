package handlers

import (
	"errors"
	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func Registration(c *gin.Context) {
	var userReg DTO.UserReg

	if err := c.ShouldBindJSON(&userReg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	user.Password = string(hash)
	user.Email = userReg.Email
	if userReg.Role != "" {
		user.Role = userReg.Role
	}

	if errors.Is(initializers.DB.Where("email = ?", userReg.Email).First(&user).Error, gorm.ErrRecordNotFound) {
		if err := initializers.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Failed to create user": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
		return
	}

	var basket models.Basket
	basket.UserID = user.ID
	if err := initializers.DB.Create(&basket).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to create basket": err.Error()})
		return
	}

	tokenString := generateToken(user.ID, user.Email, user.Role)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var userLogin DTO.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", userLogin.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Email not found": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid password": err.Error()})
		return
	}

	tokenString := generateToken(user.ID, user.Email, user.Role)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Status(http.StatusOK)
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	tokenString := generateToken(user.(models.User).ID, user.(models.User).Email, user.(models.User).Role)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.Status(http.StatusOK)
}

func generateToken(id uint, email, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"timeIssued": time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
		"user_id":    id,
		"user_role":  role,
		"user_email": email,
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString
}
