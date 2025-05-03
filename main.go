// main.go (dengan JWT di utils)
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
	// Memuat konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// Inisialisasi database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Menjalankan migrasi
	config.RunMigrations(db)

	// Inisialisasi repository
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// Inisialisasi utils
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpiration)

	// Inisialisasi service
	userService := service.NewUserService(userRepo, jwtUtil)
	eventService := service.NewEventService(eventRepo)
	ticketService := service.NewTicketService(ticketRepo, eventRepo, userRepo)
	reportService := service.NewReportService(ticketRepo, eventRepo)

	// Inisialisasi controller
	authController := controllers.NewAuthController(userService)
	eventController := controllers.NewEventController(eventService)
	ticketController := controllers.NewTicketController(ticketService)
	reportController := controllers.NewReportController(reportService)

	// Inisialisasi middleware yang membutuhkan JWT
	jwtMiddleware := middleware.NewJWTMiddleware(jwtUtil)

	// Inisialisasi router
	router := routers.SetupRouter(
		jwtMiddleware,
		authController,
		eventController,
		ticketController,
		reportController,
	)

	// Memulai server
	log.Printf("Server dimulai pada %s\n", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}
