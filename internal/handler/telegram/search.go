package telegram

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleSearch(chatID, telegramID int64, query string) bool {
	if query == "" {
		msg := tgbotapi.NewMessage(chatID, "Использование:\n/search Название книги")
		_, err := h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return false
	}

	books, err := h.bookService.SearchBooks(query)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(chatID, "Не удалось выполнить поиск. Попробуйте позже.")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return false
	}
	if len(books) == 0 {
		msg := tgbotapi.NewMessage(chatID, "По вашему запросу ничего не найдено,\nВведите другое название книги или автора ")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return false
	}

	if len(books) > 5 { // если книг будет больше 5
		books = books[:5] // новый массив не создается, создается слайс который смотрит на первые 5 элементов старого массива(слайса)
	} // если книг до 5 это не сработает

	h.searchCache[telegramID] = books // добавили результат (список книг) в кеш память
	log.Println("CACHE SAVE")
	log.Println(telegramID)
	log.Println(len(books))

	var builder strings.Builder                 // в буфере будет сохранять данные строки
	var buttons []tgbotapi.InlineKeyboardButton // хранит кнопки
	var rows [][]tgbotapi.InlineKeyboardButton  // хранит строки, мол первый элемент это первая строка

	for i, book := range books {
		var author string
		builder.WriteString("[")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString("] ")
		builder.WriteString(book.Title)
		builder.WriteString("\n")
		builder.WriteString("Автор: ")
		if len(book.Author) == 0 {
			author = "Unknown"
		} else {
			author = strings.Join(book.Author, ", ") // преобразовывает слайс в строку и ставит разделитель меж элементами
		}
		builder.WriteString(author)
		builder.WriteString("\n\n")
		// создали кнопку
		button := tgbotapi.NewInlineKeyboardButtonData("["+strconv.Itoa(i+1)+"]", "details:"+book.OpenLibraryID)
		// добавили кнопку в слайс кнопок
		buttons = append(buttons, button)
	}
	rows = append(rows, buttons)                          // добавили первый слайс в другой(получается вложенные слайсы)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...) // создали клавиатуру (внутри вложенные слайсы [][] первый отвечает за строки второй за столбцы)

	msg := tgbotapi.NewMessage(chatID, builder.String())
	msg.ReplyMarkup = keyboard // прикрепили клавиатуру к сообщению

	_, err = h.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return true
}
