package models

import (
	"api/db"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type WorkOrder struct {
	WorkOrderID		uint64	`db:"id" json:"worker_id"`
	Status			uint8	`db:"status" json:"status"`
}

type WorkOrderList struct {
	WorkOrderID		uint64	`db:"id" json:"worker_id"`
	Status			string	`db:"status" json:"status"`
}

func GetAllWorkOrder(c *[]WorkOrder, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := `SELECT work_orders.id,
			  work_orders.status 
			  FROM work_orders`
	var condition string
	// Combine where clause
	clause := false
	for key, value := range params {
		if (key != "orderBy") && (key != "orderType") {
			if clause == false {
				condition += " WHERE"
			} else {
				condition += " AND"
			}
			condition += " work_orders." + key + " = '" + value + "'"
			clause = true
		}
	}
	// Check order by
	var present bool
	var orderBy, orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY work_orders." + orderBy
		if orderType, present = params["orderType"]; present == true {
			condition += " " + orderType
		}
	}
	query += condition
	// Query limit and offset
	query += " LIMIT " + strconv.FormatUint(limit, 10)
	if offset > 0 {
		query += " OFFSET " + strconv.FormatUint(offset, 10)
	}

	// Check pagination
	var total uint64
	if pagination == true {
		countQuery := "SELECT COUNT(id) FROM work_orders" + condition
		var totalStr string
		log.Println(countQuery)
		err := db.Db.Get(&totalStr, countQuery)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		total, _ = strconv.ParseUint(totalStr, 10, 64)
	}

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if pagination == false {
		total = uint64(len(*c))
	}

	return total, nil
}