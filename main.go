package main

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"spotify-charter/db"
	"spotify-charter/model"
	"spotify-charter/spotify"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	apiClientId := os.Getenv("SPOTIFY_CHARTER_API_CLIENT_ID")
	apiClientSecret := os.Getenv("SPOTIFY_CHARTER_API_CLIENT_SECRET")
	dbFile := os.Getenv("SPOTIFY_CHARTER_DB_FILE")

	sqlDb := initDB(dbFile)

	defer sqlDb.Close()

	initCountries("countries.csv", sqlDb)

	apiClient := spotify.NewApiClient(apiClientId, apiClientSecret)
	if err := apiClient.Authorize(); err != nil {
		log.Panicln(err)
	}

	timeNow := time.Now().UTC()

	dateNow := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location()).Unix()

	reader := db.NewReader(sqlDb)
	countriesWithPlaylist := reader.GetCountriesWithPlaylist()

	tracks, err := apiClient.GetPlaylist((*countriesWithPlaylist)[0].TopPlaylistID)
	if err != nil {
		log.Panicln(err)
	}

	writer := db.NewWriter(sqlDb)
	writer.BeginTx()

	for index, track := range *tracks {
		writer.UpsertTrack(&track)
		writer.UpsertTopTrack((*countriesWithPlaylist)[0].CountryCode, track.SpotifyID, dateNow, index)
	}

	writer.CommitTx()
}

func initDB(dbPath string) *sql.DB {
	log.Printf("Initializing the DB connection with file '%s'", dbPath)

	sqlDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Initializing DB tables")

	db.CreateTables(sqlDb)

	log.Println("Successfully initialized the DB and its tables")

	return sqlDb
}

func initCountries(csvPath string, sqlDB *sql.DB) {
	log.Printf("Writing countries from the countries file '%s' to the DB\n", csvPath)

	countries, err := os.Open(csvPath)
	if err != nil {
		log.Panicln(err)
	}

	defer countries.Close()

	csvReader := csv.NewReader(countries)

	if _, err = csvReader.Read(); err != nil {
		log.Panicln(err)
	}

	record, err := csvReader.Read()

	writer := db.NewWriter(sqlDB)

	writer.BeginTx()

	for len(record) != 0 && err == nil {
		country := model.Country{
			CountryCode:   record[0],
			Name:          record[1],
			TopPlaylistID: record[2],
		}

		log.Printf("Upserting country '%s' ('%s') to the DB\n", country.Name, country.CountryCode)

		writer.UpsertCountry(&country)

		record, err = csvReader.Read()
	}

	if err != io.EOF {
		log.Panicln(err)
	}

	writer.CommitTx()

	log.Println("Successfully finished writing countries from the countries file to the DB")
}
