// service/ticket_service.go
package service

import (
	"errors"
	"ticket/models"
	"ticket/repository"
	"time"

)

type TicketService interface {
	BuyTicket(ticket *models.Ticket) error
	CancelTicket(id uint, userID uint) error
	GetTicketByID(id uint, userID uint) (*models.Ticket, error)
	GetUserTickets(userID uint, page, limit int) ([]models.Ticket, int64, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
	userRepo   repository.UserRepository
}

func NewTicketService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository, userRepo repository.UserRepository) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
		userRepo:   userRepo,
	}
}

func (s *ticketService) BuyTicket(ticket *models.Ticket) error {
	// Get event
	event, err := s.eventRepo.FindByID(ticket.EventID)
	if err != nil {
		return err
	}

	// Check if event is active
	if event.Status != "active" {
		return errors.New("tickets can only be purchased for active events")
	}

	// Check event capacity
	ticketCount, err := s.ticketRepo.CountByEvent(event.ID)
	if err != nil {
		return err
	}

	if ticketCount >= event.Capacity {
		return errors.New("event is sold out")
	}

	// Set ticket price from event
	ticket.Price = event.Price
	ticket.Status = "active"
	ticket.PurchasedAt = time.Now()

	return s.ticketRepo.Save(ticket)
}

func (s *ticketService) CancelTicket(id uint, userID uint) error {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if ticket belongs to user
	if ticket.UserID != userID {
		// Check if user is admin
		user, err := s.userRepo.FindByID(userID)
		if err != nil || user.Role != "admin" {
			return errors.New("unauthorized access to ticket")
		}
	}

	// Check if ticket is already cancelled
	if ticket.Status == "cancelled" {
		return errors.New("ticket is already cancelled")
	}

	// Check if event has already started
	if time.Now().After(ticket.Event.StartDate) {
		return errors.New("cannot cancel tickets for events that have already started")
	}

	// Cancel ticket
	ticket.Status = "cancelled"
	return s.ticketRepo.Update(ticket)
}

func (s *ticketService) GetTicketByID(id uint, userID uint) (*models.Ticket, error) {
	// Get ticket
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if ticket belongs to user
	if ticket.UserID != userID {
		// Check if user is admin
		user, err := s.userRepo.FindByID(userID)
		if err != nil || user.Role != "admin" {
			return nil, errors.New("unauthorized access to ticket")
		}
	}

	return ticket, nil
}

func (s *ticketService) GetUserTickets(userID uint, page, limit int) ([]models.Ticket, int64, error) {
	return s.ticketRepo.FindByUser(userID, page, limit)
}
