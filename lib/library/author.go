package library

import "libraryes/struct"

func NewAuthor(author string) str.Author {
	return str.Author{
		Author: author,
	}
}