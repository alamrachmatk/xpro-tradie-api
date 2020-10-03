package models

import (
	"api/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func CreateNewOrder(params map[string]string) (int, int64) {
	query := "INSERT INTO new_orders("
	var fields = ""
	var values = ""
	i := 0
	for key, value := range params {
		fields += "`" + key + "`"
		values += "'" + value + "'"
		if (len(params) - 1) > i {
			fields += ", "
			values += ", "
		}
		i++
	}

	query += fields + ", created_at) VALUES(" + values + ", NOW())"
	tx, err := db.Db.Begin()
	var lastID int64
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, lastID
	}
	result, err := tx.Exec(query)
	if err != nil {
		log.Println(err)
	}
	lastID, err = result.LastInsertId()
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, lastID
	}
	return http.StatusOK, lastID
}