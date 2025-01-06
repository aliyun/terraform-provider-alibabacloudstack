package datahub_patch

type Commonds struct {
	// StatusCode http return code

	Commond `json:"Command"`
}
type Commond struct {
	CommandId       string `json:"CommandId"`
	CommandContent  string `json:"CommandContent"`
	Description     string `json:"Description"`
	EnableParameter bool   `json:"EnableParameter"`
	Name            string `json:"Name"`
	Timeout         int64  `json:"Timeout"`
	Type            string `json:"Type"`
	WorkingDir      string `json:"WorkingDir"`
}

type EcsDescribeEcsCommandResult struct {
	Commands struct {
		Command []struct {
			CommandId       string `json:"CommandId"`
			CommandContent  string `json:"CommandContent"`
			Description     string `json:"Description"`
			EnableParameter bool   `json:"EnableParameter"`
			Name            string `json:"Name"`
			Timeout         int64  `json:"Timeout"`
			Type            string `json:"Type"`
			WorkingDir      string `json:"WorkingDir"`
		} `json:"Command"`
	} `json:"Commands"`
	RequestId string `json:"RequestId"`
}

type EcsDescribeEcsEbsStorageSetsResult struct {
	StorageSets struct {
		StorageSet []struct {
			shared                    int32  `json:"shared"`
			Department                int32  `json:"Department"`
			ZoneId                    string `json:"ZoneId"`
			CreationTime              string `json:"CreationTime"`
			DepartmentName            string `json:"DepartmentName"`
			StorageSetPartitionNumber int    `json:"StorageSetPartitionNumber"`
			RMRegionId                string `json:"RMRegionId"`
			StorageSetId              string `json:"StorageSetId"`
			RegionId                  string `json:"RegionId"`
			StorageSetName            string `json:"StorageSetName"`
			ResourceGroup             int32  `json:"ResourceGroup"`
			ResourceGroupName         string `json:"ResourceGroupName"`
		} `json:"StorageSet"`
	} `json:"StorageSets"`
	RequestId string `json:"RequestId"`
}

type EcsDescribeDeploymentSetsResult struct {
	DeploymentSets struct {
		DeploymentSet []struct {
			Granularity              string `json:"Granularity"`
			DeploymentStrategy       string `json:"DeploymentStrategy"`
			DeploymentSetDescription string `json:"DeploymentSetDescription"`
			DeploymentSetName        string `json:"DeploymentSetName"`
			Domain                   string `json:"Domain"`
		} `json:"DeploymentSet"`
	} `json:"DeploymentSets"`
}
type EcsDescribeEcsHpcClusterResult struct {
	HpcClusters struct {
		HpcCluster []struct {
			Description  string `json:"Description"`
			HpcClusterId string `json:"HpcClusterId"`
			Name         string `json:"Name"`
		} `json:"HpcCluster"`
	} `json:"HpcClusters"`
}
type EcsDeploymentSetCreateResult struct {
	DeploymentSetId string `json:"DeploymentSetId"`
}
