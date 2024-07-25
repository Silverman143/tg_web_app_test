package telegramBot

import (
	"fmt"
	"log"
	tgbotapi "project-2x/pkg/telegramBot/all"

	"github.com/sirupsen/logrus"
)

type BotHandler struct {
	bot *tgbotapi.BotAPI
}

func NewBotHandler(bot *tgbotapi.BotAPI) *BotHandler {
	return &BotHandler{bot: bot}
}

func (h *BotHandler) HandleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Hello, I'm a bot!")
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	return nil
}

func (h *BotHandler) HandlePreCheckoutQuery(preCheckoutQuery *tgbotapi.PreCheckoutQuery) error {
	log.Printf("Received PreCheckoutQuery: %+v", preCheckoutQuery)

	config := tgbotapi.PreCheckoutAnswerConfig{
		PreCheckoutQueryID: preCheckoutQuery.ID,
		OK:                 true,
		ErrorMessage:       "",
	}
	_, err := h.bot.AnswerPreCheckout(config)
	if err != nil {
		logrus.Error("Error with check payment: " + err.Error())
		return err
	}

	return nil
}

func (h *BotHandler) HandleSuccessfulPayment(message *tgbotapi.Message, successfulPayment *tgbotapi.SuccessfulPayment) error {
	log.Printf("Received SuccessfulPayment: %+v", successfulPayment)

	msg := tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("Congratulations, you successfully paid %d $$$$", successfulPayment.TotalAmount))
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

func (h *BotHandler) HandleStartBot(message *tgbotapi.Message, ref string) error {
	// log.Printf("Received SuccessfulPayment: %+v", successfulPayment)

	// msg := tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("Congratulations, you successfully paid %d $$$$", successfulPayment.TotalAmount))
	// if _, err := h.bot.Send(msg); err != nil {
	// 	log.Printf("Error sending message: %v", err)
	// 	return err
	// }
	if ref != ""{
		msg := tgbotapi.NewMessage(message.From.ID, fmt.Sprintf("Congratulations, you were invite by %d ", ref))
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
	}

	return nil
}