// controller/ticket_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"ticket/models"
	"ticket/service"

	"github.com/gin-gonic/gin"
)

type TicketController interface {
	BuyTicket(c *gin.Context)
	CancelTicket(c *gin.Context)
	GetTicketByID(c *gin.Context)
	GetUserTickets(c *gin.Context)
}

type ticketController struct {
	ticketService service.TicketService
}

func NewTicketController(ticketService service.TicketService) TicketController {
	return &ticketController{
		ticketService: ticketService,
	}
}

func (ctrl *ticketController) BuyTicket(c *gin.Context) {
	var ticketRequest struct {
		EventID uint `json:"event_id" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&ticketRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create ticket
	ticket := models.Ticket{
		EventID: ticketRequest.EventID,
		UserID:  userID.(uint),
	}

	// Save ticket
	err := ctrl.ticketService.BuyTicket(&ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket purchased successfully",
		"ticket":  ticket,
	})
}

func (ctrl *ticketController) CancelTicket(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Cancel ticket
	err = ctrl.ticketService.CancelTicket(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket cancelled successfully",
	})
}

func (ctrl *ticketController) GetTicketByID(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get ticket
	ticket, err := ctrl.ticketService.GetTicketByID(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket": ticket,
	})
}

func (ctrl *ticketController) GetUserTickets(c *gin.Context) {
	// Get query parameters for pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get tickets
	tickets, total, err := ctrl.ticketService.GetUserTickets(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination metadata
	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": tickets,
		"meta": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"limit":        limit,
		},
	})
}
