package bot

import (
	"context"
	"fmt"
	"inivoice/internal/botkit"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		startMessage := fmt.Sprintf("Привет @%s\n\nНа связи iniAI 🤖\n\nПришли мне запись встречи или лекции, а в ответ ты получишь:\n- Транскрибацию всей записи\n- Кратко содержание записи\n\n*Мы работаем в тестовом режиме, на данный момент поддерживаются аудио длиной до 1 минуты\n\nПогнали 🚀", update.Message.From.UserName)
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, startMessage)); err != nil {
			return err
		}
		return nil
	}
}

func ViewCmdHelp(cmdViews map[string]botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		var commandsList string
		for cmd := range cmdViews {
			commandsList += "/" + cmd + "\n"
		}
		helpMessage := "Доступные команды:\n" + commandsList
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)); err != nil {
			return err
		}
		return nil
	}
}

func ViewSpeechToText() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте мне аудиофайл.")); err != nil {
			return err
		}
		return nil
	}
}
