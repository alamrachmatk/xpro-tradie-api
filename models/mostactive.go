package models

import (
	"api/db"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type MostActive struct {
	Domain          *string `db:"domain" json:"domain"`
	BaseDomain      *string `db:"base_domain" json:"base_domain"`
	IPAddress       *string `db:"ip_address" json:"ip_address"`
	HasSubdomain    *uint8  `db:"has_subdomain" json:"has_subdomain"`
	TotalMostActive *uint64 `db:"total_most_active" json:"total_most_active"`
	LogDateTime     *string `db:"log_datetime" json:"log_datetime"`
	CreatedAt       *string `db:"created_at" json:"created_at"`
}

type TotalTopMostActiveList struct {
	BaseDomain string `json:"base_domain"`
	Total      uint64 `json:"total"`
}

func GetTotalTopMostActiveListQuery(c *[]MostActive, limit uint64) error {
	query := `SELECT base_domain, count(base_domain) total_most_active from dns GROUP BY base_domain`

	// Query limit
	query += " LIMIT " + strconv.FormatUint(limit, 10)

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// func GetTotalTopMostActiveListQuery(c *[]MostActive, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
// 	query := `SELECT * FROM dns`
// 	var condition string
// 	// Combine where clause
// 	clause := false
// 	for key, value := range params {
// 		if (key != "orderBy") && (key != "orderType") && (key != "groupBy") {
// 			if clause == false {
// 				condition += " WHERE"
// 			} else {
// 				condition += " AND"
// 			}
// 			condition += key + " = '" + value + "'"
// 			clause = true
// 		}
// 	}
// 	// Check order by
// 	var presentGroup bool
// 	var groupBy string
// 	if groupBy, presentGroup = params["groupBy"]; presentGroup == true {
// 		condition += " GROUP BY " + groupBy + ")"
// 		if groupBy, presentGroup = params["groupBy"]; presentGroup == true {
// 			condition += " " + groupBy
// 		}
// 	}

// 	// Check order by
// 	var presentOrder bool
// 	var orderBy, orderType string
// 	if orderBy, presentOrder = params["orderBy"]; presentOrder == true {
// 		condition += " ORDER BY " + orderBy
// 		if orderType, presentOrder = params["orderType"]; presentOrder == true {
// 			condition += " " + orderType
// 		}
// 	}
// 	query += condition
// 	// Query limit and offset
// 	query += " LIMIT " + strconv.FormatUint(limit, 10)
// 	if offset > 0 {
// 		query += " OFFSET " + strconv.FormatUint(offset, 10)
// 	}

// 	// Check pagination
// 	var total uint64
// 	if pagination == true {
// 		countQuery := "SELECT COUNT(id) FROM dns" + condition
// 		var totalStr string
// 		log.Println(countQuery)
// 		err := db.Db.Get(&totalStr, countQuery)
// 		if err != nil {
// 			log.Println(err)
// 			return 0, err
// 		}
// 		total, _ = strconv.ParseUint(totalStr, 10, 64)
// 	}

// 	// Main query
// 	log.Println(query)
// 	err := db.Db.Select(c, query)
// 	if err != nil {
// 		log.Println(err)
// 		return 0, err
// 	}
// 	if pagination == false {
// 		total = uint64(len(*c))
// 	}

// 	return total, nil
// }
