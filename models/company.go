package models

import (
	"api/db"
	"errors"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Company struct {
	CompanyID	uint64	`db:"id" json:"company_id"`
	Name		string	`db:"name" json:"name"`
}

func GetCompanyIn(c *[]Company, id ...string) (int, error) {
	if len(id) == 0 {
		return http.StatusNotFound, errors.New("no supplied IDs")
	}
	query := `SELECT 
	companies.id, 
	companies.name 
	FROM companies 
	WHERE id IN(`
	for index, value := range id {
		query += value
		if (len(id) - 1) > index {
			query += ", "
		} else {
			query += ")"
		}
	}
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return http.StatusOK, nil
}