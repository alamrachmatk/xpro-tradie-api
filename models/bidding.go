package models

import (
	"api/db"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Bidding struct {
	BiddingID 	uint64 `db:"id" json:"bidding_id"`
	OrderID    	uint64 `db:"order_id"   json:"order_id"`
	CompanyID   uint64 `db:"company_id" json:"company_id"`
	LaborTime   string `db:"labor_time" json:"labor_time"`
	Price 	 	float32 `db:"price" json:"price"`
	Description string `db:"description" json:"description"`
	Status 		string `db:"status" json:"status"`
}

type BiddingList struct {
	BiddingID 	uint64 			`json:"bidding_id"`
	OrderID    	uint64 			`json:"order_id"`
	Company   	CompanyInfo		`json:"company,omitempty"`
	LaborTime   string 			`json:"labor_time"`
	Price 	 	float32 		`json:"price"`
	Description string 			`json:"description"`
	Status 		string 			`json:"status"`
}

type CompanyInfo struct {
	CompanyID	uint64	`json:"company_id"`
	Name		string	`json:"name"`
}

func GetAllBidding(c *[]Bidding, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := `SELECT biddings.id,
			  biddings.order_id,
			  biddings.company_id,
			  biddings.labor_time,
			  biddings.price,
			  biddings.description,
			  biddings.status FROM biddings`
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
			condition += " biddings." + key + " = '" + value + "'"
			clause = true
		}
	}
	// Check order by
	var present bool
	var orderBy, orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY biddings." + orderBy
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
		countQuery := "SELECT COUNT(id) FROM biddings" + condition
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


func GetBidding(c *Bidding, id string) int {
	query := `SELECT
	biddings.id,
	biddings.company_id,
	biddings.labor_time,
	biddings.price,
	biddings.description
	FROM biddings
	WHERE biddings.id = ?`

	log.Println(query)
	err := db.Db.Get(c, query, id)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	return http.StatusOK
}

func UpdateBidding(bidding map[string]string, id string) int {
	query := "UPDATE biddings SET "
	i := 0
	for key, value := range bidding {
		query += "`" + key + "`" + " = '" + strings.Replace(value, "'", "\\'", -1) + "'"
		if (len(bidding) - 1) > i {
			query += ", "
		}
		i++
	}
	query += " WHERE id = " + id
	log.Println(query)
	tx, err := db.Db.Begin()
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway
	}
	_, err = tx.Exec(query)
	tx.Commit()
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest
	}
	return http.StatusOK
}
