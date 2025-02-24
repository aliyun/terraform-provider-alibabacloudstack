package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AdbService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AdbService) DescribeAdbCluster(id string) (instance *adb.DBClusterInDescribeDBClusters, err error) {
	request := adb.CreateDescribeDBClustersRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterIds = id
	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusters(request)
	})
	bresponse, ok := raw.(*adb.DescribeDBClustersResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.Items.DBCluster) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Cluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &bresponse.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbClusterAttribute(id string) (instance *adb.DBCluster, err error) {
	request := adb.CreateDescribeDBClusterAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAttribute(request)
	})
	bresponse, ok := raw.(*adb.DescribeDBClusterAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.Items.DBCluster) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Cluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &bresponse.Items.DBCluster[0], nil
}

func (s *AdbService) DescribeAdbAutoRenewAttribute(id string) (instance *adb.AutoRenewAttribute, err error) {
	request := adb.CreateDescribeAutoRenewAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterIds = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeAutoRenewAttribute(request)
	})
	bresponse, ok := raw.(*adb.DescribeAutoRenewAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(bresponse.Items.AutoRenewAttribute) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Cluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &bresponse.Items.AutoRenewAttribute[0], nil
}

func (s *AdbService) WaitForAdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if status != Deleted && object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *AdbService) DoAdbDescribedbclusternetinfoRequest(id string) (*adb.Address, error) {
	return s.DescribeAdbConnection(id)
}

func (s *AdbService) DescribeAdbConnection(id string) (*adb.Address, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(DefaultIntervalLong) * time.Second)
	for {
		object, err := s.DescribeAdbClusterNetInfo(parts[0])

		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return nil, errmsgs.WrapError(err)
		}

		if object != nil {
			for _, p := range object {
				if p.NetType == "Public" {
					return &p, nil
				}
			}
		}
		time.Sleep(DefaultIntervalMini * time.Second)
		if time.Now().After(deadline) {
			break
		}
	}

	return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *AdbService) DescribeAdbClusterNetInfo(id string) ([]adb.Address, error) {
	request := adb.CreateDescribeDBClusterNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterNetInfo(request)
	})
	bresponse, ok := raw.(*adb.DescribeDBClusterNetInfoResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(bresponse.Items.Address) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return bresponse.Items.Address, nil
}

func (s *AdbService) DescribeAdbClusterNetInfo2(id string) (address adb.Address, err error) {
	request := adb.CreateDescribeDBClusterNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterNetInfo(request)
	})
	bresponse, ok := raw.(*adb.DescribeDBClusterNetInfoResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return address, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return address, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(bresponse.Items.Address) < 1 {
		return address, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return bresponse.Items.Address[0], nil
}

func (s *AdbService) WaitForAdbAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbAccount(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.AccountStatus == string(status) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DoAdbDescribeaccountsRequest(id string) (ds *adb.DBAccount, err error) {
	return s.DescribeAdbAccount(id)
}

func (s *AdbService) DescribeAdbAccount(id string) (ds *adb.DBAccount, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := adb.CreateDescribeAccountsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]

	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var raw interface{}
	err = invoker.Run(func() error {
		raw, err = s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.DescribeAccounts(request)
		})
		

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	bresponse, ok := raw.(*adb.DescribeAccountsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	
	if len(bresponse.AccountList.DBAccount) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBAccount", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &bresponse.AccountList.DBAccount[0], nil
}

func (s *AdbService) WaitForAdbInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbCluster(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) setClusterTags(d *schema.ResourceData) error {
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
			request := adb.CreateUntagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "cluster"
			request.TagKey = &tagKey

			raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if len(create) > 0 {
			request := adb.CreateTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &create
			request.ResourceType = "cluster"

			raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}

	return nil
}

func (s *AdbService) diffTags(oldTags, newTags []adb.TagResourcesTag) ([]adb.TagResourcesTag, []adb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []adb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *AdbService) tagsToMap(tags []adb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *AdbService) tagsFromMap(m map[string]interface{}) []adb.TagResourcesTag {
	result := make([]adb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, adb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *AdbService) ignoreTag(t adb.TagResource) bool {
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

func (s *AdbService) DescribeTags(resourceId string, resourceType TagResourceType) (tags []adb.TagResource, err error) {
	request := adb.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}

	raw, err := s.client.WithAdbClient(func(client *adb.Client) (interface{}, error) {
		return client.ListTagResources(request)
	})
	bresponse, ok := raw.(*adb.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return bresponse.TagResources.TagResource, nil
}

func (s *AdbService) WaitForCluster(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if strings.ToLower(object.DBClusterStatus) == strings.ToLower(string(status)) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBClusterStatus, status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *AdbService) DescribeDBSecurityIps(clusterId string) (ips []string, err error) {
	request := adb.CreateDescribeDBClusterAccessWhiteListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = clusterId

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeDBClusterAccessWhiteList(request)
	})
	resp, ok := raw.(*adb.DescribeDBClusterAccessWhiteListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, clusterId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range resp.Items.IPArray {
		if ip.DBClusterIPArrayAttribute != "hidden" {
			ipstr += separator + ip.SecurityIPList
			separator = COMMA_SEPARATED
		}
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *AdbService) ModifyDBSecurityIps(clusterId, ips string) error {
	request := adb.CreateModifyDBClusterAccessWhiteListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = clusterId
	request.SecurityIps = ips

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyDBClusterAccessWhiteList(request)
	})
	bresponse, ok := raw.(*adb.ModifyDBClusterAccessWhiteListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, clusterId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *AdbService) DescribeAdbBackupPolicy(id string) (policy *adb.DescribeBackupPolicyResponse, err error) {
	request := adb.CreateDescribeBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = id

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeBackupPolicy(request)
	})
	bresponse, ok := raw.(*adb.DescribeBackupPolicyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return policy, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return bresponse, nil
}

func (s *AdbService) ModifyAdbBackupPolicy(clusterId, backupTime, backupPeriod string) error {
	request := adb.CreateModifyBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = clusterId
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyBackupPolicy(request)
	})
	bresponse, ok := raw.(*adb.ModifyBackupPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, clusterId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForCluster(clusterId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *AdbService) AdbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbClusterAttribute(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBClusterStatus == failState {
				return object, object.DBClusterStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.DBClusterStatus))
			}
		}
		return object, object.DBClusterStatus, nil
	}
}

