package model

type User struct {
	//gorm.Model
	ID        int    `gorm:"primary_key"`
	FirstName string `gorm:"type:varchar(100)"`
	Card      []Card
	Loan      []Loan `gorm:"foreignkey:user_id"`
}
