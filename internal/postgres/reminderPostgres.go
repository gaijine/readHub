package postgres

import (
	"context"
	"time"

	"readHub/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresReminderRepository struct {
	db *pgxpool.Pool
}

func NewReminderRepository(db *pgxpool.Pool) *PostgresReminderRepository {
	return &PostgresReminderRepository{
		db: db,
	}
}

func (r *PostgresReminderRepository) Create(reminder domain.Reminder) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO reminders (user_id, reminder_time) VALUES ($1, $2)",
		reminder.UserID,
		reminder.ReminderTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresReminderRepository) Update(userID int64, newTime string) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE reminders SET reminder_time=$1, is_enabled=true WHERE user_id=$2",
		newTime,
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresReminderRepository) Disable(userID int64) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE reminders SET is_enabled=false WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresReminderRepository) GetByUserID(userID int64) (domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.QueryRow(context.Background(),
		"SELECT id, user_id, reminder_time, is_enabled, last_sent_at FROM reminders WHERE user_id=$1 LIMIT 1",
		userID,
	).Scan(&reminder.ID, &reminder.UserID, &reminder.ReminderTime, &reminder.IsEnabled, &reminder.LastSentAt)
	if err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}

func (r *PostgresReminderRepository) GetAllEnabled() ([]domain.Reminder, error) {
	var reminders []domain.Reminder
	rows, err := r.db.Query(context.Background(),
		"SELECT id, user_id, reminder_time, is_enabled, last_sent_at FROM reminders WHERE is_enabled=true",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reminder domain.Reminder
		err = rows.Scan(&reminder.ID, &reminder.UserID, &reminder.ReminderTime, &reminder.IsEnabled, &reminder.LastSentAt)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *PostgresReminderRepository) UpdateLastSent(userID int64, sentAt time.Time) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE reminders SET last_sent_at=$1 WHERE user_id=$2",
		sentAt,
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}
