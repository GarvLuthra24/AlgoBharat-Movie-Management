package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"math/rand"
	"strconv"
)

type MovieServiceImpl struct{}

func (s *MovieServiceImpl) GetMovies() ([]models.Movie, error) {
	rows, err := database.DB.Query("SELECT id, title, duration_minutes FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.DurationMinutes); err != nil {
			continue
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *MovieServiceImpl) GetMovie(id string) (models.Movie, error) {
	row := database.DB.QueryRow("SELECT id, title, duration_minutes FROM movies WHERE id = ?", id)

	var movie models.Movie
	if err := row.Scan(&movie.ID, &movie.Title, &movie.DurationMinutes); err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (s *MovieServiceImpl) CreateMovie(movie models.Movie) (models.Movie, error) {
	movie.ID = strconv.Itoa(rand.Intn(1000000))

	stmt, err := database.DB.Prepare("INSERT INTO movies(id, title, duration_minutes) VALUES(?, ?, ?)")
	if err != nil {
		return models.Movie{}, err
	}
	_, err = stmt.Exec(movie.ID, movie.Title, movie.DurationMinutes)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (s *MovieServiceImpl) UpdateMovie(id string, movie models.Movie) (models.Movie, error) {
	stmt, err := database.DB.Prepare("UPDATE movies SET title = ?, duration_minutes = ? WHERE id = ?")
	if err != nil {
		return models.Movie{}, err
	}
	_, err = stmt.Exec(movie.Title, movie.DurationMinutes, id)
	if err != nil {
		return models.Movie{}, err
	}

	movie.ID = id
	return movie, nil
}

func (s *MovieServiceImpl) DeleteMovie(id string) error {
	stmt, err := database.DB.Prepare("DELETE FROM movies WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
