package entity

import (
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	id    int
	_type int
	users []string
}
