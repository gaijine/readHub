package main

import (
	"log"

	"readHub/configs"
	openlibrary "readHub/internal/client/openlibrary"
	"readHub/internal/database"
	"readHub/internal/handler/telegram"
	"readHub/internal/postgres"
	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := configs.Load() // загрузили конфигурацию
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.New(cfg) // подключились к БД
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()
	log.Println("database connected")

	bookRepo := postgres.NewBookRepository(db) // репозиторий книг, для работы с БД PostgreSQL
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)                      // репозиторий пользователей
	openLib := openlibrary.NewClient()                                    // клиент для работы с внешним openLibrary API
	bookService := service.NewBookService(bookRepo, userRepo, openLib)    // сервис бизнес-логики книг
	sessionService := service.NewSessionService(sessionRepo, bookService) // сервис сессий чтения

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	handler := telegram.NewHandler(bookService, bot, sessionService)

	handler.Run()
	// client.SearchBooks("Alice")
	// resp, err := http.Get("https://openlibrary.org")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println(resp.Status)
}
