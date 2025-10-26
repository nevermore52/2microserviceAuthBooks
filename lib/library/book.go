package library

import (
	"libraryes/struct"
	"time"
)
func NewBook(title string, author string, pages int) str.Book {
	return str.Book{
		Title: title,
		Author: author,
		Pages: pages,
		Readed: false,

		Timeadd: time.Now(),
		Timereaded: nil,
	}
}

func ReadBook(b *str.Book) {
	TimeRead := time.Now()
	b.Readed = true
	b.Timereaded = &TimeRead
}

func UnReadBook(b *str.Book) {
	b.Readed = false
	b.Timereaded = nil
}

func  BoolReadBooks(b *str.Book) bool {
	return b.Readed
}