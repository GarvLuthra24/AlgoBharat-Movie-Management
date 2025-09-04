package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type HallServiceImpl struct{}

func (s *HallServiceImpl) GetHalls(theatreID string) ([]models.Hall, error) {
	query := "SELECT id, name, theatre_id, seat_map FROM halls"
	args := []interface{}{} // Use interface{} for dynamic arguments

	if theatreID != "" {
		query += " WHERE theatre_id = ?"
		args = append(args, theatreID)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var halls []models.Hall
	for rows.Next() {
		var hall models.Hall
		var seatMapStr string
		if err := rows.Scan(&hall.ID, &hall.Name, &hall.TheatreID, &seatMapStr); err != nil {
			continue
		}
		// Unmarshal into the new map[string][]int structure
		if err := json.Unmarshal([]byte(seatMapStr), &hall.SeatMap); err != nil {
			log.Printf("Error unmarshaling seat map for hall %s: %v", hall.ID, err)
			continue
		}
		halls = append(halls, hall)
	}

	return halls, nil
}

func (s *HallServiceImpl) GetHall(id string) (models.Hall, error) {
	row := database.DB.QueryRow("SELECT id, name, theatre_id, seat_map FROM halls WHERE id = ?", id)

	var hall models.Hall
	var seatMapStr string
	if err := row.Scan(&hall.ID, &hall.Name, &hall.TheatreID, &seatMapStr); err != nil {
		return models.Hall{}, err
	}
	// Unmarshal into the new map[string][]int structure
	if err := json.Unmarshal([]byte(seatMapStr), &hall.SeatMap); err != nil {
		return models.Hall{}, fmt.Errorf("error unmarshaling seat map for hall %s: %w", hall.ID, err)
	}

	return hall, nil
}

func (s *HallServiceImpl) CreateHall(hall models.Hall) (models.Hall, error) {
	// Validation for the new seat_map structure
	for rowKey, columns := range hall.SeatMap {
		if len(columns) != 3 {
			return models.Hall{}, fmt.Errorf("row %s must have exactly 3 columns", rowKey)
		}
		for colIndex, numSeats := range columns {
			if numSeats < 2 {
				return models.Hall{}, fmt.Errorf("row %s, column %d must have at least 2 seats", rowKey, colIndex+1)
			}
		}
	}

	hall.ID = strconv.Itoa(rand.Intn(1000000))
	seatMapBytes, _ := json.Marshal(hall.SeatMap)
	seatMapStr := string(seatMapBytes)

	stmt, err := database.DB.Prepare("INSERT INTO halls(id, name, theatre_id, seat_map) VALUES(?, ?, ?, ?)")
	if err != nil {
		return models.Hall{}, err
	}
	_, err = stmt.Exec(hall.ID, hall.Name, hall.TheatreID, seatMapStr)
	if err != nil {
		return models.Hall{}, err
	}

	// Create individual seats based on the new seat_map structure
	log.Println("Creating seats for hall:", hall.ID)
	for rowKey, columns := range hall.SeatMap {
		rowNum, _ := strconv.Atoi(rowKey)
		for colIndex, numSeats := range columns {
			for i := 1; i <= numSeats; i++ {
				seat := models.Seat{
					ID:     fmt.Sprintf("%d-%d-%d-%d", rowNum, colIndex+1, i, time.Now().UnixNano()), // Generate highly unique ID
					Row:    rowNum,
					Column: colIndex + 1,
					Number: i,
					HallID: hall.ID,
				}
				seatStmt, err := database.DB.Prepare("INSERT INTO seats(id, row, number, hall_id, column) VALUES(?, ?, ?, ?, ?)")
				if err != nil {
					return models.Hall{}, err
				}
				_, err = seatStmt.Exec(seat.ID, seat.Row, seat.Number, seat.HallID, seat.Column)
				if err != nil {
					log.Printf("Error inserting seat %s: %v", seat.ID, err)
					return models.Hall{}, err
				}
			}
		}
	}

	return hall, nil
}

func (s *HallServiceImpl) GetHallSeats(hallID string) ([]models.Seat, error) {
	rows, err := database.DB.Query("SELECT id, row, number, hall_id, column FROM seats WHERE hall_id = ?", hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.Row, &seat.Number, &seat.HallID, &seat.Column); err != nil {
			continue
		}
		seats = append(seats, seat)
	}

	return seats, nil
}
