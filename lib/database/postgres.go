package database

import (
	"fmt"

	"libraryes/struct"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host 		string
	Port	 	string
	Username 	string
	Password	string
	DBName		string
	SSLMode 	string
}


type Postgres struct {
	DB 		*sqlx.DB

}


func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	db, err := sqlx.Open("postgres", config)
		if err != nil {
			fmt.Println("error to open postgres db")
			return nil, err
		}

		err = db.Ping()
		if err != nil{
			return nil, err
		}

		authors := (` CREATE TABLE IF not exists authors( 
		author VARCHAR(50) NOT NULL PRIMARY KEY)`)
		if _, err := db.Exec(authors); err != nil {
			fmt.Println(err)
			return nil, err
		}

		books := (` CREATE TABLE IF not exists books( 
	    id SERIAL PRIMARY KEY,
    	title VARCHAR(50) NOT NULL UNIQUE,
    	author VARCHAR(50) NOT NULL,
   	 	pages VARCHAR(50) NOT NULL,
   	 	readed BOOLEAN DEFAULT FALSE,
   	 	timeADD TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	timeREAD TIMESTAMP,
    	CONSTRAINT fk_author
        FOREIGN KEY(author)
            REFERENCES authors(author)
            ON DELETE CASCADE
            ON UPDATE CASCADE
		)`)
		if _, err := db.Exec(books); err != nil {
			fmt.Println(err)
			return nil, err
		}



		return db, nil
}

func (p *Postgres) DBInsertBooks(title string, author string, pages int) error {
	if _, err := p.DB.Exec(`
	INSERT INTO books (title,author,pages)
	VALUES($1, $2, $3)
	`,title, author, pages); err != nil {
		fmt.Println("Данная книга уже находится в базе данных!")
		return err
	}

	return nil
}

func (p *Postgres) DBReadBook(title string) {
	if _, err := p.DB.Exec(`
	UPDATE books SET 
	readed = $2,
	timeread = CURRENT_TIMESTAMP
	WHERE title = $1
	`,title, true); err != nil {
		fmt.Println("error update readed book: ", err)
		return
	}

}

func (p *Postgres) DBExportBooks() (map[string]str.Book, error){
	books := []str.Book{}
	err := p.DB.Select(&books, "SELECT * FROM books")
	tempMap := make(map[string]str.Book)
	tmp := str.Book{}
	for i, v := range books {
		tmp = books[i] 
		tempMap[tmp.Title] = v
	}
	return tempMap, err
}

func (p *Postgres) DBDeleteBook(title string) {
	if _, err := p.DB.Exec(`
	DELETE FROM books
	WHERE title = $1
	`, title); err != nil{
		fmt.Println("Ошибка при удалении книги", err)
		return 
	}
}

func (p *Postgres) DBExportAuthors() (map[string]str.Author, error){
	authors := []str.Author{}
	err := p.DB.Select(&authors, "SELECT * FROM authors")
	tempMap := make(map[string]str.Author)
	tmp := str.Author{}
	for i, v := range authors {
		tmp = authors[i] 
		tempMap[tmp.Author] = v
	}
	return tempMap, err
}

func (p *Postgres) DBAddAuthor(author string) error {
	if _, err := p.DB.Exec(`INSERT INTO authors (author)
	VALUES ($1)`, author); err != nil {
		fmt.Print("Ошибка при добавлении автора", err)
		return err
	}
	return nil
}

func (p *Postgres) DBDeleteAuthor(author string) error {
	if _, err := p.DB.Exec(`
	DELETE FROM authors
	WHERE author = $1
	`, author); err != nil {
		fmt.Println("Ошибка при удалении автора", err)
		return err
	}

	return nil
}