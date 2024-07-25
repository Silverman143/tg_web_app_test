package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)
const (
	authHeader = "Authorization"
	userCtx = "userId"
)
func (h *Handler) userIdentity(c *gin.Context){

	fmt.Println("identity ")
	header := c.GetHeader(authHeader)
	if header == ""{
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerPars := strings.Split(header, " ")
	if len(headerPars) != 1 {
		fmt.Println("invalid auth header format:", header)
		NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	id, err := h.service.Authorization.ParseToken(headerPars[0])

	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}