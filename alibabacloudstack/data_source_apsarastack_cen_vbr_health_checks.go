package alibabacloudstack

import (
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCenVbrHealthChecks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCenVbrHealthChecksRead,

		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vbr_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vbr_health_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vbr_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vbr_instance_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_target_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCenVbrHealthChecksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	reqQuery := make(map[string]interface{})
	reqQuery["VbrInstanceRegionId"] = client.RegionId
	if v, ok := d.GetOk("cen_id"); ok {
		reqQuery["CenId"] = v
	}
	reqQuery["VbrInstanceRegionId"] = client.RegionId
	if v, ok := d.GetOk("vbr_instance_id"); ok {
		reqQuery["VbrInstanceId"] = v
	}

	// Call API
	resp, err := client.DoTeaRequest("POST", "Cbn", "2017-09-12", "DescribeCenVbrHealthCheck", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cen_vbr_health_checks", "DescribeCenVbrHealthCheck", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Extract VbrHealthChecks array
	vbrHealthChecksData, ok := resp["VbrHealthChecks"].(map[string]interface{})
	if !ok {
		return nil
	}
	vbrHealthChecksList, ok := vbrHealthChecksData["VbrHealthCheck"].([]interface{})
	if !ok {
		return nil
	}

	// Handle ids filter
	idsMap := getIdsStringFilter(d)
	ids := make([]string, 0)
	checks := make([]map[string]interface{}, 0)

	for _, item := range vbrHealthChecksList {
		check := item.(map[string]interface{})
		key := strings.Join([]string{check["CenId"].(string), check["VbrInstanceId"].(string)}, ":")
		if len(idsMap) > 0 {
			if _, ok := idsMap[key]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{
			"cen_id":                 check["CenId"],
			"vbr_instance_id":        check["VbrInstanceId"],
			"vbr_instance_region_id": check["VbrInstanceRegionId"],
			"health_check_source_ip": check["HealthCheckSourceIp"],
			"health_check_target_ip": check["HealthCheckTargetIp"],
			"health_check_interval":  check["HealthCheckInterval"],
			"healthy_threshold":      check["HealthyThreshold"],
			"health_check_only":      check["HealthCheckOnly"],
		}
		ids = append(ids, key)
		checks = append(checks, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("vbr_health_checks", checks); err != nil {
		return errmsgs.WrapError(err)
	}

	// Write to output file if specified
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), checks)
	}

	return nil
}
