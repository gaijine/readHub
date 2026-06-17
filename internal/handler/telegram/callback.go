package telegram

import (
	"log"
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
	// data := parts[1]

	// log.Println(action)
	// log.Println(parts[1])

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
		keyboard := h.buildBookKeyboard(bookID)
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(book.CoverURL))
		photo.Caption = text
		photo.ReplyMarkup = keyboard
		// msg := tgbotapi.NewMessage(chatID, text)
		// msg.ReplyMarkup = keyboard

		_, err = h.bot.Send(photo)
		if err != nil {
			log.Println(err)
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

		text := h.buildBookCard(book)
		keyboard := h.buildBookKeyboard(bookID)

		edit := tgbotapi.NewEditMessageCaption(chatID, messageID, text)
		edit.ReplyMarkup = &keyboard

		_, err = h.bot.Send(edit)
		if err != nil {
			log.Println(err)
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

		msg := tgbotapi.NewMessage(chatID, "🗑 Книга успешно удалена")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
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
		h.progressState[telegramID] = bookID

		msg := tgbotapi.NewMessage(chatID, "Введите текущую страницу книги сообщением")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "back":

	}
}
