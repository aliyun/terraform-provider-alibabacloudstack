package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SlbService struct {
	client *connectivity.AlibabacloudStackClient
}

type SlbTag struct {
	TagKey   string
	TagValue string
}

type DescribeLoadBalancersNewResponse struct {
	*responses.BaseResponse
	RequestId     string           `json:"RequestId" xml:"RequestId"`
	PageNumber    int              `json:"PageNumber" xml:"PageNumber"`
	PageSize      int              `json:"PageSize" xml:"PageSize"`
	TotalCount    int              `json:"TotalCount" xml:"TotalCount"`
	LoadBalancers LoadBalancersNew `json:"LoadBalancers" xml:"LoadBalancers"`
}

type LoadBalancersNew struct {
	LoadBalancer []LoadBalancerNew `json:"LoadBalancer" xml:"LoadBalancer"`
}

type LoadBalancerTag struct {
	InstanceCount int    `json:"InstanceCount"`
	TagKey        string `json:"TagKey"`
	TagValue      string `json:"TagValue"`
}

type LoadBalancerNew struct {
	VpcId                        string                   `json:"VpcId" xml:"VpcId"`
	CreateTimeStamp              int64                    `json:"CreateTimeStamp" xml:"CreateTimeStamp"`
	LoadBalancerId               string                   `json:"LoadBalancerId" xml:"LoadBalancerId"`
	CreateTime                   string                   `json:"CreateTime" xml:"CreateTime"`
	PayType                      string                   `json:"PayType" xml:"PayType"`
	AddressType                  string                   `json:"AddressType" xml:"AddressType"`
	NetworkType                  string                   `json:"NetworkType" xml:"NetworkType"`
	ServiceManagedMode           string                   `json:"ServiceManagedMode" xml:"ServiceManagedMode"`
	SpecBpsFlag                  bool                     `json:"SpecBpsFlag" xml:"SpecBpsFlag"`
	AddressIPVersion             string                   `json:"AddressIPVersion" xml:"AddressIPVersion"`
	LoadBalancerName             string                   `json:"LoadBalancerName" xml:"LoadBalancerName"`
	Bandwidth                    int                      `json:"Bandwidth" xml:"Bandwidth"`
	Address                      string                   `json:"Address" xml:"Address"`
	SlaveZoneId                  string                   `json:"SlaveZoneId" xml:"SlaveZoneId"`
	MasterZoneId                 string                   `json:"MasterZoneId" xml:"MasterZoneId"`
	InternetChargeTypeAlias      string                   `json:"InternetChargeTypeAlias" xml:"InternetChargeTypeAlias"`
	LoadBalancerSpec             string                   `json:"LoadBalancerSpec" xml:"LoadBalancerSpec"`
	SpecType                     string                   `json:"SpecType" xml:"SpecType"`
	RegionId                     string                   `json:"RegionId" xml:"RegionId"`
	ModificationProtectionReason string                   `json:"ModificationProtectionReason" xml:"ModificationProtectionReason"`
	ModificationProtectionStatus string                   `json:"ModificationProtectionStatus" xml:"ModificationProtectionStatus"`
	VSwitchId                    string                   `json:"VSwitchId" xml:"VSwitchId"`
	LoadBalancerStatus           string                   `json:"LoadBalancerStatus" xml:"LoadBalancerStatus"`
	ResourceGroupId              string                   `json:"ResourceGroupId" xml:"ResourceGroupId"`
	InternetChargeType           string                   `json:"InternetChargeType" xml:"InternetChargeType"`
	BusinessStatus               string                   `json:"BusinessStatus" xml:"BusinessStatus"`
	DeleteProtection             string                   `json:"DeleteProtection" xml:"DeleteProtection"`
	RegionIdAlias                string                   `json:"RegionIdAlias" xml:"RegionIdAlias"`
	InstanceChargeType           string                   `json:"InstanceChargeType" xml:"InstanceChargeType"`
	ServiceManagedReason         string                   `json:"ServiceManagedReason" xml:"ServiceManagedReason"`
	IneffectiveOrderList         slb.IneffectiveOrderList `json:"IneffectiveOrderList" xml:"IneffectiveOrderList"`
	Tags                         []LoadBalancerTag        `json:"Tags" xml:"Tags"` // Changed to array format to match API response
}