func (s *AdbService) DescribeTask(id, taskId string) (*adb.DescribeTaskInfoResponse, error) {
	request := adb.CreateDescribeTaskInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = id
	request.TaskId = requests.Integer(taskId)

	raw, err := s.client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DescribeTaskInfo(request)
	})
	bresponse, ok := raw.(*adb.DescribeTaskInfoResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return bresponse, nil
}

func (s *AdbService) AdbTaskStateRefreshFunc(id, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTask(id, taskId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		return object, object.TaskInfo.Status, nil
	}
}

func (s *AdbService) DescribeAutoRenewAttribute(id string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["DBClusterIds"] = id
	response, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "DescribeAutoRenewAttribute", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items.AutoRenewAttribute", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.AutoRenewAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusterAccessWhiteList(id string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["DBClusterId"] = id
	response, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "DescribeDBClusterAccessWhiteList", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticdbForMysql3.0DbCluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items.IPArray", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.IPArray", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		ipList := ""
		for _, item := range v.([]interface{}) {
			if item.(map[string]interface{})["DBClusterIPArrayAttribute"] == "hidden" {
				continue
			}
			ipList += item.(map[string]interface{})["SecurityIPList"].(string) + ","
		}
		v.([]interface{})[0].(map[string]interface{})["SecurityIPList"] = strings.TrimSuffix(ipList, ",")
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		if len(removed) > 0 {
			request := make(map[string]interface{})
			request["ResourceType"] = resourceType
			request["ResourceId.1"] = d.Id()
			for i, key := range removed {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "UntagResources", "", nil, nil, request)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UntagResources", errmsgs.AlibabacloudStackSdkGoERROR)
			}
		}
		if len(added) > 0 {
			request := make(map[string]interface{})
			request["ResourceType"] = resourceType
			request["ResourceId.1"] = d.Id()
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			_, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "TagResources", "", nil, nil, request)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *AdbService) DoAdbDescribebackuppolicyRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeAdbDbCluster(id)
}

func (s *AdbService) DescribeAdbDbCluster(id string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["DBClusterId"] = id
	response, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "DescribeDBClusterAttribute", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFoundError"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AdbDbCluster", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) DescribeDBClusters(id string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["DBClusterIds"] = id
	response, err := s.client.DoTeaRequest("POST", "ADB", "2019-03-15", "DescribeDBClusters", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	addDebug("DescribeDBClusters", response, request)
	v, err := jsonpath.Get("$.Items.DBCluster", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.DBCluster", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DBClusterId"].(string) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AnalyticDBForMySQL3.0", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *AdbService) AdbDbClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAdbDbCluster(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["DBClusterStatus"].(string) == failState {
				return object, object["DBClusterStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["DBClusterStatus"].(string)))
			}
		}
		return object, object["DBClusterStatus"].(string), nil
	}
}
