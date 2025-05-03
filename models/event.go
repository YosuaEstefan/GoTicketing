package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name" binding:"required" gorm:"size:191;uniqueIndex"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Location    string    `json:"location" binding:"required" gorm:"size:191"`
	Capacity    int       `json:"capacity" binding:"required,min=1"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Status      string    `json:"status" gorm:"size:20;default:active"` // 'active','ongoing','completed'
	Tickets     []Ticket  `json:"-" gorm:"foreignKey:EventID"`
}

func (e *Event) BeforeUpdate(tx *gorm.DB) error {
	// Check if event is already completed
	var oldEvent Event
	if err := tx.First(&oldEvent, e.ID).Error; err != nil {
		return err
	}

	if oldEvent.Status == "completed" {
		return gorm.ErrInvalidTransaction
	}

	// Update status based on dates
	now := time.Now()
	if now.After(e.EndDate) {
		e.Status = "completed"
	} else if now.After(e.StartDate) {
		e.Status = "ongoing"
	}

	return nil
}
