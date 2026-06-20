package repository

import "readHub/internal/domain"

type BookRepository interface {
	Create(book domain.Book) error
	GetByID(bookID int64) (domain.Book, error)
	GetByUserID(userID int64) ([]domain.Book, error)

	UpdateStatus(bookID int64, status domain.BookStatus) error
	UpdateCurrentPage(bookID int64, page int) error
	UpdateTotalPages(bookID int64, totalPages int) error

	CountByUserID(userID int64) (int, error)
	CountByStatus(userID int64, status domain.BookStatus) (int, error)

	Delete(bookID int64) error
}
