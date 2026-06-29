package worker

import (
	"context"
	"log"
	"time"

	"readHub/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ReminderWorker struct {
	reminderService service.ReminderService
	bot             *tgbotapi.BotAPI
}

func NewReminderWorker(remiderService service.ReminderService, bot *tgbotapi.BotAPI) *ReminderWorker {
	return &ReminderWorker{
		reminderService: remiderService,
		bot:             bot,
	}
}

func (w *ReminderWorker) Run(ctx context.Context) {
	// если только запустил то пусть сначала отправит а после будет ждать по минуте
	w.process()
	// объект который имеет канал и принимает время, и каждую мин в нашем примере рантайм отправляет сигнал (кладет время) в канал
	// а в цикле читается этот канал получает сигнал и продолжает выполняться уикл а именно процесс() запускается
	// и так циклом
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		// а это чтение из канала тикера, как пройдет время рантайм отправит сигнал, отправит время и канал прочтет это и запуститься процесс()
		case <-ticker.C:
			w.process()
		// тоже своего рода канал который будет читать если закрыли канал, будет получать сигнал
		// тип канала пустая структурка, и после сигнала разблокируется горутина и ретурном завершится
		case <-ctx.Done():
			log.Println("ReminderWorker stopped")
			return
		}
	}
}

func (w *ReminderWorker) process() {
	reminders, err := w.reminderService.GetDueReminders()
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range reminders {
		// TelegramID == ChatID для личного чата с ботом
		msg := tgbotapi.NewMessage(v.ChatID, "Самое время сделать паузу на кофе и пару страниц. Какую книгу мы сегодня читаем?")
		_, err = w.bot.Send(msg)
		if err != nil {
			log.Println(err)
			continue
		}

		err = w.reminderService.UpdateLastSent(v.UserID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