const max_num_per_time = 50
const tags_max_num_per_time = 5
const tags_max_page_size = 50

func (s *SlbService) DoSlbDescribeloadbalancerattributeRequest(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	return s.DescribeSlb(id)
}
func (s *SlbService) DescribeSlb(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "InvalidLoadBalancerId.NotFound") {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Slb", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		} else {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		return response, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest)
	if bresponse.LoadBalancerId == "" {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Slb", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, err
}

func (s *SlbService) DescribeSlbRule(id string) (*slb.DescribeRuleAttributeResponse, error) {
	response := &slb.DescribeRuleAttributeResponse{}
	request := slb.CreateDescribeRuleAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RuleId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRuleAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeRuleAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "InvalidRuleId.NotFound") {
			return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbRule", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_rule", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *SlbService) DoSlbDescribevservergroupattributeRequest(id string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	return s.DescribeSlbServerGroup(id)
}
func (s *SlbService) DescribeSlbServerGroup(id string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	response := &slb.DescribeVServerGroupAttributeResponse{}
	request := slb.CreateDescribeVServerGroupAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.VServerGroupId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeVServerGroupAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "The specified VServerGroupId does not exist", "InvalidParameter") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_server_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.VServerGroupId == "" {
		return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbServerGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, err
}

func (s *SlbService) DescribeSlbMasterSlaveServerGroup(id string) (*slb.DescribeMasterSlaveServerGroupAttributeResponse, error) {
	response := &slb.DescribeMasterSlaveServerGroupAttributeResponse{}
	request := slb.CreateDescribeMasterSlaveServerGroupAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.MasterSlaveServerGroupId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeMasterSlaveServerGroupAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeMasterSlaveServerGroupAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "The specified MasterSlaveGroupId does not exist", "InvalidParameter") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultDebugMsg, "alibabacloudstack_slb_master_slave_server_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.MasterSlaveServerGroupId == "" {
		return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbMasterSlaveServerGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, err
}

func (s *SlbService) DescribeSlbBackendServer(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "InvalidLoadBalancerId.NotFound") {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbBackendServers", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		} else {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_backend_server", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		return response, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.LoadBalancerId == "" {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbBackendServers", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, err
}

func (s *SlbService) DescribeSlbListener(id string) (listener map[string]interface{}, err error) {
	parts, err := ParseSlbListenerId(id)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	protocol := parts[1]
	apiName := fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	request := s.client.NewCommonRequest("GET", "slb", "2014-05-15", apiName, "")
	request.QueryParams["LoadBalancerId"] = parts[0]
	port, _ := strconv.Atoi(parts[2])
	request.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		bresponse, err := s.client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw DescribeProjectMeta : %s", bresponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "The specified resource does not exist") {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR))
			} else if errmsgs.IsExpectedErrors(err, errmsgs.SlbIsBusy ...) {
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_slb_listener", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
			}
		}
		if err = json.Unmarshal(bresponse.GetHttpContentBytes(), &listener); err != nil {
			return resource.NonRetryableError(errmsgs.WrapError(err))
		}
		if port, ok := listener["ListenerPort"]; ok && port.(float64) > 0 {
			return nil
		} else {
			return resource.RetryableError(errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbListener", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR))
		}
	})

	return
}

