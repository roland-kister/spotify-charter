package db

import (
	"context"
	"database/sql"
)

const createCountriesQuery string = `
	CREATE TABLE IF NOT EXISTS countries (
		country_code TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		top_playlist_id TEXT
	);
`

const createArtistsQuery string = `
	CREATE TABLE IF NOT EXISTS artists (
		spotify_id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL
	);
`

const createAlbumsQuery string = `
	CREATE TABLE IF NOT EXISTS albums (
		spotify_id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL
	);
`

const createImagesQuery string = `
	CREATE TABLE IF NOT EXISTS images (
		album_id TEXT NOT NULL,
		width INTEGER NOT NULL,
		url TEXT NOT NULL,

		PRIMARY KEY(album_id, width),

		FOREIGN KEY(album_id) REFERENCES albums(spotify_id)
	);
`
const createTracksQuery string = `
	CREATE TABLE IF NOT EXISTS tracks (
		spotify_id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		album_id TEXT NOT NULL,

		FOREIGN KEY(album_id) REFERENCES albums(spotify_id)
	);
`

const createArtistsTracksQuery string = `
	CREATE TABLE IF NOT EXISTS artists_tracks (
		artist_id TEXT NOT NULL,
		track_id TEXT NOT NULL,

		PRIMARY KEY(artist_id, track_id),

		FOREIGN KEY(artist_id) REFERENCES artists(spotify_id),
		FOREIGN KEY(track_id) REFERENCES tracks(spotify_id)
	)
`

const createTopTracksQuery string = `
	CREATE TABLE IF NOT EXISTS top_tracks (
		country_code TEXT NOT NULL,
		track_id TEXT NOT NULL,
		date NUMERIC NOT NULL,
		position NUMERIC NOT NULL,

		FOREIGN KEY(country_code) REFERENCES countries(country_code),
		FOREIGN KEY(track_id) REFERENCES tracks(spotify_id)
	);
`

func CreateTables(db *sql.DB) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createCountriesQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createArtistsQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createAlbumsQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createImagesQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createTracksQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createArtistsTracksQuery); err != nil {
		panic(err)
	}

	if _, err := tx.Exec(createTopTracksQuery); err != nil {
		panic(err)
	}

	tx.Commit()
}
