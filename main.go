package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"payment_bot/handlers"
	"payment_bot/logger"
	"payment_bot/services"
	"payment_bot/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
)

func main() {
	//logging setup
	logger.InitLogger(os.Getenv("LOGLEVEL"))

	//coonecting to db
	db, err := sql.Open("sqlite3", "./payment_bot.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB sql.Open")
	}
	defer db.Close()

	if err := goose.Up(db, "./storage/migrations"); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}

	state := storage.NewStateStorage()
	storage := storage.NewSQLDBPaymentStorage(db)
	services := services.NewPaymentService(state, storage)
	handlers := handlers.NewPaymentHandler(services)

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Error().Msg("token not identified")
	}
	webhookEndpoin := os.Getenv("WEBHOOK_URL")
	if webhookEndpoin == "" {
		log.Error().Msg("token not identified")
	}
	pref := telebot.Settings{
		Token: token,
		Poller: &telebot.Webhook{
			Listen:         ":8080",
			Endpoint:       &telebot.WebhookEndpoint{PublicURL: webhookEndpoin},
			MaxConnections: 100,
		}}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal().Err(err)
	}

	input := "2024-01-01 - 2025-05-15"
	parts := strings.Split(input, "-")
	fmt.Printf("type: %T, value: %v\n", input, input)
	fmt.Printf("type: %T, value: %v\n", parts, parts)
	fmt.Println(len(parts))
	for a, b := range parts {
		fmt.Printf("a: %v, b: %v\n", a, b)
	}

	startStr := strings.TrimSpace(parts[0] + "-" + parts[1] + "-" + parts[2])
	fmt.Printf("type: %T, value: %v\n", startStr, startStr)
	endStr := strings.TrimSpace(parts[3] + "-" + parts[4] + "-" + parts[5])
	fmt.Printf("type: %T, value: %v\n", endStr, endStr)

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Printf("type: %T, value: %v\n", start, start)
	fmt.Printf("type: %T, value: %v\n", end, end)

	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	fmt.Printf("firstOfMonth - type: %T, value: %v\n", firstOfMonth, firstOfMonth)

	// Регистрация обработчиков
	bot.Handle("/start", handlers.HandleStart)
	bot.Handle("/add_payment", handlers.HandleAddPayment)
	bot.Handle("/report", handlers.HandleReport)
	bot.Handle("/export", handlers.HandleExport)
	bot.Handle(telebot.OnText, handlers.HandleTextInput)
	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		data := c.Callback().Data
		if strings.HasPrefix(data, "category:") {
			return handlers.HandleCategorySelection(c)
		}
		if strings.HasPrefix(data, "report") {
			return handlers.HandleReportSelection(c)
		}
		return nil
	})

	log.Info().Msg("The telegram bot is running")
	bot.Start()
}
