package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type EdasService struct {
	client *connectivity.AlibabacloudStackClient
}

type Hook struct {
	Exec      *Exec      `json:"exec,omitempty"`
	HttpGet   *HttpGet   `json:"httpGet,omitempty"`
	TcpSocket *TcpSocket `json:"tcpSocket,omitempty"`
}

type Exec struct {
	Command []string `json:"command"`
}

type HttpGet struct {
	Path        string       `json:"path"`
	Port        int          `json:"port"`
	Scheme      string       `json:"scheme"`
	HttpHeaders []HttpHeader `json:"httpHeaders"`
}

type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TcpSocket struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Prober struct {
	FailureThreshold    int `json:"failureThreshold"`
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	SuccessThreshold    int `json:"successThreshold"`
	TimeoutSeconds      int `json:"timeoutSeconds"`
	Hook                `json:",inline"`
}

type EdasChangeOrderInfo struct {
	Status                 int         `json:"Status" xml:"Status"`
	ChangeOrderId          string      `json:"ChangeOrderId" xml:"ChangeOrderId"`
	BatchType              string      `json:"BatchType" xml:"BatchType"`
	CoType                 string      `json:"CoType" xml:"CoType"`
	CreateTime             string      `json:"CreateTime" xml:"CreateTime"`
	ChangeOrderDescription string      `json:"ChangeOrderDescription" xml:"ChangeOrderDescription"`
	BatchCount             int         `json:"BatchCount" xml:"BatchCount"`
	CreateUserId           string      `json:"CreateUserId" xml:"CreateUserId"`
	SupportRollback        bool        `json:"SupportRollback" xml:"SupportRollback"`
	Desc                   string      `json:"Desc" xml:"Desc"`
	Targets                interface{} `json:"Targets" xml:"Targets"`
	TrafficControl         interface{} `json:"TrafficControl" xml:"TrafficControl"`
	PipelineInfoList       interface{} `json:"PipelineInfoList" xml:"PipelineInfoList"`
}

type EdasGetChangeOrderInfoResponse struct {
	EagleEyeTraceId string              `json:"eagleEyeTraceId"`
	AsapiSuccess    bool                `json:"asapiSuccess"`
	ResponseVersion string              `json:"responseVersion"`
	RequestId       string              `json:"RequestId"`
	Message         string              `json:"Message"`
	ChangeOrderInfo EdasChangeOrderInfo `json:"ChangeOrderInfo"`
	Success         bool                `json:"success"`
}

func (e *EdasService) GetChangeOrderStatus(id string) (info *EdasChangeOrderInfo, err error) {
	order := EdasChangeOrderInfo{}
	request := e.client.NewCommonRequest("POST", "Edas", "2017-08-01", "GetChangeOrderInfo", "/pop/v5/changeorder/change_order_info")
	request.QueryParams["ChangeOrderId"] = id

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := e.client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return &order, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request, request)
	response := EdasGetChangeOrderInfoResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return &order, errmsgs.WrapError(err)
	}
	order = response.ChangeOrderInfo
	return &order, nil
}

func (e *EdasService) GetDeployGroup(appId, groupId string) (groupInfo *edas.DeployGroup, err error) {
	request := edas.CreateListDeployGroupRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	rsp, ok := raw.(*edas.ListDeployGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(rsp.BaseResponse)
		}
		return groupInfo, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, appId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if rsp.Code != 200 {
		return groupInfo, errmsgs.Error("get deploy group failed for " + rsp.Message)
	}
	for _, group := range rsp.DeployGroupList.DeployGroup {
		if group.GroupId == groupId {
			return &group, nil
		}
	}
	return groupInfo, nil
}

func (e *EdasService) EdasChangeOrderStatusRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := e.GetChangeOrderStatus(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if strconv.Itoa(object.Status) == failState {
				return object, strconv.Itoa(object.Status), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, strconv.Itoa(object.Status)))
			}
		}

		return object, strconv.Itoa(object.Status), nil
	}
}

