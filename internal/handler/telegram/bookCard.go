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

	builder.WriteString("👤 Автор:		")
	builder.WriteString(book.Author)
	builder.WriteString("\n")

	builder.WriteString("📚 Статус:	 ")
	builder.WriteString(string(book.Status))
	builder.WriteString("\n")

	if book.TotalPages > 0 {
		percent := book.CurrentPage * 100 / book.TotalPages
		builder.WriteString("📄 Прогресс:     ")
		builder.WriteString(strconv.Itoa(book.CurrentPage))
		builder.WriteString(" / ")
		builder.WriteString(strconv.Itoa(book.TotalPages))
		builder.WriteString("\n")
		builder.WriteString("📈 ")
		builder.WriteString(" (")
		builder.WriteString(strconv.Itoa(percent))
		builder.WriteString("%)")
	} else {
		builder.WriteString("📄 Прогресс:     ")
		builder.WriteString(strconv.Itoa(book.CurrentPage))
		builder.WriteString(" стр.")
	}

	return builder.String()
}

func (h *Handler) buildBookKeyboard(userID int64, book domain.Book) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var readingButton tgbotapi.InlineKeyboardButton

	readingButton = tgbotapi.NewInlineKeyboardButtonData("▶ Начать чтение", "startsession:"+strconv.FormatInt(book.ID, 10))
	session, err := h.sessionService.GetActiveSession(userID)
	if err == nil {
		if session.BookID == book.ID {
			readingButton = tgbotapi.NewInlineKeyboardButtonData("⏹ Завершить чтение", "finishsession:"+strconv.FormatInt(book.ID, 10))
		}
	}

	buttonWant := tgbotapi.NewInlineKeyboardButtonData("📚 Хочу", "status:want:"+strconv.FormatInt(book.ID, 10))
	buttonReading := tgbotapi.NewInlineKeyboardButtonData("📖 Читаю", "status:reading:"+strconv.FormatInt(book.ID, 10))
	buttonCompleted := tgbotapi.NewInlineKeyboardButtonData("✅ Прочитано", "status:completed:"+strconv.FormatInt(book.ID, 10))
	buttonUpdateProgress := tgbotapi.NewInlineKeyboardButtonData("📄 Обновить прогресс", "progress:"+strconv.FormatInt(book.ID, 10))
	buttonDelete := tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", "delete:"+strconv.FormatInt(book.ID, 10))
	setPagesButton := tgbotapi.NewInlineKeyboardButtonData("📚 Указать страницы", "setpages:"+strconv.FormatInt(book.ID, 10))

	if book.Status == domain.StatusCompleted {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonWant})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonReading})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonCompleted})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonDelete})
		return tgbotapi.NewInlineKeyboardMarkup(rows...)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonWant})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonReading})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonCompleted})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{setPagesButton})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonUpdateProgress})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{readingButton})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonDelete})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (h *Handler) updateBookCard(chatID int64, messageID int, book domain.Book) error {
	text := h.buildBookCard(book)
	keyboard := h.buildBookKeyboard(book.UserID, book)

	if book.CoverURL == "" {
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		_, err := h.bot.Send(edit)
		if err != nil {
			return err
		}
	} else {
		edit := tgbotapi.NewEditMessageCaption(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		_, err := h.bot.Send(edit)
		if err != nil {
			return err
		}
	}
	return nil
}
