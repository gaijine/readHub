package repository

import "readHub/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
	GetByTelegramID(telegramID int64) (domain.User, error)
	GetByID(userID int64) (domain.User, error)
}
