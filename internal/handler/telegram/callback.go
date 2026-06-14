package telegram

import (
	"log"
	"strings"

	"readHub/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
)

func (h *Handler) handleCallback(update tgbotapi.Update) {
	text := update.CallbackQuery.Data

	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return
	}

	action := parts[0]
	openLibraryID := parts[1]
	telegramID := update.CallbackQuery.From.ID

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
		books := h.searchCache[telegramID]
		log.Println(books)

		var selectedBook domain.SearchBook

		for _, book := range books {
			if book.OpenLibraryID == openLibraryID {
				selectedBook = book
				break
			}
		}
		pp.Println(selectedBook)

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(user)
	case "back":
	}
}
