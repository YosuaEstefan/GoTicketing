// service/report_service.go
package service

import "ticket/repository"

type ReportService interface {
	GetSummaryReport() (map[string]interface{}, error)
	GetEventReport(eventID uint) (map[string]interface{}, error)
}

type reportService struct {
	ticketRepo repository.TicketRepository
	eventRepo  repository.EventRepository
}

func NewReportService(ticketRepo repository.TicketRepository, eventRepo repository.EventRepository) ReportService {
	return &reportService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}
}

func (s *reportService) GetSummaryReport() (map[string]interface{}, error) {
	return s.ticketRepo.GetSummary()
}

func (s *reportService) GetEventReport(eventID uint) (map[string]interface{}, error) {
	// Check if event exists
	_, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	return s.ticketRepo.GetEventReport(eventID)
}
