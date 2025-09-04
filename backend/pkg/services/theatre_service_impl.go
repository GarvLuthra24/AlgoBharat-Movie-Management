package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"fmt"
	"log"
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
	// Check if theatre exists
	_, err := s.GetTheatre(id)
	if err != nil {
		return fmt.Errorf("theatre not found: %w", err)
	}

	// Get all halls for this theatre
	hallRows, err := database.DB.Query("SELECT id FROM halls WHERE theatre_id = ?", id)
	if err != nil {
		return fmt.Errorf("error getting halls for theatre: %w", err)
	}
	defer hallRows.Close()

	var hallIDs []string
	for hallRows.Next() {
		var hallID string
		if err := hallRows.Scan(&hallID); err != nil {
			continue
		}
		hallIDs = append(hallIDs, hallID)
	}

	// Get all shows for halls in this theatre
	var showIDs []string
	if len(hallIDs) > 0 {
		// Build the IN clause for hall IDs
		placeholders := make([]string, len(hallIDs))
		args := make([]interface{}, len(hallIDs))
		for i, hallID := range hallIDs {
			placeholders[i] = "?"
			args[i] = hallID
		}

		placeholdersStr := ""
		for i, placeholder := range placeholders {
			if i > 0 {
				placeholdersStr += ","
			}
			placeholdersStr += placeholder
		}

		showRows, err := database.DB.Query(fmt.Sprintf("SELECT id FROM shows WHERE hall_id IN (%s)", placeholdersStr), args...)
		if err != nil {
			log.Printf("Warning: error getting shows for theatre %s: %v", id, err)
		} else {
			defer showRows.Close()
			for showRows.Next() {
				var showID string
				if err := showRows.Scan(&showID); err != nil {
					continue
				}
				showIDs = append(showIDs, showID)
			}
		}
	}

	// Delete all bookings for shows in this theatre
	if len(showIDs) > 0 {
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
			log.Printf("Warning: error deleting bookings for theatre %s: %v", id, err)
		}
	}

	// Delete all shows for halls in this theatre
	if len(hallIDs) > 0 {
		placeholders := make([]string, len(hallIDs))
		args := make([]interface{}, len(hallIDs))
		for i, hallID := range hallIDs {
			placeholders[i] = "?"
			args[i] = hallID
		}

		placeholdersStr := ""
		for i, placeholder := range placeholders {
			if i > 0 {
				placeholdersStr += ","
			}
			placeholdersStr += placeholder
		}

		_, err = database.DB.Exec(fmt.Sprintf("DELETE FROM shows WHERE hall_id IN (%s)", placeholdersStr), args...)
		if err != nil {
			log.Printf("Warning: error deleting shows for theatre %s: %v", id, err)
		}
	}

	// Delete all seats for halls in this theatre
	if len(hallIDs) > 0 {
		placeholders := make([]string, len(hallIDs))
		args := make([]interface{}, len(hallIDs))
		for i, hallID := range hallIDs {
			placeholders[i] = "?"
			args[i] = hallID
		}

		placeholdersStr := ""
		for i, placeholder := range placeholders {
			if i > 0 {
				placeholdersStr += ","
			}
			placeholdersStr += placeholder
		}

		_, err = database.DB.Exec(fmt.Sprintf("DELETE FROM seats WHERE hall_id IN (%s)", placeholdersStr), args...)
		if err != nil {
			log.Printf("Warning: error deleting seats for theatre %s: %v", id, err)
		}
	}

	// Delete all halls for this theatre
	_, err = database.DB.Exec("DELETE FROM halls WHERE theatre_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting halls for theatre: %w", err)
	}

	// Delete the theatre
	_, err = database.DB.Exec("DELETE FROM theatres WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting theatre: %w", err)
	}

	log.Printf("Theatre %s and its %d halls, %d shows deleted successfully", id, len(hallIDs), len(showIDs))
	return nil
}
