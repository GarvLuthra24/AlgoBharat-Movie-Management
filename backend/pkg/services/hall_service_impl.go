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
				seatStmt, err := database.DB.Prepare("INSERT INTO seats(id, `row`, `number`, hall_id, `column`) VALUES(?, ?, ?, ?, ?)")
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

func (s *HallServiceImpl) UpdateHall(hall models.Hall) (models.Hall, error) {
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

	seatMapBytes, _ := json.Marshal(hall.SeatMap)
	seatMapStr := string(seatMapBytes)

	stmt, err := database.DB.Prepare("UPDATE halls SET name = ?, theatre_id = ?, seat_map = ? WHERE id = ?")
	if err != nil {
		return models.Hall{}, err
	}
	_, err = stmt.Exec(hall.Name, hall.TheatreID, seatMapStr, hall.ID)
	if err != nil {
		return models.Hall{}, err
	}

	// Delete existing seats and create new ones
	_, err = database.DB.Exec("DELETE FROM seats WHERE hall_id = ?", hall.ID)
	if err != nil {
		return models.Hall{}, err
	}

	// Create individual seats based on the new seat_map structure
	log.Println("Updating seats for hall:", hall.ID)
	for rowKey, columns := range hall.SeatMap {
		rowNum, _ := strconv.Atoi(rowKey)
		for colIndex, numSeats := range columns {
			for i := 1; i <= numSeats; i++ {
				seat := models.Seat{
					ID:     fmt.Sprintf("%d-%d-%d-%d", rowNum, colIndex+1, i, time.Now().UnixNano()),
					Row:    rowNum,
					Column: colIndex + 1,
					Number: i,
					HallID: hall.ID,
				}
				seatStmt, err := database.DB.Prepare("INSERT INTO seats(id, `row`, `number`, hall_id, `column`) VALUES(?, ?, ?, ?, ?)")
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

func (s *HallServiceImpl) DeleteHall(id string) error {
	// Check if hall exists
	_, err := s.GetHall(id)
	if err != nil {
		return fmt.Errorf("hall not found: %w", err)
	}

	// Get all shows for this hall
	showRows, err := database.DB.Query("SELECT id FROM shows WHERE hall_id = ?", id)
	if err != nil {
		return fmt.Errorf("error getting shows for hall: %w", err)
	}
	defer showRows.Close()

	var showIDs []string
	for showRows.Next() {
		var showID string
		if err := showRows.Scan(&showID); err != nil {
			continue
		}
		showIDs = append(showIDs, showID)
	}

	// Delete all bookings for shows in this hall
	if len(showIDs) > 0 {
		// Build the IN clause for show IDs
		placeholders := make([]string, len(showIDs))
		args := make([]interface{}, len(showIDs))
		for i, showID := range showIDs {
			placeholders[i] = "?"
			args[i] = showID
		}

		placeholdersStr := ""
		for i, placeholder := range placeholders {
			if i > 0 {
				placeholdersStr += ","
			}
			placeholdersStr += placeholder
		}

		_, err = database.DB.Exec(fmt.Sprintf("DELETE FROM bookings WHERE show_id IN (%s)", placeholdersStr), args...)
		if err != nil {
			log.Printf("Warning: error deleting bookings for hall %s: %v", id, err)
		}
	}

	// Delete all shows for this hall
	_, err = database.DB.Exec("DELETE FROM shows WHERE hall_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting shows for hall: %w", err)
	}

	// Delete all seats for this hall
	_, err = database.DB.Exec("DELETE FROM seats WHERE hall_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting seats for hall: %w", err)
	}

	// Delete the hall
	_, err = database.DB.Exec("DELETE FROM halls WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting hall: %w", err)
	}

	log.Printf("Hall %s and its %d shows deleted successfully", id, len(showIDs))
	return nil
}

func (s *HallServiceImpl) GetHallSeats(hallID string) ([]models.Seat, error) {
	rows, err := database.DB.Query("SELECT id, `row`, `number`, hall_id, `column` FROM seats WHERE hall_id = ?", hallID)
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
