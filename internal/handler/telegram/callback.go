package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCallback(update tgbotapi.Update) {
	text := update.CallbackQuery.Data

	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	openLibraryID := parts[1]

	log.Println(action)
	log.Println(openLibraryID)

	switch action {
	case "details":
		book, err := h.bookService.GetBookDetails(openLibraryID)
		if err != nil {
			log.Println(err)
			return
		}

		var buttons []tgbotapi.InlineKeyboardButton
		var rows [][]tgbotapi.InlineKeyboardButton
		var builder strings.Builder
		builder.WriteString("📖")
		builder.WriteString(book.Title)
		builder.WriteString("\n\n")
		builder.WriteString("ID: ")
		builder.WriteString(book.OpenLibraryID)

		button := tgbotapi.NewInlineKeyboardButtonData("Добавить", "add:"+book.OpenLibraryID)
		button2 := tgbotapi.NewInlineKeyboardButtonData("Назад", "back:"+book.OpenLibraryID)
		buttons = append(buttons, button)
		buttons = append(buttons, button2)
		rows = append(rows, buttons)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, builder.String())
		msg.ReplyMarkup = keyboard

		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

	case "add":
	case "back":
	}
}
