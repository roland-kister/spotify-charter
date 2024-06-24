package db

import (
	"context"
	"database/sql"
	"spotify-charter/model"
)

const upsertCountrySql = `INSERT INTO countries (country_code, name, top_playlist_id)
							VALUES (:country_code, :name, :top_playlist_id)
							ON CONFLICT (country_code) DO UPDATE
							SET name = :name, top_playlist_id = :top_playlist_id
							WHERE country_code = :country_code;`

const upsertArtistSql = `INSERT INTO artists (spotify_id, name)
							VALUES(:spotify_id, :name)
							ON CONFLICT (spotify_id) DO UPDATE
							SET name = :name
							WHERE spotify_id = :spotify_id;`

const upsertAlbumSql = `INSERT INTO albums (spotify_id, name)
							VALUES(:spotify_id, :name)
							ON CONFLICT (spotify_id) DO UPDATE
							SET name = :name
							WHERE spotify_id = :spotify_id;`

const upsertImageSql = `INSERT INTO images (album_id, width, url)
							VALUES(:album_id, :width, :url)
							ON CONFLICT (album_id, width) DO UPDATE
							SET url = :url
							WHERE album_id = :album_id AND width = :width;`

const upsertTrackSql = `INSERT INTO tracks (spotify_id, name, album_id)
							VALUES(:spotify_id, :name, :album_id)
							ON CONFLICT (spotify_id) DO UPDATE
							SET name = :name, album_id = :album_id
							WHERE spotify_id = :spotify_id;`

type Writer struct {
	db                *sql.DB
	tx                *sql.Tx
	upsertCountryStmt *sql.Stmt
	upsertArtistStmt  *sql.Stmt
	upsertAlbumStmt   *sql.Stmt
	upsertImageStmt   *sql.Stmt
	upsertTrackStmt   *sql.Stmt
}

func NewWriter(db *sql.DB) *Writer {
	writer := Writer{
		db: db,
	}

	return &writer
}

func (writer *Writer) BeginTx() {
	if writer.tx != nil {
		panic("Trying to create a new transcation, without commiting the existing one")
	}

	if writer.upsertCountryStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert country statement")
	}

	if writer.upsertArtistStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert artist statement")
	}

	if writer.upsertAlbumStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert album statement")
	}

	if writer.upsertImageStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert image statement")
	}

	if writer.upsertTrackStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert track statement")
	}

	var err error

	writer.tx, err = writer.db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	if writer.upsertCountryStmt, err = writer.tx.Prepare(upsertCountrySql); err != nil {
		panic(err)
	}

	if writer.upsertArtistStmt, err = writer.tx.Prepare(upsertArtistSql); err != nil {
		panic(err)
	}

	if writer.upsertAlbumStmt, err = writer.tx.Prepare(upsertAlbumSql); err != nil {
		panic(err)
	}

	if writer.upsertImageStmt, err = writer.tx.Prepare(upsertImageSql); err != nil {
		panic(err)
	}

	if writer.upsertTrackStmt, err = writer.tx.Prepare(upsertTrackSql); err != nil {
		panic(err)
	}
}

func (writer *Writer) CommitTx() {
	if writer.tx == nil {
		panic("No transcation to commit")
	}

	writer.upsertCountryStmt.Close()
	writer.upsertCountryStmt = nil

	writer.upsertArtistStmt.Close()
	writer.upsertArtistStmt = nil

	writer.upsertAlbumStmt.Close()
	writer.upsertAlbumStmt = nil

	writer.upsertImageStmt.Close()
	writer.upsertImageStmt = nil

	writer.upsertTrackStmt.Close()
	writer.upsertTrackStmt = nil

	if err := writer.tx.Commit(); err != nil {
		panic(err)
	}

	writer.tx = nil
}

func (writer *Writer) UpsertCountry(country *model.Country) {
	_, err := writer.upsertCountryStmt.Exec(
		sql.Named("country_code", country.CountryCode),
		sql.Named("name", country.Name),
		sql.Named("top_playlist_id", newNullString(country.TopPlaylistID)))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) UpsertTrack(track *model.Track) {
	for _, artist := range track.Artists {
		writer.upsertArtist(&artist)
	}

	writer.upsertAlbum(&track.Album)

	for _, image := range track.Album.Images {
		writer.upsertImage(&image, track.Album.SpotifyID)
	}

	_, err := writer.upsertTrackStmt.Exec(
		sql.Named("spotify_id", track.SpotifyID),
		sql.Named("name", track.Name),
		sql.Named("album_id", track.Album.SpotifyID))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertArtist(artist *model.Artist) {
	_, err := writer.upsertArtistStmt.Exec(
		sql.Named("spotify_id", artist.SpotifyID),
		sql.Named("name", artist.Name))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertAlbum(album *model.Album) {
	_, err := writer.upsertAlbumStmt.Exec(
		sql.Named("spotify_id", album.SpotifyID),
		sql.Named("name", album.Name))

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) upsertImage(image *model.Image, albumID string) {
	_, err := writer.upsertImageStmt.Exec(
		sql.Named("album_id", albumID),
		sql.Named("width", image.Width),
		sql.Named("url", image.URL))

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
