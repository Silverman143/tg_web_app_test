package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	tgbotapi "project-2x/pkg/telegramBot/all"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

func (h *Handler) TelegramBotHandler(c *gin.Context) {
	// Чтение тела запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		NewErrorResponse(c, http.StatusBadRequest, "Error reading request body: " + err.Error())
		return
	}

	// logrus.WithField("body", string(body)).Info("Request Body")

	// Декодирование запроса в структуру tgbotapi.Update
	var update tgbotapi.Update

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&update); err != nil {
		log.Printf("Error decoding request body: %v", err)
		NewErrorResponse(c, http.StatusBadRequest, "Error decoding request body: " + err.Error())
		return
	}

	// Обработка pre-checkout запроса
	if update.PreCheckoutQuery != nil {
		if err := h.telegramBot.Handler.HandlePreCheckoutQuery(update.PreCheckoutQuery); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Error handling pre-checkout query: " + err.Error())
			return
		}
	} else if update.Message != nil && update.Message.SuccessfulPayment != nil {
		// Обработка успешного платежа
		if err := h.telegramBot.Handler.HandleSuccessfulPayment(update.Message, update.Message.SuccessfulPayment); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Error handling successful payment: " + err.Error())
			return
		}

		currency := update.Message.SuccessfulPayment.Currency

		if currency != ""{
			currency = "stars"
		}

		err = h.service.UsersData.AddPayment(update.Message.From.ID, update.Message.SuccessfulPayment.TotalAmount, update.Message.SuccessfulPayment.Currency)

		if err != nil{
			NewErrorResponse(c, http.StatusInternalServerError, "Error handling successful payment: " + err.Error())
			return
		}

	} else if update.Message != nil && update.Message.IsCommand() && strings.HasPrefix(update.Message.Text, "/start")  {

		// Извлечение реферального ключа из сообщения
		var refKey string
		parts := strings.Fields(update.Message.Text)
		if len(parts) > 1 && strings.HasPrefix(parts[1], "ref_") {
			refKey = strings.TrimPrefix(parts[1], "ref_")
		}

		// Обработка команды /start
		if err := h.telegramBot.Handler.HandleStartBot(update.Message, refKey); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Error handling start bot: "+err.Error())
			return
		}

		// Создание пользователя
		if err := h.service.UsersData.CreateUser(*update.Message, refKey); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Error creating user: "+err.Error())
			return
		}

	} else if update.Message != nil {
		// Обработка обычного сообщения
		if err := h.telegramBot.Handler.HandleMessage(update.Message); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Error handling message: " + err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}