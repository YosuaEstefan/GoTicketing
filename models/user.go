package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string    `json:"name" binding:"required" gorm:"size:191"`
	Email     string    `json:"email" binding:"required,email" gorm:"size:191;uniqueIndex"`
	Password  string    `json:"password" binding:"required,min=6" gorm:"size:191"`
	Role      string    `json:"role" gorm:"size:20;default:user"` // 'admin' or 'user'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tickets   []Ticket  `json:"-" gorm:"foreignKey:UserID"`
}
