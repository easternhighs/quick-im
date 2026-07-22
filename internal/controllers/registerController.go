package controllers

import (
	"errors"
	"quick-im-demo/internal/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterController struct{}

type UserRegisterInfo struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

func (RegisterController) DoRegister(c *gin.Context) {
	c.HTML(200, "register.html", gin.H{})
}

func (RegisterController) WriteRegisterInfo(c *gin.Context) {
	var registerInfo UserRegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser user.User
	err := user.DB.Where("username = ?", registerInfo.Username).First(&existingUser).Error
	if err == nil {
		c.JSON(409, gin.H{"error": "username already exists"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(500, gin.H{"error": "query user failed"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "password hash failed"})
		return
	}

	userInfo := user.User{
		Username:  registerInfo.Username,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}
	if err := user.DB.Create(&userInfo).Error; err != nil {
		c.JSON(500, gin.H{"error": "create user failed"})
		return
	}

	c.JSON(201, gin.H{"id": userInfo.Id, "username": userInfo.Username, "message": "register successful"})
}
