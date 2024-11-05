package database

import (
	"database/sql"
	"fmt"
	"main/domain"
	"strings"
)

type Attractions struct {
	db *sql.DB
}

func NewAttractionsRepo(db *sql.DB) *Attractions {
	return &Attractions{db: db}
}

func (r *Attractions) Find(words []string) ([]domain.Attraction, error) {
	var attractions []domain.Attraction

	for i := range words {
		words[i] = strings.TrimSpace(words[i])

		if strings.Index(words[i], " ") != -1 {
			words[i] = strings.Split(words[i], " ")[0]
		}
	}

	searchQuery := strings.Join(words, " | ")
	query := fmt.Sprintf("SELECT id, name, description, url, image, location, lat, lon FROM attractions WHERE tsv_description @@ to_tsquery('russian', '%s') AND lat != 0 AND lon != 0;", searchQuery)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var attraction domain.Attraction
		if err := rows.Scan(&attraction.Id, &attraction.Name, &attraction.Description, &attraction.Url, &attraction.ImageUrl, &attraction.Location, &attraction.Lat, &attraction.Lon); err != nil {
			return nil, err
		}
		attractions = append(attractions, attraction)
	}

	return attractions, nil
}

func (r *Attractions) GetAll() ([]domain.Attraction, error) {
	var attractions []domain.Attraction

	rows, err := r.db.Query("SELECT id, name, description, url, image, location, lat, lon FROM attractions")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var attraction domain.Attraction
		if err := rows.Scan(&attraction.Id, &attraction.Name, &attraction.Description, &attraction.Url, &attraction.ImageUrl, &attraction.Location, &attraction.Lat, &attraction.Lon); err != nil {
			return nil, err
		}
		attractions = append(attractions, attraction)
	}

	return attractions, nil
}
