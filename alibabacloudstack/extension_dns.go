package alibabacloudstack

import (
	"time"
)

type DnsRecord struct {
	AsapiSuccess bool   `json:"asapiSuccess"`
	RequestID    string `json:"RequestId"`
	PageSize     int    `json:"PageSize"`
	PageNumber   int    `json:"PageNumber"`
	TotalItems   int    `json:"TotalItems"`
	Data         []struct {
		ZoneId     int       `json:"ZoneId"`
		Name       string    `json:"Name"`
		Type       string    `json:"Type"`
		CreateTime time.Time `json:"CreateTime"`
		//RDatas          string    `json:"RDatas"`
		Id              int    `json:"Id"`
		UpdateTimestamp int64  `json:"UpdateTimestamp"`
		TTL             int    `json:"Ttl"`
		CreateTimestamp int64  `json:"CreateTimestamp"`
		Remark          string `json:"Remark"`
		LbaStrategy     string `json:"LbaStrategy"`
	} `json:"Data"`
}

type DnsDomains struct {
	AsapiSuccess    bool   `json:"asapiSuccess"`
	AsapiRequestID  string `json:"asapiRequestId"`
	PageSize        int    `json:"PageSize"`
	RequestID       string `json:"RequestId"`
	PageNumber      int    `json:"PageNumber"`
	TotalItems      int    `json:"TotalItems"`
	EagleEyeTraceId string `json:"eagleEyeTraceId"`
	Data            []struct {
		Id              string    `json:"Id"`
		VpcNumber       int       `json:"VpcNumber"`
		Name            string    `json:"Name"`
		CreateTime      time.Time `json:"CreateTime"`
		UpdateTime      time.Time `json:"UpdateTime"`
		UpdateTimestamp int64     `json:"UpdateTimestamp"`
		CreateTimestamp int64     `json:"CreateTimestamp"`
		RecordCount     int       `json:"RecordCount"`
		Remark          string    `json:"Remark,omitempty"`
	} `json:"Data"`
}

type DnsDomain struct {
	RequestID string `json:"RequestId"`
	ID        int    `json:"Id"`
}
