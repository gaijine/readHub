package telegram

import (
	"strings"

	"readHub/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) buildReminderCard(reminder domain.Reminder) string {
	var builder strings.Builder

	builder.WriteString("🔔 Напоминание\n\n")
	builder.WriteString("Статус: ")
	if reminder.IsEnabled == true {
		builder.WriteString("✅ Включено\n\n")
	} else {
		builder.WriteString("❌ Отключено\n\n")
	}

	builder.WriteString("🕒 Время: ")
	builder.WriteString(reminder.ReminderTime.Format("15:04"))

	return builder.String()
}

func (h *Handler) buildReminderKeyboard(reminder domain.Reminder) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttonStatus tgbotapi.InlineKeyboardButton

	buttonSetTime := tgbotapi.NewInlineKeyboardButtonData("🕒 Установить время", "reminder:set")

	if reminder.IsEnabled {
		buttonStatus = tgbotapi.NewInlineKeyboardButtonData("🔴 Отключить", "reminder:off")
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonSetTime})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonStatus})
	} else {
		buttonStatus = tgbotapi.NewInlineKeyboardButtonData("🟢 Включить", "reminder:on")
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonStatus})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonSetTime})
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (h *Handler) updateReminderCard(chatID int64, messageID int, reminder domain.Reminder) error {
	text := h.buildReminderCard(reminder)
	keyboard := h.buildReminderKeyboard(reminder)

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ReplyMarkup = &keyboard

	_, err := h.bot.Send(edit)
	if err != nil {
		return err
	}
	return nil
}
