package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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

var EsClusterConfig = map[string]string{
	"apack_accesslog_enabled":                               "apack.accesslog.enabled",
	"apack_accesslog_search_enabled":                        "apack.accesslog.search.enabled",
	"thread_pool_write.queue_size":                          "thread_pool.write_queue_size",
	"thread_pool_search_queue_size":                         "thread_pool.search.queue_size",
	"cluster_routing_allocation_disk_watermark_low":         "cluster.routing.allocation.disk.watermark.low",
	"cluster_routing_allocation_disk_watermark_high":        "cluster.routing.allocation.disk.watermark.high",
	"cluster_routing_allocation_disk_watermark_flood_stage": "cluster.routing.allocation.disk.watermark.flood_stage",
	"action_auto_create_index":                              "action.auto_create_index",
	"action_destructive_requires_name":                      "action_destructive_requires_name",
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := make(map[string]interface{})
	response, err = s.client.DoTeaRequest("GET", "elasticsearch-k8s", "2017-06-13", "DescribeInstance", fmt.Sprintf("/openapi/instances/%s", id), nil, nil, request)
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
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if (object["instanceId"].(string)) != id {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Elasticsearch Instance", id)), errmsgs.NotFoundWithResponse, response)
	}
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

	request["product"] = "elasticsearch"
	response, err := s.client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "TriggerNetwork", "", nil, nil, request)
	addDebug("TriggerNetwork", response, content)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"RepetitionOperationError"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "TriggerNetwork", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	request := make(map[string]interface{})

	response, err := s.client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "ModifyAclWhiteIps", fmt.Sprintf("/openapi/instances/%s/actions/modify-acl-white-ips", d.Id()), nil, nil, request)
	addDebug("ModifyWhiteIps", response, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || errmsgs.NeedRetry(err) {
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = s.client.DoTeaRequest("POST", "elasticsearch", "2017-06-13", "ModifyWhiteIps", "", nil, nil, request)
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{"InvalidAction.NotFound"}) {
						// 老版本 3.16.2不支持修改
						return nil
					}
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

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second
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

	response, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "UpdateDescription", fmt.Sprintf("/openapi/instances/%s/description", d.Id()), nil, nil, request)
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

func esSpecUpDownGrade(oldSpec, newSpec string) (string, string, error) {
	re := regexp.MustCompile(`^(\d+)C\s+(\d+)Gi$`)

	// 匹配输入字符串
	matches := re.FindStringSubmatch(oldSpec)
	if len(matches) != 3 { // 完整匹配 + 2 个捕获组
		return "", "", fmt.Errorf("格式不匹配: %s", oldSpec)
	}
	old_cpu, _ := strconv.Atoi(matches[1])
	old_memory, _ := strconv.Atoi(matches[2])

	matches = re.FindStringSubmatch(newSpec)
	if len(matches) != 3 { // 完整匹配 + 2 个捕获组
		return "", "", fmt.Errorf("格式不匹配: %s", newSpec)
	}
	new_cpu, _ := strconv.Atoi(matches[1])
	new_memory, _ := strconv.Atoi(matches[2])

	if old_cpu > new_cpu {
		old_cpu = new_cpu
	}

	if old_memory > new_memory {
		old_memory = new_memory
	}

	return fmt.Sprintf("%dC %dGi", old_cpu, old_memory), fmt.Sprintf("%dC %dGi", new_cpu, new_memory), nil
}

