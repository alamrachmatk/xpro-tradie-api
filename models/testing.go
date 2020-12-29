package models

import (
	"api/db"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type DomainExample struct {
	DomainID    uint64  `db:"id" json:"site_id"`
	DateTime    *string `db:"datetime" json:"datetime"`
	DateTimeNew *string `db:"datetime_new" json:"dadatetime_newtetime"`
	Domain      *string `db:"domain" json:"domain"`
	DomainNew   *string `db:"domain_new" json:"domain_new"`
	IpAddress   *string `db:"ip_address" json:"ip_address"`
}

type DomainExampleList struct {
	DomainID  uint64 `json:"site_id"`
	DateTime  uint64 `json:"datetime"`
	Domain    string `json:"domain"`
	IpAddress string `json:"ip_address"`
}

func GetAllDomainExample(c *[]DomainExample, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := `SELECT * FROM domainlist_example`
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
			condition += " domainlist_example." + key + " = '" + value + "'"
			clause = true
		}
	}
	// Check order by
	var present bool
	var orderBy, orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY domainlist_example." + orderBy
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
		countQuery := "SELECT COUNT(id) FROM domainlist_example" + condition
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
