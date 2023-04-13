package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func GetDB() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual), sqlmock.MonitorPingsOption(true))
	return sqlx.NewDb(db, "mysql"), mock
}

func NewRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}
