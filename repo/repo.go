package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type CreateEvent struct {
	UserID int
	Date   time.Time
	Title  string
}

type Server struct {
	db *pgx.Conn
}

type Result struct {
	Date   time.Time
	UserID int
	Title  string
}

func New(db *pgx.Conn) *Server {
	return &Server{
		db: db,
	}
}

func (event *CreateEvent) Validate() error {
	if event.UserID <= 0 {
		return fmt.Errorf("userID must be positive")
	}
	if event.Title == "" {
		return fmt.Errorf("title must not be empty")
	}
	if event.Date.IsZero() {
		return fmt.Errorf("date must be set")
	}
	return nil
}

func (s *Server) CreateEvent(ctx context.Context, event CreateEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	_, err := s.db.Exec(ctx, "insert into events (date, user_id, title) values ($1, $2, $3)", event.Date, event.UserID, event.Title)
	return err
}

func (s *Server) DeleteEvent(ctx context.Context, date time.Time, id int) error {
	if id <= 0 {
		return fmt.Errorf("id must be positive")
	}
	if date.IsZero() {
		return fmt.Errorf("date must be set")
	}
	_, err := s.db.Exec(ctx, "delete from events where id=$1 and date=$2", id, date)
	if err != nil {
		fmt.Println("error while deleting event", err)
		return err
	}

	return nil
}

func (s *Server) GetEventsForDay(ctx context.Context, id int, date time.Time) ([]Result, error) {
	var results []Result
	rows, err := s.db.Query(ctx, "select id, date, title from events where id=$1 and date=$2", id, date)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("cant find events for this day", err)
			return nil, err
		}
		fmt.Println("Error query", err)
		return nil, err
	}

	for rows.Next() {
		var result Result
		err := rows.Scan(&result.UserID, &result.Date, &result.Title)

		if err != nil {
			return nil, fmt.Errorf("getEventForDay next: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *Server) GetEventsForDates(ctx context.Context, id int, startDate time.Time, endDate time.Time) ([]Result, error) {
	var results []Result
	rows, err := s.db.Query(ctx, "select id, date, title from events where id=$1 and date between $2 and $3", id, startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("cant find events for this week", err)
			return nil, err
		}
		fmt.Println("Error query", err)
		return nil, err
	}

	for rows.Next() {
		var result Result
		err := rows.Scan(&result.UserID, &result.Date, &result.Title)

		if err != nil {
			return nil, fmt.Errorf("GetEventsForDates next: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *Server) UpdateEvent(ctx context.Context, id int, title string) (Result, error) {
	if id <= 0 {
		return Result{}, fmt.Errorf("id must be positive")
	}
	if title == "" {
		return Result{}, fmt.Errorf("title must not be empty")
	}
	var result Result
	err := s.db.QueryRow(ctx, "update events set title=$1 where id=$2 RETURNING date, id, title", title, id).Scan(&result.Date, &result.UserID, &result.Title)
	if err != nil {
		fmt.Println("Error query for updating event", err)
		return Result{}, err
	}

	return result, nil
}
