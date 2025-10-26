package http

import (
	"auth/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type HTTPHandlers struct {
	database database.Postgres
	Logged isLogged
}

func NewHTTPHandlers(db database.Postgres) *HTTPHandlers{
	return &HTTPHandlers{
		database: db,
		Logged: isLogged{},
	}
}



func (h *HTTPHandlers) HandleRegUser(w http.ResponseWriter, r *http.Request) {
	var RegDTO = RegDTO{}
		if h.Logged.isLogged{
		w.WriteHeader(http.StatusConflict) // тут редирект на главную страницу
		w.Write([]byte("Вы не можете создать новый аккаунт, так как уже залогинены"))
		return
	} 
	if err := json.NewDecoder(r.Body).Decode(&RegDTO); err != nil{
		errDTO := CreateErrDTO(err.Error(), time.Now())

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := RegDTO.ValidateToCreate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(RegDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash password: %w", err)
		return
	}
	if err := h.database.DbRegUser(RegDTO.Username, hashedPass); err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Такой аккаунт уже есть"))
		return
	}
	b, err := json.MarshalIndent(RegDTO, "", "	")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(b); err != nil {
		fmt.Print(err)
	}

	if _, err := w.Write([]byte("\nУспешная регистрация")); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}


func (h *HTTPHandlers) HandleLogUser(w http.ResponseWriter, r *http.Request) {
	var LogDTO = LogDTO{}
	if h.Logged.isLogged{
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Вы уже залогинены")) // тут редирект на главную страницу
		
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&LogDTO); err != nil{
		errDTO := CreateErrDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := LogDTO.ValidateToCreate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := h.database.DbLogUser(LogDTO.Username, LogDTO.Password)	
	if err != nil{
		fmt.Print("Ошибка при входе в аккаунт: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(LogDTO.Password))
	if err != nil {
		w.Write([]byte("Неверный пароль или логин"))
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	tokenString, err := GenerateToken(LogDTO.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	h.Logged.isLogged = true
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Authorization", "Bearer "+tokenString)
	if _, err := w.Write([]byte("Успешный вход в аккаунт")); err != nil {
		fmt.Println("failed to write http response: ", err)
		return
	}
}

func (h *HTTPHandlers) HandleVerify(w http.ResponseWriter, r *http.Request){
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if _, err := ValidateToken(token); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (h *HTTPHandlers) HandleLogout(w http.ResponseWriter, r *http.Request){
	if h.Logged.isLogged{
	h.Logged.isLogged = false
	w.Write([]byte("Успешный выход из аккаунта"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
func TestTokenHandler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
    http.Error(w, `{"error": "Неверный формат токена"}`, http.StatusUnauthorized)
    return
	}
    fmt.Fprintf(w, "Token: %s\n", tokenString)
    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        fmt.Fprintf(w, "Error: %v\n", err)
        return
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        fmt.Fprintf(w, "Valid! UserID: %v\n", claims["user_id"])
    } else {
        fmt.Fprintf(w, "Invalid token\n")
    }
}