func (e *EdasService) SyncResource(resourceType string) error {
	request := edas.CreateSynchronizeResourceRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.Type = resourceType

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.SynchronizeResource(request)
	})

	rsp, ok := raw.(*edas.SynchronizeResourceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(rsp.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "sync resource", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if rsp.Code != 200 || !rsp.Success {
		return errmsgs.WrapError(errmsgs.Error("sync resource failed for " + rsp.Message))
	}

	return nil
}

func (e *EdasService) CheckEcsStatus(instanceIds string, count int) error {
	request := ecs.CreateDescribeInstancesRequest()
	e.client.InitRpcRequest(*request.RpcRequest)
	request.Status = "Running"
	request.PageSize = requests.NewInteger(100)
	request.InstanceIds = instanceIds

	raw, err := e.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(request)
	})

	rsp, ok := raw.(*ecs.DescribeInstancesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(rsp.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceIds, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(rsp.Instances.Instance) != count {
		return errmsgs.WrapErrorf(errmsgs.Error("not enough instances"), errmsgs.DefaultErrorMsg, instanceIds, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func (e *EdasService) GetLastPackgeVersion(appId, groupId string) (string, error) {
	var versionId string
	request := edas.CreateQueryApplicationStatusRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.QueryApplicationStatus(request)
	})
	response, ok := raw.(*edas.QueryApplicationStatusResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_application_package_version", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		return "", errmsgs.WrapError(errmsgs.Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, group := range response.AppInfo.GroupList.Group {
		if group.GroupId == groupId {
			versionId = group.PackageVersionId
		}
	}

	rq := edas.CreateListHistoryDeployVersionRequest()
	e.client.InitRoaRequest(*rq.RoaRequest)
	rq.AppId = appId

	rq.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err = e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListHistoryDeployVersion(rq)
	})
	rsp, ok := raw.(*edas.ListHistoryDeployVersionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(rsp.BaseResponse)
		}
		return "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_application_package_version_list", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if rsp.Code != 200 {
		return "", errmsgs.WrapError(errmsgs.Error("QueryApplicationStatus failed for " + response.Message))
	}

	for _, version := range rsp.PackageVersionList.PackageVersion {
		if version.Id == versionId {
			return version.PackageVersion, nil
		}
	}

	return "", nil
}

func (e *EdasService) DescribeEdasApplication(appId string) (*edas.Applcation, error) {
	application := &edas.Applcation{}

	request := edas.CreateGetApplicationRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})
	response, ok := raw.(*edas.GetApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return application, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_application", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		return application, errmsgs.WrapError(errmsgs.Error("get application error :" + response.Message))
	}

	v := response.Applcation

	return &v, nil
}

func (s *EdasService) ClusterImportK8sStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEdasK8sCluster(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		status := fmt.Sprintf("%d", object.ClusterImportStatus)
		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}

		return object, status, nil
	}
}

func (e *EdasService) DescribeEdasGetCluster(clusterId string) (*edas.Cluster, error) {
	cluster := edas.Cluster{}
	request := e.client.NewCommonRequest("GET", "Edas", "2017-08-01", "GetCluster", "/pop/v5/resource/cluster")
	request.QueryParams["ResourceGroupId"] = e.client.ResourceGroup
	request.QueryParams["ClusterId"] = clusterId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	bresponse, err := e.client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return &cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request, request)

	response := edas.GetClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return &cluster, errmsgs.WrapError(err)
	}

	if response.Code != 200 {
		return &cluster, errmsgs.WrapError(errmsgs.Error("create cluster failed for " + response.Message))
	}

	cluster = response.Cluster

	return &cluster, nil
}

type EdasK8sCluster struct {
	ClusterImportStatus int    `json:"ClusterImportStatus"`
	NodeNum             int    `json:"NodeNum"`
	ClusterId           string `json:"ClusterId"`
	Cpu                 int    `json:"Cpu"`
	ClusterType         int    `json:"ClusterType"`
	NetworkMode         int    `json:"NetworkMode"`
	CsClusterId         string `json:"CsClusterId"`
	VswitchId           string `json:"VswitchId"`
	VpcId               string `json:"VpcId"`
	Mem                 int    `json:"Mem"`
	ClusterName         string `json:"ClusterName"`
	SubNetCidr          string `json:"SubNetCidr"`
	RegionId            string `json:"RegionId"`
	CsClusterStatus     string `json:"CsClusterStatus"`
	ClusterStatus       int    `json:"ClusterStatus"`
	SubClusterType      string `json:"SubClusterType"`
}

type EdasGetK8sClusterResponse struct {
	EagleEyeTraceId string `json:"eagleEyeTraceId"`
	AsapiSuccess    bool   `json:"asapiSuccess"`
	ResponseVersion string `json:"responseVersion"`
	RequestId       string `json:"RequestId"`
	Message         string `json:"Message"`
	ClusterPage     struct {
		ClusterList []EdasK8sCluster `json:"ClusterList"`
		PageSize    int              `json:"PageSize"`
		CurrentPage int              `json:"CurrentPage"`
		TotalSize   int              `json:"TotalSize"`
	} `json:"ClusterPage"`
	Success bool `json:"success"`
	Code    int  `json:"Code"`
}

