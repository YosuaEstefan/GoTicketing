package service

import (
	"errors"
	"ticket/models"
	"ticket/repository"
	"time"
)

type EventService interface {
	CreateEvent(event *models.Event) error
	UpdateEvent(event *models.Event) error
	DeleteEvent(id uint) error
	GetEventByID(id uint) (*models.Event, error)
	GetAllEvents(page, limit int, search string) ([]models.Event, int64, error)
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) CreateEvent(event *models.Event) error {
	// Validate dates
	if event.StartDate.Before(time.Now()) {
		return errors.New("start date must be in the future")
	}

	if event.EndDate.Before(event.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Validate capacity
	if event.Capacity <= 0 {
		return errors.New("capacity must be greater than zero")
	}

	// Validate price
	if event.Price < 0 {
		return errors.New("price cannot be negative")
	}

	// Set status based on dates
	now := time.Now()
	if now.After(event.EndDate) {
		event.Status = "completed"
	} else if now.After(event.StartDate) {
		event.Status = "ongoing"
	} else {
		event.Status = "active"
	}

	return s.eventRepo.Save(event)
}

func (s *eventService) UpdateEvent(event *models.Event) error {
	// Get existing event
	existingEvent, err := s.eventRepo.FindByID(event.ID)
	if err != nil {
		return err
	}

	// Check if event is completed
	if existingEvent.Status == "completed" {
		return errors.New("cannot update a completed event")
	}

	// Validate dates
	if event.StartDate.Before(time.Now()) && event.StartDate != existingEvent.StartDate {
		return errors.New("cannot change start date to a past date")
	}

	if event.EndDate.Before(event.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Validate capacity
	if event.Capacity < 0 {
		return errors.New("capacity cannot be negative")
	}

	// Update status based on dates
	now := time.Now()
	if now.After(event.EndDate) {
		event.Status = "completed"
	} else if now.After(event.StartDate) {
		event.Status = "ongoing"
	} else {
		event.Status = "active"
	}

	return s.eventRepo.Update(event)
}

func (s *eventService) DeleteEvent(id uint) error {
	// Get existing event
	existingEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if event is ongoing or completed
	if existingEvent.Status == "ongoing" || existingEvent.Status == "completed" {
		return errors.New("cannot delete an ongoing or completed event")
	}

	return s.eventRepo.Delete(id)
}

func (s *eventService) GetEventByID(id uint) (*models.Event, error) {
	return s.eventRepo.FindByID(id)
}

func (s *eventService) GetAllEvents(page, limit int, search string) ([]models.Event, int64, error) {
	return s.eventRepo.FindAll(page, limit, search)
}
