package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *Handler) buildMainMenu() tgbotapi.ReplyKeyboardMarkup {
	row1 := tgbotapi.NewKeyboardButtonRow( // две кнопки в одной строке
		tgbotapi.NewKeyboardButton("🔍 Поиск книги"), // кнопка клавиатуры
		tgbotapi.NewKeyboardButton("📚 Мои книги"),
	)

	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📊 Статистика"),
		tgbotapi.NewKeyboardButton("📖 История"),
	)

	row3 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔔 Напоминание"),
	)

	// собираем клавиатуру по строкам
	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	keyboard.ResizeKeyboard = true // для нормализации размера с пользовательскими настройками если не ошибаюсь

	return keyboard
}
