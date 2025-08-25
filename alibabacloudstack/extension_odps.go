package alibabacloudstack

type OdpsUser struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Data []struct {
		ID               int    `json:"id"`
		UserID           string `json:"userId"`
		UserPK           string `json:"aasPk"`
		UserName         string `json:"userName"`
		UserType         string `json:"userType"`
		OrganizationId   int    `json:"organizationId"`
		OrganizationName string `json:"organizationName"`
		Description      string `json:"description"`
	} `json:"data"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

type MaxComputeProjectDetailResponse struct {
	EagleEyeTraceId string `json:"eagleEyeTraceId"`
	RequestID       string `json:"RequestId"`
	Success         bool   `json:"success"`
	RequestID2      string `json:"requestId"`
	HttpStatusCode  int    `json:"HttpStatusCode"`
	Data            struct {
		TotalCount  int                 `json:"TotalCount"`
		PageSize    int                 `json:"PageSize"`
		PageNumber  int                 `json:"PageNumber"`
		CalcEngines []MaxComputeProject `json:"CalcEngines"`
	} `json:"Data"`
}

type MaxComputeProject struct {
	IsDefault      bool   `json:"IsDefault"`
	EngineId       int    `json:"EngineId"`
	DwRegion       string `json:"DwRegion"`
	CalcEngineType string `json:"CalcEngineType"`
	Data           []struct {
		Cluster   string  `json:"cluster"`
		IsDefault int     `json:"isDefault"`
		Cu        int     `json:"cu"`
		Disk      float64 `json:"disk"`
		QuotaName string  `json:"quotaName"`
		Original  struct {
			IsDefault    int         `json:"isDefault"`
			ProjectQuota interface{} `json:"projectQuota"`
			Quota        int         `json:"quota"`
			Name         string      `json:"name"`
		} `json:"original"`
		RegionId string  `json:"regionId"`
		Project  string  `json:"project"`
		QuotaId  float64 `json:"quotaId"`
	} `json:"Data,omitempty"`
	AscmCreateUser string `json:"AscmCreateUser,omitempty"`
	EnvType        string `json:"EnvType"`
	Name           string `json:"Name"`
	EngineInfo     struct {
		PubEndpoint        string `json:"pubEndpoint"`
		Specs              string `json:"specs"`
		ExternalProjectCnt int    `json:"externalProjectCnt"`
		Endpoint           string `json:"endpoint"`
		IsDisasterRecovery bool   `json:"isDisasterRecovery"`
		DefaultClusterArch string `json:"defaultClusterArch"`
		ResourceGroupType  string `json:"resourceGroupType"`
		VpcEndpoint        string `json:"vpcEndpoint"`
		ResourceGroupId    int64  `json:"resourceGroupId"`
		ProjectName        string `json:"projectName"`
		TaskSameAsOwner    bool   `json:"taskSameAsOwner"`
	} `json:"EngineInfo"`
	GmtCreate         string `json:"GmtCreate"`
	TenantId          int64  `json:"TenantId"`
	Department        int    `json:"Department"`
	Region            string `json:"Region"`
	EngineStatus      int    `json:"EngineStatus"`
	DepartmentName    string `json:"DepartmentName"`
	RegionId          string `json:"RegionId"`
	ResourceGroup     int    `json:"ResourceGroup"`
	ResourceGroupName string `json:"ResourceGroupName"`
}

type MaxComputeProjectEngineDetailResponse struct {
	EagleEyeTraceId string                      `json:"eagleEyeTraceId"`
	AsapiSuccess    bool                        `json:"asapiSuccess"`
	ResponseVersion string                      `json:"responseVersion"`
	RequestID       string                      `json:"RequestId"`
	Success         bool                        `json:"success"`
	RequestID2      string                      `json:"requestId"`
	HttpStatusCode  int                         `json:"HttpStatusCode"`
	Data            MaxComputeProjectEngineData `json:"Data"`
	Success2        bool                        `json:"Success"`
}

type MaxComputeProjectEngineData struct {
	EngineID   int `json:"EngineId"`
	EngineInfo struct {
		PubEndpoint string `json:"pubEndpoint"`
		TaskAk      struct {
			ExpireTime    string `json:"expireTime"`
			Kp            int64  `json:"kp"`
			AkType        int    `json:"akType"`
			AliyunAccount string `json:"aliyunAccount"`
			ExpTime       int64  `json:"expTime"`
			IsTempAk      bool   `json:"isTempAk"`
			ID            int    `json:"id"`
			BaseID        string `json:"baseId"`
			ProjectID     int    `json:"projectId"`
			UserBind      int    `json:"userBind"`
		} `json:"taskAk"`
		Specs              string `json:"specs"`
		Endpoint           string `json:"endpoint"`
		IsDisasterRecovery bool   `json:"isDisasterRecovery"`
		ResourceGroupType  string `json:"resourceGroupType"`
		VpcEndpoint        string `json:"vpcEndpoint"`
		ResourceGroupID    int64  `json:"resourceGroupId"`
		ProjectName        string `json:"projectName"`
		TaskSameAsOwner    bool   `json:"taskSameAsOwner"`
		OwnerAk            struct {
			ExpireTime    string `json:"expireTime"`
			Kp            int64  `json:"kp"`
			AkType        int    `json:"akType"`
			AliyunAccount string `json:"aliyunAccount"`
			ExpTime       int64  `json:"expTime"`
			IsTempAk      bool   `json:"isTempAk"`
			ID            int    `json:"id"`
			BaseID        string `json:"baseId"`
			ProjectID     int    `json:"projectId"`
			UserBind      int    `json:"userBind"`
		} `json:"ownerAk"`
	} `json:"EngineInfo"`
	McEncryptEnabled int    `json:"McEncryptEnabled"`
	TenantID         int64  `json:"TenantId"`
	Type             string `json:"Type"`
	DwRegion         string `json:"DwRegion"`
	Region           string `json:"Region"`
	EnvType          string `json:"EnvType"`
	Name             string `json:"Name"`
}
