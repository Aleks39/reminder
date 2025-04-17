package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"reminder/internal/adapters"
	"reminder/internal/config"
	"reminder/internal/entities"
	"reminder/internal/usecases"
	"time"
)

var reminders map[string]*entities.Reminder

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	ctx := context.Background()
	logger := setupLogger()

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	telegramAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		logger.Error("Не удалось инициализировать TelegramAdapter", "error", err)
		os.Exit(1)
	}

	sendReminderService := usecases.NewSendReminderService(telegramAdapter)

	c := cron.New()
	defer c.Stop()

	addReminders()

	c.AddFunc("@every 20s", func() {
		logger.Info("1")
		now := time.Now().Truncate(time.Minute).UTC()
		logger.Info("Текущее время", "time", now)

		currentReminders := GetCurrentReminders()
		for _, reminder := range currentReminders {
			logger.Info("2")
			if err = sendReminderService.SendReminder(ctx, reminder); err != nil {
				logger.Error("Ошибка отправки напоминани", "error", err)
			} else {
				logger.Info("Напоминание успешно отправлена", "remind", reminder.Text, "author", reminder.Topic)
			}
		}
	})

	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	select {}
}

func addReminders() {
	reminders = make(map[string]*entities.Reminder)
	AddReminder(
		time.Date(2025, 04, 17, 16, 55, 0, 0, time.UTC),
		"Созвон с Никитой",
		"Информативное описание",
	)
	AddReminder(
		time.Date(2025, 04, 17, 16, 59, 0, 0, time.UTC),
		"Созвон с Виталей",
		"Информативное описание",
	)
}

func AddReminder(eventTime time.Time, topic, text string) {
	reminders[eventTime.Format("2006-01-02 15:04:05")] = &entities.Reminder{
		Time:  eventTime,
		Topic: topic,
		Text:  text,
	}
}

func GetCurrentReminders() []*entities.Reminder {
	key := time.Now().Truncate(time.Minute).Format("2006-01-02 15:04:05")
	currentReminder := reminders[key]
	var currentReminders []*entities.Reminder
	if currentReminder != nil {
		currentReminders = append(currentReminders, currentReminder)
	}

	return currentReminders
}
