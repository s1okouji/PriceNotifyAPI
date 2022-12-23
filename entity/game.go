package entity

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	app_id           int
	app_name         string
	base_price       int
	final_price      int
	discount_percent int
	history          History
}