func (s *SlbService) DoSlbDescribeaccesscontrollistattributeRequest(id string) (*slb.DescribeAccessControlListAttributeResponse, error) {
	return s.DescribeSlbAcl(id)
}
func (s *SlbService) DescribeSlbAcl(id string) (*slb.DescribeAccessControlListAttributeResponse, error) {
	response := &slb.DescribeAccessControlListAttributeResponse{}
	request := slb.CreateDescribeAccessControlListAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AclId = id

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlListAttribute(request)
	})
	bresponse, ok := raw.(*slb.DescribeAccessControlListAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "AclNotExist") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_acl", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *SlbService) WaitForSlbAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbAcl(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		} else {
			return nil
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.AclId, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *SlbService) WaitForSlb(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlb(id)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		} else if strings.EqualFold(object.LoadBalancerStatus, string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoadBalancerStatus, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbListener(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbListener(id)
		if err != nil && !errmsgs.IsExpectedErrors(err, "InvalidLoadBalancerId.NotFound") {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		gotStatus := ""
		if value, ok := object["Status"]; ok {
			gotStatus = strings.ToLower(value.(string))
		}
		if gotStatus == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, gotStatus, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		_, err := s.DescribeSlbRule(id)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if status != Deleted {
			break
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.VServerGroupId == id {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.VServerGroupId, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbMasterSlaveServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbMasterSlaveServerGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.MasterSlaveServerGroupId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.MasterSlaveServerGroupId, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitSlbAttribute(id string, instanceSet *schema.Set, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

RETRY:
	object, err := s.DescribeSlb(id)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if time.Now().After(deadline) {
		return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, id, errmsgs.ProviderERROR)
	}
	servers := object.BackendServers.BackendServer
	if len(servers) > 0 {
		for _, s := range servers {
			if instanceSet.Contains(s.ServerId) {
				goto RETRY
			}
		}
	}
	return nil
}

func (s *SlbService) slbRemoveAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateRemoveAccessControlListEntryRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.RemoveAccessControlListEntry(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if !errmsgs.IsExpectedErrors(err, "AclEntryEmpty") {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_remove_access_control_list_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *SlbService) SlbRemoveAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}

		slice := list[start:end]
		if err := s.slbRemoveAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

func (s *SlbService) slbAddAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateAddAccessControlListEntryRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddAccessControlListEntry(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_add_access_control_list_entry", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *SlbService) SlbAddAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}
		slice := list[start:end]
		if err := s.slbAddAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *SlbService) FlattenSlbAclEntryMappings(list []slb.AclEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"entry":   i.AclEntryIP,
			"comment": i.AclEntryComment,
		}
		result = append(result, l)
	}

	return result
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *SlbService) flattenSlbRelatedListenerMappings(list []slb.RelatedListener) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"load_balancer_id": i.LoadBalancerId,
			"protocol":         i.Protocol,
			"frontend_port":    i.ListenerPort,
			"acl_type":         i.AclType,
		}
		result = append(result, l)
	}

	return result
}

func (s *SlbService) DoSlbDescribecacertificatesRequest(id string) (*slb.CACertificate, error) {
	return s.DescribeSlbCACertificate(id)
}
func (s *SlbService) DescribeSlbCACertificate(id string) (*slb.CACertificate, error) {
	certificate := &slb.CACertificate{}
	request := slb.CreateDescribeCACertificatesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.CACertificateId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return certificate, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_ca_certificate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeCACertificatesResponse)
	if len(response.CACertificates.CACertificate) < 1 {
		return certificate, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbCACertificate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return &response.CACertificates.CACertificate[0], nil
}

