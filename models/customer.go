package models

import (
	"api/db"
	"errors"
	"log"
	"net/http"
	"strings"

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

type CustomerDataCache struct {
	CustomerID       	string  `json:"customer_id"`
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

type CustomerData struct {
	CustomerID       	uint64  `json:"customer_id"`
	FirstName		 	string  `json:"first_name"`
	LastName		 	string  `json:"last_name"`
	Email		 	 	string  `json:"email"`
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

func GetCustomer(c *Customer, id string) int {
	query := `SELECT
	customers.id,
	customers.first_name,
	customers.last_name,
	customers.email,
	customers.password,
	customers.phone,
	customers.address,
	customers.category,
	customers.company_name,
	customers.abn_cn_number,
	customers.driving_licence,
	customers.photo_id,
	customers.avatar,
	customers.status
	FROM customers
	WHERE customers.id = ?`

	log.Println(query)
	err := db.Db.Get(c, query, id)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	return http.StatusOK
}

func GetCustomerIn(u *[]Customer, id ...string) (int, error) {
	if len(id) == 0 {
		return http.StatusNotFound, errors.New("no supplied IDs")
	}
	query := `
	SELECT
	customers.id,
	customers.first_name,
	customers.last_name,
	customers.email,
	customers.phone,
	customers.address,
	customers.category,
	customers.company_name,
	customers.abn_cn_number,
	customers.driving_licence,
	customers.photo_id,
	customers.avatar,
	customers.status
	FROM customers
	WHERE customers.id 
	IN(`
	for index, value := range id {
		query += value
		if (len(id) - 1) > index {
			query += ", "
		} else {
			query += ")"
		}
	}
	var customers []Customer
	log.Println(query)
	err := db.Db.Select(&customers, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return http.StatusOK, nil
}

func GetCustomerEmail(c *Customer, params string) int {
	query := `SELECT
	customers.id,
	customers.email,
	customers.password
	FROM customers
	WHERE customers.email = ?`

	log.Println(query)
	err := db.Db.Get(c, query, params)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}
	return http.StatusOK
}

func CreateCustomer(params map[string]string) (int, int64) {
	query := "INSERT INTO customers("
	var fields = ""
	var values = ""
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

func UpdateCustomer(customer map[string]string, id string) int {
	query := "UPDATE customers SET "
	i := 0
	for key, value := range customer {
		query += "`" + key + "`" + " = '" + strings.Replace(value, "'", "\\'", -1) + "'"
		if (len(customer) - 1) > i {
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