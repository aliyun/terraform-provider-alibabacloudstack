package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type CsService struct {
	client *connectivity.AlibabacloudStackClient
}

const (
	COMPONENT_AUTO_SCALER      = "cluster-autoscaler"
	COMPONENT_DEFAULT_VRESION  = "v1.0.0"
	SCALING_CONFIGURATION_NAME = "kubernetes_autoscaler_autogen"
	DefaultECSTag              = "k8s.aliyun.com"
	DefaultClusterTag          = "ack.aliyun.com"
	RECYCLE_MODE_LABEL         = "k8s.io/cluster-autoscaler/node-template/label/policy"
	DefaultAutoscalerTag       = "k8s.io/cluster-autoscaler"
	SCALING_GROUP_NAME         = "sg-%s-%s"
	DEFAULT_COOL_DOWN_TIME     = 300
	RELEASE_MODE               = "release"
	RECYCLE_MODE               = "recycle"

	PRIORITY_POLICY       = "PRIORITY"
	COST_OPTIMIZED_POLICY = "COST_OPTIMIZED"
	BALANCE_POLICY        = "BALANCE"

	UpgradeClusterTimeout = 30 * time.Minute
)

func (s *CsService) DoCsDescribeclusterdetailRequest(id string) (cl *cs.KubernetesClusterDetail, err error) {
	return s.DescribeCsKubernetes(id)
}

func (s *CsService) DescribeCsKubernetes(id string) (cl *cs.KubernetesClusterDetail, err error) {
	cluster := &cs.KubernetesClusterDetail{}
	cluster.ClusterId = ""

	request := s.client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClustersV1", "/api/v1/clusters")
	request.QueryParams["SignatureVersion"] = "1.0"
	request.QueryParams["ProductName"] = "CS"

	clusterdetails, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if clusterdetails == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.DenverdinoAlibabacloudStackgo)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(clusterdetails.BaseResponse)
		return cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "DescribeKubernetesCluster", errmsgs.DenverdinoAlibabacloudStackgo, errmsg)

	}

	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		addDebug("DescribeKubernetesCluster", clusterdetails, request, requestMap)
	}
	Cdetails := ClustersV1{}
	_ = json.Unmarshal(clusterdetails.GetHttpContentBytes(), &Cdetails)

	cluster = &cs.KubernetesClusterDetail{}
	for _, k := range Cdetails.Clusters {
		if k.ClusterID == id {
			cluster.Tags = k.Tags
			cluster.Name = k.Name
			cluster.State = k.State
			cluster.ClusterId = k.ClusterID
			cluster.ClusterType = cs.KubernetesClusterType(k.ClusterType)
			cluster.VpcId = k.VpcID
			cluster.ResourceGroupId = k.ResourceGroupID
			cluster.ContainerCIDR = k.SubnetCidr
			cluster.CurrentVersion = k.CurrentVersion
			cluster.DeletionProtection = k.DeletionProtection
			cluster.RegionId = common.Region(k.RegionID)
			cluster.Size = int64(k.Size)
			cluster.IngressLoadbalancerId = k.ExternalLoadbalancerID
			cluster.InitVersion = k.InitVersion
			cluster.NetworkMode = k.NetworkMode
			cluster.PrivateZone = k.PrivateZone
			cluster.Profile = k.Profile
			cluster.VSwitchIds = k.VswitchID
			break
		}
	}
	if cluster.ClusterId != id {
		return cluster, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("CsKubernetes", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return cluster, nil
}

func (s *CsService) DescribeClusterNodes(id, nodepoolid string) (pools *NodePools, err error) {
	request := s.client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClusterNodes", fmt.Sprintf("/clusters/%s/nodes", id))
	mergeMaps(request.QueryParams, map[string]string{
		"SignatureVersion": "1.0",
		"nodepool_id":      nodepoolid,
		"ClusterId":        id,
	})

	response, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if response == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorClusterNodePoolNotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "DescribeClusterNodes", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if !response.IsSuccess() {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "DescribeClusterNodes", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	var clusternodepools *NodePools
	_ = json.Unmarshal(response.GetHttpContentBytes(), &clusternodepools)
	return clusternodepools, nil
}

func (s *CsService) DescribeClusterNodePools(id string) (*NodePool, error) {
	req := s.client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClusterNodePools", fmt.Sprintf("/clusters/%s/nodepools", id))
	req.QueryParams["ProductName"] = "CS"
	req.QueryParams["ClusterId"] = id
	var nodePool *responses.CommonResponse
	nodePool, err := s.client.ProcessCommonRequest(req)
	if err != nil {
		errmsg := ""
		if nodePool != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(nodePool.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_kubernetes", "DescribeClusterNodePools", nodePool, errmsg)
	}
	var node *NodePool

	if nodePool.IsSuccess() == false {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm", "API Action", nodePool.GetHttpContentString())
	}
	err = json.Unmarshal(nodePool.GetHttpContentBytes(), &node)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes", "DescribeClusterNodePools", nodePool)
	}
	return node, nil
}

func (s *CsService) CsKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		for _, failState := range failStates {
			if object.State == failState {
				return object, object.State, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.State))
			}
		}
		return object, object.State, nil
	}
}