func (s *SlbService) WaitForSlbCACertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbCACertificate(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		} else {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.CACertificateId, id, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *SlbService) DoSlbDescribeservercertificateRequest(id string) (*slb.ServerCertificate, error) {
	return s.DescribeSlbServerCertificate(id)
}
func (s *SlbService) DescribeSlbServerCertificate(id string) (*slb.ServerCertificate, error) {
	certificate := &slb.ServerCertificate{}
	request := slb.CreateDescribeServerCertificatesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ServerCertificateId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return certificate, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_server_certificate", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeServerCertificatesResponse)
	if len(response.ServerCertificates.ServerCertificate) < 1 || response.ServerCertificates.ServerCertificate[0].ServerCertificateId != id {
		return certificate, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SlbServerCertificate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &response.ServerCertificates.ServerCertificate[0], nil
}

func (s *SlbService) WaitForSlbServerCertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerCertificate(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.ServerCertificateId == id {
			break
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ServerCertificateId, id, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func toSlbTagsString(tags []Tag) string {
	slbTags := make([]SlbTag, 0, len(tags))

	for _, tag := range tags {
		slbTag := SlbTag{
			TagKey:   tag.Key,
			TagValue: tag.Value,
		}
		slbTags = append(slbTags, slbTag)
	}

	b, _ := json.Marshal(slbTags)

	return string(b)
}

func (s *SlbService) DoSlbDescribedomainextensionattributeRequest(domainExtensionId string) (*slb.DescribeDomainExtensionAttributeResponse, error) {
	return s.DescribeDomainExtensionAttribute(domainExtensionId)
}
func (s *SlbService) DescribeDomainExtensionAttribute(domainExtensionId string) (*slb.DescribeDomainExtensionAttributeResponse, error) {
	response := &slb.DescribeDomainExtensionAttributeResponse{}
	request := slb.CreateDescribeDomainExtensionAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DomainExtensionId = domainExtensionId
	var raw interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeDomainExtensionAttribute(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.AlibabacloudStackGoClientFailure, "ServiceUnavailable", errmsgs.Throttling) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	bresponse, ok := raw.(*slb.DescribeDomainExtensionAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "InvalidParameter.DomainExtensionId", "InvalidParameter") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_domain_extension", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if bresponse.DomainExtensionId != domainExtensionId {
		return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("SLBDomainExtension", domainExtensionId)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, nil
}

func (s *SlbService) WaitForSlbDomainExtension(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeDomainExtensionAttribute(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), errmsgs.ProviderERROR)
		}
	}
}

