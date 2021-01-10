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
