package repository

import "readHub/internal/domain"

type SessionRepository interface {
	CreateSession(userID, bookID int64, startPage int) error
	FinishSession(sessionID int64, endPage int) error
	GetActiveSession(userID int64) (domain.ReadingSession, error)
}
