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
	client := Writer{
		db: db,
	}

	return &client
}

func (client *Writer) BeginTx() {
	if client.tx != nil {
		panic("Trying to create a new transcation, without commiting the existing one")
	}

	if client.upsertCountryStmt != nil {
		panic("Trying to create a new transaction, without clearing the upsert country statement")
	}

	var err error

	client.tx, err = client.db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	if client.upsertCountryStmt, err = client.tx.Prepare(upsertCountrySql); err != nil {
		panic(err)
	}
}

func (client *Writer) CommitTx() {
	if client.tx == nil {
		panic("No transcation to commit")
	}

	client.upsertCountryStmt.Close()
	client.upsertCountryStmt = nil

	if err := client.tx.Commit(); err != nil {
		panic(err)
	}

	client.tx = nil
}

func (client *Writer) UpsertCountry(country *model.Country) {
	var err error

	_, err = client.upsertCountryStmt.Exec(
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
