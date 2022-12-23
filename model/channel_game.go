package model

import "gorm.io/gorm"

type Channel_Game struct {
	gorm.Model
	AppId     int
	ChannelId string
}
