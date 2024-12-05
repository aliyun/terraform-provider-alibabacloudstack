package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ElasticsearchService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeInstance"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	request["product"] = "elasticsearch"
	request["OrganizationId"] = s.client.Department
	request["ResourceId"] = s.client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, request, &runtime)
	addDebug(action, response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if (object["instanceId"].(string)) != id {
		return object, WrapErrorf(Error(GetNotFoundMessage("Elasticsearch Instance", id)), NotFoundWithResponse, response)
	}
	return object, nil

}

func (s *ElasticsearchService) DescribeElasticsearchOnk8sInstance(id string) (object map[string]interface{}, err error) {
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		
		
		"Product":         "elasticsearch-k8s",
		"product":         "elasticsearch-k8s",
		"Action":          "DescribeInstance",
		"Version":         "2017-06-13",
		"InstanceId":      id,
		"ClientToken":     buildClientToken("DescribeInstance"),
		"OrganizationId":  s.client.Department,
	}
	request.Method = "POST"
	request.Product = "elasticsearch-k8s"
	request.Version = "2017-06-13"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DescribeInstance"
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId, "Content-Type": "application/json"}

	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, "DescribeInstance", AlibabacloudStackSdkGoERROR)
	}
	addDebug("DescribeInstance", raw, request.QueryParams)
	response, _ := raw.(*responses.CommonResponse)
	var resp map[string]interface{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		return nil, WrapErrorf(Error(GetNotFoundMessage("ElasticsearchOnK8s Instance", id)), NotFoundWithResponse, response)
	}

	if fmt.Sprint(resp["asapiSuccess"]) == "false" {
		return nil, WrapError(fmt.Errorf("%s failed, response: %v", "DescribeInstance", resp))
	}
	v, err := jsonpath.Get("$.Result", resp)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil

}

func (s *ElasticsearchService) ElasticsearchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["status"].(string) == failState {
				return object, object["status"].(string), WrapError(Error(FailedToReachTargetStatus, object["status"].(string)))
			}
		}

		return object, object["status"].(string), nil
	}
}

func (s *ElasticsearchService) ElasticsearchRetryFunc(wait func(), errorCodeList []string, do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	var raw interface{}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithElasticsearchClient(do)

		if err != nil {
			if IsExpectedErrors(err, errorCodeList) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return raw, WrapError(err)
}

func (s *ElasticsearchService) TriggerNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "TriggerNetwork"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	request["product"] = "elasticsearch"
	request["OrganizationId"] = s.client.Department
	request["ResourceId"] = s.client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, request, &runtime)
	addDebug(action, response, content)
	if err != nil {
		if IsExpectedErrors(err, []string{"RepetitionOperationError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil

}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	conn, err := s.client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ModifyWhiteIps"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	request["product"] = "elasticsearch"
	request["OrganizationId"] = s.client.Department
	request["ResourceId"] = s.client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) DescribeElasticsearchTags(id string) (tags map[string]string, err error) {
	resourceIds, err := json.Marshal([]string{id})
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapError(err)
	}

	request := elasticsearch.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{ "Product": "elasticsearch", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	request.ResourceIds = string(resourceIds)
	request.ResourceType = strings.ToUpper(string(TagResourceInstance))
	raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.ListTagResources(request)
	})

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	response, _ := raw.(*elasticsearch.ListTagResourcesResponse)
	return s.tagsToMap(response.TagResources.TagResource), nil
}

