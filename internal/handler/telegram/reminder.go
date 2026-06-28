package telegram

import (
	"errors"
	"log"
	"strings"
	"time"

	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleSetReminder(chatID, telegramID int64, query string) error {
	if query == "" {
		msg := tgbotapi.NewMessage(chatID, "Использование:\n/setreminder указываете время напоминания\nНапример, /setreminder 21:30")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return err
	}

	_, err := time.Parse("15:04", query)
	if err != nil {
		return service.ErrInvalidReminderTime
	}

	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return err
	}

	err = h.reminderService.SetReminder(user.ID, query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (h *Handler) handleReminder(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	reminder, err := h.reminderService.GetReminder(user.ID)
	if err == nil {
		var builder strings.Builder

		builder.WriteString("🔔	Напоминание ")
		if reminder.IsEnabled == true {
			builder.WriteString("включено\n\n")
		} else {
			builder.WriteString("отключено\n\n")
		}
		builder.WriteString("🕒	Время: ")
		builder.WriteString(reminder.ReminderTime.Format("15:04"))

		msg := tgbotapi.NewMessage(chatID, builder.String())
		_, sendErr := h.bot.Send(msg)
		if sendErr != nil {
			log.Println(err)
			return
		}
		return
	}

	if errors.Is(err, service.ErrReminderNotFound) {
		msg := tgbotapi.NewMessage(chatID, "У вас пока нет активного напоминания\nИспользуйте \n\n/setreminder 21:00")
		_, sendErr := h.bot.Send(msg)
		if sendErr != nil {
			log.Println(err)
			return
		}
		return
	}
}

func (h *Handler) handleDisableReminder(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.reminderService.DisableReminder(user.ID)
	if err != nil {
		if errors.Is(err, service.ErrReminderNotFound) {
			msg := tgbotapi.NewMessage(chatID, "Нечего включать — напоминание уже включено")
			_, sendErr := h.bot.Send(msg)
			if sendErr != nil {
				log.Println(err)
				return
			}
			return
		}
		return
	}
}

func (h *Handler) handleEnableReminder(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.reminderService.EnableReminder(user.ID)
	if err != nil {
		if errors.Is(err, service.ErrReminderNotFound) {
			msg := tgbotapi.NewMessage(chatID, "Нечего включать — напоминание уже включено")
			_, sendErr := h.bot.Send(msg)
			if sendErr != nil {
				log.Println(err)
				return
			}
			return
		}
		return
	}
}

func (h *Handler) handleReminderMenu(chatID, telegramID int64) {
	user, err := h.bookService.GetUserByTelegramID(telegramID)
	if err != nil {
		log.Println(err)
		return
	}

	reminder, err := h.reminderService.GetReminder(user.ID)
	if err != nil {
		if errors.Is(err, service.ErrReminderNotFound) {
			var builder strings.Builder

			builder.WriteString("🔔 Напоминание\n\n")
			builder.WriteString("У вас пока нет активного напоминания")

			keyboard := h.buildReminderKeyboard(reminder)

			msg := tgbotapi.NewMessage(chatID, builder.String())
			msg.ReplyMarkup = keyboard
			_, err = h.bot.Send(msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println(err)
		return
	}

	text := h.buildReminderCard(reminder)
	keyboard := h.buildReminderKeyboard(reminder)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}
