package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEcsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEcsDedicatedHostCreate,
		Read:   resourceAlibabacloudStackEcsDedicatedHostRead,
		Update: resourceAlibabacloudStackEcsDedicatedHostUpdate,
		Delete: resourceAlibabacloudStackEcsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_on_maintenance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stop", "Migrate"}, false),
				Default:      "Stop",
			},
			"auto_placement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "on",
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_over_commit_ratio": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"dedicated_host_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail_fee": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"network_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"slb_udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				Default:      "PostPaid",
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sale_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackEcsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AllocateDedicatedHosts"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("action_on_maintenance"); ok {
		request["ActionOnMaintenance"] = v
	}

	if v, ok := d.GetOk("auto_placement"); ok {
		request["AutoPlacement"] = v
	}

	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	if v, ok := d.GetOk("cpu_over_commit_ratio"); ok {
		request["CpuOverCommitRatio"] = v
	}

	if v, ok := d.GetOk("dedicated_host_cluster_id"); ok {
		request["DedicatedHostClusterId"] = v
	}

	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request["DedicatedHostName"] = v
	}

	request["DedicatedHostType"] = d.Get("dedicated_host_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("expired_time"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("min_quantity"); ok {
		request["MinQuantity"] = v
	}

	if v, ok := d.GetOk("network_attributes"); ok {
		networkAttributesMap := make(map[string]interface{})
		for _, networkAttributes := range v.(*schema.Set).List() {
			networkAttributesArg := networkAttributes.(map[string]interface{})
			networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
			networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
		}
		request["NetworkAttributes"] = networkAttributesMap
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}

	request["Quantity"] = 1

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("sale_cycle"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	response, err = client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	responseDedicatedHostIdSets := response["DedicatedHostIdSets"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseDedicatedHostIdSets["DedicatedHostId"].([]interface{})[0]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{"PermanentFailure"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackEcsDedicatedHostRead(d, meta)
}

func resourceAlibabacloudStackEcsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDedicatedHost(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_dedicated_host ecsService.DescribeEcsDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("action_on_maintenance", object.ActionOnMaintenance)
	d.Set("auto_placement", object.AutoPlacement)
	d.Set("auto_release_time", object.AutoReleaseTime)
	d.Set("cpu_over_commit_ratio", object.CpuOverCommitRatio)
	d.Set("dedicated_host_name", object.DedicatedHostName)
	d.Set("dedicated_host_type", object.DedicatedHostType)
	d.Set("description", object.Description)
	d.Set("expired_time", object.ExpiredTime)
	d.Set("network_attributes", object.NetworkAttributes)
	d.Set("payment_type", object.ChargeType)
	d.Set("resource_group_id", object.ResourceGroupId)
	d.Set("sale_cycle", object.SaleCycle)
	d.Set("status", object.Status)
	d.Set("tags", object.Tags)
	d.Set("zone_id", object.ZoneId)
	return nil
}

func resourceAlibabacloudStackEcsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "ddh"); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("tags")
	}
	if !d.IsNewResource() && d.HasChange("auto_release_time") {
		request := map[string]interface{}{
			"DedicatedHostId": d.Id(),
		}
		request["AutoReleaseTime"] = d.Get("auto_release_time")
		action := "ModifyDedicatedHostAutoReleaseTime"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		//d.SetPartial("auto_release_time")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"ResourceId": d.Id(),
			"ResourceType": "ddh",
			"ResourceGroupId": d.Get("resource_group_id"),
		}
		action := "JoinResourceGroup"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		//d.SetPartial("resource_group_id")
	}
	update := false
	request := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	if !d.IsNewResource() && d.HasChange("expired_time") {
		update = true
		request["Period"] = d.Get("expired_time")
	}
	if !d.IsNewResource() && d.HasChange("sale_cycle") {
		update = true
		request["PeriodUnit"] = d.Get("sale_cycle")
	}
	if update {
		action := "RenewDedicatedHosts"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		//d.SetPartial("expired_time")
		//d.SetPartial("sale_cycle")
	}
	update = false
	modifyDedicatedHostsChargeTypeReq := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
		"AutoPay":          true,
		"Period":           d.Get("expired_time"),
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		modifyDedicatedHostsChargeTypeReq["DedicatedHostChargeType"] = d.Get("payment_type")
	}
	modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = d.Get("sale_cycle")
	if update {
		if _, ok := d.GetOkExists("detail_fee"); ok {
			modifyDedicatedHostsChargeTypeReq["DetailFee"] = d.Get("detail_fee")
		}
		if _, ok := d.GetOkExists("dry_run"); ok {
			modifyDedicatedHostsChargeTypeReq["DryRun"] = d.Get("dry_run")
		}
		action := "ModifyDedicatedHostsChargeType"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, modifyDedicatedHostsChargeTypeReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		// d.SetPartial("detail_fee")
		// d.SetPartial("dry_run")
		// d.SetPartial("expired_time")
		// d.SetPartial("payment_type")
		// d.SetPartial("sale_cycle")
	}
	update = false
	modifyDedicatedHostAttributeReq := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("action_on_maintenance") {
		update = true
		modifyDedicatedHostAttributeReq["ActionOnMaintenance"] = d.Get("action_on_maintenance")
	}
	if !d.IsNewResource() && d.HasChange("auto_placement") {
		update = true
		modifyDedicatedHostAttributeReq["AutoPlacement"] = d.Get("auto_placement")
	}
	if !d.IsNewResource() && d.HasChange("cpu_over_commit_ratio") {
		update = true
		modifyDedicatedHostAttributeReq["CpuOverCommitRatio"] = d.Get("cpu_over_commit_ratio")
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_name") {
		update = true
		modifyDedicatedHostAttributeReq["DedicatedHostName"] = d.Get("dedicated_host_name")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyDedicatedHostAttributeReq["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("network_attributes") {
		update = true
		if d.Get("network_attributes") != nil {
			networkAttributesMap := make(map[string]interface{})
			for _, networkAttributes := range d.Get("network_attributes").(*schema.Set).List() {
				networkAttributesArg := networkAttributes.(map[string]interface{})
				networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
				networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
			}
			modifyDedicatedHostAttributeReq["NetworkAttributes"] = networkAttributesMap
		}
	}
	if update {
		if _, ok := d.GetOk("dedicated_host_cluster_id"); ok {
			modifyDedicatedHostAttributeReq["DedicatedHostClusterId"] = d.Get("dedicated_host_cluster_id")
		}
		action := "ModifyDedicatedHostAttribute"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, modifyDedicatedHostAttributeReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		// d.SetPartial("action_on_maintenance")
		// d.SetPartial("auto_placement")
		// d.SetPartial("cpu_over_commit_ratio")
		// d.SetPartial("dedicated_host_cluster_id")
		// d.SetPartial("dedicated_host_name")
		// d.SetPartial("description")
		// d.SetPartial("network_attributes")
	}
	d.Partial(false)
	return resourceAlibabacloudStackEcsDedicatedHostRead(d, meta)
}

func resourceAlibabacloudStackEcsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "ReleaseDedicatedHost"
	request := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
