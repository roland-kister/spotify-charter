package db

import "database/sql"

const existsCountrySql = "SELECT 1 FROM countries WHERE country_code = :country_code;"

type Reader struct {
	db                *sql.DB
	existsCountryStmt *sql.Stmt
}

func NewReader(db *sql.DB) *Reader {
	client := Reader{
		db: db,
	}

	var err error

	if client.existsCountryStmt, err = client.db.Prepare(existsCountrySql); err != nil {
		panic(err)
	}

	return &client
}

func (client *Reader) Close() {
	client.db = nil

	client.existsCountryStmt.Close()
	client.existsCountryStmt = nil
}

func (client Reader) ExistsCountry(countryCode string) bool {
	err := client.existsCountryStmt.QueryRow(sql.Named("country_code", countryCode)).Scan(new(int))

	if err == nil {
		return true
	}

	if err == sql.ErrNoRows {
		return false
	}

	panic(err)
}
