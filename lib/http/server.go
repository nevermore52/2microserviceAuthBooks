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
	router := mux.NewRouter( )
	router.Path("/login").Methods("POST").HandlerFunc(h.httpHandlers.HandleLogin)
	router.Path("/register").Methods("POST").HandlerFunc(h.httpHandlers.HandleRegister)
	router.Path("/authors").Methods("POST").HandlerFunc(authMiddleware(h.httpHandlers.HandleCreateAuthor))
	router.Path("/authors").Methods("GET").HandlerFunc(authMiddleware(h.httpHandlers.HandleListAuthors))
	router.Path("/authors/{author}").Methods("DELETE").HandlerFunc(authMiddleware(h.httpHandlers.HandleDeleteAuthor))
    router.Path("/books").Methods("POST").HandlerFunc(authMiddleware(h.httpHandlers.HandleCreateBook))
	router.Path("/books/{title}").Methods("GET").HandlerFunc(authMiddleware(h.httpHandlers.HandleGetBook))
	router.Path("/books").Methods("GET").Queries("readed", "true").HandlerFunc(authMiddleware(h.httpHandlers.HandleGetReadedBook))
	router.Path("/books").Methods("GET").Queries("readed", "false").HandlerFunc(authMiddleware(h.httpHandlers.HandleGetUnReadedBook))
	router.Path("/books").Methods("GET").Queries("author", "{author}").HandlerFunc(authMiddleware(h.httpHandlers.HandleListBookAuthor))
	router.Path("/books").Methods("GET").HandlerFunc(authMiddleware(h.httpHandlers.HandleGetAllBook))
	router.Path("/books/{title}").Methods("PATCH").HandlerFunc(authMiddleware(h.httpHandlers.HandleReadBook))
	router.Path("/books/{title}").Methods("DELETE").HandlerFunc(authMiddleware(h.httpHandlers.HandleDeleteBook))

	if err := http.ListenAndServe(":8080", router); err != nil {
		return err
	}
	return nil
}