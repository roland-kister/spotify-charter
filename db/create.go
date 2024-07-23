package db

import (
	"context"
	"database/sql"
)

const (
	crCountries = iota
	crArtists
	crAlbums
	crImages
	crTracks
	crArtistsTracks
	crChartTracks
)

var createSqls = map[int]string{
	crCountries: `
		CREATE TABLE IF NOT EXISTS countries (
			code TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			top_playlist_id TEXT
		);`,

	crArtists: `
		CREATE TABLE IF NOT EXISTS artists (
			spotify_id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL
		);`,

	crAlbums: `
		CREATE TABLE IF NOT EXISTS albums (
			spotify_id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL
		);`,

	crImages: `
		CREATE TABLE IF NOT EXISTS images (
			album_id TEXT NOT NULL,
			width INTEGER NOT NULL,
			url TEXT NOT NULL,

			PRIMARY KEY(album_id, width),

			FOREIGN KEY(album_id) REFERENCES albums(spotify_id)
		);`,

	crTracks: `
		CREATE TABLE IF NOT EXISTS tracks (
			spotify_id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			album_id TEXT NOT NULL,

			FOREIGN KEY(album_id) REFERENCES albums(spotify_id)
		);`,

	crArtistsTracks: `
		CREATE TABLE IF NOT EXISTS artists_tracks (
			artist_id TEXT NOT NULL,
			track_id TEXT NOT NULL,

			PRIMARY KEY(artist_id, track_id),

			FOREIGN KEY(artist_id) REFERENCES artists(spotify_id),
			FOREIGN KEY(track_id) REFERENCES tracks(spotify_id)
		);`,

	crChartTracks: `
		CREATE TABLE IF NOT EXISTS chart_tracks (
			country_code TEXT NOT NULL,
			track_id TEXT NOT NULL,
			chart_type TEXT NOT NULL,
			date NUMERIC NOT NULL,
			position NUMERIC NOT NULL,

			PRIMARY KEY(country_code, chart_type, date, position),

			FOREIGN KEY(country_code) REFERENCES countries(code),
			FOREIGN KEY(track_id) REFERENCES tracks(spotify_id)
		);`,
}

func CreateTables(db *sql.DB) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	for _, sql := range createSqls {
		if _, err := tx.Exec(sql); err != nil {
			panic(err)
		}
	}

	tx.Commit()
}
