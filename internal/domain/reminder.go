package domain

import "time"

type Reminder struct {
	ID           int64
	UserID       int64
	ReminderTime string
	IsEnabled    bool
	LastSentAt   *time.Time
}
