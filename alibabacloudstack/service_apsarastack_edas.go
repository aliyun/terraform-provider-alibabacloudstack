package alibabacloudstack

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"

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

func (e *EdasService) GetChangeOrderStatus(id string) (info *edas.ChangeOrderInfo, err error) {
	request := edas.CreateGetChangeOrderInfoRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.ChangeOrderId = id

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetChangeOrderInfo(request)
	})

	rsp, ok := raw.(*edas.GetChangeOrderInfoResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(rsp.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.InvalidDBClusterIdNotFound", "OperationDenied.InvalidDBClusterNameNotFound"}) {
			return info, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return info, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	return &rsp.ChangeOrderInfo, nil
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

func (e *EdasService) DescribeEdasGetCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}

	request := edas.CreateGetClusterRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	response, ok := raw.(*edas.GetClusterResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		return cluster, errmsgs.WrapError(errmsgs.Error("create cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
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
	request := edas.CreateGetApplicationRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	log.Printf("-------------------------------------------- %v", request.Headers)
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetApplication(request)
	})
	addDebug(request.GetActionName(), raw, request, request.RoaRequest)
	response, ok := raw.(*edas.GetApplicationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return "", errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_k8s_app_package_type", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if response.Code != 200 {
		return "", errmsgs.WrapError(errmsgs.Error("get application for appId:" + appId + " failed:" + response.Message))
	}
	if len(response.Applcation.ApplicationType) > 0 {
		return response.Applcation.ApplicationType, nil
	}
	return "", errmsgs.WrapError(errmsgs.Error("not package type for appId:" + appId))
}

func (e *EdasService) DescribeEdasK8sCluster(clusterId string) (*edas.Cluster, error) {
	cluster := &edas.Cluster{}

	request := edas.CreateGetClusterRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.ClusterId = clusterId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, err := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	response, ok := raw.(*edas.GetClusterResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return cluster, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, clusterId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		if strings.Contains(response.Message, "does not exist") {
			return cluster, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return cluster, errmsgs.WrapError(errmsgs.Error("create k8s cluster failed for " + response.Message))
	}

	v := response.Cluster

	return &v, nil
}

func (e *EdasService) DescribeEdasK8sApplication(appId string) (*edas.Applcation, error) {

	request := edas.CreateGetK8sApplicationRequest()
	e.client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId

	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	raw, _ := e.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetK8sApplication(request)
	})

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetK8sApplicationResponse)
	v := response.Applcation

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
