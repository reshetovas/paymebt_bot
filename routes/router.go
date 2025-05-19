package routes

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
	"os"
	"payment_bot/handlers"
	"strings"
)

type Router struct {
	handlers *handlers.PaymentHandler
}

func NewRouter(handlers *handlers.PaymentHandler) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) InitRouter() *telebot.Bot {
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
	log.Debug().Msgf("bot: %T, %+v\n", bot, bot)
	if err != nil {
		log.Fatal().Err(err)
	}

	// Регистрация обработчиков
	bot.Handle("/start", r.handlers.HandleStart)
	bot.Handle("/add_payment", r.handlers.HandleAddPayment)
	bot.Handle("/report", r.handlers.HandleReport)
	bot.Handle("/export", r.handlers.HandleExport)
	bot.Handle(telebot.OnText, r.handlers.HandleTextInput)
	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		data := c.Callback().Data
		if strings.HasPrefix(data, "category:") {
			return r.handlers.HandleCategorySelection(c)
		}
		if strings.HasPrefix(data, "report") {
			return r.handlers.HandleReportSelection(c)
		}
		return nil
	})

	return bot
}
