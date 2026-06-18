package service

import (
	"errors"

	"readHub/internal/domain"
	"readHub/internal/repository"

	"github.com/jackc/pgx/v5"
)

type SessionService interface {
	StartSession(bookID, userID int64) error
	FinishSession(userID int64) error
}

type sessionService struct {
	sessionRepo repository.SessionRepository
	bookService BookService
}

func NewSessionService(sessionRepo repository.SessionRepository, bookService BookService) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		bookService: bookService,
	}
}

func (s *sessionService) StartSession(bookID, userID int64) error {
	_, err := s.sessionRepo.GetActiveSession(userID)
	if err == nil { // активная сессия уже есть, новую не создаем
		return ErrActiveSessionIsExist
	}
	if errors.Is(err, pgx.ErrNoRows) { // проверка на то что бд не вернула ни одной строки, активной сессии нет, можно создать новую
		book, err := s.bookService.GetBookByID(bookID)
		if err != nil {
			return err
		}

		err = s.sessionRepo.CreateSession(userID, bookID, book.CurrentPage)
		if err != nil {
			return err
		}

		err = s.bookService.UpdateStatus(userID, bookID, domain.StatusReading) // автоматически после создания сессии меняем статус
		if err != nil {
			return err
		}
		return nil
	}
	return err // если есть другая ошибка (например, соединение потеряно, бд недоступна, ошибка скьюл) то возвращаем её
}

func (s *sessionService) FinishSession(userID int64) error {
	session, err := s.sessionRepo.GetActiveSession(userID)
	if errors.Is(err, pgx.ErrNoRows) { // активной сессии нет
		return ErrActiveSessionNotFound
	}
	if err != nil { // другая ошибка соединение потеряно, бд недоступна, ошибка скьюл
		return err
	}

	book, err := s.bookService.GetBookByID(session.BookID)
	if err != nil {
		return err
	}

	err = s.sessionRepo.FinishSession(session.ID, book.CurrentPage)
	if err != nil {
		return err
	}
	return nil
}
