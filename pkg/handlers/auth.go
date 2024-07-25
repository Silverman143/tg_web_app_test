package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// back "project-2x"
	"github.com/spf13/viper"
	initdata "github.com/telegram-mini-apps/init-data-golang"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) signUp (c *gin.Context){
// 	var input back.User

// 	if err := c.BindJSON(&input); err != nil {
// 		NewErrorResponse(c, http.StatusBadRequest, err.Error())
// 	}
// 	id, err := h.service.CreateUser(input);
// 	if err != nil{
// 		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
// 	}
// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"id": id,
// 	})
// }

type signInInput struct {
	InitData string `json:"initData" binding:"required"`
}

func (h *Handler) signIn (c *gin.Context){
	fmt.Println("sign-in")
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		NewErrorResponse(c, http.StatusFailedDependency, "telegram bot token not found ")
		return
	}

	expInStr := viper.GetString("exp_in")
    expIn, err := time.ParseDuration(expInStr)
    if err != nil {
        expIn = 24 * time.Hour
    }

	err = initdata.Validate(input.InitData, token, expIn)

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Validation error: "+err.Error())
		return
	}

	userInitData, err := initdata.Parse(input.InitData)

	if err != nil{
		NewErrorResponse(c, http.StatusBadRequest, "Couldn't pars initData from telegram: "+err.Error())
		return
	}

	// GetOrCreateUser
	userInfo, err := h.service.Authorization.CreateUser(userInitData);

	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, "Couldn't create user with error: "+err.Error())
		return
	}

	accessToken, err := h.service.GenerateAccessToken(int64(userInitData.User.ID));
	if err != nil{
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
		"userData": userInfo,
	})
}