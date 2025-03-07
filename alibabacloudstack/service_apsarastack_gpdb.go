package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type GpdbService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *GpdbService) GpdbAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbAccount(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AccountStatus"]) == failState {
				return object, fmt.Sprint(object["AccountStatus"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["AccountStatus"])))
			}
		}
		return object, fmt.Sprint(object["AccountStatus"]), nil
	}
}

func (s *GpdbService) DescribeGpdbAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := make(map[string]interface{})
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request["AccountName"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "gpdb", "2016-05-03", "DescribeAccounts", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GPDB:Account", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err	
	}
	v, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Accounts.DBInstanceAccount", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GPDB", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AccountName"]) != parts[1] {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GPDB", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *GpdbService) DescribeGpdbInstance(id string) (instanceAttribute gpdb.DBInstanceAttribute, err error) {
	request := gpdb.CreateDescribeDBInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.DescribeDBInstanceAttribute(request)
	})

	response, ok := raw.(*gpdb.DescribeDBInstanceAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		} else {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		return
	}

	addDebug(request.GetActionName(), response, request.RpcRequest, request)
	if len(response.Items.DBInstanceAttribute) == 0 {
		return instanceAttribute, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Gpdb Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return response.Items.DBInstanceAttribute[0], nil
}

func (s *GpdbService) DescribeGpdbSecurityIps(id string) (ips []string, err error) {
	request := gpdb.CreateDescribeDBInstanceIPArrayListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id

	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.DescribeDBInstanceIPArrayList(request)
	})
	response, ok := raw.(*gpdb.DescribeDBInstanceIPArrayListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			err = errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		} else {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		return
	}
	addDebug(request.GetActionName(), response, request.RpcRequest, request)
	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.Items.DBInstanceIPArray {
		if ip.DBInstanceIPArrayAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
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

func (s *GpdbService) ModifyGpdbSecurityIps(id, ips string) error {
	request := gpdb.CreateModifySecurityIpsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	request.SecurityIPList = ips
	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	response, ok := raw.(*gpdb.ModifySecurityIpsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), response, request.RpcRequest, request)

	return nil
}

func (s *GpdbService) DescribeGpdbConnection(id string) (*gpdb.DBInstanceNetInfo, error) {
	info := &gpdb.DBInstanceNetInfo{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return info, errmsgs.WrapError(err)
	}

	// Describe DB Instance Net Info
	request := gpdb.CreateDescribeDBInstanceNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	raw, err := s.client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
		return gpdbClient.DescribeDBInstanceNetInfo(request)
	})
	response, ok := raw.(*gpdb.DescribeDBInstanceNetInfoResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return info, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return info, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), response, request.RpcRequest, request)
	if response.DBInstanceNetInfos.DBInstanceNetInfo != nil {
		for _, o := range response.DBInstanceNetInfos.DBInstanceNetInfo {
			if strings.HasPrefix(o.ConnectionString, parts[1]) {
				return &o, nil
			}
		}
	}

	return info, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GpdbConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *GpdbService) GpdbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *GpdbService) WaitForGpdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeGpdbConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.ConnectionString != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *GpdbService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffGpdbTags(gpdbTagsFromMap(o), gpdbTagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := gpdb.CreateUntagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.TagKey = &tagKey
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(TagResourceInstance)
		raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		response, ok := raw.(*gpdb.UntagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), response, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := gpdb.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(TagResourceInstance)
		raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		response, ok := raw.(*gpdb.TagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), response, request.RpcRequest, request)
	}

	//d.SetPartial("tags")
	return nil
}

func (s *GpdbService) tagsToMap(tags []gpdb.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *GpdbService) ignoreTag(t gpdb.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}