func (s *ElasticsearchService) tagsToMap(tagSet []elasticsearch.TagResourceItem) (tags map[string]string) {
	result := make(map[string]string)
	for _, t := range tagSet {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func (s *ElasticsearchService) diffElasticsearchTags(oldTags, newTags map[string]interface{}) (remove []string, add []map[string]string) {
	for k, _ := range oldTags {
		remove = append(remove, k)
	}
	for k, v := range newTags {
		tag := map[string]string{
			"key":   k,
			"value": v.(string),
		}

		add = append(add, tag)
	}
	return
}

func (s *ElasticsearchService) getActionType(actionType bool) string {
	if actionType == true {
		return string(OPEN)
	} else {
		return string(CLOSE)
	}
}

func updateDescription(d *schema.ResourceData, meta interface{}) error {
	// var response map[string]interface{}
	// client := meta.(*connectivity.AlibabacloudStackClient)
	// action := "UpdateDescription"
	// request := map[string]interface{}{
	// 	"RegionId":    client.RegionId,
	// 	"clientToken": StringPointer(buildClientToken(action)),
	// 	"description": d.Get("description").(string),
	// }

	// request["description"] = d.Get("description").(string)
	// elasticsearchClient, err := client.NewElasticsearchClient()
	// wait := incrementalWait(3*time.Second, 5*time.Second)
	// request["product"] = "elasticsearch"
	// request["OrganizationId"] = client.Department
	// request["ResourceId"] = client.ResourceGroup
	// runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(s.client.Config.Insecure)}
	// runtime.SetAutoretry(true)
	// //response, err = elasticsearchClient.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, request, &runtime)
	// err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	// 	response, err = elasticsearchClient.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, request, &runtime)
	// 	if err != nil {
	// 		if IsExpectedErrors(err, []string{"GetCustomerLabelFail"}) || NeedRetry(err) {
	// 			wait()
	// 			return resource.RetryableError(err)
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	addDebug(action, response, nil)
	// 	return nil
	// })
	// if err != nil {
	// 	return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	// }
	// return nil
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		
		
		"Product":         "elasticsearch-k8s",
		"product":         "elasticsearch-k8s",
		"Action":          "UpdateDescription",
		"Version":         "2017-06-13",
		"InstanceId":      d.Id(),
		"description":     d.Get("description").(string),
		"ClientToken":     buildClientToken("UpdateDescription"),
		"OrganizationId":  client.Department,
	}
	request.Method = "POST"
	request.Product = "elasticsearch-k8s"
	request.Version = "2017-06-13"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "UpdateDescription"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId, "Content-Type": "application/json"}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"UpdateDescriptionFailed"}) {
			return WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Get("description").(string), "UpdateDescription", AlibabacloudStackSdkGoERROR)
	}
	addDebug("UpdateDescription", raw, request.QueryParams)
	// response, _ := raw.(*responses.CommonResponse)
	// var resp map[string]interface{}
	// err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	// if err != nil {
	// 	return nil, WrapErrorf(Error(GetNotFoundMessage("ElasticsearchOnK8s Instance", id)), NotFoundWithResponse, response)
	// }

	// if fmt.Sprint(resp["asapiSuccess"]) == "false" {
	// 	return nil, WrapError(fmt.Errorf("%s failed, response: %v", "DescribeInstance", resp))
	// }
	// v, err := jsonpath.Get("$.Result", resp)
	// if err != nil {
	// 	return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	// }
	// object = v.(map[string]interface{})
	return nil
}

