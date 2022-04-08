package model

import "gorm.io/gorm"

type Collections struct {
	gorm.Model
	Username     string `gorm:"not null;type:varchar(24);unique" json:"username"`
	Firstname    string `gorm:"not null;type:varchar(64)" json:"firstname"`
	Lastname     string `gorm:"not null;type:varchar(64)" json:"lastname"`
	Password     string `gorm:"not null" json:"password"`
	PhoneNumber  string `gorm:"not null;type:varchar(10)" json:"phone_number"`
	EmailAddress string `gorm:"not null;type:varchar(256);unique" json:"email_address"`
	FacebookLink string `gorm:"not null" json:"facebook_link"`
}

func (Collections) TableName() string { return "collections" }
