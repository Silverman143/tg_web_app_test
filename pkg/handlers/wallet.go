package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func (h *Handler) GetAllTransactions(c *gin.Context){
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "id is not exist in context")
		return
	}

	transactions, err := h.service.Wallet.GetAllTransactions(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Couldn't get transactions"+err.Error())
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) GetPositiveTransactions(c *gin.Context){
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "id is not exist in context")
		return
	}

	transactions, err := h.service.Wallet.GetPositiveTransactions(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Couldn't get transactions"+err.Error())
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) GetNegativeTransactions(c *gin.Context){
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "id is not exist in context")
		return
	}

	transactions, err := h.service.Wallet.GetNegativeTransactions(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Couldn't get transactions"+err.Error())
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) GetBalance(c *gin.Context){
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "id is not exist in context")
		return
	}

	balance, err := h.service.Wallet.GetBalance(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Couldn't get wallet balance: "+ err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"wallet_balance": balance,
	})
}