package db

import (
	"database/sql"
	"spotify-charter/model"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	reader := Reader{
		db: db,
	}

	return &reader
}

func (reader *Reader) Close() {
	reader.db = nil
}

func (reader Reader) GetCountriesWithPlaylist() *[]model.Country {
	rows, err := reader.db.Query("SELECT country_code, name, top_playlist_id FROM countries WHERE top_playlist_id IS NOT NULL;")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	countries := make([]model.Country, 0)

	for rows.Next() {
		country := model.Country{}
		if err := rows.Scan(&country.CountryCode, &country.Name, &country.TopPlaylistID); err != nil {
			panic(err)
		}

		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return &countries
}
