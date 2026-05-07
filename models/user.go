package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NamaLengkap string `gorm:"type:varchar(100);not null"`
	NIK         string `gorm:"type:varchar(16);unique;not null"`
	Email       string `gorm:"type:varchar(100);unique;not null"`
	Password    string `gorm:"type:varchar(255);not null"`
	Role        string `gorm:"type:enum('warga', 'admin');default:'warga'"`
}