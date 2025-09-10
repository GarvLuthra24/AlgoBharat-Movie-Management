package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
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

func (s *BookingServiceImpl) CreateBooking(request BookingRequest) (models.Booking, error) {
	// 1. Parse the request time
	requestTime, err := time.Parse(time.RFC3339, request.Time)
	if err != nil {
		formats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05+05:30",
		}
		for _, format := range formats {
			if requestTime, err = time.Parse(format, request.Time); err == nil {
				break
			}
		}
		if err != nil {
			return models.Booking{}, fmt.Errorf("invalid time format: %w", err)
		}
	}
	requestTime = requestTime.UTC()

	// 2. Find target show
	var targetShow models.Show
	// we should also be sending the time to get that particular show
	row := database.DB.QueryRow("SELECT id, movie_id, hall_id, time FROM shows WHERE movie_id = ? AND hall_id = ?", request.MovieID, request.HallID)
	if err := row.Scan(&targetShow.ID, &targetShow.MovieID, &targetShow.HallID, &targetShow.Time); err != nil {
		if err == sql.ErrNoRows {
			return models.Booking{}, fmt.Errorf("no show found for the given movie and hall")
		}
		return models.Booking{}, err
	}

	showTime, err := time.Parse(time.RFC3339, targetShow.Time)
	if err != nil {
		return models.Booking{}, fmt.Errorf("invalid show time format in database: %w", err)
	}
	showTime = showTime.UTC()

	if !requestTime.Truncate(time.Minute).Equal(showTime.Truncate(time.Minute)) {
		return models.Booking{}, fmt.Errorf("no show found for the given movie, hall, and time")
	}

	// 3. Get Hall and prepare seats
	hall, err := (&HallServiceImpl{}).GetHall(targetShow.HallID)
	if err != nil {
		return models.Booking{}, fmt.Errorf("could not get hall %s: %w", targetShow.HallID, err)
	}

	seatsByRow := make(map[int][]models.Seat)
	for rowKeyStr, cols := range hall.SeatMap {
		rowNum, _ := strconv.Atoi(rowKeyStr)
		for colIndex, numSeats := range cols {
			for i := 1; i <= numSeats; i++ {
				seatID := fmt.Sprintf("%d-%d-%d", rowNum, colIndex+1, i)
				seatsByRow[rowNum] = append(seatsByRow[rowNum], models.Seat{
					ID:     seatID,
					Row:    rowNum,
					Column: colIndex + 1,
					Number: i,
					HallID: hall.ID,
				})
			}
		}
	}

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
				if currentContiguousCount > 0 {
					prevSeat := potentialBlock[len(potentialBlock)-1]
					isContiguous := (seat.Column == prevSeat.Column && seat.Number == prevSeat.Number+1)
					if !isContiguous {
						currentContiguousCount = 0
						potentialBlock = []models.Seat{}
					}
				}

				currentContiguousCount++
				potentialBlock = append(potentialBlock, seat)

				if currentContiguousCount >= request.NumSeats {
					seatsToBook = potentialBlock[:request.NumSeats]
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

	seatIDsToBook := make([]string, len(seatsToBook))
	for i, seat := range seatsToBook {
		seatIDsToBook[i] = seat.ID
	}

	// 4. Transactional booking with seat_ids and new booked_seats
	tx, err := database.DB.Begin()
	if err != nil {
		return models.Booking{}, err
	}
	defer tx.Rollback()

	newBooking := models.Booking{
		ID:      strconv.Itoa(rand.Intn(1000000)),
		ShowID:  targetShow.ID,
		SeatIDs: seatIDsToBook,
	}

	seatIDsBytes, _ := json.Marshal(newBooking.SeatIDs)
	seatIDsStr := string(seatIDsBytes)

	// Insert booking with seat_ids
	stmtBooking, err := tx.Prepare("INSERT INTO bookings(id, show_id, seat_ids) VALUES(?, ?, ?)")
	if err != nil {
		return models.Booking{}, err
	}
	defer stmtBooking.Close()

	_, err = stmtBooking.Exec(newBooking.ID, newBooking.ShowID, seatIDsStr)
	if err != nil {
		return models.Booking{}, err
	}

	// Insert booked_seats with atomic constraint
	stmtSeat, err := tx.Prepare("INSERT INTO booked_seats(show_id, seat_id, booking_id) VALUES(?, ?, ?)")
	if err != nil {
		return models.Booking{}, err
	}
	defer stmtSeat.Close()

	for _, seatID := range seatIDsToBook {
		_, err := stmtSeat.Exec(newBooking.ShowID, seatID, newBooking.ID)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				return models.Booking{}, fmt.Errorf("seat %s is already booked", seatID)
			}
			return models.Booking{}, err
		}
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
		formats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05+05:30",
		}
		for _, format := range formats {
			if parsedTime, err = time.Parse(format, originalTime); err == nil {
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("invalid time format: %w", err)
		}
	}

	parsedTime = parsedTime.UTC()
	year, month, day := parsedTime.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 2. Get all shows within that day.
	rows, err := database.DB.Query(
		"SELECT id, movie_id, hall_id, time FROM shows WHERE time >= ? AND time < ?",
		startOfDay.Format(time.RFC3339),
		endOfDay.Format(time.RFC3339),
	)
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

	// 3. Check each same-day show for consecutive available seats.
	var alternatives []models.Show
	for _, show := range sameDayShows {
		hall, err := (&HallServiceImpl{}).GetHall(show.HallID)
		if err != nil {
			log.Printf("Could not get hall %s for show %s: %v", show.HallID, show.ID, err)
			continue
		}

		seatsByRow := make(map[int][]models.Seat)
		for rowKeyStr, cols := range hall.SeatMap {
			rowNum, _ := strconv.Atoi(rowKeyStr)
			for colIndex, numSeats := range cols {
				for i := 1; i <= numSeats; i++ {
					seatID := fmt.Sprintf("%d-%d-%d", rowNum, colIndex+1, i)
					seatsByRow[rowNum] = append(seatsByRow[rowNum], models.Seat{
						ID:     seatID,
						Row:    rowNum,
						Column: colIndex + 1,
						Number: i,
						HallID: hall.ID,
					})
				}
			}
		}

		for rowNum := range seatsByRow {
			sort.Slice(seatsByRow[rowNum], func(i, j int) bool {
				if seatsByRow[rowNum][i].Column != seatsByRow[rowNum][j].Column {
					return seatsByRow[rowNum][i].Column < seatsByRow[rowNum][j].Column
				}
				return seatsByRow[rowNum][i].Number < seatsByRow[rowNum][j].Number
			})
		}

		// NEW: Consistently use booked_seats table
		bookedSeatIDs, err := s.getBookedSeatIDsForShow(show.ID)
		if err != nil {
			log.Printf("Could not get booked seats for show %s: %v", show.ID, err)
			continue
		}

		foundAlternative := false
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
					if currentContiguousCount > 0 {
						prevSeat := potentialBlock[len(potentialBlock)-1]
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
	bookedMap := make(map[string]bool)

	// Prefer querying booked_seats directly to get atomic data
	rows, err := database.DB.Query("SELECT seat_id FROM booked_seats WHERE show_id = ?", showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seatID string
		if err := rows.Scan(&seatID); err != nil {
			return nil, err
		}
		bookedMap[seatID] = true
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
