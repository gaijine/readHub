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
	GetDueReminders() ([]domain.ReminderNotification, error)
	DisableReminder(userID int64) error
	EnableReminder(userID int64) error
	UpdateLastSent(userID int64) error
}

type reminderService struct {
	repo     repository.ReminderRepository
	userRepo repository.UserRepository
}

func NewReminderService(repo repository.ReminderRepository, userRepo repository.UserRepository) ReminderService {
	return &reminderService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (r *reminderService) SetReminder(userID int64, query string) error {
	t, err := time.Parse("15:04", query)
	if err != nil {
		return err
	}

	_, err = r.repo.GetByUserID(userID)

	if err == nil { // напоминание есть, новое не создаем а обновляем
		err = r.repo.Update(userID, t)
		if err != nil {
			return err
		}
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) { // бд не вернула ни одной строки, ничего нет создаем

		reminder := domain.Reminder{
			UserID:       userID,
			ReminderTime: t,
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
	if err == nil {
		return reminder, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Reminder{}, ErrReminderNotFound
	}

	return domain.Reminder{}, err
}

// кому пора отправить уведомл
func (r *reminderService) GetDueReminders() ([]domain.ReminderNotification, error) {
	var dueReminders []domain.ReminderNotification
	reminders, err := r.repo.GetAllEnabled()
	if err != nil {
		return nil, err
	}
	now := time.Now().Format("15:04")
	today := time.Now().Format("2006.01.02")
	for _, v := range reminders {

		if v.ReminderTime.Format("15:04") != now { // время не совпадает идем дальше по слайсу полученных реминдерс
			continue
		}
		// совпадает
		if v.LastSentAt == nil { // он же указатель, если нил, то увед не отправлялось, добавляем
			user, err := r.userRepo.GetByID(v.UserID)
			if err != nil {
				return nil, err
			}
			reminder := domain.ReminderNotification{
				UserID: user.ID,
				ChatID: user.TelegramID,
			}
			dueReminders = append(dueReminders, reminder)
			continue
		}
		// отправлялось, сравнивать будем даты
		lastSent := v.LastSentAt.Format("2006.01.02")
		if lastSent != today { // даты не совпадают добавляем для дальней отпр
			user, err := r.userRepo.GetByID(v.UserID)
			if err != nil {
				return nil, err
			}

			reminder := domain.ReminderNotification{
				UserID: user.ID,
				ChatID: user.TelegramID,
			}
			dueReminders = append(dueReminders, reminder)
		}
	}
	return dueReminders, nil
}

func (r *reminderService) DisableReminder(userID int64) error {
	err := r.repo.Disable(userID)
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrReminderNotFound
	}
	return err
}

func (r *reminderService) EnableReminder(userID int64) error {
	err := r.repo.Enable(userID)
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrReminderNotFound
	}
	return err
}

func (r *reminderService) UpdateLastSent(userID int64) error {
	err := r.repo.UpdateLastSent(userID, time.Now())
	if err != nil {
		return err
	}
	return nil
}
