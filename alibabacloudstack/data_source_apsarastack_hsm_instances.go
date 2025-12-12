package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackHsmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackHsmInstancesRead,

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
			"vsm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the HSM instance. Possible values: evsm, gvsm, svsm.",
			},
			"zone_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone ID where the HSM instance is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the HSM instance.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cluster ID to which the HSM instance belongs.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the HSM instance. Multiple statuses can be separated by commas.",
			},
			"vpc_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The classic network IP address registered by the HSM instance.",
			},
			"enable_details": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to retrieve detailed information for each instance.",
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
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias or remark of the HSM instance.",
						},
						"vendor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor name of the device used by the HSM instance.",
						},
						"vendor_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vendor code of the device used by the HSM instance.",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product model name of the device used by the HSM instance.",
						},
						"product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product model code of the device used by the HSM instance.",
						},
						"vsm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the HSM instance. evsm: financial data HSM, gvsm: general server HSM, svsm: signature verification server HSM.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the cluster to which the HSM instance belongs.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cluster to which the HSM instance belongs.",
						},
						"is_master": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates whether the current HSM instance is the master in the cluster. 0: no, 1: yes.",
						},
						"zone_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone ID where the HSM instance is located.",
						},
						"show_create_cluster": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether cluster creation is allowed for this HSM instance.",
						},
						"hsm_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The physical HSM ID associated with the instance.",
						},
						"release_protection": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates whether release protection is enabled. 0: disabled, 1: enabled.",
						},
						"white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The whitelist IP list of the HSM instance. If the instance belongs to a cluster, the whitelist will be consistent with the cluster's.",
						},
					},
				},
				Description: "A list of HSM instances matching the filter criteria.",
			},
		},
	}
}

func dataSourceAlibabacloudStackHsmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("vsm_type"); ok {
		request["VsmType"] = v
	}
	if v, ok := d.GetOk("zone_no"); ok {
		request["ZoneNo"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("vpc_ip"); ok {
		request["VpcIp"] = v
	}

	request["PageSize"] = 100
	request["CurrentPage"] = 1
	var allInstances []map[string]interface{}
	for true {
		response, err := client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeInstances", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		instancesRaw, ok := response["Instances"].([]interface{})
		if !ok {
			break
		}

		totalCountRaw, _ := response["TotalCount"].(json.Number).Int64()
		currentPageRaw, _ := response["CurrentPage"].(json.Number).Int64()
		pageSizeRaw, _ := response["PageSize"].(json.Number).Int64()

		totalCount := int(totalCountRaw)
		currentPage := int(currentPageRaw)
		pageSize := int(pageSizeRaw)

		for _, inst := range instancesRaw {
			allInstances = append(allInstances, inst.(map[string]interface{}))
		}

		if currentPage*pageSize >= totalCount {
			break
		}
		request["CurrentPage"] = currentPage + 1
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
		instanceId := inst["InstanceId"].(string)
		remark := inst["Remark"].(string)
		if len(idsMap) > 0 {
			if _, exists := idsMap[instanceId]; !exists {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(remark) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                  inst["InstanceId"],
			"instance_id":         inst["InstanceId"],
			"hsm_status":          inst["HsmStatus"],
			"vpc_id":              inst["VpcId"],
			"vswitch_id":          inst["VswitchId"],
			"ip":                  inst["Ip"],
			"remark":              inst["Remark"],
			"vendor_name":         inst["VendorName"],
			"vendor_code":         inst["VendorCode"],
			"product_name":        inst["ProductName"],
			"product_code":        inst["ProductCode"],
			"vsm_type":            inst["VsmType"],
			"cluster_id":          inst["ClusterId"],
			"cluster_name":        inst["ClusterName"],
			"is_master":           inst["IsMaster"],
			"zone_no":             inst["ZoneNo"],
			"show_create_cluster": inst["ShowCreateCluster"],
			"hsm_id":              inst["HsmId"],
			"release_protection":  inst["ReleaseProtection"],
			"white_list":          inst["WhiteList"],
		}
		names = append(names, remark)
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
