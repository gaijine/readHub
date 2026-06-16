package telegram

import (
	"strconv"
	"strings"

	"readHub/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) buildBookCard(book domain.Book) string {
	var builder strings.Builder
	builder.WriteString("📖 ")
	builder.WriteString(book.Title)
	builder.WriteString("\n\n")

	builder.WriteString("Автор:		")
	builder.WriteString(book.Author)
	builder.WriteString("\n")

	builder.WriteString("Статус:	 ")
	builder.WriteString(string(book.Status))
	builder.WriteString("\n")

	builder.WriteString("Прогресс:     ")
	builder.WriteString(strconv.Itoa(book.CurrentPage))
	builder.WriteString(" / ")
	builder.WriteString(strconv.Itoa(book.TotalPages))

	return builder.String()
}

func (h *Handler) buildBookKeyboard(bookID int64) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	buttonWant := tgbotapi.NewInlineKeyboardButtonData("📚 Хочу", "status:want:"+strconv.FormatInt(bookID, 10))
	buttonReading := tgbotapi.NewInlineKeyboardButtonData("📖 Читаю", "status:reading:"+strconv.FormatInt(bookID, 10))
	buttonCompleted := tgbotapi.NewInlineKeyboardButtonData("✅ Прочитано", "status:completed:"+strconv.FormatInt(bookID, 10))
	buttonUpdateProgress := tgbotapi.NewInlineKeyboardButtonData("📄 Обновить прогресс", "progress:"+strconv.FormatInt(bookID, 10))
	buttonDelete := tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", "delete:"+strconv.FormatInt(bookID, 10))

	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonWant})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonReading})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonCompleted})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonUpdateProgress})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonDelete})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
