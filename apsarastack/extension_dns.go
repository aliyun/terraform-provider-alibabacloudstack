package apsarastack

import (
	"time"
)

type DnsRecord struct {
	AsapiSuccess bool   `json:"asapiSuccess"`
	RequestID    string `json:"RequestId"`
	PageSize     int    `json:"PageSize"`
	PageNumber   int    `json:"PageNumber"`
	TotalItems   int    `json:"TotalItems"`
	Records      []struct {
		DomainID        string    `json:"Id"`
		Rr              string    `json:"Rr"`
		Type            string    `json:"Type"`
		CreateTime      time.Time `json:"CreateTime"`
		RrSet           []string  `json:"RrSet"`
		RecordID        int       `json:"RecordId"`
		UpdateTimestamp int64     `json:"UpdateTimestamp"`
		TTL             int       `json:"Ttl"`
		ZoneID          string    `json:"ZoneId"`
		CreateTimestamp int64     `json:"CreateTimestamp"`
		Remark          string    `json:"Remark,omitempty"`
	} `json:"Records"`
}

type DnsDomains struct {
	AsapiSuccess   bool   `json:"asapiSuccess"`
	AsapiRequestID string `json:"asapiRequestId"`
	PageSize       int    `json:"PageSize"`
	RequestID      string `json:"RequestId"`
	PageNumber     int    `json:"PageNumber"`
	TotalItems     int    `json:"TotalItems"`
	ZoneList       []struct {
		DomainID        int       `json:"DomainId"`
		VpcNumber       int       `json:"VpcNumber"`
		DomainName      string    `json:"DomainName"`
		CreateTime      time.Time `json:"CreateTime"`
		UpdateTime      time.Time `json:"UpdateTime"`
		UpdateTimestamp int64     `json:"UpdateTimestamp"`
		CreateTimestamp int64     `json:"CreateTimestamp"`
		RecordCount     int       `json:"RecordCount"`
		Remark          string    `json:"Remark,omitempty"`
	} `json:"ZoneList"`
}

type DnsDomain struct {
	RequestID string `json:"RequestId"`
	ID        int    `json:"Id"`
}
