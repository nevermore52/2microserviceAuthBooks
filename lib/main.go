package main

import (
	"fmt"
	"libraryes/database"
	"libraryes/http"
	"libraryes/library"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("error read env file")
	}
	db, err := database.NewPostgresDB(database.Config{
				Host: os.Getenv("DB_HOST"),
				Port: os.Getenv("DB_PORT"),
				Username: os.Getenv("DB_USERNAME"),
				Password: os.Getenv("DB_PASSWORD"),
				DBName: os.Getenv("DB_NAME"),
				SSLMode: "disable",
				})
				if err != nil {
					fmt.Print(err)
					return
				}

				defer db.Close()

	pg := database.Postgres{DB: db}
	library := library.NewLibrary(pg)
	httpHandlers := http.NewHTTPHandlers(library)
	httpServer := http.NewHTTPServer(httpHandlers)
	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server", err)
	}
}
