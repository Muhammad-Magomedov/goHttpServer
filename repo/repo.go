package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

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

func (s *Server) CreateEvent(ctx context.Context, userID int, date time.Time, title string) error {
	_, err := s.db.Exec(ctx, "insert into events (date, user_id, title) values ($1, $2, $3)", date, userID, title)
	return err
}

func (s *Server) DeleteEvent(ctx context.Context, date time.Time, id int) error {
	_, err := s.db.Exec(ctx, "delete from events where id=$1 and date=$2", id, date)
	if err != nil {
		fmt.Println("error while deleting event", err)
		return err
	}

	return nil
}

func (s *Server) GetEventsForDay(ctx context.Context, userID int, date time.Time) ([]Result, error) {
	var results []Result
	rows, err := s.db.Query(ctx, "select user_id, date, title from events where user_id=$1 and date=$2", userID, date)
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

func (s *Server) GetEventsForWeek(ctx context.Context, userID int, startDate time.Time, endDate time.Time) ([]Result, error) {
	var results []Result
	rows, err := s.db.Query(ctx, "select user_id, date from events where user_id=$1 and date between $2 and $3", userID, startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func (s *Server) UpdateEvent(ctx context.Context, id int, userId int, title string) (Result, error) {
	var result Result
	err := s.db.QueryRow(ctx, "update events set title=$1 where user_id=$2 and id=$3", title, userId, id).Scan(&result.Date, &result.UserID, &result.Title)
	if err != nil {
		fmt.Println("Error query for updating event", err)
		return Result{}, err
	}

	return result, nil
}
