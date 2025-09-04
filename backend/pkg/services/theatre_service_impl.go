package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"math/rand"
	"strconv"
)

type TheatreServiceImpl struct{}

func (s *TheatreServiceImpl) GetTheatres() ([]models.Theatre, error) {
	rows, err := database.DB.Query("SELECT id, name FROM theatres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theatres []models.Theatre
	for rows.Next() {
		var theatre models.Theatre
		if err := rows.Scan(&theatre.ID, &theatre.Name); err != nil {
			continue
		}
		theatres = append(theatres, theatre)
	}

	return theatres, nil
}

func (s *TheatreServiceImpl) GetTheatre(id string) (models.Theatre, error) {
	row := database.DB.QueryRow("SELECT id, name FROM theatres WHERE id = ?", id)

	var theatre models.Theatre
	if err := row.Scan(&theatre.ID, &theatre.Name); err != nil {
		return models.Theatre{}, err
	}

	return theatre, nil
}

func (s *TheatreServiceImpl) CreateTheatre(theatre models.Theatre) (models.Theatre, error) {
	theatre.ID = strconv.Itoa(rand.Intn(1000000))

	stmt, err := database.DB.Prepare("INSERT INTO theatres(id, name) VALUES(?, ?)")
	if err != nil {
		return models.Theatre{}, err
	}
	_, err = stmt.Exec(theatre.ID, theatre.Name)
	if err != nil {
		return models.Theatre{}, err
	}

	return theatre, nil
}

func (s *TheatreServiceImpl) UpdateTheatre(id string, theatre models.Theatre) (models.Theatre, error) {
	stmt, err := database.DB.Prepare("UPDATE theatres SET name = ? WHERE id = ?")
	if err != nil {
		return models.Theatre{}, err
	}
	_, err = stmt.Exec(theatre.Name, id)
	if err != nil {
		return models.Theatre{}, err
	}

	theatre.ID = id
	return theatre, nil
}

func (s *TheatreServiceImpl) DeleteTheatre(id string) error {
	stmt, err := database.DB.Prepare("DELETE FROM theatres WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
