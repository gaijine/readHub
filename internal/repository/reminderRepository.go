package repository

import (
	"time"

	"readHub/internal/domain"
)

type ReminderRepository interface {
	Create(reminder domain.Reminder) error
	Update(userID int64, newTime string) error
	Disable(userID int64) error
	GetByUserID(userID int64) (domain.Reminder, error)
	GetAllEnabled() ([]domain.Reminder, error)
	UpdateLastSent(userID int64, sentAt time.Time) error
}
