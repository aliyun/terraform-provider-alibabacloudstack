package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type KvstoreService struct {
	client *connectivity.AlibabacloudStackClient
}

var KVstoreInstanceStatusCatcher = Catcher{"OperationDenied.KVstoreInstanceStatus", 60, 5}

func (s *KvstoreService) DoR_KvstoreDescribeinstanceattributeRequest(id string) (*r_kvstore.DBInstanceAttribute, error) {
	return s.DescribeKVstoreInstance(id)
}
func (s *KvstoreService) DescribeKVstoreInstance(id string) (*r_kvstore.DBInstanceAttribute, error) {
	instance := &r_kvstore.DBInstanceAttribute{}
	request := r_kvstore.CreateDescribeInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeInstanceAttribute(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeInstanceAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KVstoreInstance", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.Instances.DBInstanceAttribute) <= 0 {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KVstoreInstance", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return &bresponse.Instances.DBInstanceAttribute[0], nil
}

func (s *KvstoreService) DescribeKVstoreBackupPolicy(id string) (*r_kvstore.DescribeBackupPolicyResponse, error) {
	response := &r_kvstore.DescribeBackupPolicyResponse{}
	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeBackupPolicy(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeBackupPolicyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KVstoreBackupPolicy", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *KvstoreService) WaitForKVstoreInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.InstanceStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) RdsKvstoreInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.InstanceStatus == failState {
				return object, object.InstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.InstanceStatus))
			}
		}
		return object, object.InstanceStatus, nil
	}
}

