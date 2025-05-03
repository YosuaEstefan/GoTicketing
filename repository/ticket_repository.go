package repository

import (
	"ticket/models"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Save(ticket *models.Ticket) error
	Update(ticket *models.Ticket) error
	FindByID(id uint) (*models.Ticket, error)
	FindByUser(userID uint, page, limit int) ([]models.Ticket, int64, error)
	CountByEvent(eventID uint) (int, error)
	GetSummary() (map[string]interface{}, error)
	GetEventReport(eventID uint) (map[string]interface{}, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) Save(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Event").Preload("User").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByUser(userID uint, page, limit int) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	query := r.db.Model(&models.Ticket{}).Where("user_id = ?", userID).Preload("Event")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err = query.Find(&tickets).Error
	if err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

func (r *ticketRepository) CountByEvent(eventID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Ticket{}).Where("event_id = ? AND status != ?", eventID, "cancelled").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ticketRepository) GetSummary() (map[string]interface{}, error) {
	// Calculate total tickets sold and revenue
	type SummaryResult struct {
		TotalTickets int
		TotalRevenue float64
	}

	var summary SummaryResult
	err := r.db.Model(&models.Ticket{}).
		Select("count(*) as total_tickets, sum(price) as total_revenue").
		Where("status != ?", "cancelled").
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	// Get event-wise breakdown
	type EventSummary struct {
		EventID     uint
		EventName   string
		TicketsSold int
		Revenue     float64
	}

	var eventSummaries []EventSummary
	err = r.db.Model(&models.Ticket{}).
		Select("tickets.event_id, events.name as event_name, count(*) as tickets_sold, sum(tickets.price) as revenue").
		Joins("JOIN events ON tickets.event_id = events.id").
		Where("tickets.status != ?", "cancelled").
		Group("tickets.event_id").
		Scan(&eventSummaries).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_tickets": summary.TotalTickets,
		"total_revenue": summary.TotalRevenue,
		"events":        eventSummaries,
	}

	return result, nil
}

func (r *ticketRepository) GetEventReport(eventID uint) (map[string]interface{}, error) {
	// Get event details
	var event models.Event
	err := r.db.First(&event, eventID).Error
	if err != nil {
		return nil, err
	}

	// Calculate tickets sold and revenue
	type EventStats struct {
		TicketsSold int
		Revenue     float64
	}

	var stats EventStats
	err = r.db.Model(&models.Ticket{}).
		Select("count(*) as tickets_sold, sum(price) as revenue").
		Where("event_id = ? AND status != ?", eventID, "cancelled").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// Get ticket status breakdown
	type StatusCount struct {
		Status string
		Count  int
	}

	var statusCounts []StatusCount
	err = r.db.Model(&models.Ticket{}).
		Select("status, count(*) as count").
		Where("event_id = ?", eventID).
		Group("status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"event":            event,
		"tickets_sold":     stats.TicketsSold,
		"revenue":          stats.Revenue,
		"capacity":         event.Capacity,
		"available":        event.Capacity - stats.TicketsSold,
		"status_breakdown": statusCounts,
	}

	return result, nil
}
