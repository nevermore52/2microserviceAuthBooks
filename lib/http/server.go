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

func (h *HTTPServer) StartServer() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("error read env file")
	}
	router := mux.NewRouter()
	router.Path("/login").Methods("POST").HandlerFunc(h.httpHandlers.HandleLogin)
	router.Path("/register").Methods("POST").HandlerFunc(h.httpHandlers.HandleRegister)
	router.Path("/authors").Methods("POST").HandlerFunc(h.httpHandlers.HandleCreateAuthor)
	router.Path("/authors").Methods("GET").HandlerFunc(h.httpHandlers.HandleListAuthors)
	router.Path("/authors/{author}").Methods("DELETE").HandlerFunc(h.httpHandlers.HandleDeleteAuthor)
    router.Path("/books").Methods("POST").HandlerFunc(h.httpHandlers.HandleCreateBook)
	router.Path("/books/{title}").Methods("GET").HandlerFunc(h.httpHandlers.HandleGetBook)
	router.Path("/books").Methods("GET").Queries("readed", "true").HandlerFunc(h.httpHandlers.HandleGetReadedBook)
	router.Path("/books").Methods("GET").Queries("readed", "false").HandlerFunc(h.httpHandlers.HandleGetUnReadedBook)
	router.Path("/books").Methods("GET").Queries("author", "{author}").HandlerFunc(h.httpHandlers.HandleListBookAuthor)
	router.Path("/books").Methods("GET").HandlerFunc(h.httpHandlers.HandleGetAllBook)
	router.Path("/books/{title}").Methods("PATCH").HandlerFunc(h.httpHandlers.HandleReadBook)
	router.Path("/books/{title}").Methods("DELETE").HandlerFunc(h.httpHandlers.HandleDeleteBook)
	mux := authMiddleware(router)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return err
	}
	return nil
}