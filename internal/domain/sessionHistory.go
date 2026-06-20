package domain

import "time"

type SessionHistory struct {
	BookTitle string
	PagesRead int
	Duration  time.Duration
	Date      time.Time
}

type SessionHistoryRow struct {
	BookTitle  string
	StartPage  int
	EndPage    int
	StartedAt  time.Time
	FinishedAt time.Time
}
