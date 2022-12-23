package entity

import "gorm.io/gorm"

type History struct {
	gorm.Model
	app_id           int
	day              int
	final_price      int
	discount_percent int
}
