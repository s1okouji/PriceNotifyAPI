package model

import "time"

type Game struct {
	AppId           int `gorm:"primaryKey"`
	AppName         string
	BasePrice       int
	FinalPrice      int
	DiscountPercent int
	CreatedAt       time.Time
}
