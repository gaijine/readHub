package telegram

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleSessions(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	sessions, err := h.sessionService.GetSessionHistory(user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if len(sessions) == 0 {
		msg := tgbotapi.NewMessage(chatID, "История ваших сессий чтения пока пуста")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	var builder strings.Builder
	for _, v := range sessions {
		builder.WriteString("📖 ")
		builder.WriteString(v.BookTitle)
		builder.WriteString("\n")
		builder.WriteString("📄 ")
		builder.WriteString(strconv.Itoa(v.PagesRead))
		builder.WriteString(" стр.\n")
		builder.WriteString("⏱ ")
		builder.WriteString(strconv.Itoa(int(v.Duration.Minutes())))
		builder.WriteString(" мин\n")
		builder.WriteString("📅 ")
		builder.WriteString(v.Date.Format("02.01.2006"))
		builder.WriteString("\n\n")
	}

	msg := tgbotapi.NewMessage(chatID, builder.String())
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
