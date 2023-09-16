package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stevenjack/golang-testing-reference/repository/domain"
	"github.com/stevenjack/golang-testing-reference/storage"
)

//go:generate go run go.uber.org/mock/mockgen -package repositorytest -destination=repositorytest/mocks.go -source=book.go
type BookFetcher interface {
	FetchByID(id string) (domain.Book, error)
}

type Book struct {
	storage.Fetcher
}

func NewBook(f storage.Fetcher) Book {
	return Book{
		Fetcher: f,
	}
}

func (b Book) FetchByID(id string) (domain.Book, error) {
	result, err := b.Fetch(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.Book{}, fmt.Errorf("%w: %w", err, ErrBookNotFound)
		default:
			return domain.Book{}, fmt.Errorf("unable to find book with id '%s': %w", id, err)
		}

	}

	var book domain.Book

	if err := json.Unmarshal(result, &book); err != nil {
		return domain.Book{}, fmt.Errorf("problem unmarshalling book with id '%s': %w", id, err)
	}

	return book, nil
}