func (e *EdasService) DescribeEdasK8sCluster(clusterId string) (*EdasK8sCluster, error) {
	cluster := EdasK8sCluster{}
	request := e.client.NewCommonRequest("POST", "Edas", "2017-08-01", "GetK8sCluster", "/pop/v5/k8s_clusters")
	request.QueryParams["ResourceGroupId"] = e.client.ResourceGroup
	request.QueryParams["ClusterId"] = clusterId

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := e.client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return &cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request, request)

	response := EdasGetK8sClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return &cluster, errmsgs.WrapError(err)
	}

	if response.Code != 200 {
		return &cluster, errmsgs.WrapError(errmsgs.Error("create cluster failed for " + response.Message))
	}
	ClusterList := response.ClusterPage.ClusterList
	if len(ClusterList) == 0 {
		return &cluster, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg)
	} else {
		cluster = ClusterList[0]
	}

	return &cluster, nil
}

func (e *EdasService) DescribeEdasListCluster(clusterId string) (*edas.Cluster, error) {
	cluster := edas.Cluster{}

	request := e.client.NewCommonRequest("POST", "Edas", "2017-08-01", "ListCluster", "/pop/v5/resource/cluster_list")
	request.QueryParams["ResourceGroupId"] = e.client.ResourceGroup
	request.QueryParams["LogicalRegionId"] = e.client.RegionId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	bresponse, err := e.client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return &cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request, request)

	response := edas.ListClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return &cluster, errmsgs.WrapError(err)
	}
	if response.Code != 200 {
		return &cluster, errmsgs.WrapError(errmsgs.Error("create cluster failed for " + response.Message))
	}

	for _, onecluster := range response.ClusterList.Cluster {
		if onecluster.ClusterId == clusterId {
			if onecluster.CsClusterStatus == "running" {
				cluster = onecluster
			}
		}
	}

	return &cluster, nil
}

type PortMapping struct {
	TargetPort  string `json:"targetPort"`
	ServicePort string `json:"servicePort"`
	Protocol    string `json:"protocol"`
}

type EdasK8sService struct {
	Type                  string                 `json:"type"`
	ServiceName           string                 `json:"serviceName"`
	AllowEdit             bool                   `json:"allowEdit"`
	ClusterIP             string                 `json:"clusterIP"`
	PortMappings          []PortMapping          `json:"portMappings"`
	Annotations           map[string]interface{} `json:"annotations"`
	Labels                map[string]interface{} `json:"labels"`
	ExternalTrafficPolicy string                 `json:"externalTrafficPolicy"`
	InnerEndpointer       string                 `json:"innerEndpointer"`
	Namespace             string                 `json:"namespace"`
	NodeIpList            []string               `json:"nodeIpList"`
	LoadBalancerIP        string                 `json:"loadbalancerip"`
}

type EdasK8sServiceResponse struct {
	Code      int               `json:"Code"`
	Message   string            `json:"Message"`
	RequestId string            `json:"RequestId"`
	Data      []*EdasK8sService `json:"Data"`
}

func (e *EdasService) ListEdasK8sServices(app_id string) ([]*EdasK8sService, error) {

	request := e.client.NewCommonRequest("GET", "Edas", "2017-08-01", "ListK8sServices", "/pop/v5/k8s/service/list_service")
	request.QueryParams["AppId"] = app_id

	bresponse, err := e.client.ProcessCommonRequest(request)
	addDebug("ListK8sServices", bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_service", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var response *EdasK8sServiceResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	if fmt.Sprint(response.Code) != "200" {
		return nil, errmsgs.WrapError(fmt.Errorf("List k8s services failed for %s", response.Message))
	}
	// services_map := response["Services"].([]interface{})
	// err = json.Unmarshal([]byte(fmt.Sprint(services_map)), &services)
	// if err != nil {
	// 	return nil, errmsgs.WrapError(err)
	// }
	return response.Data, nil
}

func (e *EdasService) DescribeEdasK8sService(id string) (*EdasK8sService, error) {
	parts := strings.Split(id, ":")
	app_id := parts[0]
	name := parts[1]
	services, err := e.ListEdasK8sServices(app_id)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if service.ServiceName == name {
			return service, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.NotFoundMsg, " Edas K8s cluster")
}