func (s *CsService) CsKubernetesNodePoolStateRefreshFunc(id, clusterid string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetesNodePool(id, clusterid)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.Status.State) == failState {
				return object, string(object.Status.State), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, string(object.Status.State)))
			}
		}
		return object, string(object.Status.State), nil
	}
}

func (s *CsService) DescribeCsKubernetesNodePool(id, clusterid string) (*NodePoolAlone, error) {
	req := s.client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClusterNodePoolDetail", fmt.Sprintf("/clusters/%s/nodepools/%s", clusterid, id))
	req.Headers["x-acs-asapi-gateway-version"] = "3.0"
	req.QueryParams["ClusterId"] = clusterid
	req.QueryParams["NodepoolId"] = id
	response, err := s.client.ProcessCommonRequest(req)
	if err != nil || !response.IsSuccess() {
		if response == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.IsExpectedErrors(err, []string{"<QuerySeter> no row found"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.DenverdinoAlibabacloudStackgo)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_nodepool", "DescribeNodePool", response, errmsg)
	}
	var node *NodePoolAlone

	err = json.Unmarshal(response.GetHttpContentBytes(), &node)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_nodepool", "ParsenodepoolResponse", response)
	}
	return node, nil
}

func (s *CsService) UpgradeCluster(clusterId string, args *cs.UpgradeClusterArgs) error {
	invoker := NewInvoker()
	err := invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.UpgradeCluster(clusterId, args)
		})
		if e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapError(err)
	}

	state, upgradeError := s.WaitForUpgradeCluster(clusterId, "Upgrade")
	if state == cs.Task_Status_Success && upgradeError == nil {
		return nil
	}

	// if upgrade failed cancel the task
	err = invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.CancelUpgradeCluster(clusterId)
		})
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(upgradeError)
	}

	if state, err := s.WaitForUpgradeCluster(clusterId, "CancelUpgrade"); err != nil || state != cs.Task_Status_Success {
		log.Printf("[WARN] %s ACK Cluster cancel upgrade error: %#v", clusterId, err)
	}

	return errmsgs.WrapError(upgradeError)
}

func (s *CsService) WaitForUpgradeCluster(clusterId string, action string) (string, error) {
	err := resource.Retry(UpgradeClusterTimeout, func() *resource.RetryError {
		resp, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.QueryUpgradeClusterResult(clusterId)
		})
		if err != nil || resp == nil {
			return resource.RetryableError(err)
		}

		upgradeResult := resp.(*cs.UpgradeClusterResult)
		if upgradeResult.UpgradeStep == cs.UpgradeStep_Success {
			return nil
		}

		if upgradeResult.UpgradeStep == cs.UpgradeStep_Pause && upgradeResult.UpgradeStatus.Failed == "true" {
			msg := ""
			events := upgradeResult.UpgradeStatus.Events
			if len(events) > 0 {
				msg = events[len(events)-1].Message
			}
			return resource.NonRetryableError(fmt.Errorf("faild to %s cluster, error: %s", action, msg))
		}
		return resource.RetryableError(fmt.Errorf("%s cluster state not matched", action))
	})

	if err == nil {
		log.Printf("[INFO] %s ACK Cluster %s successed", action, clusterId)
		return cs.Task_Status_Success, nil
	}

	return cs.Task_Status_Failed, errmsgs.WrapError(err)
}

