package main

import (
	"log"
	"ticket/config"
	"ticket/controllers"
	"ticket/middleware"
	"ticket/repository"
	"ticket/routers"
	"ticket/service"
	"ticket/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	config.RunMigrations(db)

	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpiration)

	userService := service.NewUserService(userRepo, jwtUtil)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo, eventRepo, userRepo)
	reportService := service.NewReportService(ticketRepo, eventRepo)

	authController := controllers.NewAuthController(userService)
	eventController := controllers.NewEventController(eventService)
	ticketController := controllers.NewTicketController(ticketService)
	reportController := controllers.NewReportController(reportService)

	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtil)

	router := routers.SetupRouter(
		jwtMiddleware,
		authController,
		eventController,
		ticketController,
		reportController,
	)

	log.Printf("Server dimulai pada %s\n", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}
