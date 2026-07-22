package controllers

import (
	"errors"
	"quick-im-demo/internal/config"
	"quick-im-demo/internal/models/user"
	"quick-im-demo/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginController struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (LoginController) ShowLogin(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

func (lc LoginController) DoLogin(c *gin.Context) {
	if err := c.ShouldBindJSON(&lc); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var storedUser user.User
	if err := user.DB.Where("username = ?", lc.Username).First(&storedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(401, gin.H{"error": "username or password is incorrect"})
			return
		}
		c.JSON(500, gin.H{"error": "query user failed"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(lc.Password)); err != nil {
		c.JSON(401, gin.H{"error": "username or password is incorrect"})
		return
	}

	token, err := utils.CreateToken(storedUser.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "create token failed"})
		return
	}

	now := time.Now()
	_ = user.DB.Model(&storedUser).Update("login_time", &now).Error
	c.SetCookie("quick_im_token", token, config.ReadConfig().JWT.ExpireHours*60*60, "/", "", false, true)
	c.JSON(200, gin.H{"message": "login successful", "token": token, "username": storedUser.Username})
}

func (LoginController) Logout(c *gin.Context) {
	c.SetCookie("quick_im_token", "", -1, "/", "", false, true)
	c.JSON(200, gin.H{"message": "logout successful"})
}
