package telegram

import (
	"strings"

	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bookService service.BookService
	bot         *tgbotapi.BotAPI
}

func NewHandler(bookService service.BookService, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		bookService: bookService,
		bot:         bot,
	}
}

func (h *Handler) Run() {
	u := tgbotapi.NewUpdate(0)
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			h.handleMessage(update)
		}
		// if update.Message.Text != "" {
		// 	h.handleMessage(update)
		// }
		if update.CallbackQuery != nil {
			h.handleCallback(update)
		}
	}
}

func (h *Handler) handleMessage(update tgbotapi.Update) {
	text := update.Message.Text      // получаем текст сообщения
	chatID := update.Message.Chat.ID // получаем айди чата

	parts := strings.Fields(text) // делим текст на слайс слов

	if len(parts) == 0 {
		return
	}

	command := parts[0]                   // берем первое слово
	query := strings.Join(parts[1:], " ") // собираем оставшуюся часть без команды в отдельную строку

	switch command {
	case "/start":
		h.handleStart(chatID)
	case "/search":
		h.handleSearch(chatID, query)
	}
}
