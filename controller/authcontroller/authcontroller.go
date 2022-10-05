package authcontroller

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go_jwt_mux/config"
	"go_jwt_mux/helper"
	"go_jwt_mux/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

func init() {
	// 12h hh:mm:ss: 2:23:20 PM
	const (
		HHMMSS12h = "3:04:05 PM"
	)
	log.SetPrefix(time.Now().UTC().Format(HHMMSS12h) + ": ")
	log.SetFlags(log.Lshortfile)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInput models.User

	// take input json into struct golang
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Printf("error decode request: %v \n", err)
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(r.Body)

	// take user data based on username
	var user models.User

	if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			log.Printf("error decode request: %v \n", err)
			helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": "username atau password salah"})
			return
		default:
			log.Printf("error decode request: %v \n", err)
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}

	fmt.Println("still running")

	// cek apakah password valid
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		log.Printf("error decode request: %v \n", err)
		helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	// generating JWT token process
	expTime := time.Now().Add(time.Minute)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// declaring algorithm that will be used
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	log.Printf("jwt key: %v \n", config.JWT_KEY)
	if err != nil {
		log.Printf("error decode request: %v \n", err)
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// set token into cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})
	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "login succeed"})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput models.User

	// take input json into struct golang
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Printf("error decode request: %v \n", err)
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(r.Body)

	// has pass with bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)

	if err := models.DB.Create(&userInput).Error; err != nil {
		log.Printf("error create with gorm: %v \n", err)
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	helper.ResponseJSON(w, http.StatusCreated, map[string]string{"message": "success register"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete token into cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "logout succeed"})
}