func (e *EdasService) DescribeEdasDeployGroup(id string) (*edas.DeployGroup, error) {
	group := &edas.DeployGroup{}

	strs := strings.Split(id, ":")

	request := edas.CreateListDeployGroupRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = strs[0]

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	response, ok := raw.(*edas.ListDeployGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return group, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_deploy_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		return group, errmsgs.WrapError(errmsgs.Error("create cluster failed for " + response.Message))
	}

	for _, v := range response.DeployGroupList.DeployGroup {
		if v.ClusterName == strs[1] {
			return &v, nil
		}
	}

	return group, nil
}

func (e *EdasService) DescribeEdasInstanceClusterAttachment(id string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasGetCluster(v[0])
	if err != nil {
		return cluster, errmsgs.WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationDeployment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, errmsgs.WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasApplicationScale(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, errmsgs.WrapError(err)
	}

	return o, nil
}

func (e *EdasService) DescribeEdasSlbAttachment(id string) (*edas.Applcation, error) {
	application := &edas.Applcation{}
	v := strings.Split(id, ":")
	o, err := e.DescribeEdasApplication(v[0])
	if err != nil {
		return application, errmsgs.WrapError(err)
	}

	return o, nil
}

type CommandArg struct {
	Argument string `json:"argument" xml:"argument"`
}

func (e *EdasService) GetK8sCommandArgs(args []interface{}) (string, error) {
	aString := make([]CommandArg, 0)
	for _, v := range args {
		aString = append(aString, CommandArg{Argument: v.(string)})
	}
	b, err := json.Marshal(aString)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

type PvcMountPaths struct {
	MountPath string `json:"mountPath" xml:"mountPath"`
	ReadOnly  bool   `json:"readOnly" xml:"readOnly"`
}

type PvcMountDescs struct {
	PvcName    string          `json:"pvcName" xml:"pvcName"`
	MountPaths []PvcMountPaths `json:"mountPaths" xml:"mountPaths"`
}

func (e *EdasService) GetK8sPvcMountDescs(pvc_mount_descs []interface{}) (string, error) {
	pvcMountDescs := make([]PvcMountDescs, 0)
	for _, pvc_mount_desc := range pvc_mount_descs {
		p := pvc_mount_desc.(map[string]interface{})
		mountPaths := p["mount_paths"].([]interface{})
		mount_paths := make([]PvcMountPaths, 0)
		for _, mountPath := range mountPaths {
			m := mountPath.(map[string]interface{})
			mount_paths = append(mount_paths, PvcMountPaths{
				MountPath: m["mount_path"].(string),
				ReadOnly:  m["read_only"].(bool),
			})
		}
		pvcMountDescs = append(pvcMountDescs, PvcMountDescs{
			PvcName:    p["pvc_name"].(string),
			MountPaths: mount_paths,
		})
	}
	b, err := json.Marshal(pvcMountDescs)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

type ConfigMaps struct {
	Type      string `json:"type" xml:"type"`
	Name      string `json:"name" xml:"name"`
	MountPath string `json:"mountPath" xml:"mountPath"`
}

func (e *EdasService) GetK8sConfigMaps(configmaps []interface{}) (string, error) {
	config_mount_descs := make([]ConfigMaps, 0)
	for _, v := range configmaps {
		m := v.(map[string]interface{})
		config_mount_descs = append(config_mount_descs, ConfigMaps{
			Type:      m["type"].(string),
			Name:      m["name"].(string),
			MountPath: m["mount_path"].(string),
		})
	}
	b, err := json.Marshal(config_mount_descs)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

type LocalVolume struct {
	Type      string `json:"type" xml:"type"`
	NodePath  string `json:"nodePath" xml:"nodePath"`
	MountPath string `json:"mountPath" xml:"mountPath"`
}

func (e *EdasService) GetK8sLocalVolumes(local_volumes []interface{}) (string, error) {
	localVolume := make([]LocalVolume, 0)
	for _, v := range local_volumes {
		m := v.(map[string]interface{})
		localVolume = append(localVolume, LocalVolume{
			Type:      m["type"].(string),
			MountPath: m["mount_path"].(string),
			NodePath:  m["node_path"].(string),
		})
	}
	b, err := json.Marshal(localVolume)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) GetK8sCommandArgsForDeploy(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

type K8sEnv struct {
	Name  string `json:"name" xml:"name"`
	Value string `json:"value" xml:"value"`
}

func (e *EdasService) GetK8sEnvs(envs map[string]interface{}) (string, error) {
	k8sEnvs := make([]K8sEnv, 0)
	for n, v := range envs {
		k8sEnvs = append(k8sEnvs, K8sEnv{Name: n, Value: v.(string)})
	}

	b, err := json.Marshal(k8sEnvs)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

func (e *EdasService) QueryK8sAppPackageType(appId string) (string, error) {
	v, err := e.DescribeEdasK8sApplication(appId)
	if err != nil {
		return "", err
	}
	return v.App.ApplicationType, nil
}

type K8sServicePorts struct {
	TargetPort  int    `json:"TargetPort"`
	ServicePort int    `json:"ServicePort"`
	Protocol    string `json:"Protocol"`
}

func (e *EdasService) GetK8sServicePorts(service_ports []interface{}) (string, error) {
	k8s_service_ports := make([]K8sServicePorts, 0)
	for _, v := range service_ports {
		s := v.(map[string]interface{})
		k8s_service_ports = append(k8s_service_ports, K8sServicePorts{
			Protocol:    s["protocol"].(string),
			ServicePort: s["service_port"].(int),
			TargetPort:  s["target_port"].(int),
		})
	}
	b, err := json.Marshal(k8s_service_ports)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(b), nil
}

type EdasK8sApplcation struct {
	Instaces      EdasK8sAppInstances `json:"Instances" xml:"Instances"`
	Conf          edas.Conf           `json:"Conf" xml:"Conf"`
	App           edas.App            `json:"App" xml:"App"`
	ImageInfo     edas.ImageInfo      `json:"ImageInfo" xml:"ImageInfo"`
	LatestVersion edas.LatestVersion  `json:"LatestVersion" xml:"LatestVersion"`
	DeployGroups  edas.DeployGroups   `json:"DeployGroups" xml:"DeployGroups"`
}

type EdasK8sAppInstances struct {
	Instance []EdasK8sAppInstance `json:"Instance"`
}

type EdasK8sAppInstance struct {
	GroupName       string `json:"GroupName"`
	Status          int    `json:"Status"`
	ContainerIp     string `json:"ContainerIp"`
	Expired         bool   `json:"Expired"`
	PackageVersion  string `json:"PackageVersion"`
	Url             string `json:"Url"`
	GroupId         string `json:"GroupId"`
	ContainerStatus string `json:"ContainerStatus"`
	AppId           string `json:"AppId"`
	GroupType       int    `json:"GroupType"`
	PrivateIp       string `json:"PrivateIp"`
	Running         bool   `json:"Running"`
	EccId           string `json:"EccId"`
	Id              string `json:"Id"`
	AppState        int    `json:"AppState"`
	Online          int    `json:"Online"`
}

type EdasGetK8sApplcationResponse struct {
	EagleEyeTraceId string            `json:"eagleEyeTraceId"`
	AsapiSuccess    bool              `json:"asapiSuccess"`
	ResponseVersion string            `json:"responseVersion"`
	RequestId       string            `json:"RequestId"`
	Message         string            `json:"Message"`
	Applcation      EdasK8sApplcation `json:"Applcation"`
	Success         bool              `json:"success"`
	Code            int               `json:"Code"`
}

func (e *EdasService) DescribeEdasK8sApplication(appId string) (*EdasK8sApplcation, error) {
	v := EdasK8sApplcation{}
	request := e.client.NewCommonRequest("GET", "Edas", "2017-08-01", "GetK8sApplication", "/pop/v5/changeorder/co_application")
	request.QueryParams["ResourceGroupId"] = e.client.ResourceGroup
	request.QueryParams["AppId"] = appId

	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := e.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	if err != nil {
		return &v, err
	}

	response := EdasGetK8sApplcationResponse{}
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	v = response.Applcation

	return &v, nil
}

func (e *EdasService) PreStopEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) PostStartEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldHook Hook
	err := json.Unmarshal([]byte(oldStr), &oldHook)
	if err != nil {
		return false
	}
	var newHook Hook
	err = json.Unmarshal([]byte(newStr), &newHook)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldHook, newHook)
}

func (e *EdasService) LivenessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}

func (e *EdasService) ReadinessEqual(old, new interface{}) bool {
	oldStr := old.(string)
	newStr := new.(string)
	var oldProber Prober
	err := json.Unmarshal([]byte(oldStr), &oldProber)
	if err != nil {
		return false
	}
	var newProber Prober
	err = json.Unmarshal([]byte(newStr), &newProber)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(oldProber, newProber)
}
func (s *EdasService) DescribeEdasNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "/pop/v5/user_region_defs"
	request := map[string]interface{}{}
	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", "ListUserDefineRegion", "/pop/v5/user_region_defs", nil, request, nil)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.UserDefineRegionList.UserDefineRegionEntity", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.UserDefineRegionList.UserDefineRegionEntity", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EDAS", id)), errmsgs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("EDAS", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}
