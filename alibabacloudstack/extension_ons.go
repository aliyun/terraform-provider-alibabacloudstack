package alibabacloudstack

type Topic struct {
	Data []struct {
		CreateTime        int64  `json:"createTime"`
		UnitFlag          bool   `json:"unitFlag"`
		Remark            string `json:"remark"`
		StatusName        string `json:"statusName"`
		Applier           string `json:"applier"`
		ID                int    `json:"id"`
		RelationName      string `json:"relationName"`
		IndependentNaming bool   `json:"independentNaming"`
		RegionID          string `json:"regionId"`
		Topic             string `json:"topic"`
		ChannelName       string `json:"channelName"`
		NamespaceID       string `json:"namespaceId"`
		UpdateTime        int64  `json:"updateTime"`
		Status            int    `json:"status"`
		ChannelID         int    `json:"channelId"`
		OrderType         int    `json:"orderType"`
		Relation          int    `json:"relation"`
		RegionName        string `json:"regionName"`
		PageStart         int    `json:"pageStart"`
		PageSize          int    `json:"pageSize"`
		Owner             string `json:"owner"`
	} `json:"Data"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      int    `json:"Code"`
	Total     int    `json:"Total"`
}

type OnsInstance struct {
	Data struct {
		Cluster            string `json:"cluster"`
		InstanceName       string `json:"instanceName"`
		NamespaceRulesType bool   `json:"namespaceRulesType"`
		TpsReceiveMax      int    `json:"tpsReceiveMax"`
		InstanceType       int    `json:"instanceType"`
		IndependentNaming  bool   `json:"independentNaming"`
		InstanceStatus     int    `json:"instanceStatus"`
		TopicCapacity      int    `json:"topicCapacity"`
		Department         int    `json:"Department"`
		InstanceID         string `json:"instanceId"`
		CreateTime         int64  `json:"createTime"`
		RegionID           string `json:"regionId"`
		DepartmentName     string `json:"DepartmentName"`
		TpsMax             int    `json:"tpsMax"`
		Remark             string `json:"remark"`
		SpInstanceID       string `json:"spInstanceId"`
		ResourceGroup      int    `json:"ResourceGroup"`
		ResourceGroupName  string `json:"ResourceGroupName"`
	} `json:"Data"`
	PageNumber   int    `json:"PageNumber"`
	PageSize     int    `json:"PageSize"`
	Total        int    `json:"Total"`
	Code         int    `json:"code"`
	Cost         int    `json:"cost"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}
type OInstance struct {
	Data []struct {
		Cluster            string `json:"cluster"`
		InstanceName       string `json:"instanceName"`
		NamespaceRulesType bool   `json:"namespaceRulesType"`
		TpsReceiveMax      int    `json:"tpsReceiveMax"`
		InstanceType       int    `json:"instanceType"`
		IndependentNaming  bool   `json:"independentNaming"`
		InstanceStatus     int    `json:"instanceStatus"`
		TopicCapacity      int    `json:"topicCapacity"`
		Department         int    `json:"Department"`
		InstanceID         string `json:"instanceId"`
		CreateTime         int64  `json:"createTime"`
		RegionID           string `json:"regionId"`
		DepartmentName     string `json:"DepartmentName"`
		TpsMax             int    `json:"tpsMax"`
		SpInstanceID       string `json:"spInstanceId"`
		ResourceGroup      int    `json:"ResourceGroup"`
		ResourceGroupName  string `json:"ResourceGroupName"`
	} `json:"Data"`
	PageNumber   int    `json:"PageNumber"`
	PageSize     int    `json:"PageSize"`
	Total        int    `json:"Total"`
	Code         string `json:"code"`
	Cost         int    `json:"cost"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

type OGroup struct {
	Data struct {
		CreateTime        int64  `json:"createTime"`
		UnitFlag          bool   `json:"unitFlag"`
		Remark            string `json:"remark"`
		StatusName        string `json:"statusName"`
		Applier           string `json:"applier"`
		ID                int    `json:"id"`
		RelationName      string `json:"relationName"`
		IndependentNaming bool   `json:"independentNaming"`
		RegionID          string `json:"regionId"`
		GroupType         int    `json:"groupType"`
		ChannelName       string `json:"channelName"`
		NamespaceID       string `json:"namespaceId"`
		Status            int    `json:"status"`
		ChannelID         int    `json:"channelId"`
		UpdateTime        int64  `json:"updateTime"`
		Relation          int    `json:"relation"`
		RegionName        string `json:"regionName"`
		ConsumerID        string `json:"consumerId"`
		GroupID           string `json:"groupId"`
		Owner             string `json:"owner"`
	} `json:"Data"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      int    `json:"Code"`
	Total     int    `json:"Total"`
}
type OnsGroup struct {
	Data []struct {
		CreateTime        int64  `json:"createTime"`
		UnitFlag          bool   `json:"unitFlag"`
		Remark            string `json:"remark"`
		StatusName        string `json:"statusName"`
		Applier           string `json:"applier"`
		ID                int    `json:"id"`
		RelationName      string `json:"relationName"`
		IndependentNaming bool   `json:"independentNaming"`
		RegionID          string `json:"regionId"`
		GroupType         int    `json:"groupType"`
		ChannelName       string `json:"channelName"`
		NamespaceID       string `json:"namespaceId"`
		Status            int    `json:"status"`
		ChannelID         int    `json:"channelId"`
		UpdateTime        int64  `json:"updateTime"`
		Relation          int    `json:"relation"`
		RegionName        string `json:"regionName"`
		ConsumerID        string `json:"consumerId"`
		GroupID           string `json:"groupId"`
		Owner             string `json:"owner"`
	} `json:"Data"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      int    `json:"Code"`
	Total     int    `json:"Total"`
}
