package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const startMessage = `📚 Добро пожаловать в ReadHub!

ReadHub поможет сохранять книги, которые вы хотите прочитать, отслеживать прогресс чтения и вести свою личную библиотеку.

Находите книги, добавляйте их в список и отмечайте свой прогресс по мере чтения.

Приятного чтения! 📖
	`

func (h *Handler) handleStart(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, startMessage)

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
