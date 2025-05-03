// repository/event_repository.go
package repository

import (
	"ticket/models"

	"gorm.io/gorm"
)

type EventRepository interface {
	Save(event *models.Event) error
	Update(event *models.Event) error
	Delete(id uint) error
	FindByID(id uint) (*models.Event, error)
	FindAll(page, limit int, search string) ([]models.Event, int64, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

func (r *eventRepository) Save(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	// Check if event has tickets
	var count int64
	if err := r.db.Model(&models.Ticket{}).Where("event_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrInvalidTransaction
	}

	return r.db.Delete(&models.Event{}, id).Error
}

func (r *eventRepository) FindByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindAll(page, limit int, search string) ([]models.Event, int64, error) {
	var events []models.Event
	var total int64

	query := r.db.Model(&models.Event{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR location LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
