package telegram

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleStats(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	stats, err := h.statsService.GetStats(user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	var builder strings.Builder

	builder.WriteString("📊 Статистика чтения\n\n")
	builder.WriteString("📚 Всего книг: ")
	builder.WriteString(strconv.Itoa(stats.TotalBooks))
	builder.WriteString("\n")
	builder.WriteString("📖 Сейчас читаю: ")
	builder.WriteString(strconv.Itoa(stats.ReadingBooks))
	builder.WriteString("\n")
	builder.WriteString("✅ Прочитано: ")
	builder.WriteString(strconv.Itoa(stats.CompletedBooks))
	builder.WriteString("\n\n")
	builder.WriteString("⏱ Сессий чтения: ")
	builder.WriteString(strconv.Itoa(stats.TotalSessions))
	builder.WriteString("\n")
	builder.WriteString("📄 Страниц прочитано: ")
	builder.WriteString(strconv.Itoa(stats.PagesRead))

	msg := tgbotapi.NewMessage(chatID, builder.String())
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