type Cluster struct {
	_                      string `json:"-"`
	Department             int64  `json:"Department"`
	DepartmentName         string `json:"DepartmentName"`
	ResourceGroup          int64  `json:"ResourceGroup"`
	ResourceGroupName      string `json:"ResourceGroupName"`
	ClusterHealthy         string `json:"cluster_healthy"`
	ClusterID              string `json:"cluster_id"`
	ClusterType            string `json:"cluster_type"`
	Created                string `json:"created"`
	CurrentVersion         string `json:"current_version"`
	DataDiskCategory       string `json:"data_disk_category"`
	DataDiskSize           int64  `json:"data_disk_size"`
	DeletionProtection     bool   `json:"deletion_protection"`
	DockerVersion          string `json:"docker_version"`
	EnabledMigration       bool   `json:"enabled_migration"`
	ErrMsg                 string `json:"err_msg"`
	ExternalLoadbalancerID string `json:"external_loadbalancer_id"`
	GwBridge               string `json:"gw_bridge"`
	InitVersion            string `json:"init_version"`
	InstanceType           string `json:"instance_type"`
	MasterURL              string `json:"master_url"`
	MetaData               string `json:"meta_data"`
	Name                   string `json:"name"`
	NeedUpdateAgent        bool   `json:"need_update_agent"`
	NetworkMode            string `json:"network_mode"`
	NodeStatus             string `json:"node_status"`
	Outputs                []struct {
		Description string      `json:"Description"`
		OutputKey   string      `json:"OutputKey"`
		OutputValue interface{} `json:"OutputValue"`
	} `json:"outputs"`
	Parameters struct {
		ALIYUN__AccountID        string `json:"ALIYUN::AccountId"`
		ALIYUN__NoValue          string `json:"ALIYUN::NoValue"`
		ALIYUN__Region           string `json:"ALIYUN::Region"`
		ALIYUN__StackID          string `json:"ALIYUN::StackId"`
		ALIYUN__StackName        string `json:"ALIYUN::StackName"`
		AdjustmentType           string `json:"AdjustmentType"`
		AuditFlags               string `json:"AuditFlags"`
		BetaVersion              string `json:"BetaVersion"`
		Ca                       string `json:"CA"`
		ClientCA                 string `json:"ClientCA"`
		CloudMonitorFlags        string `json:"CloudMonitorFlags"`
		CloudMonitorVersion      string `json:"CloudMonitorVersion"`
		ContainerCIDR            string `json:"ContainerCIDR"`
		DockerVersion            string `json:"DockerVersion"`
		Eip                      string `json:"Eip"`
		EipAddress               string `json:"EipAddress"`
		ElasticSearchHost        string `json:"ElasticSearchHost"`
		ElasticSearchPass        string `json:"ElasticSearchPass"`
		ElasticSearchPort        string `json:"ElasticSearchPort"`
		ElasticSearchUser        string `json:"ElasticSearchUser"`
		EtcdVersion              string `json:"EtcdVersion"`
		ExecuteVersion           string `json:"ExecuteVersion"`
		GPUFlags                 string `json:"GPUFlags"`
		HealthCheckType          string `json:"HealthCheckType"`
		IPVSEnable               string `json:"IPVSEnable"`
		ImageID                  string `json:"ImageId"`
		K8SMasterPolicyDocument  string `json:"K8SMasterPolicyDocument"`
		K8sWorkerPolicyDocument  string `json:"K8sWorkerPolicyDocument"`
		Key                      string `json:"Key"`
		KeyPair                  string `json:"KeyPair"`
		KubernetesVersion        string `json:"KubernetesVersion"`
		LoggingType              string `json:"LoggingType"`
		MasterAutoRenew          string `json:"MasterAutoRenew"`
		MasterAutoRenewPeriod    string `json:"MasterAutoRenewPeriod"`
		MasterDataDisk           string `json:"MasterDataDisk"`
		MasterDataDiskCategory   string `json:"MasterDataDiskCategory"`
		MasterDataDiskDevice     string `json:"MasterDataDiskDevice"`
		MasterDataDiskSize       string `json:"MasterDataDiskSize"`
		MasterImageID            string `json:"MasterImageId"`
		MasterInstanceChargeType string `json:"MasterInstanceChargeType"`
		MasterInstanceType       string `json:"MasterInstanceType"`
		MasterKeyPair            string `json:"MasterKeyPair"`
		MasterLoginPassword      string `json:"MasterLoginPassword"`
		MasterPeriod             string `json:"MasterPeriod"`
		MasterPeriodUnit         string `json:"MasterPeriodUnit"`
		MasterSystemDiskCategory string `json:"MasterSystemDiskCategory"`
		MasterSystemDiskSize     string `json:"MasterSystemDiskSize"`
		NatGateway               string `json:"NatGateway"`
		NatGatewayID             string `json:"NatGatewayId"`
		Network                  string `json:"Network"`
		NodeCIDRMask             string `json:"NodeCIDRMask"`
		NumOfNodes               string `json:"NumOfNodes"`
		Password                 string `json:"Password"`
		ProtectedInstances       string `json:"ProtectedInstances"`
		PublicSLB                string `json:"PublicSLB"`
		RemoveInstanceIds        string `json:"RemoveInstanceIds"`
		SLSProjectName           string `json:"SLSProjectName"`
		SNatEntry                string `json:"SNatEntry"`
		SSHFlags                 string `json:"SSHFlags"`
		ServiceCIDR              string `json:"ServiceCIDR"`
		SnatTableID              string `json:"SnatTableId"`
		UserCA                   string `json:"UserCA"`
		VSwitchID                string `json:"VSwitchId"`
		VpcID                    string `json:"VpcId"`
		WillReplace              string `json:"WillReplace"`
		WorkerAutoRenew          string `json:"WorkerAutoRenew"`
		WorkerAutoRenewPeriod    string `json:"WorkerAutoRenewPeriod"`
		WorkerDataDisk           string `json:"WorkerDataDisk"`
		WorkerDataDiskCategory   string `json:"WorkerDataDiskCategory"`
		WorkerDataDiskDevice     string `json:"WorkerDataDiskDevice"`
		WorkerDataDiskSize       string `json:"WorkerDataDiskSize"`
		WorkerImageID            string `json:"WorkerImageId"`
		WorkerInstanceChargeType string `json:"WorkerInstanceChargeType"`
		WorkerInstanceType       string `json:"WorkerInstanceType"`
		WorkerKeyPair            string `json:"WorkerKeyPair"`
		WorkerLoginPassword      string `json:"WorkerLoginPassword"`
		WorkerPeriod             string `json:"WorkerPeriod"`
		WorkerPeriodUnit         string `json:"WorkerPeriodUnit"`
		WorkerSystemDiskCategory string `json:"WorkerSystemDiskCategory"`
		WorkerSystemDiskSize     string `json:"WorkerSystemDiskSize"`
		ZoneID                   string `json:"ZoneId"`
	} `json:"parameters"`
	Port              int64  `json:"port"`
	PrivateZone       bool   `json:"private_zone"`
	Profile           string `json:"profile"`
	RegionID          string `json:"region_id"`
	ResourceGroupID   string `json:"resource_group_id"`
	SecurityGroupID   string `json:"security_group_id"`
	Size              int64  `json:"size"`
	State             string `json:"state"`
	SubnetCidr        string `json:"subnet_cidr"`
	SwarmMode         bool   `json:"swarm_mode"`
	Updated           string `json:"updated"`
	UpgradeComponents struct {
		Kubernetes struct {
			CanUpgrade     bool   `json:"can_upgrade"`
			Changed        string `json:"changed"`
			ComponentName  string `json:"component_name"`
			Exist          bool   `json:"exist"`
			Force          bool   `json:"force"`
			Message        string `json:"message"`
			NextVersion    string `json:"next_version"`
			Policy         string `json:"policy"`
			ReadyToUpgrade string `json:"ready_to_upgrade"`
			Required       bool   `json:"required"`
			Version        string `json:"version"`
		} `json:"Kubernetes"`
	} `json:"upgrade_components"`
	VpcID       string `json:"vpc_id"`
	VswitchCidr string `json:"vswitch_cidr"`
	VswitchID   string `json:"vswitch_id"`
	ZoneID      string `json:"zone_id"`
}

