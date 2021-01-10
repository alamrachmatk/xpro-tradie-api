package models

import (
	"api/db"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type TotalDns struct {
	DnsAll uint64 `json:"dns_all"`
}

type Dns struct {
	DnsID        uint64 `db:"id" json:"dns_id"`
	Domain       string `db:"domain" json:"domain"`
	BaseDomain   string `db:"base_domain" json:"base_domain"`
	IpAddress    string `db:"ip_address" json:"ip_address"`
	HasSubdomain uint8  `db:"has_subdomain" json:"has_subdomain"`
	LogDatetime  string `db:"log_datetime" json:"log_datetime"`
	CreatedAt    string `db:"created_at" json:"created_at"`
}

type DnsList struct {
	DnsID        uint64 `json:"dns_id"`
	Domain       string `json:"domain"`
	BaseDomain   string `json:"base_domain"`
	IpAddress    string `json:"ip_address"`
	HasSubdomain string `json:"has_subdomain"`
	LogDatetime  string `json:"log_datetime"`
	CreatedAt    string `json:"created_at"`
}

func GetAllDns(c *[]Dns, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := `SELECT * FROM dns`
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
			condition += " dns." + key + " = '" + value + "'"
			clause = true
		}
	}
	// Check order by
	var present bool
	var orderBy, orderType string
	if orderBy, present = params["orderBy"]; present == true {
		condition += " ORDER BY dns." + orderBy
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
		countQuery := "SELECT COUNT(id) FROM dns" + condition
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
