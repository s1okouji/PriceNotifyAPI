package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	id string
}
