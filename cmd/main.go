package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())

	go reminderWorker.Run(ctx)
	go handler.Run()

	quit := make(chan os.Signal, 1)

	// если пользователь нажмет ctrl+c или процесс получит sigterm, положи этот сигнал в канал quit
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// блокирует мэин, она будет ждать сигнала
	<-quit

	log.Println("Received shutdown signal")
	cancel()
	log.Println("Application stopped")
}
