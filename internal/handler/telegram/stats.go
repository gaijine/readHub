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
	builder.WriteString("✅ Прочитано: ")
	builder.WriteString(strconv.Itoa(stats.CompletedBooks))
	builder.WriteString(" из ")
	builder.WriteString(strconv.Itoa(stats.TotalBooks))
	builder.WriteString(" книг\n")
	builder.WriteString("📖 Читаю: ")
	builder.WriteString(strconv.Itoa(stats.ReadingBooks))
	builder.WriteString("\n\n")
	builder.WriteString("⏱ Сессий чтения: ")
	builder.WriteString(strconv.Itoa(stats.TotalSessions))
	builder.WriteString("\n")
	builder.WriteString("📄 Страниц прочитано: ")
	builder.WriteString(strconv.Itoa(stats.PagesRead))
	builder.WriteString("\n")
	builder.WriteString("📈 Среднее за сессию: ")
	builder.WriteString(strconv.Itoa(stats.AveragePagesPerSession))
	builder.WriteString(" стр.\n\n")

	totalMinutes := int(stats.TotalReadingTime.Minutes())
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	if hours > 0 {
		builder.WriteString("⌛ Время чтения: ")
		builder.WriteString(strconv.Itoa(hours))
		builder.WriteString(" ч ")
		builder.WriteString(strconv.Itoa(minutes))
		builder.WriteString(" мин\n")
	} else {
		builder.WriteString("⌛ Время чтения: ")
		builder.WriteString(strconv.Itoa(minutes))
		builder.WriteString(" мин\n")
	}

	totalMinutesAverage := int(stats.AverageSessionDuration.Minutes())
	hoursAverage := totalMinutesAverage / 60
	minutesAverage := totalMinutesAverage % 60
	if hoursAverage > 0 {
		builder.WriteString("⏱ Средняя сессия: ")
		builder.WriteString(strconv.Itoa(hoursAverage))
		builder.WriteString(" ч ")
		builder.WriteString(strconv.Itoa(minutesAverage))
		builder.WriteString(" мин\n")
	} else {
		builder.WriteString("⏱ Средняя сессия: ")
		builder.WriteString(strconv.Itoa(minutesAverage))
		builder.WriteString(" мин\n\n")
	}

	builder.WriteString("🔥 Процент завершения: ")
	builder.WriteString(strconv.Itoa(stats.CompletionRate))
	builder.WriteString("%")

	msg := tgbotapi.NewMessage(chatID, builder.String())
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
