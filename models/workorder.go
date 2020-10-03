package models

import (
	"api/db"
	"log"
	"net/http"
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

type WorkOrderData struct {
	WorkOrderID		uint64	`json:"worker_id"`
	Status			string	`json:"status"`
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

func GetWorkOrder(c *WorkOrder, id string) int {
	query := "SELECT work_orders.id, work_orders.status FROM work_orders WHERE id = " + id
	log.Println(query)
	err := db.Db.Get(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}
	return http.StatusOK
}

func CreateWorkOrder(workorder map[string]string) (int, int64) {
	query := "INSERT INTO work_orders("
	var fields = ""
	var values = ""
	i := 0
	for key, value := range workorder {
		fields += "`" + key + "`"
		values += "'" + value + "'"
		if (len(workorder) - 1) > i {
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