func (s *SlbService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.TagKey)
		}
		request := slb.CreateRemoveTagsRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.LoadBalancerId = d.Id()
		bytes, _ := json.Marshal(remove)
		s2 := string(bytes)
		request.Tags = fmt.Sprint(s2)
		request.RegionId = s.client.RegionId

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.RemoveTags(request)
			})
			bresponse, ok := raw.(*responses.CommonResponse)
			if err != nil {
				if errmsgs.IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_remove_tags", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return err
		}
	}

	if len(create) > 0 {
		request := slb.CreateAddTagsRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.LoadBalancerId = d.Id()
		bytes, _ := json.Marshal(create)
		s2 := string(bytes)
		request.Tags = fmt.Sprint(s2)
		request.RegionId = s.client.RegionId

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.AddTags(request)
			})
			bresponse, ok := raw.(*responses.CommonResponse)
			if err != nil {
				if errmsgs.IsThrottling(err) {
					wait()
					return resource.RetryableError(err)
				}
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_add_tags", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SlbService) tagsToMap(tags []slb.TagSet) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *SlbService) ignoreTag(t slb.TagSet) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *SlbService) diffTags(oldTags, newTags []slb.TagSet) ([]slb.TagSet, []slb.TagSet) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.TagKey] = t.TagValue
	}

	// Build the list of what to remove
	var remove []slb.TagSet
	for _, t := range oldTags {
		old, ok := create[t.TagKey]
		if !ok || old != t.TagValue {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *SlbService) tagsFromMap(m map[string]interface{}) []slb.TagSet {
	result := make([]slb.TagSet, 0, len(m))
	for k, v := range m {
		result = append(result, slb.TagSet{
			TagKey:   k,
			TagValue: v.(string),
		})
	}

	return result
}

func (s *SlbService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []slb.TagSet, err error) {
	request := slb.CreateDescribeTagsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = resourceId
	if len(resourceTags) > 0 {
		var reqTags []slb.TagSet
		for key, value := range resourceTags {
			reqTags = append(reqTags, slb.TagSet{
				TagKey:   key,
				TagValue: value.(string),
			})
		}
		bytes, _ := json.Marshal(reqTags)
		s2 := string(bytes)
		request.Tags = fmt.Sprint(s2)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithSlbClient(func(Client *slb.Client) (interface{}, error) {
			return Client.DescribeTags(request)
		})
		bresponse, ok := raw.(*slb.DescribeTagsResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_describe_tags", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		tags = bresponse.TagSets.TagSet
		return nil
	})
	if err != nil {
		return nil, err
	}

	return
}

func (s *SlbService) SetAccessLogsDownloadAttribute(logs_attr_str string, load_balancer_id string) error {
	request := s.client.NewCommonRequest("POST", "Slb", "2014-05-15", "SetAccessLogsDownloadAttribute", "")
	request.QueryParams = map[string]string{
		"LogsDownloadAttributes": logs_attr_str,
		"loadBalancerId":         load_balancer_id,
	}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_access_log", "SetAccessLogsDownloadAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func (s *SlbService) DeleteAccessLogsDownloadAttribute(id string) error {
	request := s.client.NewCommonRequest("POST", "Slb", "2014-05-15", "DeleteAccessLogsDownloadAttribute", "")
	parts := strings.Split(id, ":")
	logs_download_attributes := fmt.Sprintf("[{\"LoadBalancerId\":\"%s\",}]", parts[0])
	request.QueryParams = map[string]string{
		"LogsDownloadAttributes": logs_download_attributes,
		"loadBalancerId":         parts[0],
	}
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_access_log", "DeleteAccessLogsDownloadAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	return nil
}

type Slblogsdownloadattribute struct {
	LogProject     string `json:"LogProject"`
	LogStore       string `json:"LogStore"`
	LoadBalancerId string `json:"LoadBalancerId"`
	LogType        string `json:"LogType"`
	RoleArn        string `json:"RoleArn"`
	Region         string `json:"Region"`
}

type SlbDescribeaccesslogsdownloadattributeResponse struct {
	LogsDownloadAttributes struct {
		LogsDownloadAttribute []Slblogsdownloadattribute `json:"LogsDownloadAttribute"`
	} `json:"LogsDownloadAttributes"`
	RequestId  string `json:"RequestId"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
	TotalCount int    `json:"TotalCount"`
	Count      int    `json:"Count"`
}

func (s *SlbService) DescribeAccessLogsDownloadAttributes(id string) (logsattr []Slblogsdownloadattribute, err error) {
	request := s.client.NewCommonRequest("POST", "Slb", "2014-05-15", "DescribeAccessLogsDownloadAttribute", "")
	request.QueryParams = map[string]string{
		"LoadBalancerId": id,
		"LogType":        "layer7",
	}
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_slb_access_log", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if !bresponse.IsSuccess() {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_slb_access_log", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	slblogs := &SlbDescribeaccesslogsdownloadattributeResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &slblogs)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	logsattr = []Slblogsdownloadattribute{}
	for _, logattr  :=  range slblogs.LogsDownloadAttributes.LogsDownloadAttribute {
		if logattr.LoadBalancerId == id {
			logsattr = append(logsattr, logattr)
		}
	}
	return logsattr, nil
}

func (s *SlbService) DescribeAccessLogsDownloadAttribute(id string) (logsattr *Slblogsdownloadattribute, err error) {
	parts := strings.Split(id, ":")
	logsDownloadAttributes, err := s.DescribeAccessLogsDownloadAttributes(parts[0])
	if err != nil {
		return logsattr, nil
	}
	if len(logsDownloadAttributes) > 0 {
		for _, log := range logsDownloadAttributes {
			if log.LogProject == parts[1] && log.LogStore == parts[2] {
				logsattr = &log
			}
		}

	}
	return logsattr, nil
}
