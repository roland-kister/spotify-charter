package db

import (
	"database/sql"
	"spotify-charter/model"
)

const (
	selPlaylistCountries = iota
	selChartTracks
	selArtistsByTrack
)

var readerSqls = map[int]string{
	selPlaylistCountries: `
		SELECT country_code, name, top_playlist_id
			FROM countries
		WHERE top_playlist_id IS NOT NULL;`,

	selChartTracks: `
		SELECT ct.country_code, ct.position, ct.track_id, t.name AS track_name, t.album_id, a.name AS album_name
			FROM chart_tracks ct
			RIGHT JOIN tracks t ON t.spotify_id = ct.track_id 
			RIGHT JOIN albums a ON a.spotify_id = t.album_id 
		WHERE ct.chart_type = :chart_type AND ct.date = :date;`,

	selArtistsByTrack: `
		SELECT at.artist_id, a.name 
			FROM artists_tracks at
			INNER JOIN artists a ON a.spotify_id = at.artist_id 
		WHERE at.track_id = :track_id;`,
}

type Reader struct {
	db    *sql.DB
	stmts map[int]*sql.Stmt
}

func NewReader(db *sql.DB) *Reader {
	var err error

	reader := &Reader{
		db:    db,
		stmts: make(map[int]*sql.Stmt),
	}

	for index, sql := range readerSqls {
		if reader.stmts[index], err = reader.db.Prepare(sql); err != nil {
			panic(err)
		}
	}

	return reader
}

func (reader *Reader) Close() {
	reader.db = nil

	for index := range readerSqls {
		if err := reader.stmts[index].Close(); err != nil {
			panic(err)
		}
	}

}

func (reader *Reader) GetCountriesWithPlaylist() []*model.Country {
	rows, err := reader.stmts[selPlaylistCountries].Query()
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	countries := make([]*model.Country, 0)

	for rows.Next() {
		country := model.Country{}
		if err := rows.Scan(&country.CountryCode, &country.Name, &country.TopPlaylistID); err != nil {
			panic(err)
		}

		countries = append(countries, &country)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return countries
}

func (reader *Reader) GetChartTracks(date int64) []*model.ChartTrack {
	rows, err := reader.stmts[selChartTracks].Query(
		sql.Named("chart_type", model.DailyTopTrack),
		sql.Named("date", date))

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	chartTracks := make([]*model.ChartTrack, 0)

	for rows.Next() {
		chartTrack := model.ChartTrack{
			Country:   &model.Country{},
			Track:     &model.Track{},
			ChartType: model.DailyTopTrack,
			Date:      date,
		}

		if err := rows.Scan(&chartTrack.Country.CountryCode, &chartTrack.Position, &chartTrack.Track.SpotifyID,
			&chartTrack.Track.Name, &chartTrack.Track.Album.SpotifyID, &chartTrack.Track.Album.Name); err != nil {
			panic(err)
		}

		chartTrack.Track.Artists = reader.getArtistsForTrack(chartTrack.Track.SpotifyID)

		chartTracks = append(chartTracks, &chartTrack)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return chartTracks
}

func (reader *Reader) getArtistsForTrack(trackID string) []model.Artist {
	rows, err := reader.stmts[selArtistsByTrack].Query(sql.Named("track_id", trackID))
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	artists := make([]model.Artist, 0)

	for rows.Next() {
		artist := model.Artist{}

		if err := rows.Scan(&artist.SpotifyID, &artist.Name); err != nil {
			panic(err)
		}

		artists = append(artists, artist)
	}

	return artists
}
