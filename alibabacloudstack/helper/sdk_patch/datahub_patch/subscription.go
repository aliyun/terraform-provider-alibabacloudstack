package datahub_patch

type SubscriptionCreate struct {
	SubId           string `json:"SubscriptionId"`
	TopicName       string `json:"TopicName"`
	eagleEyeTraceId string `json:"eagleEyeTraceId"`
	asapiSuccess    bool   `json:"asapiSuccess"`
	serverRole      string `json:"serverRole"`
	asapiRequestId  string `json:"asapiRequestId"`
	RequestId       string `json:"RequestId"`
	domain          string `json:"domain"`
	api             string `json:"api"`
	Success         bool   `json:"Success"`
	ProjectName     string `json:"ProjectName"`
}
type EcsCreate struct {
	CommandId       string `json:"CommandId"`
	TopicName       string `json:"TopicName"`
	eagleEyeTraceId string `json:"eagleEyeTraceId"`
	asapiSuccess    bool   `json:"asapiSuccess"`
	serverRole      string `json:"serverRole"`
	asapiRequestId  string `json:"asapiRequestId"`
	RequestId       string `json:"RequestId"`
	domain          string `json:"domain"`
	api             string `json:"api"`
	Success         bool   `json:"Success"`
	ProjectName     string `json:"ProjectName"`
}

type EcsStorageSetsCreate struct {
	RequestId    string `json:"RequestId"`
	StorageSetId string `json:"StorageSetId"`
}
