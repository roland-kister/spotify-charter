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
	"sync"
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

	dateNow := model.TimeToDatestamp(time.Now())

	reader := db.NewReader(sqlDb)
	countriesWithPlaylist := reader.GetCountriesWithPlaylist()

	wg := new(sync.WaitGroup)

	writer := db.NewWriter(sqlDb)

	for _, country := range countriesWithPlaylist {
		wg.Add(1)

		go getCountry(wg, apiClient, country, writer, dateNow, model.DailyTopTrack)
	}

	wg.Wait()

	writer.Commit()
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

	for len(record) != 0 && err == nil {
		country := model.Country{
			CountryCode:   record[0],
			Name:          record[1],
			TopPlaylistID: record[2],
		}

		log.Printf("Upserting country '%s' ('%s') to the DB\n", country.Name, country.CountryCode)

		writer.SaveCountry(&country)

		record, err = csvReader.Read()
	}

	if err != io.EOF {
		log.Panicln(err)
	}

	log.Println("Successfully finished writing countries from the countries file to the DB")

	writer.Commit()
}

func getCountry(wg *sync.WaitGroup, apiClient *spotify.ApiClient, country *model.Country, writer *db.Writer, date int64, chartType model.ChartType) {
	tracks, err := apiClient.GetPlaylist(country.TopPlaylistID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, track := range tracks {
		chartTrack := &model.ChartTrack{
			Country:   country,
			Track:     track,
			ChartType: chartType,
			Date:      date,
			Position:  index,
		}

		writer.SaveChartTrack(chartTrack)

		log.Printf("[%s] %d: %s\n", country.CountryCode, index, track.Name)
	}

	wg.Done()
}
