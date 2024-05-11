package botkit

import (
	"context"
	"github.com/google/uuid"
	"inivoice/internal/botkit/markup"
	"inivoice/internal/constants"
	"inivoice/internal/openai"
	"inivoice/libs"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
	openAI   *openai.OpenAI
	logger   libs.Logger
}

func New(api *tgbotapi.BotAPI, logger libs.Logger, openAI *openai.OpenAI, cmdViews map[string]ViewFunc) *Bot {
	return &Bot{
		api:      api,
		logger:   logger,
		openAI:   openAI,
		cmdViews: cmdViews,
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update *tgbotapi.Update) {
	if update.Message == nil {
		b.logger.Errorf("error: %v", "not message")
		return
	}
	switch {

	case update.Message.IsCommand():
		if viewFunc, ok := b.cmdViews[update.Message.Command()]; ok {
			if err := viewFunc(ctx, b.api, update); err != nil {
				b.logger.Errorf("viewFunc error: %v", err.Error())
				if _, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())); err != nil {
					b.logger.Errorf("send error: %v", err.Error())
				}
			}
		} else {
			b.logger.Errorf("error: %v", "wrong command")
			if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, constants.NOT_COMMAND_ERROR)); err != nil {
				b.logger.Errorf("send error: %v", err.Error())
			}
			return
		}
	case update.Message.Voice != nil || update.Message.Audio != nil:
		var fileID string
		var duration int
		var msg tgbotapi.MessageConfig
		if update.Message.Voice != nil {
			fileID = update.Message.Voice.FileID
			duration = update.Message.Voice.Duration
		} else if update.Message.Audio != nil {
			fileID = update.Message.Audio.FileID
			duration = update.Message.Audio.Duration
		}

		if duration > 60 || duration == 0 {
			if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, constants.ERROR_FORMAT)); err != nil {
				b.logger.Errorf("send error: %v", err.Error())
			}
			return
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Идет обработка аудио...")
		if _, err := b.api.Send(msg); err != nil {
			b.logger.Errorf("send error: %v", err.Error())
		}

		file, err := b.api.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			b.logger.Errorf("error getting file: %v", err.Error())
			return
		}

		fileURL := file.Link(b.api.Token)
		filePath, err := libs.UploadFile(fileURL, uuid.New().String(), file.FilePath)
		if err != nil {
			log.Printf("error downloading file: %v", err.Error())
			return
		}

		result, err := b.openAI.SpeechToText(filePath)
		if err != nil {
			b.logger.Errorf("error converting speech to text: %v", err.Error())
			return
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "👉 Транскрибация вашей записи: \n"+markup.EscapeForMarkdown(result))
		if _, err = b.api.Send(msg); err != nil {
			b.logger.Errorf("send error: %v", err.Error())
		}

		summary, err := b.openAI.Summarize(result)
		if err != nil {
			b.logger.Errorf("error summarizing text: %v", err.Error())
			return
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "👉 Краткое содержание вашей записи: \n"+markup.EscapeForMarkdown(summary))
		if _, err = b.api.Send(msg); err != nil {
			b.logger.Errorf("send error: %v", err.Error())
		}

		if err = libs.DeleteFile(filePath); err != nil {
			b.logger.Errorf("error deleting file: %v", err.Error())
		}
	default:
		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, constants.ERROR_FORMAT)); err != nil {
			b.logger.Errorf("send error: %v", err.Error())
		}
		return
	}
}

func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 20*time.Second)
			b.handleUpdate(updateCtx, &update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error
