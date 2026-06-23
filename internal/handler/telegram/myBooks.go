package telegram

import (
	"log"
	"math"
	"strconv"
	"strings"

	"readHub/internal/domain"

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

	text, keyboard := h.buildBooksPage(books, 0)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	sentMessage, err := h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

	h.libraryState[telegramID] = LibraryState{
		MessageID: sentMessage.MessageID,
		Page:      0,
	}
}

func (h *Handler) buildBooksPage(books []domain.Book, page int) (string, tgbotapi.InlineKeyboardMarkup) {
	pageSize := 4
	totalPages := int(math.Ceil(float64(len(books)) / float64(pageSize))) // округляем общее кол-во стр, даже 3,1=4
	start := page * pageSize
	end := start + pageSize
	if end > len(books) {
		end = len(books)
	}
	pageBooks := books[start:end]

	var builder strings.Builder
	var buttons []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for i, book := range pageBooks {
		builder.WriteString("[")
		builder.WriteString(strconv.Itoa(start + i + 1))
		builder.WriteString("]		")
		builder.WriteString(book.Title)
		builder.WriteString("\n")
		builder.WriteString("Автор:		")
		builder.WriteString(book.Author)
		builder.WriteString("\n")
		builder.WriteString("Статус:   	 ")
		builder.WriteString(string(book.Status))
		builder.WriteString("\n\n")

		button := tgbotapi.NewInlineKeyboardButtonData("["+strconv.Itoa(start+i+1)+"]", "mybook:"+strconv.FormatInt(book.ID, 10))
		buttons = append(buttons, button)

		if len(buttons) == 4 { // если кнопки 9/10 будут последними к примеру они по условию непридут и не добавяться
			rows = append(rows, buttons)
			buttons = nil
		}
	}

	if len(buttons) > 0 { // поэтому после цикла добавляем оставшиеся
		rows = append(rows, buttons)
	}

	var pageButtons []tgbotapi.InlineKeyboardButton

	if totalPages > 1 {
		if page > 0 {
			backButton := tgbotapi.NewInlineKeyboardButtonData("⬅️", "books:page:"+strconv.Itoa(page-1))
			pageButtons = append(pageButtons, backButton)
		}

		pagebutton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(page+1)+"/"+strconv.Itoa(totalPages), "noop")
		pageButtons = append(pageButtons, pagebutton)

		if page < totalPages-1 {
			nextButton := tgbotapi.NewInlineKeyboardButtonData("➡️", "books:page:"+strconv.Itoa(page+1))
			pageButtons = append(pageButtons, nextButton)

		}
		rows = append(rows, pageButtons)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return builder.String(), keyboard
}
