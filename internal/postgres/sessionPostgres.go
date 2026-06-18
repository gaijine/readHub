package postgres

import (
	"context"

	"readHub/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{
		db: db,
	}
}

func (r *PostgresSessionRepository) CreateSession(userID, bookID int64, startPage int) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO reading_sessions (book_id, user_id, start_page) VALUES ($1, $2, $3)",
		bookID,
		userID,
		startPage,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSessionRepository) FinishSession(sessionID int64, endPage int) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE reading_sessions SET end_page=$2, finished_at=NOW() WHERE id=$1",
		sessionID,
		endPage,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSessionRepository) GetActiveSession(userID int64) (domain.ReadingSession, error) {
	var session domain.ReadingSession
	err := r.db.QueryRow(context.Background(),
		"SELECT id, book_id, user_id, started_at, finished_at, start_page, end_page FROM reading_sessions WHERE user_id=$1 AND finished_at IS NULL",
		userID,
	).Scan(&session.ID, &session.BookID, &session.UserID, &session.StartedAt, &session.FinishedAt, &session.StartPage, &session.EndPage)
	if err != nil {
		return domain.ReadingSession{}, err
	}
	return session, nil
}
