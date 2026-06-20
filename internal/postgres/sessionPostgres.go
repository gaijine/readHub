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

func (r *PostgresSessionRepository) CountByUserID(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM reading_sessions WHERE user_id=$1",
		userID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostgresSessionRepository) GetPagesRead(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(end_page - start_page), 0) FROM reading_sessions WHERE user_id=$1 AND end_page IS NOT NULL",
		userID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostgresSessionRepository) GetListSessions(userID int64) ([]domain.SessionHistoryRow, error) {
	var sessionsRow []domain.SessionHistoryRow

	rows, err := r.db.Query(context.Background(),
		`SELECT books.title, reading_sessions.start_page, reading_sessions.end_page, reading_sessions.started_at, reading_sessions.finished_at
FROM books
INNER JOIN reading_sessions ON books.id=reading_sessions.book_id 
WHERE reading_sessions.user_id=$1
AND reading_sessions.finished_at IS NOT NULL
ORDER BY reading_sessions.started_at DESC
LIMIT 10`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row domain.SessionHistoryRow
		err = rows.Scan(&row.BookTitle, &row.StartPage, &row.EndPage, &row.StartedAt, &row.FinishedAt)
		if err != nil {
			return nil, err
		}
		sessionsRow = append(sessionsRow, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return sessionsRow, nil
}
