package usecases

import (
	"context"
	"fmt"

	"reminder/internal/entities"
	"reminder/internal/interfaces"
)

type SendReminderService struct {
	telegram interfaces.TelegramSender
}

func NewSendReminderService(telegram interfaces.TelegramSender) *SendReminderService {
	return &SendReminderService{telegram: telegram}
}

func (s *SendReminderService) SendReminder(ctx context.Context, reminder *entities.Reminder) error {
	localTime := reminder.Time.Local().Format("02.01.2006 15:04")
	message := fmt.Sprintf("Напоминание: \n\n ⌚ %s\n 📌 %s\n ✉ %s", localTime, reminder.Topic, reminder.Text)

	err := s.telegram.SendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	return nil
}