func updateInstanceTags(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}

	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	remove, add := elasticsearchService.diffElasticsearchTags(o, n)

	// 对系统 Tag 进行过滤
	removeTagKeys := make([]string, 0)
	for _, v := range remove {
		if !elasticsearchTagIgnored(v, "") {
			removeTagKeys = append(removeTagKeys, v)
		}
	}
	if len(removeTagKeys) > 0 {
		tagKeys, err := json.Marshal(removeTagKeys)
		if err != nil {
			return WrapError(err)
		}

		resourceIds, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return WrapError(err)
		}
		request := elasticsearch.CreateUntagResourcesRequest()
		request.RegionId = client.RegionId
		request.TagKeys = string(tagKeys)
		request.ResourceType = strings.ToUpper(string(TagResourceInstance))
		request.ResourceIds = string(resourceIds)
		request.SetContentType("application/json")
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{ "Product": "elasticsearch", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.ResourceIds = string(resourceIds)
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UntagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}

	if len(add) > 0 {
		content := make(map[string]interface{})
		content["ResourceIds"] = []string{d.Id()}
		content["ResourceType"] = strings.ToUpper(string(TagResourceInstance))
		content["Tags"] = add
		data, err := json.Marshal(content)
		if err != nil {
			return WrapError(err)
		}

		request := elasticsearch.CreateTagResourcesRequest()
		request.RegionId = client.RegionId
		request.SetContent(data)
		request.SetContentType("application/json")

		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{ "Product": "elasticsearch", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.TagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func updateInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	var response map[string]interface{}
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchClient, err := client.NewElasticsearchClient()
	action := "UpdateInstanceChargeType"
	content := make(map[string]interface{})
	content["paymentType"] = strings.ToLower(d.Get("instance_charge_type").(string))
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = string(Year)
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = string(Month)
		}

		content["paymentInfo"] = paymentInfo
	}
	content["product"] = "elasticsearch"
	content["clientToken"] = StringPointer(buildClientToken(action))
	content["RegionId"] = client.RegionId
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err = elasticsearchClient.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
	time.Sleep(10 * time.Second)
	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func renewInstance(d *schema.ResourceData, meta interface{}) error {
	var response map[string]interface{}
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchClient, err := client.NewElasticsearchClient()
	action := "RenewInstance"
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	if d.Get("period").(int) >= 12 {
		content["duration"] = d.Get("period").(int) / 12
		content["pricingCycle"] = string(Year)
	} else {
		content["duration"] = d.Get("period").(int)
		content["pricingCycle"] = string(Month)
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err = elasticsearchClient.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
	time.Sleep(10 * time.Second)

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func updateDataNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"
	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	content["nodeAmount"] = d.Get("data_node_amount").(int)
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateDataNodeSpec(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	spec := make(map[string]interface{})
	spec["spec"] = d.Get("data_node_spec")
	spec["disk"] = d.Get("data_node_disk_size")
	spec["diskType"] = d.Get("data_node_disk_type")
	content["nodeSpec"] = spec
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	if d.Get("master_node_spec") != nil {
		master := make(map[string]interface{})
		master["spec"] = d.Get("master_node_spec").(string)
		master["amount"] = "3"
		master["diskType"] = "cloud_ssd"
		master["disk"] = "20"
		content["masterConfiguration"] = master
		content["advancedDedicateMaster"] = true
	} else {
		content["advancedDedicateMaster"] = false
	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	addDebug(action, response, content)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateAdminPassword"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {

		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		content["esAdminPassword"] = decryptResp


	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

func filterWhitelist(destIPs []string, localIPs *schema.Set) []string {
	var whitelist []string
	if destIPs != nil {
		for _, ip := range destIPs {
			if (ip == "::1" || ip == "::/0" || ip == "127.0.0.1" || ip == "0.0.0.0/0") && !localIPs.Contains(ip) {
				continue
			}
			whitelist = append(whitelist, ip)
		}
	}
	return whitelist
}

func updateClientNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "UpdateInstance"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content["isHaveClientNode"] = true

	spec := make(map[string]interface{})
	spec["spec"] = d.Get("client_node_spec")
	if d.Get("client_node_amount") == nil {
		spec["amount"] = "2"
	} else {
		spec["amount"] = d.Get("client_node_amount")
	}
	spec["disk"] = "20"
	spec["diskType"] = "cloud_efficiency"
	content["clientNodeConfiguration"] = spec
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func openHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "OpenHttps"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	action := "CloseHttps"

	var response map[string]interface{}
	content := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content["product"] = "elasticsearch"
	content["OrganizationId"] = client.Department
	content["ResourceId"] = client.ResourceGroup
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-13"), StringPointer("AK"), nil, content, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
