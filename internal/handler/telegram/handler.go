package telegram

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"readHub/internal/domain"
	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReminderState struct {
	SentMessageID int
	MessageID     int
}

type LibraryState struct {
	MessageID int
	Page      int
}

type DeleteState struct {
	SentMessageID int
}

type SearchState struct{}

type PagesState struct {
	BookID        int64
	MessageID     int // айди карточки книги
	SentMessageID int // айди сообщения которое отправляем пользователю
}

type ProgressState struct {
	BookID        int64
	MessageID     int
	SentMessageID int
}

type Handler struct {
	bookService   service.BookService
	bot           *tgbotapi.BotAPI
	searchCache   map[int64][]domain.SearchBook // будет хранится результат поиска в кеше, мол список книг после поиска
	progressState map[int64]ProgressState       // будет хранится телеграм айди как ключь, айди книги и сообщения как значение, чтоб понимать у кого какую книгу обновлять
	// для того чтоб после нажатия "обновить прогресс" в памяти сохранялось h.progressState[8798127434] = {3, 23}
	// что значит Пользователь 8798127434 сейчас вводит прогресс для книги 3
	sessionService  service.SessionService
	pagesState      map[int64]PagesState
	statsService    service.StatsService
	searchState     map[int64]SearchState
	deleteState     map[int64]DeleteState
	libraryState    map[int64]LibraryState
	reminderService service.ReminderService
	reminderState   map[int64]ReminderState
}

func NewHandler(bookService service.BookService, bot *tgbotapi.BotAPI, sessionService service.SessionService, statsService service.StatsService, reminderService service.ReminderService) *Handler {
	return &Handler{
		bookService:     bookService,
		bot:             bot,
		searchCache:     make(map[int64][]domain.SearchBook),
		progressState:   make(map[int64]ProgressState),
		sessionService:  sessionService,
		pagesState:      make(map[int64]PagesState),
		statsService:    statsService,
		searchState:     make(map[int64]SearchState),
		deleteState:     make(map[int64]DeleteState),
		libraryState:    make(map[int64]LibraryState),
		reminderService: reminderService,
		reminderState:   make(map[int64]ReminderState),
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

		deleteConfig := tgbotapi.NewDeleteMessage(chatID, state.SentMessageID)
		_, err = h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}

		deleteConfig1 := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		_, err = h.bot.Request(deleteConfig1)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}

		delete(h.progressState, telegramID) // удаляем состояние из мапы (очищаем)

		return
	}

	statePage, exists := h.pagesState[telegramID]
	if exists {
		totalPage, err := strconv.Atoi(text)
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

		err = h.bookService.UpdateTotalPages(user.ID, statePage.BookID, totalPage)
		if err != nil {
			log.Println(err)
			return
		}

		book, err := h.bookService.GetBookByID(statePage.BookID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateBookCard(chatID, statePage.MessageID, book)
		if err != nil {
			log.Println(err)
			return
		}

		deleteConfig := tgbotapi.NewDeleteMessage(chatID, statePage.SentMessageID)
		_, err = h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}
		deleteConfig1 := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		_, err = h.bot.Request(deleteConfig1)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}

		delete(h.pagesState, telegramID)

		return
	}

	_, exists = h.searchState[telegramID]
	if exists {
		found := h.handleSearch(chatID, telegramID, text)
		if found {
			delete(h.searchState, telegramID)
		}
		return
	}

	remState, exists := h.reminderState[telegramID]
	if exists {
		err := h.handleSetReminder(chatID, telegramID, text)
		if err != nil {
			if errors.Is(err, service.ErrInvalidReminderTime) {
				msg := tgbotapi.NewMessage(chatID, "Формат времени неверный. \nНапишите время в формате ЧЧ:ММ, \nнапример, *23:21* или *6:06*")
				_, sendErr := h.bot.Send(msg)
				if sendErr != nil {
					log.Println(sendErr)
					return
				}
				return
			}
			log.Println(err)
			return
		}

		deleteConfig := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		_, err = h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}

		deleteConfig2 := tgbotapi.NewDeleteMessage(chatID, remState.SentMessageID)
		_, err = h.bot.Request(deleteConfig2)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
		}

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		reminder, err := h.reminderService.GetReminder(user.ID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateReminderCard(chatID, remState.MessageID, reminder)
		if err != nil {
			log.Println(err)
			return
		}

		delete(h.reminderState, telegramID)
		return
	}

	switch text {
	case "🔍 Поиск книги":
		h.searchState[telegramID] = SearchState{}

		msg := tgbotapi.NewMessage(chatID, "Введите название книги или автора")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
		return
	case "📚 Мои книги":
		h.handleMyBooks(chatID, telegramID)
		return
	case "📊 Статистика":
		h.handleStats(chatID, telegramID)
		return
	case "📖 История":
		h.handleSessions(chatID, telegramID)
		return
	case "🔔 Напоминание":
		h.handleReminderMenu(chatID, telegramID)
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
	case "/stats":
		h.handleStats(chatID, telegramID)
	case "/sessions":
		h.handleSessions(chatID, telegramID)
	case "/setreminder":
		h.handleSetReminder(chatID, telegramID, query)
	case "/reminder":
		h.handleReminder(chatID, telegramID)
	case "/reminderoff":
		h.handleDisableReminder(chatID, telegramID)
	}
}
