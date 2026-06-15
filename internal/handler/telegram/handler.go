package telegram

import (
	"strings"

	"readHub/internal/domain"
	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bookService service.BookService
	bot         *tgbotapi.BotAPI
	searchCache map[int64][]domain.SearchBook // будет хранится результат поиска в кеше, мол список книг после поиска
}

func NewHandler(bookService service.BookService, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		bookService: bookService,
		bot:         bot,
		searchCache: make(map[int64][]domain.SearchBook),
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
	text := update.Message.Text          // получаем текст сообщения
	chatID := update.Message.Chat.ID     // получаем айди чата
	telegramID := update.Message.From.ID // получаем телеграм id
	username := update.Message.From.UserName

	parts := strings.Fields(text) // делим текст на слайс слов

	if len(parts) == 0 {
		return
	}

	command := parts[0]                   // берем первое слово
	query := strings.Join(parts[1:], " ") // собираем оставшуюся часть без команды в отдельную строку

	switch command {
	case "/start":
		h.handleStart(chatID, telegramID, username)
	case "/search":
		h.handleSearch(chatID, telegramID, query)
	case "/mybooks":
		h.handleMyBooks(chatID, telegramID)
	}
}
