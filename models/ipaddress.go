package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type TotalIpAddress struct {
	IpAddressAll uint64 `json:"ip_address_all"`
}

type TotalIpAddressDayList struct {
	DayName string `json:"day_name"`
	Total   uint64 `json:"total"`
}
