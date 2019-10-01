package model

import "github.com/jinzhu/gorm"

type Card struct {
	gorm.Model
	//User      User
	CreatedAt string `gorm:"type:timestamp"`
	UpdatedAt string `gorm:"type:timestamp"`
	Signature string `gorm:"type:varchar(100)"`
	Reusable  int
	ExpYear   int
	ExpMonth  int
}

func (Card) TableName() string {
	return "user_cards"
}
