package telegram

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleMyBooks(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	books, err := h.bookService.GetUserBooks(user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if len(books) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📚 Ваша библиотека пока пуста.\n\nИспользуйте /search чтобы найти и добавить книгу.")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	var builder strings.Builder
	var buttons []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for i, book := range books {
		builder.WriteString("[")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString("]		")
		builder.WriteString(book.Title)
		builder.WriteString("\n")
		builder.WriteString("Автор:		")
		builder.WriteString(book.Author)
		builder.WriteString("\n")
		builder.WriteString("Статус:   	 ")
		builder.WriteString(string(book.Status))
		builder.WriteString("\n\n")

		button := tgbotapi.NewInlineKeyboardButtonData("["+strconv.Itoa(i+1)+"]", "mybook:"+strconv.FormatInt(book.ID, 10))
		buttons = append(buttons, button)
	}
	rows = append(rows, buttons)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, builder.String())
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
