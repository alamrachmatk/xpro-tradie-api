package models

import (
	"api/db"
	"log"
	"strconv"
	"errors"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

type News struct {
	NewsID		 uint64	 `db:"id" json:"news_id"`
	TopicID		 uint64  `db:"topic_id" json:"topic_id"`
	News         string  `db:"name" json:"news"`
	Description  string  `db:"description" json:"description"`
	Status       string  `db:"status" json:"status"`
	DateCreated  string  `db:"date_created" json:"date_created"`
}

type NewsList struct {
	NewsID		 uint64	 `json:"news_id"`
	Topic        string  `json:"topic"`
	News         string  `json:"news"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	DateCreated  string  `json:"date_created"`
}

func GetAllNews(u *[]NewsList, limit uint64, offset uint64, pagination bool, params map[string]string) (uint64, error) {
	query := "SELECT * FROM news"

	var condition string
	// Combine where clause
	clause := false
	for key, value := range params {
		if clause == false {
			condition += " WHERE"
		} else {
			condition += " AND"
		}
		condition += " news." + key + " = '" + value + "'"
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
		countQuery := "SELECT COUNT(id) FROM news" + condition
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

func APIGetNewsIn(u *[]News, id ...string) (int, error) {
	if len(id) == 0 {
		return http.StatusNotFound, errors.New("no supplied IDs")
	}
	query := "SELECT * FROM news WHERE topic_id IN("
	for index, value := range id {
		query += value
		if (len(id) - 1) > index {
			query += ", "
		} else {
			query += ")"
		}
	}
	log.Println(query)
	err := db.Db.Select(u, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return http.StatusOK, nil
}

func CreateNews(params map[string]string) int {
	query := "INSERT INTO news("
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

func UpdateNews(params map[string]string) int {
	query := "UPDATE news SET "
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

func DeleteNews(id string) int {
	query := "UPDATE news SET status = 'deleted'"
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
