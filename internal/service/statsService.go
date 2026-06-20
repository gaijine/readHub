package service

import (
	"time"

	"readHub/internal/domain"
	"readHub/internal/repository"
)

type StatsService interface {
	GetStats(userID int64) (domain.ReadingStats, error)
}

type statsService struct {
	bookRepo    repository.BookRepository
	sessionRepo repository.SessionRepository
}

func NewStatsService(bookRepo repository.BookRepository, sessionRepo repository.SessionRepository) StatsService {
	return &statsService{
		bookRepo:    bookRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *statsService) GetStats(userID int64) (domain.ReadingStats, error) {
	var (
		completionRate int
		averagePage    int
		averageSession time.Duration
		totalDuration  time.Duration
	)

	totalBooks, err := s.bookRepo.CountByUserID(userID)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	completedBooks, err := s.bookRepo.CountByStatus(userID, domain.StatusCompleted)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	readingBooks, err := s.bookRepo.CountByStatus(userID, domain.StatusReading)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	totalSessions, err := s.sessionRepo.CountByUserID(userID)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	pagesRead, err := s.sessionRepo.GetPagesRead(userID)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	sessionsRow, err := s.sessionRepo.GetListSessions(userID)
	if err != nil {
		return domain.ReadingStats{}, err
	}

	for _, v := range sessionsRow {
		duration := v.FinishedAt.Sub(v.StartedAt)
		totalDuration += duration
	}

	if totalSessions > 0 {
		averageSession = totalDuration / time.Duration(totalSessions)
	}

	if totalBooks > 0 {
		completionRate = completedBooks * 100 / totalBooks
	}

	if totalSessions > 0 {
		averagePage = pagesRead / totalSessions
	}

	stats := domain.ReadingStats{
		TotalBooks:             totalBooks,
		CompletedBooks:         completedBooks,
		ReadingBooks:           readingBooks,
		TotalSessions:          totalSessions,
		PagesRead:              pagesRead,
		CompletionRate:         completionRate,
		AveragePagesPerSession: averagePage,
		TotalReadingTime:       totalDuration,
		AverageSessionDuration: averageSession,
	}

	return stats, nil
}
