package http

import (
	"encoding/json"
	"errors"
	"time"
)
type RegDTO struct {
	Username string 	`db:"username" json:"Username"`
	Password string		`db:"hashedpass" json:"Password"`
}

type isLogged struct {
	isLogged bool
}

type LogDTO struct {
	Username string 	`json:"Username"`
	Password string 	`json:"Password"`
}

type Token struct{
	Token string
}

func (r RegDTO) ValidateToCreate() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

func (l LogDTO) ValidateToCreate() error {
	if l.Username == "" {
		return errors.New("username is empty")
	}
	if l.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

type ErrorDTO struct {
	Message string
	time    time.Time
}

func CreateErrDTO(message string, time time.Time) ErrorDTO {
	return ErrorDTO{
		Message: message,
		time: time,
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}