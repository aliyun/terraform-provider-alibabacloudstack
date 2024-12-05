package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type CloudApiService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CloudApiService) DoCloudapiDescribeapigroupRequest(id string) (*cloudapi.DescribeApiGroupResponse, error) {
	return s.DescribeApiGatewayGroup(id)
}

func (s *CloudApiService) DescribeApiGatewayGroup(id string) (*cloudapi.DescribeApiGroupResponse, error) {
	request := cloudapi.CreateDescribeApiGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.GroupId = id

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApiGroup(request)
	})
	bresponse, ok := raw.(*cloudapi.DescribeApiGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApiGroup"}) {
			return bresponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return bresponse, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.GroupId == "" {
		return bresponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ApiGatewayGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, nil
}

func (s *CloudApiService) WaitForApiGatewayGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if string(object.GroupId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, string(object.GroupId), id, errmsgs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApp(id string) (*cloudapi.DescribeAppResponse, error) {
	request := cloudapi.CreateDescribeAppRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.AppId = requests.Integer(id)

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApp(request)
	})
	bresponse, ok := raw.(*cloudapi.DescribeAppResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApp"}) {
			return bresponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return bresponse, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *CloudApiService) WaitForApiGatewayApp(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApp(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if fmt.Sprint(object.AppId) == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object.AppId), id, errmsgs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayApi(id string) (*cloudapi.DescribeApiResponse, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	request := cloudapi.CreateDescribeApiRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ApiId = parts[1]
	request.GroupId = parts[0]

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeApi(request)
	})
	bresponse, ok := raw.(*cloudapi.DescribeApiResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
			return bresponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return bresponse, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse.ApiId == "" {
		return bresponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ApiGatewayApi", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return bresponse, nil
}

func (s *CloudApiService) WaitForApiGatewayApi(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeApiGatewayApi(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errmsgs.WrapError(err)
		}
		respId := fmt.Sprintf("%s%s%s", object.GroupId, COLON_SEPARATED, object.ApiId)
		if respId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, respId, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *CloudApiService) DescribeApiGatewayAppAttachment(id string) (*cloudapi.AuthorizedApp, error) {
	app := &cloudapi.AuthorizedApp{}
	request := cloudapi.CreateDescribeAuthorizedAppsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return app, errmsgs.WrapError(err)
	}
	request.GroupId = parts[0]
	request.ApiId = parts[1]
	request.StageName = parts[3]
	appId, _ := strconv.ParseInt(parts[2], 10, 64)

	var allApps []cloudapi.AuthorizedApp

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeAuthorizedApps(request)
		})
		bresponse, ok := raw.(*cloudapi.DescribeAuthorizedAppsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi"}) {
				return app, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return app, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		allApps = append(allApps, bresponse.AuthorizedApps.AuthorizedApp...)

		if len(allApps) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return app, errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredAppsTemp []cloudapi.AuthorizedApp
	for _, app := range allApps {
		if app.AppId == appId {
			filteredAppsTemp = append(filteredAppsTemp, app)
		}
	}

	if len(filteredAppsTemp) < 1 {
		return app, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ApigatewayAppAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	app = &filteredAppsTemp[0]
	return app, nil
}

func (s *CloudApiService) DescribeApiGatewayVpcAccess(id string) (*cloudapi.VpcAccessAttribute, error) {
	vpc := &cloudapi.VpcAccessAttribute{}
	request := cloudapi.CreateDescribeVpcAccessesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return vpc, errmsgs.WrapError(err)
	}
	var allVpcs []cloudapi.VpcAccessAttribute

	for {
		raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeVpcAccesses(request)
		})
		bresponse, ok := raw.(*cloudapi.DescribeVpcAccessesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return vpc, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		allVpcs = append(allVpcs, bresponse.VpcAccessAttributes.VpcAccessAttribute...)

		if len(allVpcs) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return vpc, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredVpcsTemp []cloudapi.VpcAccessAttribute
	for _, vpc := range allVpcs {
		iPort, _ := strconv.Atoi(parts[3])
		if vpc.Port == iPort && vpc.InstanceId == parts[2] && vpc.VpcId == parts[1] && vpc.Name == parts[0] {
			filteredVpcsTemp = append(filteredVpcsTemp, vpc)
		}
	}

	if len(filteredVpcsTemp) < 1 {
		return vpc, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ApiGatewayVpcAccess", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return &filteredVpcsTemp[0], nil
}

func (s *CloudApiService) WaitForApiGatewayAppAttachment(id string, status Status, timeout int) (err error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	appIds := parts[2]
	for {
		object, err := s.DescribeApiGatewayAppAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if strconv.FormatInt(object.AppId, 10) == appIds && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatInt(object.AppId, 10), appIds, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CloudApiService) DescribeDeployedApi(id string, stageName string) (*cloudapi.DescribeDeployedApiResponse, error) {
	request := cloudapi.CreateDescribeDeployedApiRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DescribeDeployedApi(request)
	})
	bresponse, ok := raw.(*cloudapi.DescribeDeployedApiResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return bresponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return bresponse, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return bresponse, nil
}

func (s *CloudApiService) DeployedApi(id string, stageName string) (err error) {
	request := cloudapi.CreateDeployApiRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName
	request.Description = DeployCommonDescription

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.DeployApi(request)
	})
	bresponse, ok := raw.(*cloudapi.DeployApiResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) AbolishApi(id string, stageName string) (err error) {
	request := cloudapi.CreateAbolishApiRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.ApiId = parts[1]
	request.GroupId = parts[0]
	request.StageName = stageName

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.AbolishApi(request)
	})
	bresponse, ok := raw.(*cloudapi.AbolishApiResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"NotFoundApiGroup", "NotFoundApi", "NotFoundStage"}) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return
}

func (s *CloudApiService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []cloudapi.TagResource, err error) {
	request := cloudapi.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}

	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []cloudapi.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, cloudapi.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}

	raw, err := s.client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.ListTagResources(request)
	})
	bresponse, ok := raw.(*cloudapi.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, resourceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	tags = bresponse.TagResources.TagResource
	return
}

func (s *CloudApiService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := cloudapi.CreateUntagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(resourceType)
		request.TagKey = &tagKey

		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		bresponse, ok := raw.(*cloudapi.UntagResourcesResponse)
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
		request := cloudapi.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(resourceType)

		raw, err := s.client.WithCloudApiClient(func(client *cloudapi.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		bresponse, ok := raw.(*cloudapi.TagResourcesResponse)
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

	return nil
}

func (s *CloudApiService) tagsToMap(tags []cloudapi.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *CloudApiService) ignoreTag(t cloudapi.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Apsara Stack Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *CloudApiService) diffTags(oldTags, newTags []cloudapi.TagResourcesTag) ([]cloudapi.TagResourcesTag, []cloudapi.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []cloudapi.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *CloudApiService) tagsFromMap(m map[string]interface{}) []cloudapi.TagResourcesTag {
	result := make([]cloudapi.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, cloudapi.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}
