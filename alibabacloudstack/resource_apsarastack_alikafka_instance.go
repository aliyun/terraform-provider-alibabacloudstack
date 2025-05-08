package alibabacloudstack

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	"encoding/json"


	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAlikafkaInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deploy_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 5}),
			},
			"partition_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"partition_num", "topic_quota"},
			},
			"topic_quota": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, _ := strconv.Atoi(old)
					partitionNum := d.Get("partition_num").(int)
					if o > 0 {
						return o-1000 == partitionNum
					}
					return false
				},
				Deprecated: "Attribute `topic_quota` has been deprecated since 1.194.0 and it will be removed in the next future. Using new attribute `partition_num` instead.",
			},
			"io_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"io_max_spec": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"io_max", "io_max_spec"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"paid_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Default:      PostPaid,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "normal",
			},
			"eip_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("deploy_type").(int) == 5
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"service_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					if old == "" {
						return false
					}
				
					oldMap := make(map[string]string)
					err := json.Unmarshal([]byte(old), &oldMap)
					if err != nil {
						return false
					}
				
					newMap := make(map[string]string)
					err = json.Unmarshal([]byte(new), &newMap)
					if err != nil {
						return false
					}
				
					// key exist in oldMap && found new value item different with old item
					for k, newValueItem := range newMap {
						oldValueItem, ok := oldMap[k]
						if ok && newValueItem != oldValueItem {
							return false
						}
					}
				
					return true
				},
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_auto_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_topic": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, false),
			},
			"default_topic_partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"selected_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
			"end_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sasl_domain_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_num_of_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"topic_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"partition_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_left": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_partition_buy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAlikafkaInstanceCreate, resourceAlibabacloudStackAlikafkaInstanceRead,
		resourceAlibabacloudStackAlikafkaInstanceUpdate, resourceAlibabacloudStackAlikafkaInstanceDelete)
	return resource
}

func resourceAlibabacloudStackAlikafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	vpcService := VpcService{client}
	var err error

	// 1. Create order
	var createOrderAction string
	createOrderResponse := make(map[string]interface{})
	createOrderReq := make(map[string]interface{})
	createOrderReq["RegionId"] = client.RegionId

	if v, ok := d.GetOk("partition_num"); ok {
		createOrderReq["PartitionNum"] = v
	} else if v, ok := d.GetOk("topic_quota"); ok {
		createOrderReq["TopicQuota"] = v
	}

	createOrderReq["DiskType"] = d.Get("disk_type")

	createOrderReq["DiskSize"] = d.Get("disk_size")

	createOrderReq["DeployType"] = d.Get("deploy_type")

	if v, ok := d.GetOk("io_max"); ok {
		createOrderReq["IoMax"] = v
	}

	if v, ok := d.GetOk("io_max_spec"); ok {
		createOrderReq["IoMaxSpec"] = v
	}

	if v, ok := d.GetOk("spec_type"); ok {
		createOrderReq["SpecType"] = v
	}

	if v, ok := d.GetOkExists("eip_max"); ok {
		createOrderReq["EipMax"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		createOrderReq["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("paid_type"); ok {
		switch v.(string) {
		case "PostPaid":
			createOrderAction = "CreatePostPayOrder"
		case "PrePaid":
			createOrderAction = "CreatePrePayOrder"
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		createOrderResponse, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", createOrderAction, "", nil, nil, createOrderReq)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL", "ONS_SYSTEM_ERROR"}) || errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(createOrderAction, createOrderResponse, createOrderReq)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alicloud_alikafka_instance", createOrderAction, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(createOrderResponse["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", createOrderAction, createOrderResponse))
	}

	alikafkaInstanceVO, err := alikafkaService.DescribeAlikafkaInstanceByOrderId(fmt.Sprint(createOrderResponse["OrderId"]), 60)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.SetId(fmt.Sprint(alikafkaInstanceVO.InstanceId))

	// 2. Start instance
	startInstanceAction := "StartInstance"
	startInstanceResponse := make(map[string]interface{})
	startInstanceReq := make(map[string]interface{})
	startInstanceReq["RegionId"] = client.RegionId
	startInstanceReq["InstanceId"] = alikafkaInstanceVO.InstanceId
	startInstanceReq["VSwitchId"] = d.Get("vswitch_id")

	if v, ok := d.GetOk("vpc_id"); ok {
		startInstanceReq["VpcId"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		startInstanceReq["ZoneId"] = v
	}

	if startInstanceReq["VpcId"] == nil {
		vsw, err := vpcService.DescribeVSwitch(startInstanceReq["VSwitchId"].(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if v, ok := startInstanceReq["VpcId"].(string); !ok || v == "" {
			startInstanceReq["VpcId"] = vsw.VpcId
		}
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		startInstanceReq["VSwitchIds"] = v
	}

	if _, ok := d.GetOkExists("eip_max"); ok {
		startInstanceReq["DeployModule"] = "eip"
		startInstanceReq["IsEipInner"] = true
	}

	if v, ok := d.GetOk("name"); ok {
		startInstanceReq["Name"] = v
	}

	if v, ok := d.GetOk("security_group"); ok {
		startInstanceReq["SecurityGroup"] = v
	}

	if v, ok := d.GetOk("service_version"); ok {
		startInstanceReq["ServiceVersion"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		startInstanceReq["Config"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		startInstanceReq["KMSKeyId"] = v
	}

	if v, ok := d.GetOk("selected_zones"); ok {
		startInstanceReq["SelectedZones"] = formatSelectedZonesReq(v.([]interface{}))
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		startInstanceResponse, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", startInstanceAction, "", nil, nil, startInstanceReq)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(startInstanceAction, startInstanceResponse, startInstanceReq)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alicloud_alikafka_instance", startInstanceAction, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(startInstanceResponse["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", startInstanceAction, startInstanceResponse))
	}

	// 3. wait until running
	stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	
	return nil

}

func resourceAlibabacloudStackAlikafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
	if err != nil {
		// Handle exceptions
		if !d.IsNewResource() && errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecp_instance alikafkaService.DescribeAliKafkaInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("disk_type", object.DiskType)
	d.Set("disk_size", object.DiskSize)
	d.Set("deploy_type", object.DeployType)
	d.Set("io_max", object.IoMax)
	d.Set("io_max_spec", object.IoMaxSpec)
	d.Set("eip_max", object.EipMax)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("zone_id", object.ZoneId)
	d.Set("paid_type", PostPaid)
	d.Set("spec_type", object.SpecType)
	d.Set("security_group", object.SecurityGroup)
	d.Set("end_point", object.EndPoint)
	d.Set("ssl_endpoint", object.SslEndPoint)
	d.Set("domain_endpoint", object.DomainEndpoint)
	d.Set("ssl_domain_endpoint", object.SslDomainEndpoint)
	d.Set("sasl_domain_endpoint", object.SaslDomainEndpoint)
	d.Set("status", object.ServiceStatus)
	// object.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO[0].Current2OpenSourceVersion can guaranteed not to be null
	d.Set("service_version", object.UpgradeServiceDetailInfo.Current2OpenSourceVersion)
	d.Set("config", object.AllConfig)
	d.Set("kms_key_id", object.KmsKeyId)
// 	d.Set("enable_auto_group", object.AutoCreateGroupEnable)
// 	d.Set("enable_auto_topic", convertAliKafkaAutoCreateTopicEnableResponse(object.AutoCreateTopicEnable))
// 	d.Set("default_topic_partition_num", formatInt(object.DefaultPartitionNum))
// 
// 	if vSwitchIds, ok := object.VSwitchIds; ok {
// 		vSwitchIdsArg := vSwitchIds.(map[string]interface{})
// 
// 		if vSwitchIdsList, ok := vSwitchIdsArg["VSwitchIds"]; ok {
// 			d.Set("vswitch_ids", vSwitchIdsList)
// 		}
// 	}

	if fmt.Sprint(object.PaidType) == "0" {
		d.Set("paid_type", PrePaid)
	}

	quota, err := alikafkaService.GetQuotaTip(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("topic_quota", quota["TopicQuota"])
	d.Set("partition_num", quota["PartitionNumOfBuy"])
	d.Set("topic_num_of_buy", quota["TopicNumOfBuy"])
	d.Set("topic_used", quota["TopicUsed"])
	d.Set("topic_left", quota["TopicLeft"])
	d.Set("partition_used", quota["PartitionUsed"])
	d.Set("partition_left", quota["PartitionLeft"])
	d.Set("group_used", quota["GroupUsed"])
	d.Set("group_left", quota["GroupLeft"])
	d.Set("is_partition_buy", quota["IsPartitionBuy"])

	tags, err := alikafkaService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("tags", alikafkaService.tagsToMap(tags))

	return nil
}

func resourceAlibabacloudStackAlikafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if err := alikafkaService.setInstanceTags(d, TagResourceInstance); err != nil {
		return errmsgs.WrapError(err)
	}

	// Process change instance name.
	if !d.IsNewResource() && d.HasChange("name") {
		action := "ModifyInstanceName"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("name"); ok {
			request["InstanceName"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}	}

	// Process paid type change, note only support change from post to pre pay.
	if !d.IsNewResource() && d.HasChange("paid_type") {
		o, n := d.GetChange("paid_type")
		oldPaidType := o.(string)
		newPaidType := n.(string)
		oldPaidTypeInt := 1
		newPaidTypeInt := 1
		if oldPaidType == string(PrePaid) {
			oldPaidTypeInt = 0
		}
		if newPaidType == string(PrePaid) {
			newPaidTypeInt = 0
		}
		if oldPaidTypeInt == 1 && newPaidTypeInt == 0 {
			action := "ConvertPostPayOrder"
			request := map[string]interface{}{
				"RegionId":   client.RegionId,
				"InstanceId": d.Id(),
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
				if err != nil {
					if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)

			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
			}

			stateConf := BuildStateConf([]string{}, []string{strconv.Itoa(newPaidTypeInt)}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "PaidType", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}

		} else {
			return errmsgs.WrapError(errors.New("paid type only support change from post pay to pre pay"))
		}	
	}

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}
	// updating topic_quota only by updating partition_num
	if !d.IsNewResource() && d.HasChange("partition_num") {
		update = true
	}
	request["PartitionNum"] = d.Get("partition_num")

	if !d.IsNewResource() && d.HasChange("disk_size") {
		update = true
	}
	request["DiskSize"] = d.Get("disk_size")

	if !d.IsNewResource() && d.HasChange("io_max") {
		update = true

		if v, ok := d.GetOk("io_max"); ok {
			request["IoMax"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("io_max_spec") {
		update = true

		if v, ok := d.GetOk("io_max_spec"); ok {
			request["IoMaxSpec"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("spec_type") {
		update = true
	}
	request["SpecType"] = d.Get("spec_type")

	if !d.IsNewResource() && d.HasChange("deploy_type") {
		update = true
	}
	if d.Get("deploy_type").(int) == 4 {
		request["EipModel"] = true
	} else {
		request["EipModel"] = false
	}
	if !d.IsNewResource() && d.HasChange("eip_max") {
		update = true
	}
	request["EipMax"] = d.Get("eip_max").(int)

	if update {
		action := "UpgradePostPayOrder"

		if d.Get("paid_type").(string) == string(PrePaid) {
			action = "UpgradePrePayOrder"
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
			if err != nil {
				if errmsgs.NeedRetry(err) || errmsgs.IsExpectedErrors(err, []string{"ONS_SYSTEM_FLOW_CONTROL"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("disk_size"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "DiskSize", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if d.HasChange("io_max") {
			stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("io_max"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "IoMax", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}

		stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("eip_max"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "EipMax", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		stateConf = BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("spec_type"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "SpecType", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		stateConf = BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}	}

	if !d.IsNewResource() && d.HasChange("service_version") {
		action := "UpgradeInstanceVersion"
		request := map[string]interface{}{
			"InstanceId": d.Id(),
			"RegionId":   client.RegionId,
		}

		if v, ok := d.GetOk("service_version"); ok {
			request["TargetVersion"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				// means no need to update version
				if errmsgs.IsExpectedErrors(err, []string{"ONS_INIT_ENV_ERROR"}) {
					return nil
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		time.Sleep(60 * time.Second)
		// upgrade service may be last a long time
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}	}

	if !d.IsNewResource() && d.HasChange("config") {
		action := "UpdateInstanceConfig"
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": d.Id(),
		}

		if v, ok := d.GetOk("config"); ok {
			request["Config"] = v
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		// wait for upgrade task be invoke
		stateConf := BuildStateConf([]string{}, []string{"5"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}	}

	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ResourceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ChangeResourceGroup"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, changeResourceGroupReq)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}	}

	update = false
	enableAutoGroupCreationReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": d.Id(),
	}

	if d.HasChange("enable_auto_group") {
		update = true

		if v, ok := d.GetOkExists("enable_auto_group"); ok {
			enableAutoGroupCreationReq["Enable"] = v
		}
	}

	if update {
		action := "EnableAutoGroupCreation"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, enableAutoGroupCreationReq)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, enableAutoGroupCreationReq)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}	}

	update = false
	enableAutoTopicCreationReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"InstanceId": d.Id(),
	}

	if d.HasChange("enable_auto_topic") {
		update = true
	}
	if v, ok := d.GetOk("enable_auto_topic"); ok {
		enableAutoTopicCreationReq["Operate"] = v
	}

	if update {
		action := "EnableAutoTopicCreation"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, enableAutoTopicCreationReq)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, enableAutoTopicCreationReq)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}	}

	update = false
	updateTopicPartitionNumReq := map[string]interface{}{
		"RegionId":        client.RegionId,
		"Operate":         "updatePartition",
		"UpdatePartition": true,
		"InstanceId":      d.Id(),
	}

// 	object, err := alikafkaService.DescribeAlikafkaInstance(d.Id())
// 	if err != nil {
// 		return errmsgs.WrapError(err)
// 	}

// 	defaultTopicPartitionNum, ok := d.GetOkExists("default_topic_partition_num")
// 	if ok && fmt.Sprint(object.DefaultPartitionNum) != fmt.Sprint(defaultTopicPartitionNum) {
// 		update = true
// 		updateTopicPartitionNumReq["PartitionNum"] = defaultTopicPartitionNum
// 	}

	if update {
		action := "EnableAutoTopicCreation"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, updateTopicPartitionNumReq)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action+" updateTopicPartitionNum", response, updateTopicPartitionNumReq)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}	}

	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackAlikafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	action := "ReleaseInstance"
	var err error
	var response map[string]interface{}
	request := map[string]interface{}{
		"InstanceId":          d.Id(),
		"RegionId":            client.RegionId,
		"ForceDeleteInstance": true,
	}

	// Pre paid instance can not be release.
	if d.Get("paid_type").(string) == string(PrePaid) {
		return nil
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) || errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"15"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	action = "DeleteInstance"
	request = map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
	}

	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "alikafka", "2019-09-16", action, "", nil, nil, request)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf = BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, alikafkaService.AliKafkaInstanceStateRefreshFunc(d.Id(), "ServiceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func formatSelectedZonesReq(configured []interface{}) string {
	doubleList := make([][]interface{}, len(configured))
	for i, v := range configured {
		doubleList[i] = []interface{}{v}
	}

	if len(doubleList) < 1 {
		return ""
	}

	if len(doubleList) == 1 {
		return "[[\"" + doubleList[0][0].(string) + "\"],[]]"
	}

	result := "[["

	for i := 0; i < len(doubleList); i++ {
		switch i {
		case len(doubleList) - 2:
			result += "\"" + doubleList[i][0].(string) + "\""
		case len(doubleList) - 1:
			result += "],[\"" + doubleList[i][0].(string) + "\"]"
		default:
			result += "\"" + doubleList[i][0].(string) + "\","
		}
	}

	result += "]"

	return result
}

func convertAliKafkaAutoCreateTopicEnableResponse(source interface{}) interface{} {
	switch source {
	case true:
		return "enable"
	case false:
		return "disable"
	}

	return source
}
