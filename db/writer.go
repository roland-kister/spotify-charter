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

type Writer struct {
	db                *sql.DB
	tx                *sql.Tx
	upsertCountryStmt *sql.Stmt
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

	var err error

	writer.tx, err = writer.db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	if writer.upsertCountryStmt, err = writer.tx.Prepare(upsertCountrySql); err != nil {
		panic(err)
	}
}

func (writer *Writer) CommitTx() {
	if writer.tx == nil {
		panic("No transcation to commit")
	}

	writer.upsertCountryStmt.Close()
	writer.upsertCountryStmt = nil

	if err := writer.tx.Commit(); err != nil {
		panic(err)
	}

	writer.tx = nil
}

func (writer *Writer) UpsertCountry(country *model.Country) {
	var err error

	_, err = writer.upsertCountryStmt.Exec(
		sql.Named("country_code", country.CountryCode),
		sql.Named("name", country.Name),
		sql.Named("top_playlist_id", newNullString(country.TopPlaylistID)))

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
