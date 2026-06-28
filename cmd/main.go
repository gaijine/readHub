package main

import (
	"context"
	"log"

	"readHub/configs"
	openlibrary "readHub/internal/client/openlibrary"
	"readHub/internal/database"
	"readHub/internal/handler/telegram"
	"readHub/internal/postgres"
	"readHub/internal/service"
	"readHub/internal/worker"

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
	userRepo := postgres.NewUserRepository(db) // репозиторий пользователей
	sessionRepo := postgres.NewSessionRepository(db)
	openLib := openlibrary.NewClient()                                    // клиент для работы с внешним openLibrary API
	bookService := service.NewBookService(bookRepo, userRepo, openLib)    // сервис бизнес-логики книг
	sessionService := service.NewSessionService(sessionRepo, bookService) // сервис сессий чтения
	statsService := service.NewStatsService(bookRepo, sessionRepo)
	reminderRepo := postgres.NewReminderRepository(db)
	reminderService := service.NewReminderService(reminderRepo, userRepo)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	reminderWorker := worker.NewReminderWorker(reminderService, bot)

	handler := telegram.NewHandler(bookService, bot, sessionService, statsService, reminderService)

	ctx, _ := context.WithCancel(context.Background())
	go reminderWorker.Run(ctx)

	handler.Run()
	// client.SearchBooks("Alice")
	// resp, err := http.Get("https://openlibrary.org")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println(resp.Status)
}
