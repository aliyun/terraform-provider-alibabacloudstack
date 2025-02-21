package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	Hb_LAUNCHING            = "LAUNCHING"
	Hb_CREATING             = "CREATING"
	Hb_ACTIVATION           = "ACTIVATION"
	Hb_DELETING             = "DELETING"
	Hb_CREATE_FAILED        = "CREATE_FAILED"
	Hb_NODE_RESIZING        = "HBASE_SCALE_OUT"
	Hb_NODE_RESIZING_FAILED = "NODE_RESIZE_FAILED"
	Hb_DISK_RESIZING        = "HBASE_EXPANDING"
	Hb_DISK_RESIZE_FAILED   = "DISK_RESIZING_FAILED"
	Hb_LEVEL_MODIFY         = "INSTANCE_LEVEL_MODIFY"
	Hb_LEVEL_MODIFY_FAILED  = "INSTANCE_LEVEL_MODIFY_FAILED"
	Hb_HBASE_COLD_EXPANDING = "HBASE_COLD_EXPANDING"
)

type HBaseService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *HBaseService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := hbase.CreateUnTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.TagKey = &tagKey
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.UnTagResources(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*hbase.UnTagResourcesResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := hbase.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		raw, err := s.client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.TagResources(request)
		})
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*hbase.TagResourcesResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	//d.SetPartial("tags")
	return nil
}

func (s *HBaseService) diffTags(oldTags, newTags []hbase.TagResourcesTag) ([]hbase.TagResourcesTag, []hbase.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []hbase.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *HBaseService) tagsFromMap(m map[string]interface{}) []hbase.TagResourcesTag {
	result := make([]hbase.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, hbase.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *HBaseService) tagsToMap(tags []hbase.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *HBaseService) ignoreTag(t hbase.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Apsara Stack Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *HBaseService) DescribeHBaseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ClusterId": id,
		"PageSize":  PageSizeLarge,
		"PageNumber": 1,
	}
	response, err = s.client.DoTeaRequest("POST", "HBase", "2019-01-01", "DescribeInstance", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Hbase:Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *HBaseService) DescribeIpWhitelist(id string) (instance hbase.DescribeIpWhitelistResponse, err error) {
	request := hbase.CreateDescribeIpWhitelistRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ClusterId = id
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeIpWhitelist(request)
	})
	response, ok := raw.(*hbase.DescribeIpWhitelistResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"Instance.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) DescribeSecurityGroups(id string) (object hbase.DescribeSecurityGroupsResponse, err error) {
	request := hbase.CreateDescribeSecurityGroupsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeSecurityGroups(request)
	})
	response, ok := raw.(*hbase.DescribeSecurityGroupsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}

func (s *HBaseService) HBaseClusterStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHBaseInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *HBaseService) ModifyClusterDeletionProtection(clusterId string, protection bool) error {
	request := hbase.CreateModifyClusterDeletionProtectionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ClusterId = clusterId
	request.Protection = requests.NewBoolean(protection)
	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.ModifyClusterDeletionProtection(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*hbase.ModifyClusterDeletionProtectionResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, clusterId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *HBaseService) DescribeEndpoints(id string) (object hbase.DescribeEndpointsResponse, err error) {
	request := hbase.CreateDescribeEndpointsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeEndpoints(request)
	})
	response, ok := raw.(*hbase.DescribeEndpointsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	return *response, nil
}

func (s *HBaseService) DescribeClusterConnection(id string) (object hbase.DescribeClusterConnectionResponse, err error) {
	request := hbase.CreateDescribeClusterConnectionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ClusterId = id

	raw, err := s.client.WithHbaseClient(func(client *hbase.Client) (interface{}, error) {
		return client.DescribeClusterConnection(request)
	})
	response, ok := raw.(*hbase.DescribeClusterConnectionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *response, nil
}
