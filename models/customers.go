package models

import (
	"api/db"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Customer struct {
	CustomerID       	uint64  `db:"id" json:"customer_id"`
	FirstName		 	string  `db:"first_name" json:"first_name"`
	LastName		 	string  `db:"last_name" json:"last_name"`
	Email		 	 	string  `db:"email" json:"email"`
	Password		 	string  `db:"password" json:"password"`
	Phone		 	 	string  `db:"phone" json:"phone"`
	Address		 	 	string  `db:"address" json:"address"`
	Category		 	uint64  `db:"category" json:"category"`
	CompanyName		 	*string  `db:"company_name" json:"company_name"`
	AbnCnNumber		 	*string  `db:"abn_cn_number" json:"abn_cn_number"`
	DrivingLicence	 	string  `db:"driving_licence" json:"driving_licence"`
	PhotoId	 			string  `db:"photo_id" json:"photo_id"`
	Avatar		 		*string  `db:"avatar" json:"avatar"`
	Status		 		string  `db:"status" json:"status"`
}

type CustomerData struct {
	CustomerID       	uint64  `json:"customer_id"`
	FirstName		 	string  `json:"first_name"`
	LastName		 	string  `json:"last_name"`
	Email		 	 	string  `json:"email"`
	Password		 	string  `json:"password"`
	Phone		 	 	string  `json:"phone"`
	Address		 	 	string  `json:"address"`
	Category		 	string  `json:"category"`
	CompanyName		 	*string  `json:"company_name"`
	AbnCnNumber		 	*string  `json:"abn_cn_number"`
	DrivingLicence	 	string  `json:"driving_licence"`
	PhotoId	 			string  `json:"photo_id"`
	Avatar		 		*string  `json:"avatar"`
	Status		 		string  `json:"status"`
}

func CreateCustomer(params map[string]string) (int, int64) {
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