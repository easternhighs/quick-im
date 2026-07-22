package utils

import (
	"encoding/base64"
	"quick-im-demo/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 将yaml中的密钥进行base64编码，返回编码后的字符串
func GenerateSecretKey() string {
	b := []byte(config.ReadConfig(config.ConfigFileName).JWT.Secret)
	return base64.RawURLEncoding.EncodeToString(b)
}

// 签发Token
func CreateToken(username string) (string, error) {
	//创建声明，为map类型，存储一些用户信息和token的过期时间等
	claims := jwt.MapClaims{
		"username": username,
		"name":     "im-server-signed-token",
		"iat":      jwt.NewNumericDate(jwt.TimeFunc()),
		"exp":      jwt.NewNumericDate(jwt.TimeFunc().Add(time.Hour * time.Duration(config.ReadConfig(config.ConfigFileName).JWT.ExpireHours))),
	}

	//使用HS256算法和密钥生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GenerateSecretKey()))
}

// 验证Token
func VerifyToken(tokenString string) (*jwt.Token, error) {
	//验证token时，直接调用jwt.Parse()方法，传入token字符串和一个回调函数，回调函数中返回密钥
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(GenerateSecretKey()), nil
	})
}
