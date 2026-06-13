package postgres

import (
	"context"

	"readHub/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *PostgresBookRepository {
	repo := PostgresBookRepository{
		db: db,
	}
	return &repo
}

func (r *PostgresBookRepository) Create(book domain.Book) error {
	// Exec обычно используют для тех запросов которые ничего не возвращают, это INSERT, UPDATE, DELETE/ в остальных случаев query
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO books (user_id, open_library_id, title, author, total_pages, status, cover_url) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		book.UserID,
		book.OpenLibraryID,
		book.Title,
		book.Author,
		book.TotalPages,
		book.Status,
		book.CoverURL,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresBookRepository) GetByID(bookID int64) (domain.Book, error) {
	var book domain.Book
	// для записи одной строки берется QueryRow, и err здесь возвращает не queryrow а Scan
	err := r.db.QueryRow(context.Background(),
		"SELECT id, user_id, open_library_id, title, author, total_pages, current_page, status, cover_url, created_at FROM books WHERE id=$1",
		bookID,
	).Scan(&book.ID, &book.UserID, &book.OpenLibraryID, &book.Title, &book.Author, &book.TotalPages, &book.CurrentPage, &book.Status, &book.CoverURL, &book.CreatedAt)
	if err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func (r *PostgresBookRepository) GetByUserID(userID int64) ([]domain.Book, error) {
	var books []domain.Book
	// Query нужен для записи нескольких строк, слайса данных так сказать
	rows, err := r.db.Query(context.Background(),
		"SELECT id, user_id, open_library_id, title, author, total_pages, current_page, status, cover_url, created_at FROM books WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() { // если след строка есть, выполнить тело цикла
		var book domain.Book // создаем внутри с каждым циклом новую переменную, еслиб была вне цикла то данные перезаписывались бы в одну и ту же переменную
		err = rows.Scan(&book.ID, &book.UserID, &book.OpenLibraryID, &book.Title, &book.Author, &book.TotalPages, &book.CurrentPage, &book.Status, &book.CoverURL, &book.CreatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	err = rows.Err() // проверка на ошибку, так как даже если цикл завершился в моменте его выполнения могла быть ошибка, вот мы её и отлоавливаем
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *PostgresBookRepository) UpdateCurrentPage(bookID int64, page int) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE books SET current_page=$2 WHERE id=$1",
		bookID,
		page,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresBookRepository) UpdateStatus(bookID int64, status domain.BookStatus) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE books SET status=$2 WHERE id=$1",
		bookID,
		status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresBookRepository) Delete(bookID int64) error {
	_, err := r.db.Exec(context.Background(),
		"DELETE FROM books WHERE id=$1",
		bookID,
	)
	if err != nil {
		return err
	}
	return nil
}
