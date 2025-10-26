package database

import (
	"errors"
	"fmt"

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

		users := ` CREATE TABLE IF NOT EXISTS users( 
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			hashedpass VARCHAR(100) NOT NULL)`
		if _, err := db.Exec(users); err != nil {
			fmt.Println(err)
			return nil, err
		}
		
		return db, nil
}


func (p *Postgres) DbRegUser(username string, password []byte) error {
	if _, err := p.DB.Exec(`
	INSERT INTO users (username, hashedpass)
	VALUES ($1, $2)`, username, password); err != nil {
		return err
	}
	return nil  
}

func (p *Postgres) DbLogUser(username string, password string) (string, error) {
	var passhash string
	err := p.DB.QueryRow( `SELECT hashedpass FROM users WHERE username = $1
	`, username).Scan(&passhash)
	if err != nil {
		return "", errors.New("такого пользователя нету в системе")
	}
	return passhash, nil
}