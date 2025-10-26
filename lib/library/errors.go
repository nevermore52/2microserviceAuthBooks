package library

import "errors"

var ErrBookNotFound = errors.New("book not found")
var ErrBookAlreadyExists = errors.New("book already exists")
var ErrAuthorNotFound = errors.New("author not found")
var ErrAuthorAlreadyExists = errors.New("author already exists")
