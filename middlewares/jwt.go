package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
	"go_jwt_mux/config"
	"go_jwt_mux/helper"
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

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("error decode request: %v \n", err)
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
				return
			}
		}
		// get token value
		tokenString := c.Value
		claims := &config.JWTClaim{}

		// ========== parsing token jwt ==========
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		// check error
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				// token invalid
				log.Printf("error decode request: %v \n", jwt.ValidationErrorSignatureInvalid)
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			case jwt.ValidationErrorExpired:
				log.Printf("error decode request: %v \n", jwt.ValidationErrorExpired)
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized, Token Expired!"})
				return
			default:
				log.Printf("error decode request: %v \n", http.StatusText(http.StatusUnauthorized))
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			}

		}

		// check token is valid
		if !token.Valid {
			log.Printf("error decode request: %v \n", http.StatusText(http.StatusUnauthorized))
			helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
