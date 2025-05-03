// controllers/event_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"ticket/models"
	"ticket/service"

	"github.com/gin-gonic/gin"

)

type EventController interface {
	CreateEvent(c *gin.Context)
	UpdateEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
	GetAllEvents(c *gin.Context)
}

type eventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) EventController {
	return &eventController{
		eventService: eventService,
	}
}

func (ctrl *eventController) CreateEvent(c *gin.Context) {
	var event models.Event

	// Validate input
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create event
	err := ctrl.eventService.CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

func (ctrl *eventController) UpdateEvent(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if event exists
	event, err := ctrl.eventService.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Bind input data
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID is preserved
	event.ID = uint(id)

	// Update event
	err = ctrl.eventService.UpdateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   event,
	})
}

func (ctrl *eventController) DeleteEvent(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Delete event
	err = ctrl.eventService.DeleteEvent(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted successfully",
	})
}

func (ctrl *eventController) GetAllEvents(c *gin.Context) {
	// Get query parameters for pagination and search
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	// Get events
	events, total, err := ctrl.eventService.GetAllEvents(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination metadata
	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": events,
		"meta": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"limit":        limit,
		},
	})
}
