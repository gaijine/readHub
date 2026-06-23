package telegram

import (
	"log"
	"math"
	"strconv"
	"strings"

	"readHub/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCallback(update tgbotapi.Update) {
	text := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	telegramID := update.CallbackQuery.From.ID

	parts := strings.Split(text, ":")
	if len(parts) < 2 {
		return
	}

	action := parts[0]

	switch action {
	case "details":
		book, err := h.bookService.GetBookDetails(parts[1])
		if err != nil {
			log.Println(err)
			return
		}

		var buttons []tgbotapi.InlineKeyboardButton
		var rows [][]tgbotapi.InlineKeyboardButton
		var builder strings.Builder
		builder.WriteString("📖")
		builder.WriteString(book.Title)
		builder.WriteString("\n\n")
		builder.WriteString("ID: ")
		builder.WriteString(book.OpenLibraryID)

		button := tgbotapi.NewInlineKeyboardButtonData("Добавить", "add:"+book.OpenLibraryID)
		button2 := tgbotapi.NewInlineKeyboardButtonData("Назад", "back:"+book.OpenLibraryID)
		buttons = append(buttons, button)
		buttons = append(buttons, button2)
		rows = append(rows, buttons)

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(chatID, builder.String())
		msg.ReplyMarkup = keyboard

		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

	case "add":
		books := h.searchCache[telegramID]

		var selectedBook domain.SearchBook

		for _, book := range books {
			if book.OpenLibraryID == parts[1] {
				selectedBook = book
				break
			}
		}

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%+v\n", selectedBook)
		err = h.bookService.AddBook(user.ID, selectedBook)
		if err != nil {
			log.Println(err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "✅ Книга успешно добавлена в библиотеку")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
	case "back":
		deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)

		_, err := h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
			return
		}
	case "mybook":
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		book, err := h.bookService.GetBookByID(bookID)
		if err != nil {
			log.Println(err)
			return
		}

		text := h.buildBookCard(book)
		keyboard := h.buildBookKeyboard(book.UserID, book)

		var sentMessage tgbotapi.Message
		if book.CoverURL == "" {
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = keyboard

			sentMessage, err = h.bot.Send(msg)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(book.CoverURL))
			photo.Caption = text
			photo.ReplyMarkup = keyboard

			sentMessage, err = h.bot.Send(photo)
			if err != nil {
				log.Println(err)
				return
			}
		}

		h.deleteState[telegramID] = DeleteState{
			SentMessageID: sentMessage.MessageID,
		}

	case "status":
		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		status := parts[1]
		bookID, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.bookService.UpdateStatus(user.ID, bookID, domain.BookStatus(status))
		if err != nil {
			log.Println(err)
			return
		}

		book, err := h.bookService.GetBookByID(bookID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateBookCard(chatID, messageID, book)
		if err != nil {
			log.Println(err)
			return
		}
	case "delete":
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		buttonYes := tgbotapi.NewInlineKeyboardButtonData("✅ Да", "confirmdelete:"+strconv.FormatInt(bookID, 10))
		buttonNo := tgbotapi.NewInlineKeyboardButtonData("❌ Нет", "canceldelete:"+strconv.FormatInt(bookID, 10))

		var rows [][]tgbotapi.InlineKeyboardButton
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonYes})
		rows = append(rows, []tgbotapi.InlineKeyboardButton{buttonNo})

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(chatID, "Вы действительно хотите удалить книгу?")
		msg.ReplyMarkup = keyboard

		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "confirmdelete":
		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.bookService.DeleteBook(user.ID, bookID)
		if err != nil {
			log.Println(err)
			return
		}

		deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)

		_, err = h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
			return
		}

		state, exist := h.deleteState[telegramID]
		if exist {
			deleteConfig := tgbotapi.NewDeleteMessage(chatID, state.SentMessageID)

			_, err = h.bot.Request(deleteConfig)
			if err != nil {
				log.Printf("Ошибка удаления сообщения: %v", err)
				return
			}
			delete(h.deleteState, telegramID)
		}

		stateLibrary, exist := h.libraryState[telegramID]
		if exist {
			books, err := h.bookService.GetUserBooks(user.ID)
			if err != nil {
				log.Println(err)
				return
			}

			pageSize := 4
			totalPages := int(math.Ceil(float64(len(books)) / float64(pageSize)))
			currentPage := stateLibrary.Page

			if totalPages > 0 && currentPage >= totalPages {
				currentPage = totalPages - 1
			}
			text, keyboard := h.buildBooksPage(books, currentPage)

			edit := tgbotapi.NewEditMessageText(chatID, stateLibrary.MessageID, text)
			edit.ReplyMarkup = &keyboard

			_, err = h.bot.Send(edit)
			if err != nil {
				log.Println(err)
				return
			}
		}

		msg := tgbotapi.NewMessage(chatID, "🗑 Книга успешно удалена")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
	case "canceldelete":
		deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageID)

		_, err := h.bot.Request(deleteConfig)
		if err != nil {
			log.Printf("Ошибка удаления сообщения: %v", err)
			return
		}
	case "progress":
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "Введите текущую страницу книги сообщением")
		sentMessage, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}

		h.progressState[telegramID] = ProgressState{
			BookID:        bookID,
			MessageID:     messageID,
			SentMessageID: sentMessage.MessageID,
		}

	case "startsession":
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.sessionService.StartSession(bookID, user.ID)
		if err != nil {
			log.Println(err)

			msg := tgbotapi.NewMessage(chatID, "❌ У вас уже есть активная сессия чтения")
			_, _ = h.bot.Send(msg)

			return
		}

		book, err := h.bookService.GetBookByID(bookID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateBookCard(chatID, messageID, book)
		if err != nil {
			log.Println(err)
			return
		}
	case "finishsession":
		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.sessionService.FinishSession(user.ID)
		if err != nil {
			log.Println(err)
			return
		}
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		book, err := h.bookService.GetBookByID(bookID)
		if err != nil {
			log.Println(err)
			return
		}

		err = h.updateBookCard(chatID, messageID, book)
		if err != nil {
			log.Println(err)
			return
		}
	case "setpages":
		bookID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "Введите количество страниц в книге")
		sentMessage, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}

		h.pagesState[telegramID] = PagesState{
			BookID:        bookID,
			MessageID:     messageID,
			SentMessageID: sentMessage.MessageID,
		}
	case "books":
		if len(parts) < 3 {
			return
		}
		page, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
			return
		}

		user, err := h.bookService.GetUserByTelegramID(telegramID)
		if err != nil {
			log.Println(err)
			return
		}

		books, err := h.bookService.GetUserBooks(user.ID)
		if err != nil {
			log.Println(err)
			return
		}

		text, keyboard := h.buildBooksPage(books, page)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		_, err = h.bot.Send(edit)
		if err != nil {
			log.Println(err)
			return
		}

		h.libraryState[telegramID] = LibraryState{
			MessageID: messageID,
			Page:      page,
		}

	case "sessions":
		if len(parts) < 3 {
			return
		}

		page, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Println(err)
			return
		}

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

		text, keyboard := h.buildSessionsPage(sessions, page)

		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		_, err = h.bot.Send(edit)
		if err != nil {
			log.Println(err)
			return
		}

	case "noop":
		return
	}
}
