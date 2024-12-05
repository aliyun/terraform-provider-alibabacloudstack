package datahub

type SubscriptionEntry struct {
	SubId          string            `json:"SubId"`
	TopicName      string            `json:"TopicName"`
	IsOwner        bool              `json:"IsOwner"`
	Type           SubscriptionType  `json:"Type"`
	State          SubscriptionState `json:"State,omitempty"`
	Comment        string            `json:"Comment,omitempty"`
	CreateTime     int64             `json:"CreateTime"`
	LastModifyTime int64             `json:"LastModifyTime"`
}
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
	RequestId       string `json:"RequestId"`
	StorageSetId       string `json:"StorageSetId"`
}

type SubscriptionOffset struct {
	Timestamp int64  `json:"Timestamp"`
	Sequence  int64  `json:"Sequence"`
	VersionId int64  `json:"Version"`
	SessionId *int64 `json:"SessionId"`
}
