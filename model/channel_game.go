package model

type Channel_Game struct {
	Id        int `gorm:"primaryKey"`
	AppId     int
	ChannelId string
}
