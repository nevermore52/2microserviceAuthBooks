package library

import (
	"fmt"
	"libraryes/database"
	"libraryes/struct"
	"sync"
)

type Library struct {
	books    map[string]str.Book
	authors  map[string]str.Author
	mtx      sync.RWMutex
	postgres database.Postgres
}

func NewLibrary(pg database.Postgres) *Library {
	tempBooks, err := pg.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	tempAuthors, err := pg.DBExportAuthors()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Library{
		books:    tempBooks,
		postgres: pg,
		authors:  tempAuthors,
	}
}

func (l *Library) AddBook(book str.Book) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if err := l.postgres.DBInsertBooks(book.Title, book.Author, book.Pages); err != nil {
		fmt.Println(err)
		return err
	}
	
	if _, ok := l.books[book.Title]; ok {
		return ErrBookAlreadyExists
	}

	l.books[book.Title] = book

	return nil
}

func (l *Library) GetBook(title string) (str.Book, error) {	
	l.mtx.RLock()
	defer l.mtx.RUnlock()


	book, ok := l.books[title]
	if !ok {
		return str.Book{}, ErrBookNotFound
	}
	return book, nil
}

func (l *Library) ListBooks() map[string]str.Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tempBooks
}

func (l *Library) ListUnReadedBooks() map[string]str.Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	UnReadedBooks := make(map[string]str.Book)

	for title, book := range tempBooks {
		if !book.Readed {
			UnReadedBooks[title] = book
		}
	}
	return UnReadedBooks
}

func (l *Library) ListReadedBooks() map[string]str.Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ReadedBooks := make(map[string]str.Book)

	for title, book := range tempBooks {
		if book.Readed {
			ReadedBooks[title] = book
		}
	}
	return ReadedBooks
}

func (l *Library) ReadBook(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	book := tempBooks[title]

	ReadBook(&book)
	l.postgres.DBReadBook(title)
	l.books[title] = book

	return nil
}

func (l *Library) UnReadBook(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	book, ok := tempBooks[title]
	if !ok {
		return ErrBookNotFound
	}

	UnReadBook(&book)

	l.books[title] = book

	return nil
}

func (l *Library) DeleteBook(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, ok := tempBooks[title]
	if !ok {
		return ErrBookNotFound
	}
	l.postgres.DBDeleteBook(title)
	delete(l.books, title)

	return nil
}

func (l *Library) BoolReadBook(title string) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	book, ok := l.books[title]
	if !ok {
		return false
	}

	b := BoolReadBooks(&book)

	return b
}

func (l *Library) ListBooksAuthor(author string) map[string]str.Book {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	tempBooks, err := l.postgres.DBExportBooks()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ListBooksAuthor := make(map[string]str.Book)
	for title, book := range tempBooks {
		if author == book.Author {
			ListBooksAuthor[title] = book
		}
	}
	return ListBooksAuthor
}

func (l *Library) ListAuthors() map[string]str.Author {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	

	return l.authors
}

func (l *Library) AddAuthor(author str.Author) error {
	if err := l.postgres.DBAddAuthor(author.Author); err != nil {
		return err
	}

	l.authors[author.Author] = author

	return nil
}

func (l *Library) DeleteAuthor(author string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if err := l.postgres.DBDeleteAuthor(author); err != nil {
		return err
	}
	_, ok := l.authors[author]
	fmt.Println(l.authors)
	if !ok {
		return ErrAuthorNotFound
	}
	delete(l.authors,author)
	
	return nil
}