func (s *KvstoreService) WaitForKVstoreInstanceVpcAuthMode(id string, status string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreInstance(id)
		if err != nil && !errmsgs.NotFoundError(err) {
			return err
		}
		if object.VpcAuthMode == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.VpcAuthMode, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DescribeParameters(id string) (*r_kvstore.DescribeParametersResponse, error) {
	response := &r_kvstore.DescribeParametersResponse{}
	request := r_kvstore.CreateDescribeParametersRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id

	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeParameters(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeParametersResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return response, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Parameters", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *KvstoreService) ModifyInstanceConfig(id string, config string) error {
	request := r_kvstore.CreateModifyInstanceConfigRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	request.Config = config

	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ModifyInstanceConfig(request)
	})
	bresponse, ok := raw.(*r_kvstore.ModifyInstanceConfigResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *KvstoreService) setInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

		if len(remove) > 0 {
			var tagKey []string
			for _, v := range remove {
				tagKey = append(tagKey, v.Key)
			}
			request := r_kvstore.CreateUntagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = strings.ToUpper(string(TagResourceInstance))
			request.TagKey = &tagKey
			raw, err := s.client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			bresponse, ok := raw.(*r_kvstore.UntagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := r_kvstore.CreateTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = strings.ToUpper(string(TagResourceInstance))
			raw, err := s.client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			bresponse, ok := raw.(*r_kvstore.TagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		//d.SetPartial("tags")
	}

	return nil
}

func (s *KvstoreService) tagsToMap(tags []r_kvstore.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *KvstoreService) tagsFromMap(m map[string]interface{}) []r_kvstore.TagResourcesTag {
	result := make([]r_kvstore.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, r_kvstore.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *KvstoreService) ignoreTag(t r_kvstore.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagValue)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *KvstoreService) diffTags(oldTags, newTags []r_kvstore.TagResourcesTag) ([]r_kvstore.TagResourcesTag, []r_kvstore.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []r_kvstore.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *KvstoreService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []r_kvstore.TagResource, err error) {
	request := r_kvstore.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = strings.ToUpper(string(resourceType))
	request.ResourceId = &[]string{resourceId}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.ListTagResources(request)
	})
	bresponse, ok := raw.(*r_kvstore.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse.TagResources.TagResource, nil
}

func (s *KvstoreService) WaitForKVstoreAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeKVstoreAccount(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *KvstoreService) DoR_KvstoreDescribeaccountsRequest(id string) (*r_kvstore.Account, error) {
	return s.DescribeKVstoreAccount(id)
}
func (s *KvstoreService) DescribeKVstoreAccount(id string) (*r_kvstore.Account, error) {
	ds := &r_kvstore.Account{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return ds, errmsgs.WrapError(err)
	}
	request := r_kvstore.CreateDescribeAccountsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(KVstoreInstanceStatusCatcher)
	var raw interface{}
	err = invoker.Run(func() error {
		raw, err = s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeAccounts(request)
		})

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return err
	})
	response, ok := raw.(*r_kvstore.DescribeAccountsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return ds, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return ds, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if len(response.Accounts.Account) < 1 {
		return ds, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KVstoreAccount", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return &response.Accounts.Account[0], nil
}

func (s *KvstoreService) DescribeKVstoreSecurityGroupId(id string) (*r_kvstore.DescribeSecurityGroupConfigurationResponse, error) {
	response := &r_kvstore.DescribeSecurityGroupConfigurationResponse{}
	request := r_kvstore.CreateDescribeSecurityGroupConfigurationRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return response, errmsgs.WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeSecurityGroupConfiguration(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeSecurityGroupConfigurationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *KvstoreService) DescribeDBInstanceNetInfo(id string) (*r_kvstore.NetInfoItemsInDescribeDBInstanceNetInfo, error) {
	response := &r_kvstore.DescribeDBInstanceNetInfoResponse{}
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id
	if err := s.WaitForKVstoreInstance(id, Normal, DefaultLongTimeout); err != nil {
		return &response.NetInfoItems, errmsgs.WrapError(err)
	}
	raw, err := s.client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeDBInstanceNetInfo(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return &response.NetInfoItems, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return &bresponse.NetInfoItems, nil
}

func (s *KvstoreService) DoR_KvstoreDescribedbinstancenetinfoRequest(id string) (object []r_kvstore.InstanceNetInfo, err error) {
	return s.DescribeKvstoreConnection(id)
}
func (s *KvstoreService) DescribeKvstoreConnection(id string) (object []r_kvstore.InstanceNetInfo, err error) {
	request := r_kvstore.CreateDescribeDBInstanceNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id

	raw, err := s.client.WithRkvClient(func(r_kvstoreClient *r_kvstore.Client) (interface{}, error) {
		return r_kvstoreClient.DescribeDBInstanceNetInfo(request)
	})
	bresponse, ok := raw.(*r_kvstore.DescribeDBInstanceNetInfoResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KvstoreConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.NetInfoItems.InstanceNetInfo) < 1 {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("KvstoreConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, bresponse.RequestId)
		return
	}
	return bresponse.NetInfoItems.InstanceNetInfo, nil
}

func (s *KvstoreService) InstanceSslStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstanceSSL(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["SSLEnabled"].(string) == failState {
				return object, object["SSLEnabled"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["SSLEnabled"].(string)))
			}
		}
		return object, object["SSLEnabled"].(string), nil
	}
}

func (s *KvstoreService) DescribeInstanceSSL(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("POST", "R-kvstore", "2015-01-01", "DescribeInstanceSSL", "")
	request.QueryParams["InstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeInstanceSSL", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

// DescribeInstanceTDEStatus

func (s *KvstoreService) InstanceTDEStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInstanceTDEStatus(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["TDEStatus"].(string) == failState {
				return object, object["TDEStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *KvstoreService) DescribeInstanceTDEStatus(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("POST", "R-kvstore", "2015-01-01", "DescribeInstanceTDEStatus", "")
	request.QueryParams["InstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeInstanceTDEStatus", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

type GetKVInstanceClassResponse struct {
	*responses.BaseResponse
	Code      any               `json:"Code"`
	Message   string            `json:"Message"`
	RequestId string            `json:"RequestId"`
	Success   bool              `json:"Success"`
	Data      []KVInstanceClass `json:"data"`
}

type KVInstanceClass struct {
	Architecture   string `json:"architecture"`
	Cpu            int    `json:"cpu"`
	EngineVersion  string `json:"engineVersion"`
	InstanceClass  string `json:"instanceClass"`
	MaxBandWidth   int    `json:"maxBandWidth"`
	MaxConnections int    `json:"maxConnections"`
	Memory         int    `json:"memory"`
	NodeType       string `json:"nodeType"`
	Product        string `json:"product"`
	Series         string `json:"series"`
	Status         string `json:"status"`
}
