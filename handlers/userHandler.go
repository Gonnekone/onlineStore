package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Gonnekone/onlineStore/DTO"
	"github.com/Gonnekone/onlineStore/initializers"
	"github.com/Gonnekone/onlineStore/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	var userReg DTO.UserReg

	if err := json.NewDecoder(r.Body).Decode(&userReg); err != nil {
		http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := userReg.Validate(); err != nil {
		http.Error(w, "error validating user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), 10)
	if err != nil {
		http.Error(w, "error generating password hash: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var user models.User

	user.Password = string(hash)
	user.Email = userReg.Email
	if userReg.Role != "" {
		user.Role = userReg.Role
	}

	var existingUser models.User
	if err := initializers.DB.Where("email = ?", userReg.Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := initializers.DB.Create(&user).Error; err != nil {
				http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "error checking existing email: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	var basket models.Basket
	basket.UserID = user.ID
	if err := initializers.DB.Create(&basket).Error; err != nil {
		http.Error(w, "Failed to create basket: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := generateToken(user.ID, user.Email, user.Role)

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userLogin DTO.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := userLogin.Validate(); err != nil {
		http.Error(w, "error validating user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", userLogin.Email).First(&user).Error; err != nil {
		http.Error(w, "Email not found: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		http.Error(w, "Invalid password: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := generateToken(user.ID, user.Email, user.Role)

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func Validate(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	tokenString := generateToken(user.ID, user.Email, user.Role)

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func generateToken(id uint, email, role string) string {
	claims := jwt.MapClaims{
		"timeIssued": time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
		"user_id":    id,
		"user_role":  role,
		"user_email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString
}