type NodePools struct {
	Nodes []struct {
		CreationTime       time.Time `json:"creation_time"`
		ErrorMessage       string    `json:"error_message"`
		InstanceName       string    `json:"instance_name"`
		NodeStatus         string    `json:"node_status"`
		IsAliyunNode       bool      `json:"is_aliyun_node"`
		NodeName           string    `json:"node_name"`
		ExpiredTime        time.Time `json:"expired_time"`
		IPAddress          []string  `json:"ip_address"`
		Source             string    `json:"source"`
		InstanceTypeFamily string    `json:"instance_type_family"`
		InstanceID         string    `json:"instance_id"`
		InstanceChargeType string    `json:"instance_charge_type"`
		InstanceRole       string    `json:"instance_role"`
		State              string    `json:"state"`
		InstanceStatus     string    `json:"instance_status"`
		ImageID            string    `json:"image_id"`
		InstanceType       string    `json:"instance_type"`
		NodepoolID         string    `json:"nodepool_id"`
		HostName           string    `json:"host_name"`
	} `json:"nodes"`
	Page struct {
		PageNumber int `json:"page_number"`
		TotalCount int `json:"total_count"`
		PageSize   int `json:"page_size"`
	} `json:"page"`
}

type NodePool struct {
	Nodepools []NodePoolAlone `json:"nodepools"`
}

