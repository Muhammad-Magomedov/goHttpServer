package main

import (
	"context"
	"event/handler"
	"event/repo"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("DB_URL")
	if err != nil {
		log.Fatal("DB environment variable is not set")
	}
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB environment variable is not set")
	}
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к бд", err)
	}

	server := repo.New(db)
	eventHandler := handler.New(server)

	http.HandleFunc("/create_event", eventHandler.CreateEventHandler)
	http.HandleFunc("/events_for_day", eventHandler.GetEventsForDayHandler)
	http.HandleFunc("/events_for_week", eventHandler.GetEventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventHandler.GetEventsForMonthHandler)
	http.HandleFunc("/delete_event", eventHandler.RemoveEventHandler)
	http.HandleFunc("/update_event", eventHandler.UpdateEventHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
