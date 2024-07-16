package db

import (
	"context"
	"database/sql"
	"spotify-charter/model"
)

const (
	upsCountry = iota
	upsArtist
	upsAlbum
	upsImage
	upsTrack
	upsArtistTrack
	upsChartTrack
)

var writerSqls = map[int]string{
	upsCountry: `
		INSERT INTO countries (country_code, name, top_playlist_id)
			VALUES (:country_code, :name, :top_playlist_id)
		ON CONFLICT (country_code) DO UPDATE
			SET name = :name, top_playlist_id = :top_playlist_id
		WHERE country_code = :country_code;`,

	upsArtist: `
		INSERT INTO artists (spotify_id, name)
			VALUES(:spotify_id, :name)
		ON CONFLICT (spotify_id) DO UPDATE
			SET name = :name
		WHERE spotify_id = :spotify_id;`,

	upsAlbum: `
		INSERT INTO albums (spotify_id, name)
			VALUES(:spotify_id, :name)
		ON CONFLICT (spotify_id) DO UPDATE
			SET name = :name
		WHERE spotify_id = :spotify_id;`,

	upsImage: `
		INSERT INTO images (album_id, width, url)
			VALUES(:album_id, :width, :url)
		ON CONFLICT (album_id, width) DO UPDATE
			SET url = :url
		WHERE album_id = :album_id AND width = :width;`,

	upsTrack: `
		INSERT INTO tracks (spotify_id, name, album_id)
			VALUES(:spotify_id, :name, :album_id)
		ON CONFLICT (spotify_id) DO UPDATE
			SET name = :name, album_id = :album_id
		WHERE spotify_id = :spotify_id;`,

	upsArtistTrack: `
		INSERT INTO artists_tracks (artist_id, track_id)
			VALUES(:artist_id, :track_id)
		ON CONFLICT (artist_id, track_id) DO NOTHING;`,

	upsChartTrack: `
		INSERT INTO chart_tracks (country_code, track_id, chart_type, date, position)
			VALUES(:country_code, :track_id, :chart_type, :date, :position)
		ON CONFLICT (country_code, chart_type, date, position) DO UPDATE
			SET track_id = :track_id
		WHERE country_code = :country_code AND chart_type = :chart_type AND date = :date AND position = :position;`,
}

type Writer struct {
	db               *sql.DB
	tx               *sql.Tx
	stmts            map[int]*sql.Stmt
	done             chan bool
	countryToSave    chan *model.Country
	chartTrackToSave chan *model.ChartTrack
}

func NewWriter(db *sql.DB) *Writer {
	var err error

	writer := &Writer{
		db:               db,
		stmts:            make(map[int]*sql.Stmt),
		done:             make(chan bool),
		countryToSave:    make(chan *model.Country),
		chartTrackToSave: make(chan *model.ChartTrack),
	}

	if writer.tx, err = writer.db.BeginTx(context.Background(), nil); err != nil {
		panic(err)
	}

	for index, sql := range writerSqls {
		if writer.stmts[index], err = writer.tx.Prepare(sql); err != nil {
			panic(err)
		}
	}

	go writer.writingRoutine()

	return writer
}

func (writer *Writer) writingRoutine() {
	for {
		select {
		case country := <-writer.countryToSave:
			writer.upsertCountry(country)
		case chartTrack := <-writer.chartTrackToSave:
			writer.upsertChartTrack(chartTrack)
		case <-writer.done:
			return
		}
	}
}

func (writer *Writer) Commit() {
	var err error

	writer.done <- true

	close(writer.countryToSave)
	close(writer.chartTrackToSave)
	close(writer.done)

	for index := range writerSqls {
		if err = writer.stmts[index].Close(); err != nil {
			panic(err)
		}
	}

	if err = writer.tx.Commit(); err != nil {
		panic(err)
	}

	writer.tx = nil
}

func (writer *Writer) SaveCountry(country *model.Country) {
	writer.countryToSave <- country
}

func (writer *Writer) SaveChartTrack(chartTrack *model.ChartTrack) {
	writer.chartTrackToSave <- chartTrack
}

func (writer *Writer) upsertCountry(country *model.Country) {
	_, err := writer.stmts[upsCountry].Exec(
		sql.Named("country_code", country.CountryCode),
		sql.Named("name", country.Name),
		sql.Named("top_playlist_id", newNullString(country.TopPlaylistID)))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertChartTrack(chartTrack *model.ChartTrack) {
	writer.upsertTrack(chartTrack.Track)

	_, err := writer.stmts[upsChartTrack].Exec(
		sql.Named("country_code", chartTrack.Country.CountryCode),
		sql.Named("track_id", chartTrack.Track.SpotifyID),
		sql.Named("chart_type", chartTrack.ChartType),
		sql.Named("date", chartTrack.Date),
		sql.Named("position", chartTrack.Position))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertTrack(track *model.Track) {
	for _, artist := range track.Artists {
		writer.upsertArtist(&artist)
	}

	writer.upsertAlbum(&track.Album)

	for _, image := range track.Album.Images {
		writer.upsertImage(&image, track.Album.SpotifyID)
	}

	_, err := writer.stmts[upsTrack].Exec(
		sql.Named("spotify_id", track.SpotifyID),
		sql.Named("name", track.Name),
		sql.Named("album_id", track.Album.SpotifyID))

	if err != nil {
		panic(err)
	}

	for _, artist := range track.Artists {
		writer.upsertArtistTrack(artist.SpotifyID, track.SpotifyID)
	}
}

func (writer *Writer) upsertArtist(artist *model.Artist) {
	_, err := writer.stmts[upsArtist].Exec(
		sql.Named("spotify_id", artist.SpotifyID),
		sql.Named("name", artist.Name))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertAlbum(album *model.Album) {
	_, err := writer.stmts[upsAlbum].Exec(
		sql.Named("spotify_id", album.SpotifyID),
		sql.Named("name", album.Name))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertImage(image *model.Image, albumID string) {
	_, err := writer.stmts[upsImage].Exec(
		sql.Named("album_id", albumID),
		sql.Named("width", image.Width),
		sql.Named("url", image.URL))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertArtistTrack(artistID string, trackID string) {
	_, err := writer.stmts[upsArtistTrack].Exec(
		sql.Named("artist_id", artistID),
		sql.Named("track_id", trackID))

	if err != nil {
		panic(err)
	}
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
