package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ShowServiceImpl struct{}

func (s *ShowServiceImpl) GetShows() ([]models.Show, error) {
	rows, err := database.DB.Query("SELECT id, movie_id, hall_id, time, price FROM shows")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []models.Show
	for rows.Next() {
		var show models.Show
		if err := rows.Scan(&show.ID, &show.MovieID, &show.HallID, &show.Time, &show.Price); err != nil {
			log.Println(err)
			continue
		}
		shows = append(shows, show)
	}

	return shows, nil
}

func (s *ShowServiceImpl) CreateShow(show models.Show) (models.Show, error) {
	// 1. Get movie duration
	movieService := &MovieServiceImpl{}
	movie, err := movieService.GetMovie(show.MovieID)
	if err != nil {
		return models.Show{}, fmt.Errorf("could not get movie details: %w", err)
	}

	// 2. Parse show time and calculate end time
	showStartTime, err := time.Parse(time.RFC3339, show.Time)
	if err != nil {
		return models.Show{}, fmt.Errorf("invalid show time format: %w", err)
	}
	showEndTime := showStartTime.Add(time.Duration(movie.DurationMinutes) * time.Minute)

	// 3. Check for overlaps with existing shows in the same hall
	existingShowsRows, err := database.DB.Query("SELECT id, movie_id, time FROM shows WHERE hall_id = ?", show.HallID)
	if err != nil {
		return models.Show{}, fmt.Errorf("could not query existing shows: %w", err)
	}
	defer existingShowsRows.Close()

	for existingShowsRows.Next() {
		var existingShow models.Show
		if err := existingShowsRows.Scan(&existingShow.ID, &existingShow.MovieID, &existingShow.Time); err != nil {
			log.Printf("Error scanning existing show: %v", err)
			continue
		}

		// Get existing movie duration
		existingMovie, err := movieService.GetMovie(existingShow.MovieID)
		if err != nil {
			log.Printf("Could not get existing movie details for show %s: %v", existingShow.ID, err)
			continue
		}

		existingShowStartTime, err := time.Parse(time.RFC3339, existingShow.Time)
		if err != nil {
			log.Printf("Invalid time format for existing show %s: %v", existingShow.ID, err)
			continue
		}
		existingShowEndTime := existingShowStartTime.Add(time.Duration(existingMovie.DurationMinutes) * time.Minute)

		// Check for overlap: (start1 < end2 && end1 > start2)
		if showStartTime.Before(existingShowEndTime) && showEndTime.After(existingShowStartTime) {
			return models.Show{}, fmt.Errorf("show overlaps with existing show %s in the same hall", existingShow.ID)
		}
	}

	// 4. If no overlap, proceed with insertion
	show.ID = strconv.Itoa(rand.Intn(1000000))
	stmt, err := database.DB.Prepare("INSERT INTO shows(id, movie_id, hall_id, time, price) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return models.Show{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(show.ID, show.MovieID, show.HallID, show.Time, show.Price)
	if err != nil {
		return models.Show{}, err
	}

	return show, nil
}
