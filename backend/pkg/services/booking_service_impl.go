package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// ErrSeatsAlreadyBooked is a custom error type for when seats are already booked.
type ErrSeatsAlreadyBooked struct{}

func (e *ErrSeatsAlreadyBooked) Error() string {
	return "one or more seats are already booked"
}

// ErrNoContiguousSeats is a custom error type for when no contiguous seats are available.
type ErrNoContiguousSeats struct{}

func (e *ErrNoContiguousSeats) Error() string {
	return "no contiguous seats available for the requested show"
}

type BookingServiceImpl struct{}

// CreateBooking attempts to find and book a contiguous block of seats for a specific show.
func (s *BookingServiceImpl) CreateBooking(request BookingRequest) (models.Booking, error) {
	// 1. Find the specific show using a direct and efficient query.
	var targetShow models.Show
	row := database.DB.QueryRow("SELECT id, movie_id, hall_id, time FROM shows WHERE movie_id = ? AND hall_id = ? AND time = ?", request.MovieID, request.HallID, request.Time)
	if err := row.Scan(&targetShow.ID, &targetShow.MovieID, &targetShow.HallID, &targetShow.Time); err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, fmt.Errorf("no show found for the given movie, hall, and time")
		}
		return models.Booking{}, err
	}

	// 2. Find a contiguous block of available seats.
	hall, err := (&HallServiceImpl{}).GetHall(targetShow.HallID)
	if err != nil {
		return models.Booking{}, fmt.Errorf("could not get hall %s: %w", targetShow.HallID, err)
	}

	// Re-organize hall seats by row and then by column and number for easier contiguous check
	seatsByRow := make(map[int][]models.Seat)
	for rowKeyStr, cols := range hall.SeatMap {
		rowNum, _ := strconv.Atoi(rowKeyStr)
		for colIndex, numSeats := range cols {
			for i := 1; i <= numSeats; i++ {
				seatID := fmt.Sprintf("%d-%d-%d", rowNum, colIndex+1, i)
				seatsByRow[rowNum] = append(seatsByRow[rowNum], models.Seat{
					ID: seatID, Row: rowNum, Column: colIndex + 1, Number: i, HallID: hall.ID,
				})
			}
		}
	}

	// Sort seats within each row for consistent contiguous checking
	for rowNum := range seatsByRow {
		sort.Slice(seatsByRow[rowNum], func(i, j int) bool {
			if seatsByRow[rowNum][i].Column != seatsByRow[rowNum][j].Column {
				return seatsByRow[rowNum][i].Column < seatsByRow[rowNum][j].Column
			}
			return seatsByRow[rowNum][i].Number < seatsByRow[rowNum][j].Number
		})
	}

	bookedSeatIDs, err := s.getBookedSeatIDsForShow(targetShow.ID)
	if err != nil {
		return models.Booking{}, fmt.Errorf("could not get booked seats for show %s: %w", targetShow.ID, err)
	}

	var seatsToBook []models.Seat
	foundBlock := false

	// Iterate through rows in ascending order
	var rowNums []int
	for rNum := range seatsByRow {
		rowNums = append(rowNums, rNum)
	}
	sort.Ints(rowNums)

	for _, rowNum := range rowNums {
		seatsInRow := seatsByRow[rowNum]

		currentContiguousCount := 0
		potentialBlock := []models.Seat{}

		for _, seat := range seatsInRow {
			if bookedSeatIDs[seat.ID] {
				currentContiguousCount = 0
				potentialBlock = []models.Seat{}
			} else {
				// Check for contiguity within the same column group
				// A seat is contiguous if it's in the same row, same column, and is the next number
				if currentContiguousCount > 0 {
					prevSeat := potentialBlock[len(potentialBlock)-1]
					// Check if current seat is in the same column and is the next number
					isContiguous := (seat.Column == prevSeat.Column && seat.Number == prevSeat.Number+1)

					if !isContiguous {
						currentContiguousCount = 0
						potentialBlock = []models.Seat{}
					}
				}

				currentContiguousCount++
				potentialBlock = append(potentialBlock, seat)

				if currentContiguousCount >= request.NumSeats {
					seatsToBook = potentialBlock[:request.NumSeats] // Take only the required number of seats
					foundBlock = true
					break
				}
			}
		}
		if foundBlock {
			break
		}
	}

	if !foundBlock {
		return models.Booking{}, &ErrNoContiguousSeats{}
	}

	// 3. Atomically book the found seats.
	var seatIDsToBook []string
	for _, seat := range seatsToBook {
		seatIDsToBook = append(seatIDsToBook, seat.ID)
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return models.Booking{}, err
	}
	defer tx.Rollback()

	// Double-check seats within the transaction to prevent race conditions.
	for _, seatID := range seatIDsToBook {
		var count int
		// Check if any existing booking contains this exact seatID
		// Using INSTR for partial match within the JSON string of seat_ids
		err := tx.QueryRow("SELECT COUNT(*) FROM bookings WHERE show_id = ? AND INSTR(seat_ids, ?)", targetShow.ID, `"`+seatID+`"`).Scan(&count)
		if err != nil {
			return models.Booking{}, err
		}
		if count > 0 {
			return models.Booking{}, fmt.Errorf("seat %s is already booked", seatID)
		}
	}

	newBooking := models.Booking{
		ID:      strconv.Itoa(rand.Intn(1000000)),
		ShowID:  targetShow.ID,
		SeatIDs: seatIDsToBook,
	}

	seatIDsBytes, _ := json.Marshal(newBooking.SeatIDs)
	seatIDsStr := string(seatIDsBytes)

	stmt, err := tx.Prepare("INSERT INTO bookings(id, show_id, seat_ids) VALUES(?, ?, ?)")
	if err != nil {
		return models.Booking{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newBooking.ID, newBooking.ShowID, seatIDsStr)
	if err != nil {
		return models.Booking{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Booking{}, err
	}

	return newBooking, nil
}

// FindAlternativeShows performs a global search for shows on the same day that have enough consecutive seats.
func (s *BookingServiceImpl) FindAlternativeShows(originalTime string, numSeats int) ([]models.Show, error) {
	// 1. Determine the date range for the same day.
	parsedTime, err := time.Parse(time.RFC3339, originalTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}
	year, month, day := parsedTime.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, parsedTime.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 2. Get all shows within that day.
	rows, err := database.DB.Query("SELECT id, movie_id, hall_id, time FROM shows WHERE time >= ? AND time < ?", startOfDay.Format(time.RFC3339), endOfDay.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sameDayShows []models.Show
	for rows.Next() {
		var show models.Show
		if err := rows.Scan(&show.ID, &show.MovieID, &show.HallID, &show.Time); err != nil {
			log.Println(err)
			continue
		}
		sameDayShows = append(sameDayShows, show)
	}

	// 3. Check each same-day show for consecutive seats.
	var alternatives []models.Show
	for _, show := range sameDayShows {
		hall, err := (&HallServiceImpl{}).GetHall(show.HallID)
		if err != nil {
			log.Printf("Could not get hall %s for show %s: %v", show.HallID, show.ID, err)
			continue
		}

		// Re-organize hall seats by row and column for easier contiguous check
		seatsByRow := make(map[int][]models.Seat)
		for rowKeyStr, cols := range hall.SeatMap {
			rowNum, _ := strconv.Atoi(rowKeyStr)
			for colIndex, numSeats := range cols {
				for i := 1; i <= numSeats; i++ {
					seatID := fmt.Sprintf("%d-%d-%d", rowNum, colIndex+1, i)
					seatsByRow[rowNum] = append(seatsByRow[rowNum], models.Seat{
						ID: seatID, Row: rowNum, Column: colIndex + 1, Number: i, HallID: hall.ID,
					})
				}
			}
		}

		// Sort seats within each row for consistent contiguous checking
		for rowNum := range seatsByRow {
			sort.Slice(seatsByRow[rowNum], func(i, j int) bool {
				if seatsByRow[rowNum][i].Column != seatsByRow[rowNum][j].Column {
					return seatsByRow[rowNum][i].Column < seatsByRow[rowNum][j].Column
				}
				return seatsByRow[rowNum][i].Number < seatsByRow[rowNum][j].Number
			})
		}

		bookedSeatIDs, err := s.getBookedSeatIDsForShow(show.ID)
		if err != nil {
			log.Printf("Could not get booked seats for show %s: %v", show.ID, err)
			continue
		}

		foundAlternative := false
		var rowNumsAlt []int
		for rNum := range seatsByRow {
			rowNumsAlt = append(rowNumsAlt, rNum)
		}
		sort.Ints(rowNumsAlt)

		for _, rowNum := range rowNumsAlt {
			seatsInRow := seatsByRow[rowNum]

			currentContiguousCount := 0
			potentialBlock := []models.Seat{}

			for _, seat := range seatsInRow {
				if bookedSeatIDs[seat.ID] {
					currentContiguousCount = 0
					potentialBlock = []models.Seat{}
				} else {
					// Check for contiguity within the same column group
					if currentContiguousCount > 0 {
						prevSeat := potentialBlock[len(potentialBlock)-1]
						// Check if current seat is in the same column and is the next number
						isContiguous := (seat.Column == prevSeat.Column && seat.Number == prevSeat.Number+1)

						if !isContiguous {
							currentContiguousCount = 0
							potentialBlock = []models.Seat{}
						}
					}

					currentContiguousCount++
					potentialBlock = append(potentialBlock, seat)

					if currentContiguousCount >= numSeats {
						foundAlternative = true
						break
					}
				}
			}
			if foundAlternative {
				alternatives = append(alternatives, show)
				break
			}
		}
	}

	return alternatives, nil
}

// getBookedSeatIDsForShow is a helper to get all booked seat IDs for a given show.
func (s *BookingServiceImpl) getBookedSeatIDsForShow(showID string) (map[string]bool, error) {
	rows, err := database.DB.Query("SELECT id, show_id, seat_ids FROM bookings WHERE show_id = ?", showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookedMap := make(map[string]bool)
	for rows.Next() {
		var booking models.Booking
		var seatIDsStr string
		if err := rows.Scan(&booking.ID, &booking.ShowID, &seatIDsStr); err != nil {
			return nil, err
		}

		var seatIDs []string
		if err := json.Unmarshal([]byte(seatIDsStr), &seatIDs); err != nil {
			return nil, err
		}

		for _, id := range seatIDs {
			bookedMap[id] = true
		}
	}
	return bookedMap, nil
}

// GetBookingsByShowID retrieves all bookings for a specific show.
func (s *BookingServiceImpl) GetBookingsByShowID(showID string) ([]models.Booking, error) {
	rows, err := database.DB.Query("SELECT id, show_id, seat_ids FROM bookings WHERE show_id = ?", showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		var seatIDsStr string
		if err := rows.Scan(&booking.ID, &booking.ShowID, &seatIDsStr); err != nil {
			return nil, err
		}

		var seatIDs []string
		if err := json.Unmarshal([]byte(seatIDsStr), &seatIDs); err != nil {
			return nil, err
		}
		booking.SeatIDs = seatIDs
		bookings = append(bookings, booking)
	}
	return bookings, nil
}
