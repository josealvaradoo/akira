package model

import "gorm.io/gorm"

type Member struct {
	*gorm.Model
	ID   string
	Name string
}

type Winner struct {
	*gorm.Model
	ID   string
	Name string
}