type NodePoolAlone struct {
	TeeConfig struct {
		TeeEnable bool   `json:"tee_enable"`
		TeeType   string `json:"tee_type"`
	} `json:"tee_config"`
	ScalingGroup struct {
		InstanceTypes                    []string              `json:"instance_types"`
		PeriodUnit                       string                `json:"period_unit"`
		SecurityGroupID                  string                `json:"security_group_id"`
		MultiAzPolicy                    string                `json:"multi_az_policy"`
		Platform                         string                `json:"platform"`
		WorkerHpcClusterID               string                `json:"worker_hpc_cluster_id"`
		DataDisks                        []cs.NodePoolDataDisk `json:"data_disks"`
		RAMPolicy                        string                `json:"ram_policy"`
		LoginPassword                    string                `json:"login_password"`
		InstanceChargeType               string                `json:"instance_charge_type"`
		VswitchIds                       []string              `json:"vswitch_ids"`
		ScalingGroupID                   string                `json:"scaling_group_id"`
		Period                           int                   `json:"period"`
		AutoRenewPeriod                  int                   `json:"auto_renew_period"`
		WorkerDeploymentsetID            string                `json:"worker_deploymentset_id"`
		KeyPair                          string                `json:"key_pair"`
		SpotStrategy                     string                `json:"spot_strategy"`
		SystemDiskSize                   int                   `json:"system_disk_size"`
		Tags                             []cs.Tag              `json:"tags"`
		SpotPriceLimit                   []cs.SpotPrice        `json:"spot_price_limit"`
		AutoRenew                        bool                  `json:"auto_renew"`
		SystemDiskCategory               string                `json:"system_disk_category"`
		RdsInstances                     []interface{}         `json:"rds_instances"`
		WorkerSystemDiskSnapshotPolicyID string                `json:"worker_system_disk_snapshot_policy_id"`
		ImageID                          string                `json:"image_id"`
		ScalingPolicy                    string                `json:"scaling_policy"`
	} `json:"scaling_group"`
	KubernetesConfig struct {
		RuntimeVersion    string     `json:"runtime_version"`
		CPUPolicy         string     `json:"cpu_policy"`
		CmsEnabled        bool       `json:"cms_enabled"`
		Runtime           string     `json:"runtime"`
		OverwriteHostname bool       `json:"overwrite_hostname"`
		UserData          string     `json:"user_data"`
		NodeNameMode      string     `json:"node_name_mode"`
		Unschedulable     bool       `json:"unschedulable"`
		Taints            []cs.Taint `json:"taints"`
		Labels            []cs.Label `json:"labels"`
	} `json:"kubernetes_config"`
	AutoScaling  cs.AutoScaling `json:"auto_scaling"`
	NodepoolInfo struct {
		ResourceGroupID string    `json:"resource_group_id"`
		Created         time.Time `json:"created"`
		RegionID        string    `json:"region_id"`
		Name            string    `json:"name"`
		IsDefault       bool      `json:"is_default"`
		Type            string    `json:"type"`
		NodepoolID      string    `json:"nodepool_id"`
		Updated         time.Time `json:"updated"`
	} `json:"nodepool_info"`
	Status struct {
		ServingNodes  int    `json:"serving_nodes"`
		TotalNodes    int    `json:"total_nodes"`
		State         string `json:"state"`
		OfflineNodes  int    `json:"offline_nodes"`
		RemovingNodes int    `json:"removing_nodes"`
		InitialNodes  int    `json:"initial_nodes"`
		FailedNodes   int    `json:"failed_nodes"`
		HealthyNodes  int    `json:"healthy_nodes"`
	} `json:"status"`
}

