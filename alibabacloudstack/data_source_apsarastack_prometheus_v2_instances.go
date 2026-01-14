package alibabacloudstack

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPrometheusV2Instances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPrometheusV2InstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of instance IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by cluster name.",
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Prometheus V2 instance.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Prometheus V2 instance.",
						},
						"http_api": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The HTTP API endpoint of the Prometheus V2 instance.",
						},
						"remote_write_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remote write URL of the Prometheus V2 instance.",
						},
						"push_gateway_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The push gateway URL of the Prometheus V2 instance.",
						},
						"objid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The internal object ID of the Prometheus V2 instance.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the Prometheus V2 instance.",
						},
						"security_level_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security level tag of the instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current status of the instance.",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of tags assigned to the instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPrometheusV2InstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	prometheusService := PrometheusService{client}
	now := time.Now()
	endTime := now.Format("2006-01-02 15:04:05")
	startTime := now.Add(-5 * time.Minute).Format("2006-01-02 15:04:05")
	period := fmt.Sprintf("%s~%s", startTime, endTime)
	reqQuery := map[string]interface{}{
		"period":     period,
		"sumMetrics": false,
		"timeType":   2,
	}
	respBody, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "ListCluster", "/log/api/v2/cloud-native/cluster/list", nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	success, ok := respBody["success"].(bool)
	if !ok || !success {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_prometheus_v2_instances", "ListCluster", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data, ok := respBody["data"].([]interface{})
	if !ok {
		data = []interface{}{}
	}

	// Filter by ids
	idsMap := getIdsStringFilter(d)

	// Filter by name_regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredInstances []interface{}
	for _, item := range data {
		instance, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		clusterName, _ := instance["clusterName"].(string)
		clusterId, _ := instance["clusterId"].(string)

		// Apply name_regex filter
		if nameRegex != nil && !nameRegex.MatchString(clusterName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[clusterId]; !exists {
				continue
			}
		}

		filteredInstances = append(filteredInstances, instance)
	}

	// Prepare result
	instances := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range filteredInstances {
		instance := item.(map[string]interface{})

		mapping := map[string]interface{}{
			"id":                 instance["clusterId"],
			"cluster_id":         instance["clusterId"],
			"cluster_name":       instance["clusterName"],
			"http_api":           instance["httpApi"],
			"remote_write_url":   instance["remoteWriteUrl"],
			"push_gateway_url":   instance["pushGatewayUrl"],
			"objid":              instance["id"],
			"cluster_type":       instance["clusterType"],
			"security_level_tag": instance["SecurityLevelTag"],
			"status":             instance["status"],
		}
		real, err := prometheusService.DescribePrometheusV2InstanceReal(instance["clusterId"].(string))
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DescribePrometheusV2InstanceReal", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if tagSet, ok := real["tagSet"].([]interface{}); ok {
			mapping["tag_set"] = tagSet
		}

		instances = append(instances, mapping)
		id, _ := instance["clusterId"].(string)
		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", instances); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
