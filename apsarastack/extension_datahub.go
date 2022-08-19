package apsarastack

type GetTopicResult struct {
	RequestId      string
	Success        bool
	AsapiSuccess   bool   `json:"asapiSuccess"`
	ProjectName    string `json:"ProjectName"`
	TopicName      string `json:"TopicName"`
	ShardCount     int    `json:"ShardCount"`
	LifeCycle      int    `json:"LifeCycle"`
	RecordType     string `json:"RecordType"`
	RecordSchema   string `json:"RecordSchema"`
	Comment        string `json:"Comment"`
	CreateTime     int64  `json:"CreateTime"`
	LastModifyTime int64  `json:"LastModifyTime"`
	Storage        int
	//TopicStatus    TopicStatus   `json:"Status"`
	//ExpandMode     ExpandMode    `json:"ExpandMode"`
}
