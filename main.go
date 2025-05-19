package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"os"
	"payment_bot/handlers"
	"payment_bot/logger"
	"payment_bot/routes"
	"payment_bot/services"
	"payment_bot/state_machine"
	"payment_bot/storage"
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
	paymentServices := services.NewPaymentService(state, storage)
	reportService := services.NewReportService(state, storage)
	exportService := services.NewExportService(state, storage)
	stateMachine := state_machine.NewStateMachine(paymentServices, reportService, exportService)
	handlers := handlers.NewPaymentHandler(paymentServices, reportService, exportService, stateMachine)
	router := routes.NewRouter(handlers)

	bot := router.InitRouter()

	log.Info().Msg("The telegram bot is running")
	bot.Start()
}
