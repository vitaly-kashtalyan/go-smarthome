package models

import "time"

type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"type:timestamp" json:"-"`
	UpdatedAt time.Time  `gorm:"type:timestamp" json:"-"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"-" sql:"index"`
}
