package app

import (
	"FriendsAdvice/internal/database/postgresql"
	"FriendsAdvice/internal/services"
	transport "FriendsAdvice/internal/transport/rest"
	"log"
	"net/http"
	"os"
)

func Run() {
	// 1.Storage manager - LAYER OF STORING
	connectDTO := createConnectionDTO()
	storageManager, err := postgresql.InitManager(connectDTO)
	if err != nil {
		log.Fatal(err)
	}

	// 2.Controller - LAYER OF BUSINESS-LOGIC
	controller := services.InitController(storageManager)

	// 3.Router - LAYER OF THE TRANSPOT
	router := transport.CreateRouter(controller)

	http.ListenAndServe(":8000", router.MuxRouter)
}

func createConnectionDTO() *postgresql.ConnectionDTO {
	return &postgresql.ConnectionDTO{
		HostName: os.Getenv("DATABASE_HOST_NAME"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Port:     os.Getenv("DATABASE_PORT"),
		DBName:   os.Getenv("DATABASE_NAME")}
}