type ClustersV1 struct {
	Redirect        bool   `json:"redirect"`
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	Code            string `json:"code"`
	Cost            int    `json:"cost"`
	Message         string `json:"message"`
	ServerRole      string `json:"serverRole"`
	AsapiRequestID  string `json:"asapiRequestId"`
	Success         bool   `json:"success"`
	PageInfo        struct {
		PageNumber int `json:"page_number"`
		TotalCount int `json:"total_count"`
		PageSize   int `json:"page_size"`
	} `json:"page_info"`
	Domain       string `json:"domain"`
	PureListData bool   `json:"pureListData"`
	API          string `json:"api"`
	Clusters     []struct {
		Tags                   []cs.Tag  `json:"tags"`
		ResourceGroupID        string    `json:"resource_group_id"`
		PrivateZone            bool      `json:"private_zone"`
		VpcID                  string    `json:"vpc_id"`
		NetworkMode            string    `json:"network_mode"`
		SecurityGroupID        string    `json:"security_group_id"`
		ClusterType            string    `json:"cluster_type"`
		DockerVersion          string    `json:"docker_version"`
		DataDiskCategory       string    `json:"data_disk_category"`
		NextVersion            string    `json:"next_version"`
		ZoneID                 string    `json:"zone_id"`
		ClusterID              string    `json:"cluster_id"`
		Department             int       `json:"Department"`
		ExternalLoadbalancerID string    `json:"external_loadbalancer_id"`
		VswitchID              string    `json:"vswitch_id"`
		SwarmMode              bool      `json:"swarm_mode"`
		RMRegionID             string    `json:"RMRegionId"`
		State                  string    `json:"state"`
		ResourceGroup          int       `json:"ResourceGroup"`
		InitVersion            string    `json:"init_version"`
		NodeStatus             string    `json:"node_status"`
		NeedUpdateAgent        bool      `json:"need_update_agent"`
		Created                time.Time `json:"created"`
		DeletionProtection     bool      `json:"deletion_protection"`
		SubnetCidr             string    `json:"subnet_cidr"`
		Profile                string    `json:"profile"`
		RegionID               string    `json:"region_id"`
		MasterURL              string    `json:"master_url"`
		CurrentVersion         string    `json:"current_version"`
		NAMING_FAILED          string    `json:"-"`
		VswitchCidr            string    `json:"vswitch_cidr"`
		ClusterHealthy         string    `json:"cluster_healthy"`
		ClusterSpec            string    `json:"cluster_spec"`
		Size                   int       `json:"size"`
		DataDiskSize           int       `json:"data_disk_size"`
		Port                   int       `json:"port"`
		EnabledMigration       bool      `json:"enabled_migration"`
		Name                   string    `json:"name"`
		DepartmentName         string    `json:"DepartmentName"`
		Updated                time.Time `json:"updated"`
		InstanceType           string    `json:"instance_type"`
		WorkerRAMRoleName      string    `json:"worker_ram_role_name"`
		ResourceGroupName      string    `json:"ResourceGroupName"`
	} `json:"clusters"`
}
