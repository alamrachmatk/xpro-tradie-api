package models

import (
	"api/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	OrderID		uint64 `db:"id" json:"order_id"`
	NewOrderID	uint64 `db:"new_order_id" json:"new_order_id"`
	WorkOrderID	uint64 `db:"work_order_id" json:"work_order_id"`
}

func GetOrder(c *Order, id string) int {
	query := `SELECT
	orders.id,
	orders.new_order_id,
	orders.work_order_id
	FROM orders
	WHERE orders.id = ?`

	log.Println(query)
	err := db.Db.Get(c, query, id)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	return http.StatusOK
}