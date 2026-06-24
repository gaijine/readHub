package service

import (
	"errors"
	"time"

	"readHub/internal/domain"
	"readHub/internal/repository"

	"github.com/jackc/pgx/v5"
)

type ReminderService interface {
	SetReminder(userID int64, time string) error
	GetReminder(userID int64) (domain.Reminder, error)
	GetDueReminders() ([]domain.Reminder, error)
	DisableReminder(userID int64) error
	UpdateLastSent(userID int64, sentAt time.Time) error
}

type reminderService struct {
	repo repository.ReminderRepository
}

func NewReminderService(repo repository.ReminderRepository) ReminderService {
	return &reminderService{
		repo: repo,
	}
}

func (r *reminderService) SetReminder(userID int64, time string) error {
	_, err := r.repo.GetByUserID(userID)

	if err == nil { // напоминание есть, новое не создаем а обновляем
		err = r.repo.Update(userID, time)
		if err != nil {
			return err
		}
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) { // бд не вернула ни одной строки, ничего нет создаем

		reminder := domain.Reminder{
			UserID:       userID,
			ReminderTime: time,
		}

		err = r.repo.Create(reminder)
		if err != nil {
			return err
		}
		return nil
	}
	return err // любая другая ошибка выведется здесь
}

func (r *reminderService) GetReminder(userID int64) (domain.Reminder, error) {
	reminder, err := r.repo.GetByUserID(userID)
	if err != nil {
		return domain.Reminder{}, err
	}
	return reminder, nil
}

func (r *reminderService) GetDueReminders() ([]domain.Reminder, error) {
	var dueReminders []domain.Reminder
	reminders, err := r.repo.GetAllEnabled()
	if err != nil {
		return nil, err
	}
	now := time.Now().Format("15:04")
	today := time.Now().Format("2006.01.02")
	for _, v := range reminders {
		if v.ReminderTime != now { // время не совпадает идем дальше
			continue
		}
		// совпадает
		if v.LastSentAt == nil { // он же указатель, если нил, то увед не отправлялось, добавляем
			dueReminders = append(dueReminders, v)
			continue
		}
		// отправлялось, сравнивать будем даты
		lastSent := v.LastSentAt.Format("2006.01.02")
		if lastSent != today { // даты не совпадают добавляем для дальней отпр
			dueReminders = append(dueReminders, v)
		}
	}
	return dueReminders, nil
}

func (r *reminderService) DisableReminder(userID int64) error {
	err := r.repo.Disable(userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *reminderService) UpdateLastSent(userID int64, sentAt time.Time) error {
	err := r.repo.UpdateLastSent(userID, sentAt)
	if err != nil {
		return err
	}
	return nil
}
