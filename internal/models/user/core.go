package user

import (
	"log"
	"quick-im-demo/internal/config"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
}

func initDB() (db *gorm.DB) {
	config := config.ReadConfig(config.ConfigFileName)
	dsn := config.MySQL.Username + ":" + config.MySQL.Password + "@tcp(" + config.MySQL.Host + ":" + strconv.Itoa(config.MySQL.Port) + ")/" + config.MySQL.Database + "?charset=" + config.MySQL.Charset + "&parseTime=True&loc=Local"

	//根据配置信息，连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect mysql failed:", err)
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal("auto migrate failed:", err)
		panic("failed to auto migrate database")
	}

	return
}

var DB = initDB()
