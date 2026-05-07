package models

import "gorm.io/gorm"

type PengajuanSurat struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	User       User   `gorm:"foreignKey:UserID"`
	JenisSurat string `gorm:"type:varchar(50);not null"`
	Keperluan  string `gorm:"type:text;not null"`
	FileKTPUrl string `gorm:"type:varchar(255);not null"`
	FileKKUrl  string `gorm:"type:varchar(255);not null"`
	Status     string `gorm:"type:enum('Pending', 'Diproses', 'Selesai', 'Ditolak');default:'Pending'"`
	Keterangan string `gorm:"type:text"`
}