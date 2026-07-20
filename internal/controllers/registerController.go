package controllers

import (
	"fmt"
	"quick-im-demo/internal/models/user"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	// DefaultController
}

type UserRegisterInfo struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (rc RegisterController) DoRegister(c *gin.Context) {
	c.HTML(200, "register.html", gin.H{})
}

func (rc RegisterController) WriteRegisterInfo(c *gin.Context) {
	var registerInfo UserRegisterInfo
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 处理注册逻辑，例如验证输入、保存用户信息等
	username := registerInfo.Username
	//对密码进行哈希加密后存储到数据库中
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	userInfo := user.User{
		Username:  username,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
	}

	user.DB.Create(&userInfo)

	var test user.User
	user.DB.Find(&test)
	fmt.Println(test)
}
