package domain

import "time"

type BookStatus string

const (
	StatusWantToRead BookStatus = "want"
	StatusReading    BookStatus = "reading"
	StatusCompleted  BookStatus = "completed"
)

type Book struct {
	ID            int64
	UserID        int64
	OpenLibraryID string
	Title         string
	Author        string
	TotalPages    int
	CurrentPage   int
	Status        BookStatus
	CoverURL      string
	CreatedAt     time.Time
}
