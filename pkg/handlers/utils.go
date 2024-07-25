package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTelegramID(c *gin.Context) (int64, error) {
    telegramId, exists := c.Get(userCtx)
    if !exists {
        return 0, fmt.Errorf("id does not exist in context")
    }

    telegramIdInt, ok := telegramId.(int64)
    if !ok {
        return 0, fmt.Errorf("id in context is not an int")
    }

    return telegramIdInt, nil
}