// controllers/report_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"ticket/service"

	"github.com/gin-gonic/gin"
)

type ReportController interface {
	GetSummaryReport(c *gin.Context)
	GetEventReport(c *gin.Context)
}

type reportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) ReportController {
	return &reportController{
		reportService: reportService,
	}
}

func (ctrl *reportController) GetSummaryReport(c *gin.Context) {
	// Get summary report
	data, err := ctrl.reportService.GetSummaryReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ctrl *reportController) GetEventReport(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get event report
	data, err := ctrl.reportService.GetEventReport(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
