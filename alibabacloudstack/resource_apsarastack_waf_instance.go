package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
			},
			"wafinstance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_make_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudstackWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// var err error
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
	// vpcvswitchitem := map[string]interface{}{}
	// // if value, ok := d.GetOk("vpc_vswitch"); ok {
	// for _, vpcvswitchinfo := range d.Get("vpc_vswitch").([]interface{}) {
	// 	vswitchMap := vpcvswitchinfo.(map[string]interface{})
	// 	vpcvswitchitem["vswitch_name"] = vswitchMap["vswitch_name"]
	// 	vpcvswitchitem["vswitch"] = vswitchMap["vswitch"]
	// 	vpcvswitchitem["cidr_block"] = vswitchMap["cidr_block"]
	// 	vpcvswitchitem["available_zone"] = vswitchMap["available_zone"]
	// 	vpcvswitchitem["vpc"] = vswitchMap["vpc"]
	// 	vpcvswitchitem["vpc_name"] = vswitchMap["vpc_name"]
	// }
	// }
	// }
	if value, ok := d.GetOk("vpc_vswitch"); ok {
		vpcvswitchsMappings := value.([]interface{})
		if vpcvswitchsMappings != nil && len(vpcvswitchsMappings) > 0 {
			mappings := make([]string, 0, len(vpcvswitchsMappings))
			for _, diskDeviceMapping := range vpcvswitchsMappings {
				mapping := diskDeviceMapping.(map[string]interface{})
				vpcvswithMapping := WafVPCVSwitch{
					VSwitchName:   mapping["vswitch_name"].(string),
					VSwitch:       mapping["vswitch"].(string),
					CIDRBlock:     mapping["cidr_block"].(string),
					AvailableZone: mapping["available_zone"].(string),
					VPC:           mapping["vpc"].(string),
					VPCName:       mapping["vpc_name"].(string),
				}
				jsonBytes, err := json.Marshal(vpcvswithMapping)
				if err != nil {
					fmt.Println("Error marshalling to JSON:", err)
					return errmsgs.WrapError(err)
				}
				jsonStr := string(jsonBytes)
				mappings = append(mappings, jsonStr)
			}
			request["VpcVswitch"] = mappings
		}
	}
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
	// jsonBytes, err := json.Marshal(vpcvswitchitem)
	// if err != nil {
	// 	fmt.Println("Error marshalling to JSON:", err)
	// 	return errmsgs.WrapError(err)
	// }
	// jsonStr := string(jsonBytes)
	// vpcvswitch := []string{jsonStr}

	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, nil, request)
	addDebug(action, response, request)
	if err != nil {
		return err
	}
	response = response["Result"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["instance_id"]))
	WafInstanceService := WafOpenapiService{client}
	stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, WafInstanceService.Wafv3InstanceStateRefreshFunc(d.Id(), []string{"faild"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, fmt.Sprint(response["instance_id"]))
	}

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
	d.Set("name", object["name"])
	d.Set("instance_status", object["instance_status"])
	d.Set("vpc_vswitch", object["vpc_vswitch"])
	d.Set("detector_version", object["detector_version"])
	d.Set("instance_make_status", object["instance_make_status"])
	d.Set("detector_specs", object["detector_specs"])
	d.Set("detector_nodenum", object["detector_node_num"])
	return nil
}
func resourceAlibabacloudstackWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"WafInstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("detector_nodenum") {
		update = true
	}
	if update {
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
		request["WafInstanceId"] = object["instance_id"]
		request["InstanceId"] = object["instance_id"]
		vpcvswitchsMappings := object["vpc_vswitch"].([]interface{})
		if vpcvswitchsMappings != nil && len(vpcvswitchsMappings) > 0 {
			mappings := make([]string, 0, len(vpcvswitchsMappings))
			for _, diskDeviceMapping := range vpcvswitchsMappings {
				mapping := diskDeviceMapping.(map[string]interface{})
				vpcvswithMapping := WafVPCVSwitch{
					VSwitchName:   mapping["vswitch_name"].(string),
					VSwitch:       mapping["vswitch"].(string),
					CIDRBlock:     mapping["cidr_block"].(string),
					AvailableZone: mapping["available_zone"].(string),
					VPC:           mapping["vpc"].(string),
					VPCName:       mapping["vpc_name"].(string),
				}
				jsonBytes, err := json.Marshal(vpcvswithMapping)
				if err != nil {
					fmt.Println("Error marshalling to JSON:", err)
					return errmsgs.WrapError(err)
				}
				jsonStr := string(jsonBytes)
				mappings = append(mappings, jsonStr)
			}
			request["VpcVswitch"] = mappings
		}
		request["DetectorNodeNum"] = d.Get("detector_nodenum")
		detectorNodeNumObject, ok := object["detector_node_num"].(json.Number)
		if !ok {
			return fmt.Errorf("failed to convert object[\"detector_node_num\"] to json.Number")
		}
		detectorNodeNumInt, err := detectorNodeNumObject.Int64()
		if err != nil {
			return fmt.Errorf("failed to convert json.Number to int: %v", err)
		}
		dgetnodenumberobject := int64(d.Get("detector_nodenum").(float64))
		if !ok {
			return fmt.Errorf("failed to convert d.Get(\"detector_node_num\") to int")
		}
		log.Printf("detectorNodeNumObject is %d  %d", detectorNodeNumInt, dgetnodenumberobject)
		if detectorNodeNumInt > dgetnodenumberobject {
			prescaledown := "CreatePreScaleDownInstance"
			scaledown := "CreateScaleDownInstance"
			response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", prescaledown, "", nil, nil, request)
			if err != nil {
				return err
			}
			addDebug(prescaledown, response, request)
			response, err = client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", scaledown, "", nil, nil, request)
			if err != nil {
				return err
			}
			addDebug(scaledown, response, request)
		} else {
			scaleup := "CreateScaleUpInstance"
			response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", scaleup, "", nil, nil, request)
			if err != nil {
				return err
			}
			addDebug(scaleup, response, request)
		}
		stateConf := BuildStateConf([]string{}, []string{"success"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, waf_openapiService.Wafv3InstanceStateRefreshFunc(d.Id(), []string{"faild"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, fmt.Sprint(d.Id()))
		}
	}
	return resourceAlibabacloudstackWafInstanceRead(d, meta)
}
func resourceAlibabacloudstackWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteWAFInstance"
	var err error
	request := map[string]interface{}{
		"WafInstanceId": d.Id(),
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	response, err := client.DoTeaRequest("POST", "waf-onecs", "2020-07-01", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	addDebug(action, response, request)
	return nil
}
