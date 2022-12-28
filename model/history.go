package model

import "gorm.io/gorm"

type History struct {
	gorm.Model
	AppId           int
	Day             string
	FinalPrice      int
	DiscountPercent int
}
