package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"


	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackElasticsearchOnk8s() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Basic instance information
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w\-.]{0,30}$`), "be 0 to 30 characters in length and can contain numbers, letters, underscores, (_) and hyphens (-). It must start with a letter, a number or Chinese character."),
				Computed:     true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"version": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: esVersionDiffSuppressFunc,
				ForceNew:         true,
			},
			"tags": tagsSchema(),

			// Life cycle
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Default:      PostPaid,
				Optional:     true,
			},

			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},

			// Data node configuration
			"data_node_amount": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(2, 50),
			},

			"data_node_spec": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_node_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"data_node_disk_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"data_node_disk_encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"private_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"enable_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnablePublicDiffSuppressFunc,
			},

			"master_node_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Client node configuration
			"client_node_amount": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 25),
			},

			"client_node_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTP",
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Kibana node configuration
			"kibana_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kibana_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"enable_kibana_public_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"kibana_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnableKibanaPublicDiffSuppressFunc,
			},

			"enable_kibana_private_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kibana_private_whitelist": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Computed:         true,
				DiffSuppressFunc: elasticsearchEnableKibanaPrivateDiffSuppressFunc,
			},

			"zone_count": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3),
				Default:      1,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"setting_config": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackElasticsearchOnk8sCreate, resourceAlibabacloudStackElasticsearchOnk8sRead, resourceAlibabacloudStackElasticsearchOnk8sUpdate, resourceAlibabacloudStackElasticsearchOnk8sDelete)
	return resource
}

func resourceAlibabacloudStackElasticsearchOnk8sCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	action := "createInstance"
	requestBody, _ := buildElasticsearchOnk8sCreateRequestBody(d, meta)
	body, _ := json.Marshal(requestBody)
	bodydata := string(body)

	request := client.NewCommonRequest("POST", "elasticsearch-k8s", "2017-06-13", "createInstance", "")
	request.Headers["Content-Type"] = "application/json"
	request.QueryParams["X-acs-body"] = bodydata
	request.QueryParams["ClientToken"] = buildClientToken("createInstance")

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_elasticsearch_k8s_instance", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(action, raw, request.QueryParams)

	var resp CreateElasticsearchOnk8sResponse

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "apsarastack_elasticsearch_on_k8s_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	Result := resp.Result
	log.Printf("############## Result : %s", Result)

	d.SetId(Result["instanceId"])

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, elasticsearchService.ElasticsearchOnK8sStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackElasticsearchOnk8sRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}

	object, err := elasticsearchService.DescribeElasticsearchOnk8sInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if (object["instanceId"].(string)) != d.Id() {
		return errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ElasticsearchOnK8s Instance", d.Id())), errmsgs.NotFoundWithResponse)
	}

	d.Set("description", object["description"])
	d.Set("status", object["status"])
	d.Set("vswitch_id", object["networkConfig"].(map[string]interface{})["vswitchId"])

	esIPWhitelist := object["esIPWhitelist"].([]interface{})
	publicIpWhitelist := object["publicIpWhitelist"].([]interface{})
	d.Set("private_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(esIPWhitelist), d.Get("private_whitelist").(*schema.Set)))
	d.Set("public_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(publicIpWhitelist), d.Get("public_whitelist").(*schema.Set)))
	d.Set("enable_public", object["enablePublic"])
	d.Set("version", object["esVersion"])
	d.Set("instance_charge_type", getChargeType(object["paymentType"].(string)))

	d.Set("domain", object["domain"])
	d.Set("port", object["port"])

	// Kibana configuration
	d.Set("enable_kibana_public_network", object["enableKibanaPublicNetwork"])
	kibanaIPWhitelist := object["kibanaIPWhitelist"].([]interface{})
	d.Set("kibana_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(kibanaIPWhitelist), d.Get("kibana_whitelist").(*schema.Set)))
	if object["enableKibanaPublicNetwork"].(bool) {
		d.Set("kibana_domain", object["kibanaDomain"])
		d.Set("kibana_port", object["kibanaPort"])
	}

	d.Set("enable_kibana_private_network", object["enableKibanaPrivateNetwork"])
	kibanaPrivateIPWhitelist := object["kibanaPrivateIPWhitelist"].([]interface{})
	d.Set("kibana_private_whitelist", filterWhitelist(convertArrayInterfaceToArrayString(kibanaPrivateIPWhitelist), d.Get("kibana_private_whitelist").(*schema.Set)))

	// Data node configuration
	d.Set("data_node_amount", object["nodeAmount"])
	d.Set("data_node_spec", object["nodeSpec"].(map[string]interface{})["spec"])
	d.Set("data_node_disk_size", object["nodeSpec"].(map[string]interface{})["disk"])
	d.Set("data_node_disk_type", object["nodeSpec"].(map[string]interface{})["diskType"])
	d.Set("data_node_disk_encrypted", object["nodeSpec"].(map[string]interface{})["diskEncryption"])
	d.Set("master_node_spec", object["masterConfiguration"].(map[string]interface{})["spec"])
	// Client node configuration
	d.Set("client_node_amount", object["clientNodeConfiguration"].(map[string]interface{})["amount"])
	d.Set("client_node_spec", object["clientNodeConfiguration"].(map[string]interface{})["spec"])
	// Protocol: HTTP/HTTPS
	d.Set("protocol", object["protocol"])

	// Cross zone configuration
	d.Set("zone_count", object["zoneCount"])
	d.Set("resource_group_id", object["resourceGroup"])

	esConfig := object["esConfig"].(map[string]interface{})
	if esConfig != nil {
		d.Set("setting_config", esConfig)
	}

	// tags
	//tags, err := elasticsearchService.DescribeElasticsearchTags(d.Id())
	//if err != nil {
	//	return errmsgs.WrapError(err)
	//}
	//if len(tags) > 0 {
	//	d.Set("tags", tags)
	//}

	return nil
}

func resourceAlibabacloudStackElasticsearchOnk8sUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchOnK8sStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if d.HasChange("description") {
		if err := updateDescription(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("description")
	}

	if d.HasChange("private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("private_whitelist")
	}

	// if d.HasChange("enable_public") {
	// 	content := make(map[string]interface{})
	// 	content["networkType"] = string(PUBLIC)
	// 	content["nodeType"] = string(WORKER)
	// 	content["actionType"] = elasticsearchService.getActionType(d.Get("enable_public").(bool))
	// 	if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
	// 		return errmsgs.WrapError(err)
	// 	}

	// 	//d.SetPartial("enable_public")
	// }

	if d.Get("enable_public").(bool) && d.HasChange("public_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(WORKER)
		content["whiteIpList"] = d.Get("public_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("public_whitelist")
	}

	// if d.HasChange("enable_kibana_public_network") || d.IsNewResource() {
	// 	content := make(map[string]interface{})
	// 	content["networkType"] = string(PUBLIC)
	// 	content["nodeType"] = string(KIBANA)
	// 	content["actionType"] = elasticsearchService.getActionType(d.Get("enable_kibana_public_network").(bool))
	// 	if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
	// 		return errmsgs.WrapError(err)
	// 	}

	// 	//d.SetPartial("enable_kibana_public_network")
	// }

	if d.Get("enable_kibana_public_network").(bool) && d.HasChange("kibana_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PUBLIC)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("kibana_whitelist")
	}

	// if d.HasChange("enable_kibana_private_network") {
	// 	content := make(map[string]interface{})
	// 	content["networkType"] = string(PRIVATE)
	// 	content["nodeType"] = string(KIBANA)
	// 	content["actionType"] = elasticsearchService.getActionType(d.Get("enable_kibana_private_network").(bool))
	// 	if err := elasticsearchService.TriggerNetwork(d, content, meta); err != nil {
	// 		return errmsgs.WrapError(err)
	// 	}

	// 	//d.SetPartial("enable_kibana_private_network")
	// }

	if d.Get("enable_kibana_private_network").(bool) == true && d.HasChange("kibana_private_whitelist") {
		content := make(map[string]interface{})
		content["networkType"] = string(PRIVATE)
		content["nodeType"] = string(KIBANA)
		content["whiteIpList"] = d.Get("kibana_private_whitelist").(*schema.Set).List()
		if err := elasticsearchService.ModifyWhiteIps(d, content, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("kibana_private_whitelist")
	}

	if d.HasChange("tags") {
		if err := updateInstanceTags(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("tags")
	}

	if d.HasChanges("client_node_spec", "client_node_amount") {

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := updateClientNode(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("client_node_spec")
		//d.SetPartial("client_node_amount")
	}

	// if d.HasChange("protocol") {

	// 	if _, err := stateConf.WaitForState(); err != nil {
	// 		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	// 	}

	// 	var https func(*schema.ResourceData, interface{}) error

	// 	if d.Get("protocol") == "HTTPS" {
	// 		https = openHttps
	// 	} else if d.Get("protocol") == "HTTP" {
	// 		https = closeHttps
	// 	}

	// 	if nil != https {
	// 		if err := https(d, meta); err != nil {
	// 			return errmsgs.WrapError(err)
	// 		}
	// 	}

	// 	//d.SetPartial("protocol")
	// }
	if d.HasChange("setting_config") {
		action := "UpdateInstanceSettings"
		content := map[string]interface{}{
			"RegionId":    client.RegionId,
			"clientToken": StringPointer(buildClientToken(action)),
		}
		config := d.Get("setting_config").(map[string]interface{})
		content["esConfig"] = config
		_, err = client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", action, "", nil, nil, content)
		if err != nil && !errmsgs.IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
			return err
		}
		stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchOnK8sStateRefreshFunc(d.Id(), []string{"inactive"}))
		stateConf.PollInterval = 5 * time.Second
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		//d.SetPartial("setting_config")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	if d.HasChange("instance_charge_type") {
		if err := updateInstanceChargeType(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("instance_charge_type")
		//d.SetPartial("period")
	} else if d.Get("instance_charge_type").(string) == string(PrePaid) && d.HasChange("period") {
		if err := renewInstance(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("period")
	}

	if d.HasChange("data_node_amount") {

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := updateDataNodeAmount(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("data_node_amount")
	}

	if d.HasChanges("data_node_spec", "data_node_disk_size", "data_node_disk_type") {

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := updateDataNodeSpec(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("data_node_spec")
		//d.SetPartial("data_node_disk_size")
		//d.SetPartial("data_node_disk_type")
		//d.SetPartial("data_node_disk_encrypted")
	}

	if d.HasChange("master_node_spec") {

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := updateMasterNode(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("master_node_spec")
	}

	if d.HasChanges("password", "kms_encrypted_password") {

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := updatePassword(d, meta); err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("password")
	}

	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackElasticsearchOnk8sDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	elasticsearchService := ElasticsearchService{client}
	action := "DeleteInstance"

	if strings.ToLower(d.Get("instance_charge_type").(string)) == strings.ToLower(string(PrePaid)) {
		return errmsgs.WrapError(errmsgs.Error("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically"))
	}
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"clientToken": StringPointer(buildClientToken(action)),
	}
	// retry


	_, err = client.DoTeaRequest("POST", "elasticsearch-k8s", "2017-06-13", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}
		return err
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutDelete), 5*time.Minute, elasticsearchService.ElasticsearchOnK8sStateRefreshFunc(d.Id(), []string{}))
	stateConf.PollInterval = 5 * time.Second

	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	// Instance will be completed deleted in 5 minutes, so deleting vswitch is available after the time.
	time.Sleep(5 * time.Minute)

	return nil
}

func buildElasticsearchOnk8sCreateRequestBody(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	content := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		content["ResourceGroup"] = v.(string)
	} else {
		content["ResourceGroup"] = client.ResourceGroup
	}
	content["OrganizationId"] = client.Department
	content["region"] = client.RegionId
	content["Department"] = client.Department
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

	//content["nodeAmount"] = d.Get("data_node_amount")
	content["esVersion"] = d.Get("version")
	content["description"] = d.Get("description")
	content["instanceCategory"] = "x-pack"
	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return nil, errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return content, errmsgs.WrapError(err)
		}
		content["esAdminPassword"] = decryptResp
	}

	// Data node configuration
	if d.Get("data_node_spec") != nil && d.Get("data_node_spec") != "" {
		dataNodeSpec := make(map[string]interface{})
		dataNodeSpec["spec"] = d.Get("data_node_spec")
		dataNodeSpec["disk"] = d.Get("data_node_disk_size")
		dataNodeSpec["diskType"] = "fast-disks"
		content["dataNode"] = true
		content["nodeSpec"] = dataNodeSpec
		content["nodeAmount"] = 3
	}
	//dataNodeSpec := make(map[string]interface{})
	//dataNodeSpec["spec"] = d.Get("data_node_spec")
	//dataNodeSpec["disk"] = d.Get("data_node_disk_size")
	//dataNodeSpec["diskType"] = d.Get("data_node_disk_type")
	//dataNodeSpec["diskEncryption"] = d.Get("data_node_disk_encrypted")
	//content["nodeSpec"] = dataNodeSpec

	// Master node configuration
	if d.Get("master_node_spec") != nil && d.Get("master_node_spec") != "" {
		masterNode := make(map[string]interface{})
		masterNode["spec"] = d.Get("master_node_spec")
		masterNode["amount"] = "3"
		masterNode["disk"] = "100"
		masterNode["diskType"] = "yoda-lvm"
		content["advancedDedicateMaster"] = true
		content["masterConfiguration"] = masterNode
	}

	// kinba node configuration
	if d.Get("kibana_node_spec") != nil && d.Get("kibana_node_spec") != "" {
		kibanaNode := make(map[string]interface{})
		kibanaNode["spec"] = d.Get("kibana_node_spec")
		kibanaNode["amount"] = "1"
		content["haveKibana"] = true
		content["kibanaConfiguration"] = kibanaNode
	}

	// Client node configuration
	if d.Get("client_node_spec") != nil && d.Get("client_node_spec") != "" {
		clientNode := make(map[string]interface{})
		clientNode["spec"] = d.Get("client_node_spec")
		clientNode["amount"] = 2
		content["haveClientNode"] = true
		content["clientNodeConfiguration"] = clientNode
	}

	// Network configuration
	vswitchId := d.Get("vswitch_id")
	vsw, err := vpcService.DescribeVSwitch(vswitchId.(string))
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	network := make(map[string]interface{})
	network["type"] = "vpc"
	network["vpcId"] = vsw.VpcId
	network["vswitchId"] = vswitchId
	network["vsArea"] = vsw.ZoneId

	content["networkConfig"] = network

	if d.Get("zone_count") != nil && d.Get("zone_count") != "" {
		content["zoneCount"] = d.Get("zone_count")
	}
	var zones []map[string]interface{}
	zones = append(zones, map[string]interface{}{
		"zoneName": vsw.ZoneId,
	})
	content["zones"] = zones
	content["zoneName"] = vsw.ZoneId

	return content, nil
}

type CreateElasticsearchOnk8sResponse struct {
	*responses.BaseResponse
	RequestId      string `json:"RequestId" xml:"RequestId"`
	AsapiRequestId string `json:"asapiRequestId" xml:"asapiRequestId"`
	Result         map[string]string
}

// {"eagleEyeTraceId":"0a1c830b16989984038808531d0052","asapiSuccess":true,"asapiRequestId":"D0D646B0-0205-4B73-9E83-B198419B88E7","RequestId":"D0D646B0-0205-4B73-9E83-B198419B88E7","Result":{"instanceId":"es-cn-re90hg799q019d250"}}

func (s *ElasticsearchService) ElasticsearchOnK8sStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchOnk8sInstance(id)
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
		log.Printf("DescribeElasticsearchOnk8sInstance result status: %s", object["status"])
		return object, object["status"].(string), nil
	}
}
