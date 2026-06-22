package telegram

import (
	"log"

	"readHub/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const startMessage = `📚 Добро пожаловать в ReadHub!

ReadHub поможет сохранять книги, которые вы хотите прочитать, отслеживать прогресс чтения и вести свою личную библиотеку.

Находите книги, добавляйте их в список и отмечайте свой прогресс по мере чтения.

Приятного чтения! 📖
	`

func (h *Handler) handleStart(chatID, telegramID int64, username string) {
	user := domain.User{
		TelegramID: telegramID,
		Username:   username,
	}

	_, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		err = h.bookService.CreateUser(user)
		if err != nil {
			log.Println(err)
			return
		}
	}

	msg := tgbotapi.NewMessage(chatID, startMessage)
	keyboard := h.buildMainMenu()
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
