package main

import (
	"database/sql"
	"os"
	"time"

	"payment_bot/handlers"
	"payment_bot/logger"
	"payment_bot/services"
	"payment_bot/storage"

	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
)

func main() {
	//logging setup
	logger.InitLogger(os.Getenv("LogLevel"))

	//coonecting to db
	db, err := sql.Open("sqlite3", "./payment_bot.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB sql.Open")
	}
	defer db.Close()

	state := storage.NewStateStorage()
	storage := storage.NewSQLDBPaymentStorage(db)
	services := services.NewPaymentService(state, storage)
	handlers := handlers.NewPaymentHandler(services)

	pref := telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Регистрация обработчиков
	bot.Handle("/start", handlers.HandleStart)
	//bot.Handle("/add_payment", handlers.HandleAddPayment)
	bot.Handle("/report", handlers.HandleReport)
	bot.Handle("/export", handlers.HandleExport)
	bot.Handle(telebot.OnText, handlers.HandleTextInput)

	log.Info().Msg("The telegram bot is running")
	bot.Start()
}
