package alibabacloudstack

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/PaesslerAG/jsonpath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
// 	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackElasticsearch() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Basic instance information
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: esVersionDiffSuppressFunc,
				ForceNew:         true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w\-.]{0,30}$`), "be 0 to 30 characters in length and can contain numbers, letters, underscores, (_) and hyphens (-). It must start with a letter, a number or Chinese character."),
				Computed:     true,
			},
			"scene": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"high", "normal", "log"}, false),
			},

			// Data node configuration

			"data_node_amount": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(3, 50),
			},

			"data_node_spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+C \d+Gi`), "Spec format mast be like '\\d+C \\d+Gi'"),
			},

			"data_node_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(500, 20480),
			},

			"data_node_disk_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"yoda-lvm", "fast-disks", "fast-disks-ssd"}, false),
			},

			"data_node_affinity": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Kibana node configuration
			"kibana_node_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"kibana_password"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+C \d+Gi`), "Spec format mast be like '\\d+C \\d+Gi'"),
			},

			"kibana_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("kibana_node_spec"); !ok || v.(string) == "" {
						return true
					}
					return oldValue == newValue
				},
			},

			// Master node configuration
			"master_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 50),
				RequiredWith: []string{"master_node_spec", "master_node_disk_size", "master_node_disk_type"},
			},

			"master_node_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"master_node_amount", "master_node_disk_size", "master_node_disk_type"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+C \d+Gi`), "Spec format mast be like '\\d+C \\d+Gi'"),
			},

			"master_node_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(100),
				RequiredWith: []string{"master_node_amount", "master_node_spec", "master_node_disk_type"},
			},

			"master_node_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"yoda-lvm", "fast-disks", "fast-disks-ssd"}, false),
				RequiredWith: []string{"master_node_amount", "master_node_spec", "master_node_disk_size"},
			},

			// Client node configuration
			"client_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 25),
				RequiredWith: []string{"client_node_spec"},
			},

			"client_node_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"client_node_amount"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+C \d+Gi`), "Spec format mast be like '\\d+C \\d+Gi'"),
			},

			// network info
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"password": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(12, 32),
			},
			"monitor_password": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(12, 32),
			},


			// cluster config
			"apack_accesslog_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"apack_accesslog_search_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"thread_pool_write_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"thread_pool_search_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_routing_allocation_disk_watermark_low": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_routing_allocation_disk_watermark_high": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_routing_allocation_disk_watermark_flood_stage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action_auto_create_index": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"action_destructive_requires_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			// 只读属性
			"slb_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"kibana_slb_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kibana_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kibana_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kibana_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// 3.16.2不支持修改参数

			"private_whitelist": {
				Type: schema.TypeSet,
				// 				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"enable_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"public_whitelist": {
				Type: schema.TypeSet,
				// 				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"kibana_whitelist": {
				Type: schema.TypeSet,
				// 				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"kibana_private_whitelist": {
				Type: schema.TypeSet,
				// 				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"setting_config": {
				Type: schema.TypeMap,
				// 				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackElasticsearchCreate,
		resourceAlibabacloudStackElasticsearchRead, resourceAlibabacloudStackElasticsearchUpdate, resourceAlibabacloudStackElasticsearchDelete)
	return resource
}

func resourceAlibabacloudStackElasticsearchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	action := "CreateInstance"

	requestBody, err := buildElasticsearchCreateRequestBody(d, meta)
	var response map[string]interface{}

	// retry

	response, err = client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", action, "/openapi/instances", nil, nil, requestBody)
	if err != nil {
		return err
	}

	resp, err := jsonpath.Get("$.Result.instanceId", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Result.instanceId", response)
	}
	d.SetId(resp.(string))

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}

	object, err := elasticsearchService.DescribeElasticsearchInstance(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("version", object["version"].(string))
	d.Set("description", object["description"].(string))
	d.Set("cpu_type", object["cpuType"].(map[string]interface{})["cpuBrand"].(string))
	d.Set("zone_id", object["zoneInfos"].([]interface{})[0].(map[string]interface{})["zoneId"].(string))

	if object["dataNode"].(bool) {
		if v, err := object["nodeAmount"].(json.Number).Int64(); err == nil {
			d.Set("data_node_amount", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
		nodeSpec := object["nodeSpec"].(map[string]interface{})
		d.Set("data_node_spec", nodeSpec["spec"].(string))
		if v, err := nodeSpec["disk"].(json.Number).Int64(); err == nil {
			d.Set("data_node_disk_size", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
		d.Set("data_node_disk_type", nodeSpec["storageClassName"].(string))
	} else {
		d.Set("data_node_amount", nil)
		d.Set("data_node_spec", nil)
		d.Set("data_node_disk_size", nil)
		d.Set("data_node_disk_type", nil)
	}

	if object["haveKibana"].(bool) {
		d.Set("kibana_node_spec", object["kibanaConfiguration"].(map[string]interface{})["spec"].(string))
		d.Set("kibana_slb_address", object["kibanaSlbAddress"].(string))
		d.Set("kibana_domain", object["kibanaDomain"].(string))
		d.Set("kibana_protocol", object["kibanaProtocol"].(string))
		if v, err := object["kibanaPort"].(json.Number).Int64(); err == nil {
			d.Set("kibana_port", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
	} else {
		d.Set("kibana_node_spec", nil)
		d.Set("kibana_slb_address", nil)
		d.Set("kibana_domain", nil)
		d.Set("kibana_protocol", nil)
		d.Set("kibana_port", nil)
	}

	if object["advancedDedicateMaster"].(bool) {
		masterConfiguration := object["masterConfiguration"].(map[string]interface{})
		if v, err := masterConfiguration["amount"].(json.Number).Int64(); err == nil {
			d.Set("master_node_amount", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
		d.Set("master_node_spec", masterConfiguration["spec"].(string))
		if v, err := masterConfiguration["disk"].(json.Number).Int64(); err == nil {
			d.Set("master_node_disk_size", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
		d.Set("master_node_disk_type", masterConfiguration["storageClassName"].(string))
	} else {
		d.Set("master_node_amount", nil)
		d.Set("master_node_disk_size", nil)
		d.Set("master_node_disk_type", nil)
	}

	if object["haveClientNode"].(bool) {
		clientNodeConfiguration := object["clientNodeConfiguration"].(map[string]interface{})
		if v, err := clientNodeConfiguration["amount"].(json.Number).Int64(); err == nil {
			d.Set("client_node_amount", int(v))
		} else {
			return errmsgs.WrapError(err)
		}
		d.Set("client_node_spec", clientNodeConfiguration["spec"].(string))
	} else {
		d.Set("client_node_amount", nil)
		d.Set("client_node_spec", nil)
	}

	d.Set("vswitch_id", object["networkConfig"].(map[string]interface{})["vswitchId"])

	d.Set("slb_address", object["slbAddress"].(string))
	d.Set("domain", object["domain"])
	d.Set("port", object["port"])
	d.Set("status", object["status"])

	esIPWhitelist := object["esIPWhitelist"].([]interface{})
	publicIpWhitelist := object["publicIpWhitelist"].([]interface{})
	d.Set("private_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(esIPWhitelist), d.Get("private_whitelist").(*schema.Set)))
	d.Set("public_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(publicIpWhitelist), d.Get("public_whitelist").(*schema.Set)))
	d.Set("enable_public", object["enablePublic"])

	// Kibana configuration
	kibanaIPWhitelist := object["kibanaIPWhitelist"].([]interface{})
	d.Set("kibana_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(kibanaIPWhitelist), d.Get("kibana_whitelist").(*schema.Set)))
	if object["enableKibanaPublicNetwork"].(bool) {
		d.Set("kibana_domain", object["kibanaDomain"])
		d.Set("kibana_port", object["kibanaPort"])
	}

	kibanaPrivateIPWhitelist := object["kibanaPrivateIPWhitelist"].([]interface{})
	d.Set("kibana_private_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(kibanaPrivateIPWhitelist), d.Get("kibana_private_whitelist").(*schema.Set)))

	// Protocol: HTTP/HTTPS
	d.Set("protocol", object["protocol"])

	esConfig := object["esConfig"].(map[string]interface{})
	if esConfig != nil {
		d.Set("setting_config", esConfig)
	}

	var response map[string]interface{}
	action := "GetElasticsearchSettings"
	response, err = client.DoTeaRequest("GET", "elasticsearch-k8s", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/es-settings", d.Id()), nil, nil, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}
		return err
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, d.Id(), "$", response)
	}
	for _, item := range v.([]interface{}) {
		info := item.(map[string]interface{})
		key := info["key"].(string)
		if info["type"] == "Boolean" {
			if reflect.TypeOf(info["value"]) == reflect.TypeOf(true) {
				d.Set(strings.Replace(key, ".", "_", -1), info["value"])
			} else {
				value := strings.ToLower(info["value"].(string)) == "true"
				d.Set(strings.Replace(key, ".", "_", -1), value)
			}
		} else if info["type"] == "Integer" {
			var valueString string
			if reflect.TypeOf(info["value"]) == reflect.TypeOf("0") {
				valueString = info["value"].(string)
			} else {
				valueString = info["value"].(json.Number).String()
			}
			value, _ := strconv.Atoi(valueString)
			d.Set(strings.Replace(key, ".", "_", -1), value)
		} else {
			d.Set(strings.Replace(key, ".", "_", -1), info["value"])
		}
	}
	return nil
}

func resourceAlibabacloudStackElasticsearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 10 * time.Second

	if d.HasChange("private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("public_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("public_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("kibana_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("kibana_private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("protocol") {
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		var https func(*schema.ResourceData, interface{}) error
		if d.Get("protocol") == "HTTPS" {
			https = openHttps
		} else if d.Get("protocol") == "HTTP" {
			https = closeHttps
		}
		if nil != https {
			if err := https(d, meta); err != nil && !errmsgs.IsExpectedErrors(err, []string{"InvalidAction.NotFound"}) {
				// 3162 老版本不支持HTTPS -> HTTP
				return errmsgs.WrapError(err)
			}
		}
	}

	if d.HasChange("setting_config") {
		action := "UpdateInstanceSettings"
		content := map[string]interface{}{
			"RegionId":    client.RegionId,
			"clientToken": StringPointer(buildClientToken(action)),
		}
		config := d.Get("setting_config").(map[string]interface{})
		content["esConfig"] = config
		_, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/instance-settings", d.Id()), nil, nil, content)

		if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
			return err
		}
		stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
		stateConf.PollInterval = 5 * time.Second
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	cluster_configs := []string{"apack_accesslog_enabled", "apack_accesslog_search_enabled", "thread_pool_write_queue_size", "thread_pool_search_queue_size", "cluster_routing_allocation_disk_watermark_low", "cluster_routing_allocation_disk_watermark_high", "cluster_routing_allocation_disk_watermark_flood_stage", "action_auto_create_index", "action_destructive_requires_name"}
	if d.HasChanges(cluster_configs...) {
		esConfig := map[string]interface{}{}
		for _, config := range cluster_configs {
			if d.HasChange(config) {
				esConfig[EsClusterConfig[config]] = d.Get(config)
			}
		}
		request := map[string]interface{}{
			"esConfig": esConfig,
		}
		_, err := client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", "UpdateElasticsearchSettings", fmt.Sprintf("/openapi/instances/%s/actions/es-settings", d.Id()), nil, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
				return nil
			}
			return err
		}
		stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
		stateConf.PollInterval = 5 * time.Second
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if d.IsNewResource() {
		return nil
	}

	if d.HasChange("description") {
		if err := updateDescription(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChanges("data_node_amount", "data_node_spec", "master_node_spec", "master_node_amount", "client_node_spec", "client_node_amount", "kibana_node_spec") {
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		if err := updateNodes(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChanges("password", "kibana_password", "monitor_password") {
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		if err := updatePassword(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func resourceAlibabacloudStackElasticsearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	action := "DeleteInstance"

	if _, ok := d.GetOk("vswitch_id"); ok {
		// Instance will be completed deleted in 5 minutes, so deleting vswitch is available after the time.
		defer time.Sleep(3 * time.Minute)
	}

	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	_, err := client.DoTeaRequest("DELETE", "elasticsearch-k8s", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}
		return err
	}

	stateConf := BuildStateConf([]string{"activating", "inactive", "active"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{}))
	stateConf.PollInterval = 5 * time.Second

	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func buildElasticsearchCreateRequestBody(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	content := make(map[string]interface{})
	content["ClientToken"] = buildClientToken("createInstance")

	content["nodeAmount"] = connectivity.GetResourceData(d, "data_node_amount", "data_node.amount")
	content["esVersion"] = d.Get("version")
	content["description"] = d.Get("description")
	content["cpuType"] = map[string]interface{}{
		"cpuBrand":    d.Get("cpu_type"),
		"defaultType": false,
	}

	content["esAdminPassword"] = d.Get("password")
	content["monitorPassword"] = d.Get("monitor_password")
	content["scene"] = d.Get("scene")

	// Data node configuration
	content["dataNode"] = true
	dataNodeSpec := make(map[string]interface{})
	dataNodeSpec["spec"] = d.Get("data_node_spec")
	dataNodeSpec["disk"] = d.Get("data_node_disk_size")
	dataNodeSpec["diskType"] = d.Get("data_node_disk_type")
	content["nodeSpec"] = dataNodeSpec
	content["dataNodeAffinity"] = d.Get("data_node_affinity")

	// Kibana node configure
	if v, ok := d.GetOk("kibana_node_spec"); ok {
		content["haveKibana"] = true
		kibanaConfiguration := make(map[string]interface{})
		kibanaConfiguration["amount"] = 1
		kibanaConfiguration["spec"] = v
		content["kibanaConfiguration"] = kibanaConfiguration
		content["kibanaPassword"] = d.Get("kibana_password")
	} else {
		content["haveKibana"] = false
	}

	// Master node configuration
	if _, ok := d.GetOk("master_node_spec"); ok {
		content["advancedDedicateMaster"] = true
		masterConfiguration := make(map[string]interface{})
		masterConfiguration["spec"] = d.Get("master_node_spec")
		masterConfiguration["amount"] = d.Get("master_node_amount")
		masterConfiguration["disk"] = d.Get("master_node_disk_size")
		masterConfiguration["diskType"] = d.Get("master_node_disk_type")
		content["masterConfiguration"] = masterConfiguration
	} else {
		content["advancedDedicateMaster"] = false
	}

	// Client node configuration
	if _, ok := d.GetOk("client_node_spec"); ok {
		clientNode := make(map[string]interface{})
		clientNode["spec"] = d.Get("client_node_spec")
		clientNode["amount"] = d.Get("client_node_amount")

		content["haveClientNode"] = true
		content["clientNodeConfiguration"] = clientNode
	} else {
		content["haveClientNode"] = false
	}

	// Network configuration
	network := make(map[string]interface{})
	network["vsArea"] = d.Get("zone_id")
	if _, ok := d.GetOk("vswitch_id"); ok {
		vswitchId := d.Get("vswitch_id")
		vsw, err := vpcService.DescribeVSwitch(vswitchId.(string))
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		network["type"] = "vpc"
		network["vpcId"] = vsw.VpcId
		network["vswitchId"] = vswitchId

	} else {
		network["type"] = "anytunnel"
	}
	content["networkConfig"] = network

	return content, nil
}
