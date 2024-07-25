package telegramBot

import (
	"fmt"
	"log"
	tgbotapi "project-2x/pkg/telegramBot/all"
)

type Handlers interface{
	HandleMessage(message *tgbotapi.Message) error
	HandlePreCheckoutQuery(preCheckoutQuery *tgbotapi.PreCheckoutQuery) error
	HandleSuccessfulPayment(message *tgbotapi.Message, successfulPayment *tgbotapi.SuccessfulPayment) error
	HandleStartBot(message *tgbotapi.Message, ref string) error 
}

type Bot struct {
	API *tgbotapi.BotAPI
	Handler Handlers
}

func NewBot(token string) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	return &Bot{
		API: botAPI,
		Handler: NewBotHandler(botAPI),
		}, nil
}

func (bot *Bot) Start() error {
	url := fmt.Sprintf("https://proj-2x-78a0ca7fa5b0.herokuapp.com/web_hook/%s", bot.API.Self.UserName)
	wh := tgbotapi.NewWebhook(url)

	_, err := bot.API.SetWebhook(wh)
	if err != nil {
		return err
	}

	info, err := bot.API.GetWebhookInfo()
	if err != nil {
		return err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return nil
}

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	From           *User  `json:"from"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}

// User struct
type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}
