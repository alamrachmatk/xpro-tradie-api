package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type TotalBlock struct {
	BlockAll uint64 `json:"block_all"`
}

type TotalDnsBlock struct {
	DnsBlock uint64 `json:"dns_block_all"`
}

type TotalIpAddressBlockCategoryDayList struct {
	CategoryName string `json:"category_name"`
	Total        uint64 `json:"total"`
}

type TotalDnsBlockCategoryDayList struct {
	CategoryName string `json:"category_name"`
	Total        uint64 `json:"total"`
}
