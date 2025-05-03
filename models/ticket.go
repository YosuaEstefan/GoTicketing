package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	EventID     uint      `json:"event_id" binding:"required"`
	Event       Event     `json:"event" gorm:"foreignKey:EventID"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Status      string    `json:"status" gorm:"size:20;default:active"` // 'active','used','cancelled'
	Price       float64   `json:"price"`
	PurchasedAt time.Time `json:"purchased_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
