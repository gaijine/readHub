package domain

import "time"

type Reminder struct {
	ID           int64
	UserID       int64
	ReminderTime time.Time
	IsEnabled    bool
	LastSentAt   *time.Time
}

type ReminderNotification struct {
	UserID int64
	ChatID int64
}