func updateNodes(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	downgradeContent := make(map[string]interface{})
	downgradeContent["orderActionType"] = "downgrade"
	upgradeContent := make(map[string]interface{})
	if old, new := d.GetChange("data_node_amount"); old.(int) < new.(int) {
		downgradeContent["nodeAmount"] = old
		upgradeContent["nodeAmount"] = new
	} else {
		downgradeContent["nodeAmount"] = new
		upgradeContent["nodeAmount"] = new
	}

	d_spec := make(map[string]interface{})
	u_spec := make(map[string]interface{})
	old, new := d.GetChange("data_node_spec")
	var err error
	if d_spec["spec"], u_spec["spec"], err = esSpecUpDownGrade(old.(string), new.(string)); err != nil {
		return err
	}
	d_spec["disk"] = d.Get("data_node_disk_size")
	u_spec["disk"] = d.Get("data_node_disk_size")
	d_spec["diskType"] = d.Get("data_node_disk_type")
	u_spec["diskType"] = d.Get("data_node_disk_type")
	downgradeContent["nodeSpec"] = d_spec
	upgradeContent["nodeSpec"] = u_spec

	d_masterConfiguration := make(map[string]interface{})
	u_masterConfiguration := make(map[string]interface{})
	old, new = d.GetChange("master_node_spec")
	if d_masterConfiguration["spec"], u_masterConfiguration["spec"], err = esSpecUpDownGrade(old.(string), new.(string)); err != nil {
		return err
	}
	if old, new := d.GetChange("master_node_amount"); old.(int) < new.(int) {
		d_masterConfiguration["amount"] = old
		u_masterConfiguration["amount"] = new
	} else {
		d_masterConfiguration["amount"] = new
		u_masterConfiguration["amount"] = new
	}
	d_masterConfiguration["disk"] = d.Get("master_node_disk_size")
	u_masterConfiguration["disk"] = d.Get("master_node_disk_size")
	d_masterConfiguration["diskType"] = d.Get("master_node_disk_type")
	u_masterConfiguration["diskType"] = d.Get("master_node_disk_type")
	downgradeContent["masterConfiguration"] = d_masterConfiguration
	upgradeContent["masterConfiguration"] = u_masterConfiguration

	if v, ok := d.GetOk("client_node_spec"); ok && v.(string) != "" {
		d_clientNode := make(map[string]interface{})
		u_clientNode := make(map[string]interface{})
		old, new = d.GetChange("client_node_spec")
		if d_clientNode["spec"], u_clientNode["spec"], err = esSpecUpDownGrade(old.(string), new.(string)); err != nil {
			return err
		}

		if old, new := d.GetChange("client_node_amount"); old.(int) < new.(int) {
			d_clientNode["amount"] = old
			u_clientNode["amount"] = new
		} else {
			d_clientNode["amount"] = new
			u_clientNode["amount"] = new
		}

		downgradeContent["haveClientNode"] = true
		upgradeContent["haveClientNode"] = true
		downgradeContent["clientNodeConfiguration"] = d_clientNode
		upgradeContent["clientNodeConfiguration"] = u_clientNode
	} else {
		downgradeContent["haveClientNode"] = false
		upgradeContent["haveClientNode"] = false
	}

	if v, ok := d.GetOk("kibana_node_spec"); ok && v.(string) != "" {
		downgradeContent["haveKibana"] = true
		upgradeContent["haveKibana"] = true
		d_kibanaConfiguration := make(map[string]interface{})
		u_kibanaConfiguration := make(map[string]interface{})
		d_kibanaConfiguration["amount"] = 1
		u_kibanaConfiguration["amount"] = 1
		old, new = d.GetChange("kibana_node_spec")
		if d_kibanaConfiguration["spec"], u_kibanaConfiguration["spec"], err = esSpecUpDownGrade(old.(string), new.(string)); err != nil {
			return err
		}
		downgradeContent["kibanaConfiguration"] = d_kibanaConfiguration
		upgradeContent["kibanaConfiguration"] = u_kibanaConfiguration
	} else {
		downgradeContent["haveKibana"] = false
		upgradeContent["haveKibana"] = false
	}

	for _, content := range []map[string]interface{}{downgradeContent, upgradeContent} {
		response, err := client.DoTeaRequest("PUT", "elasticsearch-k8s", "2017-06-13", "UpdateInstance", fmt.Sprintf("/openapi/instances/%s", d.Id()), nil, nil, content)
		addDebug("UpdateInstance", response, content)
		if err != nil && errmsgs.IsExpectedErrors(err, []string{"UpdateInstanceNoChange"}) {
			// 当前配置没有修改，忽略错误
			continue
		}
		if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
		stateConf.PollInterval = 10 * time.Second

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	contents := make([]map[string]interface{}, 0)

	if d.HasChange("password") {
		change := map[string]interface{}{
			"userName": "elastic",
			"type":     "PASSWORD",
			"password": d.Get("password").(string),
		}
		contents = append(contents, change)
	}
	if d.HasChange("monitor_password") {
		change := map[string]interface{}{
			"userName": "monitoring_collector",
			"type":     "PASSWORD",
			"password": d.Get("monitor_password").(string),
		}
		contents = append(contents, change)
	}
	if v, ok := d.GetOk("kibana_node_spec"); ok && v.(string) != "" && d.HasChange("kibana_password") {
		change := map[string]interface{}{
			"userName": "kibanaserver",
			"type":     "PASSWORD",
			"password": d.Get("kibana_password").(string),
		}
		contents = append(contents, change)
	}
	for _, content := range contents {
		response, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "UpdateUserSecurityInfo", fmt.Sprintf("/openapi/instances/%s/securitys", d.Id()), nil, nil, content)
		addDebug("UpdateAdminPassword", response, content)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "UpdateUserSecurityInfo", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
		stateConf.PollInterval = 10 * time.Second

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
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

func openHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})

	response, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "EnableHttps", fmt.Sprintf("/openapi/instances/%s/actions/enable-https", d.Id()), nil, nil, content)
	addDebug("OpenHttps", response, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "OpenHttps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	content := make(map[string]interface{})

	response, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "DisableHttps", fmt.Sprintf("/openapi/instances/%s/actions/disable-https", d.Id()), nil, nil, content)
	addDebug("CloseHttps", response, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "CloseHttps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
