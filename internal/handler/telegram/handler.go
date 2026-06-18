package telegram

import (
	"log"
	"strconv"
	"strings"

	"readHub/internal/domain"
	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProgressState struct {
	BookID    int64
	MessageID int
}

type Handler struct {
	bookService   service.BookService
	bot           *tgbotapi.BotAPI
	searchCache   map[int64][]domain.SearchBook // будет хранится результат поиска в кеше, мол список книг после поиска
	progressState map[int64]ProgressState       // будет хранится телеграм айди как ключь, айди книги и сообщения как значение, чтоб понимать у кого какую книгу обновлять
	// для того чтоб после нажатия "обновить прогресс" в памяти сохранялось h.progressState[8798127434] = {3, 23}
	// что значит Пользователь 8798127434 сейчас вводит прогресс для книги 3
	sessionService service.SessionService
}

func NewHandler(bookService service.BookService, bot *tgbotapi.BotAPI, sessionService service.SessionService) *Handler {
	return &Handler{
		bookService:    bookService,
		bot:            bot,
		searchCache:    make(map[int64][]domain.SearchBook),
		progressState:  make(map[int64]ProgressState),
		sessionService: sessionService,
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

	state, exists := h.progressState[telegramID]
	if exists {
		page, err := strconv.Atoi(text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Введите число")
			_, err = h.bot.Send(msg)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.bookService.UpdateProgress(user.ID, state.BookID, page)
		if err != nil {
			log.Println(err)

			msg := tgbotapi.NewMessage(chatID, err.Error())
			_, err = h.bot.Send(msg)
			if err != nil {
				log.Println(err)
				return
			}

			return
		}

		book, err := h.bookService.GetBookByID(state.BookID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateBookCard(chatID, state.MessageID, book)
		if err != nil {
			log.Println(err)
			return
		}

		delete(h.progressState, telegramID) // удаляем состояние из мапы (очищаем)

		return
	}

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
