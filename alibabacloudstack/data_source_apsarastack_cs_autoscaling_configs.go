package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCSAutoscalingConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCSAutoscalingConfigsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cool_down_duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unneeded_duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"utilization_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu_utilization_threshold": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCSAutoscalingConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	var clusterIds []string
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			clusterIds = append(clusterIds, id.(string))
		}
	}

	// If no ids provided, get all clusters
	if len(clusterIds) == 0 {
		request := client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClustersV1", "/api/v1/clusters")
		request.QueryParams["SignatureVersion"] = "1.0"
		request.QueryParams["ProductName"] = "CS"

		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := ""
			if response != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_autoscaling_configs", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		var clusters ClustersV1
		err = json.Unmarshal(response.GetHttpContentBytes(), &clusters)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		for _, cluster := range clusters.Clusters {
			clusterIds = append(clusterIds, cluster.ClusterID)
		}
	}

	var ids []string
	var configs []map[string]interface{}

	for _, clusterId := range clusterIds {
		addonConfigRaw, err := csService.GetAutoscalingConfig(clusterId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				log.Printf("[DEBUG] GetAutoscalingConfig for cluster %s returned not found, skipping", clusterId)
				continue
			}
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_autoscaling_configs", "GetAutoscalingConfig", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		config := map[string]interface{}{
			"id":         clusterId,
			"cluster_id": clusterId,
		}

		if v, ok := addonConfigRaw["ScaleDownUnneededTime"]; ok && v != nil {
			config["cool_down_duration"] = fmt.Sprint(v)
		}
		if v, ok := addonConfigRaw["ScaleDownDelayAfterAdd"]; ok && v != nil {
			config["unneeded_duration"] = fmt.Sprint(v)
		}
		if v, ok := addonConfigRaw["ScaleDownUtilizationThreshold"]; ok && v != nil {
			config["utilization_threshold"] = fmt.Sprint(v)
		}
		if v, ok := addonConfigRaw["ScaleDownGpuUtilizationThreshold"]; ok && v != nil {
			config["gpu_utilization_threshold"] = fmt.Sprint(v)
		}

		ids = append(ids, clusterId)
		configs = append(configs, config)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("configs", configs); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
