package main

import (
	"context"
	"event/handler"
	"event/repo"
	"event/utils"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := "postgres://postgres:postgres@localhost:5432/events"
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к бд", err)
	}

	server := repo.New(db)
	utils := utils.New()
	eventHandler := handler.New(server, utils)

	http.HandleFunc("/create_event", eventHandler.CreateEventHandler)
	http.HandleFunc("/events_for_day", eventHandler.GetEventsForDayHandler)
	http.HandleFunc("/events_for_week", eventHandler.GetEventsForWeekHandler)
	http.HandleFunc("/events_for_month", eventHandler.GetEventsForMonthHandler)
	http.HandleFunc("/delete_event", eventHandler.RemoveEventHandler)
	http.HandleFunc("/update_event", eventHandler.UpdateEventHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
