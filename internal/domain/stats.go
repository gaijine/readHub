package domain

import "time"

type ReadingStats struct {
	TotalBooks             int
	CompletedBooks         int
	ReadingBooks           int
	TotalSessions          int
	PagesRead              int
	CompletionRate         int
	AveragePagesPerSession int
	TotalReadingTime       time.Duration
	AverageSessionDuration time.Duration
}
