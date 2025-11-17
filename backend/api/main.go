package main

import (
	"backend/utils"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from system:", err)
	}

	client := utils.Open(os.Getenv("DATABASE_URL"))
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("gagal membuat skema: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Tambahkan CORS (PENTING untuk frontend React Anda)
	e.Use(middleware.CORS())

	h := &Handler{DB: client}

	e.GET("/listings", h.getListings)
	e.GET("/events", h.getEvents)
	e.GET("/profiles/:address", h.getUserProfile)
	e.GET("/accessories", h.getAccessories)
	e.GET("/moments", h.getMoments)

	e.POST("/moment/free", h.freeMintMoment)
	e.POST("/moment/with-event-pass", h.mintMomentWithEventPass)

	log.Println("Server API dimulai di http://localhost:8000")
	e.Logger.Fatal(e.Start(":8000"))
}
