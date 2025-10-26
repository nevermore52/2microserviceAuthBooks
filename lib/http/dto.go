package http

import (
	"encoding/json"
	"errors"
	"time"
)
type BookDTO struct {
	Title		string
	Author		string
	Pages		int
}

func (b BookDTO) ValidateToCreate() error {
	if b.Title == "" {
		return errors.New("title is empty")
	}
	if b.Author == "" {
		return errors.New("author is empty")
	}
	if b.Pages == 0 {
		return errors.New("pages is 0")
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

type CompleteBookDTO struct {
	Complete bool
}

type AuthorDTO struct {
	Author string
}

type AuthClient struct {
	BaseURL string
}

type VerifyResponse struct {
   		Valid  bool   `json:"valid"`
    	UserName string `json:"user_name"`
	}