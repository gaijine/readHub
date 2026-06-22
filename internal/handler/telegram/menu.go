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

	// собираем клавиатуру по строкам
	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	keyboard.ResizeKeyboard = true

	return keyboard
}
