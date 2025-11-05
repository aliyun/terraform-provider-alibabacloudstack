package alibabacloudstack

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
	EnableSchemaRegistry bool `json:"EnableSchemaRegistry"`
	ExpandMode           bool `json:"ExpandMode"`
}

type DataHubRecordSchema struct {
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	AllowNull bool   `json:"AllowNull"`
	Comment   string `json:"Comment"`
}

type ListTopicResult struct {
	RequestId    string `json:"RequestId"`
	Success      bool   `json:"Success"`
	AsapiSuccess bool   `json:"asapiSuccess"`
	List         struct {
		Topic []struct {
			TopicName  string `json:"TopicName"`
			Comment    string `json:"Comment"`
			ShardCount int    `json:"ShardCount"`
			RecordType string `json:"RecordType"`
			LifeCycle  int    `json:"LifeCycle"`
			CreateTime int64  `json:"CreateTime"`
			UpdateTime int64  `json:"UpdateTime"`
			Storage    int    `json:"Storage"`
		} `json:"Topic"`
	} `json:"List"`
}
