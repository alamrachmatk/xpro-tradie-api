package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type TotalBlok struct {
	BlokAll uint64 `json:"blok_all"`
}

type TotalDnsBlok struct {
	DnsBlok uint64 `json:"dns_blok_all"`
}
