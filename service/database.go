package service

import (
	"os"

	"github.com/s1okouji/pnabot_client/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUp() {
	var err error
	pass := os.Getenv("dbpass")
	user := os.Getenv("dbuser")
	host := os.Getenv("dbhost")
	dsn := user + ":" + pass + "@tcp(" + host + ")/price_notify?parseTime=true"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&model.Game{})
	DB.AutoMigrate(&model.History{})
	DB.AutoMigrate(&model.Channel{})
	DB.AutoMigrate(&model.Admin{})
	DB.AutoMigrate(&model.Channel_Game{})
}
