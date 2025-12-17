package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCspprivateHsmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCspprivateHsmInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by instance name (Remark).",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance names (Remark) corresponding to the returned instances.",
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the resource.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the HSM instance.",
						},
						"hsm_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the HSM instance. 1: uninitialized, 3: released, 4: failed, 5: running, 6: syncing, 7: resetting, 8: disabled.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID assigned to the HSM instance.",
						},
						"vswitch_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vSwitch ID assigned to the HSM instance.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The classic network IP address associated with the HSM instance.",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias or alias_name of the HSM instance.",
						},
						"product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product model code of the device used by the HSM instance.",
						},
						"vendor_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor code of the device used by the HSM instance.",
						},
						"vsm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the HSM instance. evsm: financial data HSM, gvsm: general server HSM, svsm: signature verification server HSM.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone ID where the HSM instance is located.",
						},
					},
				},
				Description: "A list of HSM instances matching the filter criteria.",
			},
		},
	}
}

func dataSourceAlibabacloudStackCspprivateHsmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	pageNum := 1
	request["PageSize"] = 100
	request["PageNumber"] = pageNum
	var allInstances []map[string]interface{}
	for true {
		response, err := client.DoTeaRequest("GET", "Cspprivate", "2022-02-17", "DescribeHsms", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		instancesRaw, err := jsonpath.Get("$.HsmInfos.HsmInfo", response)
		if err != nil || len(instancesRaw.([]interface{})) == 0 {
			return errmsgs.WrapError(err)
		}

		totalCountRaw, _ := response["TotalCount"].(json.Number).Int64()

		totalCount := int(totalCountRaw)

		for _, inst := range instancesRaw.([]interface{}) {
			allInstances = append(allInstances, inst.(map[string]interface{}))
		}

		if pageNum*100 >= totalCount {
			break
		}
		request["PageNumber"] = pageNum + 1
	}

	// Process filtering
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	names := make([]string, 0)
	ids := make([]string, 0)
	result := make([]map[string]interface{}, 0)
	for _, inst := range allInstances {
		instanceId := inst["HsmId"].(string)
		alias_name := inst["AliasName"].(string)
		if len(idsMap) > 0 {
			if _, exists := idsMap[instanceId]; !exists {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(alias_name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":           instanceId,
			"instance_id":  instanceId,
			"hsm_status":   inst["InstanceStatus"],
			"vpc_id":       inst["VpcId"],
			"vswitch_id":   inst["VswitchId"],
			"ip":           inst["HsmIp"],
			"alias_name":   inst["AliasName"],
			"vendor_code":  inst["VendorCode"],
			"product_code": inst["ProductCode"],
			"vsm_type":     inst["VsmType"],
			"zone_id":      inst["ZoneNo"],
		}
		names = append(names, alias_name)
		ids = append(ids, instanceId)
		result = append(result, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("instances", result); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
