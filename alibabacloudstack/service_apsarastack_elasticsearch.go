package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ElasticsearchService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := make(map[string]interface{})
	response, err = s.client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "DescribeInstance", "", nil, nil, request)
	addDebug("DescribeInstance", response, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "DescribeInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "DescribeInstance", response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if (object["instanceId"].(string)) != id {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Elasticsearch Instance", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *ElasticsearchService) DescribeElasticsearchOnk8sInstance(id string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["InstanceId"] = id
	request["ClientToken"] = buildClientToken("DescribeInstance")
	
	response, err := s.client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "DescribeInstance", "", nil, nil, request)
	addDebug("DescribeInstance", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "DescribeInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["asapiSuccess"]) == "false" {
		return nil, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "DescribeInstance", response))
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *ElasticsearchService) ElasticsearchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["status"].(string) == failState {
				return object, object["status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["status"].(string)))
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
			if errmsgs.IsExpectedErrors(err, errorCodeList) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return raw, errmsgs.WrapError(err)
}

func (s *ElasticsearchService) TriggerNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	request := make(map[string]interface{})
	
	request["clientToken"] = buildClientToken("TriggerNetwork")
	request["product"] = "elasticsearch"
	response, err := s.client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "TriggerNetwork", "", nil, nil, request)
	addDebug("TriggerNetwork", response, content)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"RepetitionOperationError"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "TriggerNetwork", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	request := make(map[string]interface{})
	
	request["clientToken"] = buildClientToken("ModifyWhiteIps")
	request["product"] = "elasticsearch"
	response, err := s.client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "ModifyWhiteIps", "", nil, nil, request)
	addDebug("ModifyWhiteIps", response, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || errmsgs.NeedRetry(err) {
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = s.client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "ModifyWhiteIps", "", nil, nil, request)
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug("ModifyWhiteIps", response, nil)
				return nil
			})
		}
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ModifyWhiteIps", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) DescribeElasticsearchTags(id string) (tags map[string]string, err error) {
	resourceIds, err := json.Marshal([]string{id})
	if err != nil {
		tmp := make(map[string]string)
		return tmp, errmsgs.WrapError(err)
	}

	request := elasticsearch.CreateListTagResourcesRequest()
	s.client.InitRoaRequest(*request.RoaRequest)
	request.ResourceIds = string(resourceIds)
	request.ResourceType = strings.ToUpper(string(TagResourceInstance))
	raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.ListTagResources(request)
	})

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	response, ok := raw.(*elasticsearch.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*responses.BaseResponse))
		}
		tmp := make(map[string]string)
		return tmp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

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
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["description"] = d.Get("description").(string)
	request["ClientToken"] = buildClientToken("UpdateDescription")
	
	response, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "UpdateDescription", "", nil, nil, request)
	addDebug("UpdateDescription", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"UpdateDescriptionFailed"}) {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Get("description").(string), "UpdateDescription", errmsgs.AlibabacloudStackSdkGoERROR)
	}
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
			return errmsgs.WrapError(err)
		}

		resourceIds, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request := elasticsearch.CreateUntagResourcesRequest()
		client.InitRoaRequest(*request.RoaRequest)
		request.TagKeys = string(tagKeys)
		request.ResourceType = strings.ToUpper(string(TagResourceInstance))
		request.ResourceIds = string(resourceIds)
		request.SetContentType("application/json")

		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UntagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, ok := raw.(elasticsearch.UntagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	if len(add) > 0 {
		content := make(map[string]interface{})
		content["ResourceIds"] = []string{d.Id()}
		content["ResourceType"] = strings.ToUpper(string(TagResourceInstance))
		content["Tags"] = add
		data, err := json.Marshal(content)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		request := elasticsearch.CreateTagResourcesRequest()
		client.InitRoaRequest(*request.RoaRequest)
		request.SetContent(data)
		request.SetContentType("application/json")

		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.TagResources(request)
		})

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, ok := raw.(*elasticsearch.TagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	return nil
}

func updateInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	var response map[string]interface{}
	client := meta.(*connectivity.AlibabacloudStackClient)
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
	content["clientToken"] = buildClientToken("UpdateInstanceChargeType")
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateInstanceChargeType", "", nil, nil, content)
	time.Sleep(10 * time.Second)
	addDebug("UpdateInstanceChargeType", response, content)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), response, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func renewInstance(d *schema.ResourceData, meta interface{}) error {
	var response map[string]interface{}
	client := meta.(*connectivity.AlibabacloudStackClient)
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("RenewInstance")
	if d.Get("period").(int) >= 12 {
		content["duration"] = d.Get("period").(int) / 12
		content["pricingCycle"] = string(Year)
	} else {
		content["duration"] = d.Get("period").(int)
		content["pricingCycle"] = string(Month)
	}
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "RenewInstance", "", nil, nil, content)
	time.Sleep(10 * time.Second)

	addDebug("RenewInstance", response, content)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), response, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func updateDataNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("UpdateInstance")
	content["nodeAmount"] = d.Get("data_node_amount").(int)
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateInstance", "", nil, nil, content)
	addDebug("UpdateInstance", response, content)
	if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func updateDataNodeSpec(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("UpdateInstance")
	spec := make(map[string]interface{})
	spec["spec"] = d.Get("data_node_spec")
	spec["disk"] = d.Get("data_node_disk_size")
	spec["diskType"] = d.Get("data_node_disk_type")
	content["nodeSpec"] = spec
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateInstance", "", nil, nil, content)
	addDebug("UpdateInstance", response, content)
	if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("UpdateInstance")
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
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateInstance", "", nil, nil, content)
	addDebug("UpdateInstance", response, content)
	if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("UpdateAdminPassword")
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}
	if password != "" {
		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		content["esAdminPassword"] = decryptResp
	}
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateAdminPassword", "", nil, nil, content)
	addDebug("UpdateAdminPassword", response, content)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateAdminPassword", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
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
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("UpdateInstance")
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
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "UpdateInstance", "", nil, nil, content)
	addDebug("UpdateInstance", response, content)
	if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func openHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("OpenHttps")
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "OpenHttps", "", nil, nil, content)
	addDebug("OpenHttps", response, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "OpenHttps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})
	
	content["clientToken"] = buildClientToken("CloseHttps")
	response, err := client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "CloseHttps", "", nil, nil, content)
	addDebug("CloseHttps", response, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "CloseHttps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
