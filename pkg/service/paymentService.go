package service

import (
	"encoding/json"
	"project-2x/pkg/database"
	"project-2x/pkg/telegramBot"
	tgbotapi "project-2x/pkg/telegramBot/all"

	"github.com/sirupsen/logrus"
)

type PaymentService struct {
	db database.Payment
	telegramBot telegramBot.Bot
}

func NewPaymentService(db database.Payment, telegramBot telegramBot.Bot) *PaymentService{
	return &PaymentService{
		db: db,
		telegramBot: telegramBot,
	}
}

func (s *PaymentService) CreateStarsInvoice(amount int) (string, error) {
	var url string
	prices := []tgbotapi.LabeledPrice{
		{Label: "Product", Amount: amount}, // Указываем цену продукта
	}
	config := tgbotapi.InvoiceLinkConfig{
		Title:       "Your Product Title",     // Название продукта
		Description: "Your Product Description", // Описание продукта
		Payload:     "YourPayload",             // Ваш payload
		ProviderToken: "",                     // Токен платежного провайдера, пустая строка для Telegram Stars
		Currency:    "XTR",                    // Валюта, "XTR" для Telegram Stars
		Prices: &prices,
	}
	logrus.Info("!!!!!! Send request to telegram API")
	response, err := s.telegramBot.API.CreateInvoiceLink(config)
	if err != nil {
		logrus.Info("!!!!!! eeerrrrrooooorrrr")
		return "", err
	}
	
	err = json.Unmarshal(response.Result, &url)
	if err != nil {
		logrus.Info("!!!!!! Unmarshal error")
		return "", err
	}
	return url, nil
}