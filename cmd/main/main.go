package main

import (
	"context"
	"errors"
	"inivoice/config"
	"inivoice/internal/bot"
	"inivoice/internal/botkit"
	"inivoice/internal/db"
	"inivoice/internal/openai"
	"inivoice/libs"
	"os"
	"os/signal"
	"syscall"

	"math/rand"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	logger := libs.NewApiLogger(cfg)
	logger.InitLogger()

	botapi, err := botapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Panicf("failed to create botAPI: %v", err)
	}
	logger.Info("Connecting to database...")

	dbClient, err := db.New(cfg.DB)
	if err != nil {
		logger.Panicf("failed to connect to db: %v", err)
	}
	logger.Info("Database connected successfully")
	
	defer dbClient.Close()
	logger.Info("Bot started successfully")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	openAI := openai.New(cfg.OpenAIKey, logger)
	viewFuncs := map[string]botkit.ViewFunc{
		"start":  bot.ViewCmdStart(),
		"speech": bot.ViewSpeechToText(),
	}
	iniBot := botkit.New(botapi, logger, openAI,
		map[string]botkit.ViewFunc{
			"start":  bot.ViewCmdStart(),
			"speech": bot.ViewSpeechToText(),
			"help":   bot.ViewCmdHelp(viewFuncs),
		})

	if err = iniBot.Start(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Panicf("failed to start bot: %v", err)
		}
	}
	logger.Info("Bot successfully stopped")
}
