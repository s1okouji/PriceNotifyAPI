package model

import (
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	id    string
	_type int
}
