package models

import (
	"api/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type NewOrder struct {
	NewOrderID			uint64	`db:"id" json:"new_order_id"`
	CustomerID			uint64	`db:"customer_id" json:"customer_id"`
	CompanySettingID	uint64	`db:"customer_id" json:"customer_id"`
}

func GetNewOrder(c *NewOrder, id string) int {
	query := `SELECT
	new_orders.id,
	new_orders.name,
	new_orders.customer_id,
	new_orders.company_setting_id,
	new_orders.due_date
	FROM new_orders
	WHERE new_orders.id = ?`

	log.Println(query)
	err := db.Db.Get(c, query, id)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	return http.StatusOK
}

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