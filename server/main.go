package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"spotify-charter/db"
	"spotify-charter/model"
	"spotify-charter/server"
	"spotify-charter/spotify"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	apiClientID := os.Getenv("SPOTIFY_CHARTER_API_CLIENT_ID")
	apiClientSecret := os.Getenv("SPOTIFY_CHARTER_API_CLIENT_SECRET")
	dbFile := os.Getenv("SPOTIFY_CHARTER_DB_FILE")

	sqlDB := initDB(dbFile)

	defer sqlDB.Close()

	initCountries("countries.csv", sqlDB)

	apiClient := spotify.NewAPIClient(apiClientID, apiClientSecret)
	if err := apiClient.Authorize(); err != nil {
		log.Panicln(err)
	}

	dateNow := model.TimeToDatestamp(time.Now())

	reader := db.NewReader(sqlDB)
	countriesWithPlaylist := reader.GetCountriesWithPlaylist()

	wg := new(sync.WaitGroup)

	writer := db.NewWriter(sqlDB)

	for _, country := range countriesWithPlaylist {
		wg.Add(1)

		go getCountry(wg, apiClient, country, writer, dateNow, model.DailyTopTrack)
	}

	wg.Wait()

	writer.Commit()

	chartTracks := reader.GetChartTracksExt(dateNow)
	for countryCode := range *chartTracks {
		if countryCode != "SK" {
			continue
		}

		fmt.Println((*chartTracks)[countryCode][0].Name)

		for _, artist := range (*chartTracks)[countryCode][0].Artists {
			fmt.Println("\t", artist.Name)
		}

		for _, image := range (*chartTracks)[countryCode][0].Album.Images {
			fmt.Println("\t\t", image.URL)
		}
	}

	server := server.Server{
		Reader: reader,
	}

	fmt.Println("listening on: http://localhost:8080/test")

	http.HandleFunc("/test", server.GetPlaylists)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB(dbPath string) *sql.DB {
	log.Printf("Initializing the DB connection with file '%s'", dbPath)

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Initializing DB tables")

	db.CreateTables(sqlDB)

	log.Println("Successfully initialized the DB and its tables")

	return sqlDB
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
			Code:          record[0],
			Name:          record[1],
			TopPlaylistID: record[2],
		}

		log.Printf("Upserting country '%s' ('%s') to the DB\n", country.Name, country.Code)

		writer.SaveCountry(&country)

		record, err = csvReader.Read()
	}

	if err != io.EOF {
		log.Panicln(err)
	}

	log.Println("Successfully finished writing countries from the countries file to the DB")

	writer.Commit()
}

func getCountry(wg *sync.WaitGroup, apiClient *spotify.APICLient, country *model.Country, writer *db.Writer, date int64, chartType model.ChartType) {
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

		log.Printf("[%s] %d: %s\n", country.Code, index, track.Name)
	}

	wg.Done()
}
