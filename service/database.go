package service

import (
	"github.com/s1okouji/price_notify_api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUp() {
	var err error
	dsn := "root:price_notify@tcp(127.0.0.1:3306)/price_notify?parseTime=true"
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
