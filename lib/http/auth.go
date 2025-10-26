package http

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func NewAuthClient() *AuthClient{
	return &AuthClient{
		BaseURL: "http://auth-service:8081",
		}
	}

func (c *AuthClient) VerifyToken(tokenString string) (*Claims, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, token.Valid, nil
	}

	return nil, false, fmt.Errorf("инвалид токен")

}

func authMiddleware(next http.Handler) http.Handler {
	authClient := NewAuthClient()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.URL.Path == "/login" {
			next.ServeHTTP(w,r)
			return
		}
		if r.URL.Path == "/register" {
			next.ServeHTTP(w,r)
			return
		}
		if r.URL.Path == "/logout" {
			next.ServeHTTP(w,r)
			return
		}
	

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный токен нету префикса Bearer"))
		return
	}

	_, valid, err := authClient.VerifyToken(token)
	if err != nil || !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный токен"))
		return
	}

	next.ServeHTTP(w, r)
	})
}
