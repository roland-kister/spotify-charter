package db

import "database/sql"

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

const createArtistsAlbumsQuery string = `
	CREATE TABLE IF NOT EXISTS artists_tracks (
		artist_id TEXT NOT NULL,
		album_id TEXT NOT NULL,

		FOREIGN KEY(artist_id) REFERENCES artists(spotify_id),
		FOREIGN KEY(album_id) REFERENCES albums(spotify_id)
	)
`

const createImagesQuery string = `
	CREATE TABLE IF NOT EXISTS images (
		url TEXT NOT NULL PRIMARY KEY,
		height INTEGER NOT NULL,
		width INTEGER NOT NULL,
		album_id TEXT NOT NULL,

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
	writer := NewWriter(db)

	writer.BeginTx()

	if _, err := writer.db.Exec(createCountriesQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createArtistsQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createAlbumsQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createArtistsAlbumsQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createImagesQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createTracksQuery); err != nil {
		panic(err)
	}

	if _, err := writer.db.Exec(createTopTracksQuery); err != nil {
		panic(err)
	}

	writer.CommitTx()
}
