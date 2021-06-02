package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type TotalRequestList struct {
	Time  string `json:"time"`
	Total uint64 `json:"total"`
}
