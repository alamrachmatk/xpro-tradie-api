package db

import (
	"api/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

func MariaDBInit() *sqlx.DB {
	db := sqlx.MustConnect("mysql", config.MariaDBUser+":"+config.MariaDBPassword+"@tcp("+config.MariaDBHost+":"+config.MariaDBPort+")/"+config.MariaDBDB)
	return db
}
