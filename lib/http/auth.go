package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func NewAuthClient() *AuthClient{
	return &AuthClient{
		BaseURL: "http://auth-service:8081",
		}
	}

func (c *AuthClient) VerifyToken(Token string) (bool, string, error) {
	requestBody, err := json.Marshal(map[string]string{"Token": Token})
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(c.BaseURL+"/auth/verify", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()



	var result VerifyResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Valid, result.UserName, nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	authClient := NewAuthClient()

	return func(w http.ResponseWriter, r *http.Request){
		if r.URL.Path == "/login" {
			next(w,r)
			return
		}
	

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный токен нету префикса Bearer"))
		return
	}

	valid, username, err := authClient.VerifyToken(token)
	if err != nil || !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный токен"))
		return
	}

	r.Header.Set("X-User-Name", username)
	next(w, r)
	}
}
