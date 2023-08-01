package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OAuthMain(c *gin.Context) {
	var body struct {
		GrantType    string `json:"grant_type" binding:"required"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		ExchangeCode string `json:"exchange_code"`
		RefreshToken string `json:"refresh_token"`
	}

	authHeader := c.GetHeader("Authorization")

	fmt.Println(body, authHeader)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}