package services

import (
	"algoBharat/backend/pkg/database"
	"encoding/json"
	"log"
)

type AnalyticsServiceImpl struct{}

func (s *AnalyticsServiceImpl) GetMovieRevenue(movieID string) (float64, error) {
	// Get all shows for the movie, including their price
	rows, err := database.DB.Query("SELECT id, price FROM shows WHERE movie_id = ?", movieID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalRevenue float64
	for rows.Next() {
		var showID string
		var showPrice float64
		if err := rows.Scan(&showID, &showPrice); err != nil {
			log.Println(err)
			continue
		}

		bookingRows, err := database.DB.Query("SELECT seat_ids FROM bookings WHERE show_id = ?", showID)
		if err != nil {
			log.Println(err)
			continue
		}
		defer bookingRows.Close()

		var seatsBookedInShow int
		for bookingRows.Next() {
			var seatIDsStr string
			if err := bookingRows.Scan(&seatIDsStr); err != nil {
				log.Println(err)
				continue
			}

			var seatIDs []string
			if err := json.Unmarshal([]byte(seatIDsStr), &seatIDs); err != nil {
				return 0, err
			}

			seatsBookedInShow += len(seatIDs)
		}
		totalRevenue += float64(seatsBookedInShow) * showPrice
	}

	return totalRevenue, nil
}
