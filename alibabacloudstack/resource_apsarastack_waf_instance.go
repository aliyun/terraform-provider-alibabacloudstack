package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudstackWafInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudstackWafInstanceCreate,
		Read:   resourceAlibabacloudstackWafInstanceRead,
		Update: resourceAlibabacloudstackWafInstanceUpdate,
		Delete: resourceAlibabacloudstackWafInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// "vpc_vswitch": {
			// 	Type:     schema.TypeList,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// 	Required: true,
			// 	ForceNew: true,
			// },
			"vpc_vswitch": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"vswitch": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"available_zone": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"vpc": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},
			"detector_specs": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"detector_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"detector_nodenum": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"wafinstance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudstackWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var err error
	action := "CreateWAFInstance"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}

	request["ProductCode"] = "waf"
	request["ProductType"] = "waf"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}

	if v, ok := d.GetOk("arch"); ok {
		request["Arch"] = v
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		request["CpuType"] = v
	}
	// if v, ok := d.GetOk("vpcVswitch"); ok {
	// 	request["VpcVswitch"] = v
	// }
	if v, ok := d.GetOk("detector_specs"); ok {
		request["DetectorSpecs"] = v
	}
	if v, ok := d.GetOk("detector_version"); ok {
		request["DetectorVersion"] = v
	}
	if v, ok := d.GetOk("detector_nodenum"); ok {
		request["DetectorNodeNum"] = v
	}
	vpcvswitchitem := map[string]interface{}{}
	// if value, ok := d.GetOk("vpc_vswitch"); ok {
	for _, vpcvswitchinfo := range d.Get("vpc_vswitch").([]interface{}) {
		vswitchMap := vpcvswitchinfo.(map[string]interface{})
		vpcvswitchitem["vswitch_name"] = vswitchMap["vswitch_name"]
		vpcvswitchitem["vswitch"] = vswitchMap["vswitch"]
		vpcvswitchitem["cidr_block"] = vswitchMap["cidr_block"]
		vpcvswitchitem["available_zone"] = vswitchMap["available_zone"]
		vpcvswitchitem["vpc"] = vswitchMap["vpc"]
		vpcvswitchitem["vpc_name"] = vswitchMap["vpc_name"]
	}
	// }
	// }
	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	// if v, ok := d.GetOk("region"); ok && v.(string) != "" {
	// 	region = v.(string)
	// }
	// request["SubscriptionType"] = d.Get("subscription_type")
	// request["Parameter"] = []map[string]string{
	// 	{
	// 		"Code":  "BigScreen",
	// 		"Value": d.Get("big_screen").(string),
	// 	},
	// 	{
	// 		"Code":  "ExclusiveIpPackage",
	// 		"Value": d.Get("exclusive_ip_package").(string),
	// 	},
	// 	{
	// 		"Code":  "ExtBandwidth",
	// 		"Value": d.Get("ext_bandwidth").(string),
	// 	},
	// 	{
	// 		"Code":  "ExtDomainPackage",
	// 		"Value": d.Get("ext_domain_package").(string),
	// 	},
	// 	{
	// 		"Code":  "LogStorage",
	// 		"Value": d.Get("log_storage").(string),
	// 	},
	// 	{
	// 		"Code":  "LogTime",
	// 		"Value": d.Get("log_time").(string),
	// 	},
	// 	{
	// 		"Code":  "PackageCode",
	// 		"Value": d.Get("package_code").(string),
	// 	},
	// 	{
	// 		"Code":  "PrefessionalService",
	// 		"Value": d.Get("prefessional_service").(string),
	// 	},
	// {
	// 	"Code":  "Region",
	// 	"Value": region,
	// },
	// 	{
	// 		"Code":  "WafLog",
	// 		"Value": d.Get("waf_log").(string),
	// 	},
	// }
	if v, ok := d.GetOk("vpc_id"); ok {
		request["Vpc"] = v
		// vpcService := VpcService{client}
		// object_vpc, err := vpcService.DescribeVpc(v.(string))
		// if err != nil {
		// 	if errmsgs.NotFoundError(err) {
		// 		return nil
		// 	}
		// 	return errmsgs.WrapError(err)
		// }
		// vpcvswitchitem["vpc_name"] = object_vpc.VpcName
		// vpcvswitchitem["vpc"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["Vswitch1"] = v
		// vpcService := VpcService{client}
		// object_vswitch, err := vpcService.DescribeVSwitch(v.(string))
		// if err != nil {
		// 	if errmsgs.NotFoundError(err) {
		// 		return nil
		// 	}
		// 	return errmsgs.WrapError(err)
		// }
		// vpcvswitchitem["vswitch_name"] = object_vswitch.VSwitchName
		// vpcvswitchitem["cidr_block"] = object_vswitch.CidrBlock
		// vpcvswitchitem["vswitch"] = v.(string)
		// vpcvswitchitem["available_zone"] = object_vswitch.ZoneId
	}
	jsonBytes, err := json.Marshal(vpcvswitchitem)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return errmsgs.WrapError(err)
	}
	jsonStr := string(jsonBytes)
	vpcvswitch := []string{jsonStr}
	request["vpcVswitch"] = vpcvswitch
	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
	if err != nil {
		return err
	}
	addDebug(action, response, request)
	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	return resourceAlibabacloudstackWafInstanceUpdate(d, meta)
}
func resourceAlibabacloudstackWafInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	waf_openapiService := WafOpenapiService{client}
	object, err := waf_openapiService.DescribeWafInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_waf_instance waf_openapiService.DescribeWafInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("status", formatInt(object["Status"]))
	d.Set("subscription_type", object["SubscriptionType"])
	return nil
}
func resourceAlibabacloudstackWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("subscription_type") {
		update = true
	}
	request["ProductType"] = "waf"
	request["SubscriptionType"] = d.Get("subscription_type")
	request["ProductCode"] = "waf"
	request["ModifyType"] = d.Get("modify_type")
	request["Parameter"] = []map[string]string{
		{
			"Code":  "BigScreen",
			"Value": d.Get("big_screen").(string),
		},
		{
			"Code":  "ExclusiveIpPackage",
			"Value": d.Get("exclusive_ip_package").(string),
		},
		{
			"Code":  "ExtBandwidth",
			"Value": d.Get("ext_bandwidth").(string),
		},
		{
			"Code":  "ExtDomainPackage",
			"Value": d.Get("ext_domain_package").(string),
		},
		{
			"Code":  "LogStorage",
			"Value": d.Get("log_storage").(string),
		},
		{
			"Code":  "LogTime",
			"Value": d.Get("log_time").(string),
		},
		{
			"Code":  "PackageCode",
			"Value": d.Get("package_code").(string),
		},
		{
			"Code":  "PrefessionalService",
			"Value": d.Get("prefessional_service").(string),
		},
		{
			"Code":  "WafLog",
			"Value": d.Get("waf_log").(string),
		},
	}
	if update {
		action := "EditWAFInstance"
		response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
		if err != nil {
			return err
		}
		addDebug(action, response, request)
	}
	return resourceAlibabacloudstackWafInstanceRead(d, meta)
}
func resourceAlibabacloudstackWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteWAFInstance"
	var err error
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, request)
	if err != nil {
		return err
	}
	addDebug(action, response, request)
	return nil
}
