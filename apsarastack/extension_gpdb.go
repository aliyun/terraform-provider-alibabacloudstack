package apsarastack

type GpdbInstance struct {
	Items struct {
		DBInstance []struct {
			EngineVersion         string `json:"EngineVersion"`
			ZoneID                string `json:"ZoneId"`
			DBInstanceStatus      string `json:"DBInstanceStatus"`
			DBInstanceNetType     string `json:"DBInstanceNetType"`
			CreateTime            string `json:"CreateTime"`
			VSwitchID             string `json:"VSwitchId,omitempty"`
			PayType               string `json:"PayType"`
			LockMode              string `json:"LockMode"`
			InstanceNetworkType   string `json:"InstanceNetworkType"`
			Department            int    `json:"Department"`
			VpcID                 string `json:"VpcId,omitempty"`
			DBInstanceID          string `json:"DBInstanceId"`
			DepartmentName        string `json:"DepartmentName"`
			RegionID              string `json:"RegionId"`
			LockReason            string `json:"LockReason"`
			DBInstanceDescription string `json:"DBInstanceDescription"`
			Engine                string `json:"Engine"`
			Tags                  struct {
				Tag []interface{} `json:"Tag"`
			} `json:"Tags"`
			ResourceGroup     int    `json:"ResourceGroup"`
			ResourceGroupName string `json:"ResourceGroupName"`
		} `json:"DBInstance"`
	} `json:"Items"`
	PageNumber       int    `json:"PageNumber"`
	PageSize         int    `json:"PageSize"`
	TotalRecordCount int    `json:"TotalRecordCount"`
	Code             string `json:"code"`
	Cost             int    `json:"cost"`
	Message          string `json:"message"`
	PureListData     bool   `json:"pureListData"`
	Redirect         bool   `json:"redirect"`
	Success          bool   `json:"success"`
}
