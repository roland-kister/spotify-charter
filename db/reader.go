package db

import "database/sql"

const existsCountrySql = "SELECT 1 FROM countries WHERE country_code = :country_code;"

type Reader struct {
	db                *sql.DB
	existsCountryStmt *sql.Stmt
}

func NewReader(db *sql.DB) *Reader {
	reader := Reader{
		db: db,
	}

	var err error

	if reader.existsCountryStmt, err = reader.db.Prepare(existsCountrySql); err != nil {
		panic(err)
	}

	return &reader
}

func (reader *Reader) Close() {
	reader.db = nil

	reader.existsCountryStmt.Close()
	reader.existsCountryStmt = nil
}

func (reader Reader) ExistsCountry(countryCode string) bool {
	err := reader.existsCountryStmt.QueryRow(sql.Named("country_code", countryCode)).Scan(new(int))

	if err == nil {
		return true
	}

	if err == sql.ErrNoRows {
		return false
	}

	panic(err)
}
