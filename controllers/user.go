package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/helpers"
	"github.com/zombman/server/models"
)

func UserCreate(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:  body.Username,
		Password:  body.Password,
	}

	result := helpers.Postgres.Create(&user)
	
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}