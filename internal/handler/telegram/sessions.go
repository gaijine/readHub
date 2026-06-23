package telegram

import (
	"log"
	"math"
	"strconv"
	"strings"

	"readHub/internal/domain"

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

	text, keyboard := h.buildSessionsPage(sessions, 0)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h *Handler) buildSessionsPage(sessions []domain.SessionHistory, page int) (string, tgbotapi.InlineKeyboardMarkup) {
	pageSize := 5
	totalPages := int(math.Ceil(float64(len(sessions)) / float64(pageSize)))
	start := page * pageSize
	end := start + pageSize
	if end > len(sessions) {
		end = len(sessions)
	}
	pageSession := sessions[start:end]

	var builder strings.Builder
	var buttons []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, v := range pageSession {
		builder.WriteString("📖 ")
		builder.WriteString(v.BookTitle)
		builder.WriteString("\n")
		builder.WriteString("📄 ")
		builder.WriteString(strconv.Itoa(v.PagesRead))
		builder.WriteString(" стр.\n")

		totalMinutes := int(v.Duration.Minutes())
		hours := totalMinutes / 60
		minutes := totalMinutes % 60
		if hours > 0 {
			builder.WriteString("⏱ ")
			builder.WriteString(strconv.Itoa(hours))
			builder.WriteString(" ч ")
			builder.WriteString(strconv.Itoa(minutes))
			builder.WriteString(" мин\n")
		} else {
			builder.WriteString("⏱ ")
			builder.WriteString(strconv.Itoa(minutes))
			builder.WriteString(" мин\n")
		}

		builder.WriteString("📅 ")
		builder.WriteString(v.Date.Format("02.01.2006"))
		builder.WriteString("\n\n")
	}

	if totalPages > 1 {
		if page > 0 {
			buttonBack := tgbotapi.NewInlineKeyboardButtonData("⬅️", "sessions:page:"+strconv.Itoa(page-1))
			buttons = append(buttons, buttonBack)
		}

		pageButton := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(page+1)+"/"+strconv.Itoa(totalPages), "noop") // noop- no operations
		buttons = append(buttons, pageButton)

		if page < totalPages-1 {
			buttonNext := tgbotapi.NewInlineKeyboardButtonData("➡️", "sessions:page:"+strconv.Itoa(page+1))
			buttons = append(buttons, buttonNext)
		}
		rows = append(rows, buttons)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return builder.String(), keyboard
}
