package domain

import "time"

type ReadingSession struct {
	ID         int64
	BookID     int64
	UserID     int64
	StartedAt  time.Time
	FinishedAt *time.Time
	StartPage  int
	EndPage    *int
}
