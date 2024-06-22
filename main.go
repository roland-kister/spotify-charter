package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"spotify-charter/db"
	"spotify-charter/model"
	"spotify-charter/spotify"

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
		panic(err)
	}

	tracks, err := apiClient.GetPlaylist("37i9dQZEVXbKIVTPX9a2Sb")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range *tracks {
		fmt.Println(track.Name)
	}
}

func initDB(dbPath string) *sql.DB {
	log.Printf("Initializing the DB connection with file '%s'", dbPath)

	sqlDb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	defer countries.Close()

	csvReader := csv.NewReader(countries)

	if _, err = csvReader.Read(); err != nil {
		panic(err)
	}

	record, err := csvReader.Read()

	writer := db.NewWriter(sqlDB)

	writer.BeginTx()
	defer writer.CommitTx()

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
		panic(err)
	}

	log.Println("Successfully finished writing countries from the countries file to the DB")
}
