package models

import (
	"api/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerList struct {
	CustomerID       	uint64  `db:"id" json:"customer_id"`
	FirstName		 	string  `db:"first_name" json:"first_name"`
	LastName		 	string  `db:"last_name" json:"last_name"`
	Email		 	 	string  `json:"email"`
	Phone		 	 	string  `json:"phone"`
	Address		 	 	string  `json:"address"`
	Category		 	uint64  `json:"category"`
	CompanyName		 	*string  `db:"company_name" json:"company_name"`
	AbnCnNumber		 	*string  `db:"abn_cn_number" json:"abn_cn_number"`
	DrivingLicence	 	string  `db:"driving_licence" json:"driving_licence"`
	PhotoId	 			string  `db:"photo_id" json:"photo_id"`
	Avatar		 		*string  `json:"avatar"`
	Status		 		string  `json:"status"`
}

func CreateCustomer(params map[string]string) int {
	query := "INSERT INTO customers("
	// Get params
	var fields, values string
	i := 0
	for key, value := range params {
		fields += "`" + key + "`"
		values += "'" + value + "'"
		if (len(params) - 1) > i {
			fields += ", "
			values += ", "
		}
		i++
	}
	// Combile params to build query
	query += fields + ") VALUES(" + values + ")"
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