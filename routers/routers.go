package routers

import (
	"ticket/controllers"
	"ticket/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter menginisialisasi semua rute aplikasi dengan base path /api
func SetupRouter(
	jwtMiddleware middleware.JWTMiddleware,
	authController controllers.AuthController,
	eventController controllers.EventController,
	ticketController controllers.TicketController,
	reportController controllers.ReportController,
) *gin.Engine {
	router := gin.Default()

	// Global CORS
	router.Use(middleware.CORS())

	// Grupkan semuanya di bawah /api
	api := router.Group("/api")
	{
		// — Public routes —
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)
		api.GET("/events", eventController.GetAllEvents)

		// — Protected routes —
		protected := api.Group("")
		protected.Use(jwtMiddleware.JWTAuth())

		// Routes untuk user & admin
		userGroup := protected.Group("/tickets")
		userGroup.Use(middleware.Authorize("user", "admin"))
		{
			userGroup.GET("", ticketController.GetUserTickets)
			userGroup.POST("", ticketController.BuyTicket)
			userGroup.GET("/:id", ticketController.GetTicketByID)
			userGroup.PATCH("/:id", ticketController.CancelTicket)
		}

		// Routes khusus admin
		adminGroup := protected.Group("/events")
		adminGroup.Use(middleware.Authorize("admin"))
		{
			adminGroup.POST("", eventController.CreateEvent)
			adminGroup.PUT(":id", eventController.UpdateEvent)
			adminGroup.DELETE(":id", eventController.DeleteEvent)

			adminGroup.GET("summary", reportController.GetSummaryReport)
			adminGroup.GET("/reports/:id", reportController.GetEventReport)
		}
	}

	return router
}
