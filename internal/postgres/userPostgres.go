package postgres

import (
	"context"

	"readHub/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	repo := PostgresUserRepository{
		db: db,
	}
	return &repo
}

func (r *PostgresUserRepository) Create(user domain.User) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO users (telegram_id, username) VALUES($1, $2)",
		user.TelegramID,
		user.Username,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetByTelegramID(telegramID int64) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(context.Background(),
		"SELECT id, telegram_id, username, created_at FROM users WHERE telegram_id=$1",
		telegramID,
	).Scan(&user.ID, &user.TelegramID, &user.Username, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByID(userID int64) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(context.Background(),
		"SELECT id, telegram_id, username, created_at FROM users WHERE id=$1",
		userID,
	).Scan(&user.ID, &user.TelegramID, &user.Username, &user.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
