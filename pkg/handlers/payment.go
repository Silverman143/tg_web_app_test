package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type paymentInStarsInput struct {
	Amount int `json:"payment_amount" binding:"required"`
}

func (h *Handler) GetStarsPaymentURL(c *gin.Context){
	var input paymentInStarsInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Payment.CreateStarsInvoice(input.Amount)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"url": response,
	})
}