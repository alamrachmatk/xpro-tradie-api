package models

import (
	"api/db"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Site struct {
	SiteID    uint64  `db:"id" json:"site_id"`
	Ip        string  `db:"ip" json:"ip"`
	SiteName  string  `db:"name" json:"site_name"`
	Status    uint8   `db:"status" json:"status"`
	Deleted   uint8   `db:"deleted" json:"deleted"`
	CreatedAt *string `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
	CreatedBy *string `db:"created_by" json:"created_by"`
	UpdatedBy *string `db:"updated_by" json:"updated_by"`
}

type SiteList struct {
	SiteID   uint64 `json:"site_id"`
	Ip       string `json:"ip"`
	SiteName string `json:"sitename"`
	Status   string `json:"status"`
}

func GetTotalSite(params map[string]string) (uint64, error) {
	var total uint64
	var totalStr string

	err := db.Db.Get(&totalStr, "SELECT COUNT(id) FROM dns")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	total, _ = strconv.ParseUint(totalStr, 10, 64)

	return total, nil
}

func GetAllSite(c *[]Site, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := `SELECT * FROM site`
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
			condition += " site." + key + " = '" + value + "'"
			clause = true
		}
	}
	// Check order by
	var present bool
	var orderBy, orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY site." + orderBy
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
		countQuery := "SELECT COUNT(id) FROM site" + condition
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
