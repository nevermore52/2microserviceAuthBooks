package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(HTTPHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: HTTPHandlers,
	}
}

func (h *HTTPServer) StartServer() error{
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("error read env file")
	}
	router := mux.NewRouter()
	router.Path("/register").Methods("POST").HandlerFunc(h.httpHandlers.HandleRegUser)
	router.Path("/auth/login").Methods("POST").HandlerFunc(h.httpHandlers.HandleLogUser)
	router.Path("/auth/verify").Methods("POST").HandlerFunc(h.httpHandlers.HandleVerify)
	router.Path("/logout").Methods("POST").HandlerFunc(h.httpHandlers.HandleLogout)
	router.Path("/test-token").Methods("POST").HandlerFunc(TestTokenHandler)
	if err := http.ListenAndServe(":8081", router); err != nil {
		return err
	}
	return nil
}