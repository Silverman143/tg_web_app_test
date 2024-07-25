package handler

import (
	"fmt"
	"net/http"
	back "project-2x"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetDailyBonys(c *gin.Context){

	fmt.Println("get daily bonus")
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "id is not exist in context")
		return
	}

	dailyBonusInfo, err := h.service.UsersData.GetDailyBonusInfo(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Error with getting daily bonus info:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, dailyBonusInfo)
}

func (h *Handler) ClaimDailyBonys(c *gin.Context){

	fmt.Println("Claim bonus")
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "id is not exist in context")
		return
	}

	amount, err := h.service.UsersData.ClaimDailyBonus(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Error with getting daily bonus info:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"claimed": amount,
	})
}

func (h *Handler) GetReferralURL(c *gin.Context){

	fmt.Println("Get ref url ")
	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "id is not exist in context")
		return
	}

	refKey, err := h.service.UsersData.GetReferralCode(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Error with getting referral code: " + err.Error())
		return
	}

	url := fmt.Sprintf("https://t.me/%s?start=ref_%s", h.telegramBot.API.Self.UserName, refKey)
	c.JSON(http.StatusOK, map[string]interface{}{
		"ref_url": url,
	})
}

type getReferralsInput struct{
	Offset int `json:"offset"`
	PageSize int `json:"pageSize" binding:"required"`
}

func (h *Handler) GetUserReferrals(c *gin.Context){
	fmt.Println("Get user referrals ")

	var input getReferralsInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "id is not exist in context")
		return
	}

	referrals, err := h.service.UsersData.GetUserReferrals(telegramId, input.Offset, input.PageSize )

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Error with getting referrals: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, referrals)
}

func (h *Handler) GetUserData(c *gin.Context){
	fmt.Println("Get user data")

	var userData back.User


	telegramId, exists := GetTelegramID(c)
	if exists != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "id is not exist in context")
		return
	}

	userData, err := h.service.UsersData.GetUserProfil(telegramId)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Error with getting user data: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, userData)
}
