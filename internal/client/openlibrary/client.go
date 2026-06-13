package openlibrary

import (
	"net/http"
	"readHub/internal/domain"
)

type OpenLibraryClient interface {
	SearchBooks(query string) ([]domain.SearchBook, error)
	// GetByOpenLibraryID(openLibraryID string) (domain.SearchBook, error)
	GetBookDetails(openLibraryID string) (domain.BookDetails, error)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    "https://openlibrary.org",
	}
}

// TODO:
// Исследовать получение TotalPages через Editions API.
