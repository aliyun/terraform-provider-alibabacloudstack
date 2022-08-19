package apsarastack

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

type MaxComputeProject struct {
	AsapiSuccess bool   `json:"asapiSuccess"`
	Code         string `json:"code"`
	Data         struct {
		TotalCount  int `json:"otalCount"`
		PageSize    int `json:"PageSize"`
		PageNumber  int `json"PageNumber"`
		CalcEngines []struct {
			IsDefault      bool   `json:"IsDefault"`
			EngineId       int    `json:"EngineId"`
			DwRegion       string `json:"DwRegion"`
			CalcEngineType string `json:"CalcEngineType"`
			EnvType        string `json:"EnvType"`
			Name           string `json:"Name"`
			EngineInfo     struct {
				PubEndpoint        string `json:"pubEndpoint"`
				Specs              string `json:"specs"`
				ExternalProjectCnt int    `json:"externalProjectCnt"`
				Endpoint           string `json:"endpoint"`
				DefaultClusterArch string `json:"defaultClusterArch"`
				ResourceGroupType  string `json:"ODPS"`
				VpcEndpoint        string `json:"vpcEndpoint"`
				ProjectName        string `json:"projectName"`
				TaskSameAsOwner    bool   `json:"taskSameAsOwner"`
			} `json:"EngineInfo"`
			Department        int    `json:"Department"`
			Region            string `json:"Region"`
			DepartmentName    string `json:"DepartmentName"`
			RMRegionId        string `json:"RMRegionId"`
			ResourceGroup     int    `json:"ResourceGroup"`
			ResourceGroupName string `json:ResourceGroupName`
		} `json:"CalcEngines"`
	} `json:"Data"`
	Message      string `json:"message"`
	Success      bool   `json:"success"`
	Domain       string `json:"domain"`
	PureListData bool   `json:"pureListData"`
	TotalCount   int    `json:"innerTotalCount"`
}
