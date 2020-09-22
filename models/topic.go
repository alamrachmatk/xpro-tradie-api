package models

import (
	"api/db"
	"log"
	"strconv"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

type Topic struct {
	TopicID		 uint64	 `db:"id" json:"topic_id"`
	Name         string  `db:"name" json:"name"`
	Status       string  `db:"status" json:"status"`
	DateCreated  string  `db:"date_created" json:"date_created"`
}

func GetAllTopics(u *[]Topic, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := "SELECT * FROM topics"

	var condition string
	// Combine where clause
	clause := false
	for key, value := range params {
		if clause == false {
			condition += " WHERE"
		} else {
			condition += " AND"
		}
		condition += " topics." + key + " = '" + value + "'"
		clause = true
	}

	query += condition

	if limit > 0 {
		query += " LIMIT " + strconv.FormatUint(limit, 10)
	}
	if offset > 0 {
		query += " OFFSET " + strconv.FormatUint(offset, 10)
	}

	// Check pagination
	var total uint64
	if pagination == true {
		countQuery := "SELECT COUNT(id) FROM topics" + condition
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
	err := db.Db.Select(u, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if pagination == false {
		total = uint64(len(*u))
	}

	return total, nil
}

func CreateTopic(params map[string]string) int {
	query := "INSERT INTO topics("
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
	query += fields + ", date_created) VALUES(" + values + ", NOW())"
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

func UpdateTopic(params map[string]string) int {
	query := "UPDATE topics SET "
	// Get params
	i := 0
	for key, value := range params {
		if key != "id" {
			query += key + " = '" + value + "'"
			if (len(params) - 2) > i {
				query += ", "
			}
			i++
		}
	}
	query += " WHERE id = " + params["id"]
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

func DeleteTopic(id string) int {
	query := "UPDATE topics SET status = 'deleted'"
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
