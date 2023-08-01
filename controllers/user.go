package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zombman/server/common"
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

	user, err := common.CreateUser(body.Username, body.Password)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserLogin(c *gin.Context) {
	var body struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := common.GetUserByUsernameAndPlainPassword(body.Username, body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}