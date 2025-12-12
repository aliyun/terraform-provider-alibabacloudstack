package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackApfsZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApfsZonesRead,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by zone ID.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by cluster ID.",
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID.",
						},
						"clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID.",
									},
									"available_capacity": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Available capacity of the cluster.",
									},
									"storage_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Storage type of the instance.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackApfsZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	request := make(map[string]interface{})

	// Call API
	resp, err := client.DoTeaRequest("POST", "EFS", "2017-06-26", "DescribeZones", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_apfs_zones", "DescribeZones", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	zoneList, err := jsonpath.Get("$.Zones.Zone", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_apfs_zones", "$.Zones.Zone", resp)
	}

	var ids []string
	var s []map[string]interface{}

	// Filter parameters
	filterZoneId := d.Get("zone_id").(string)
	filterClusterId := d.Get("cluster_id").(string)

	for _, zone := range zoneList.([]interface{}) {
		zoneMap, ok := zone.(map[string]interface{})
		if !ok {
			continue
		}

		zoneId, _ := zoneMap["ZoneId"].(string)

		// Apply zone_id filter
		if filterZoneId != "" && filterZoneId != zoneId {
			continue
		}

		clusters := make([]map[string]interface{}, 0)
		if clusterList, ok := zoneMap["Clusters"].(map[string]interface{})["Cluster"].([]interface{}); ok {
			for _, cluster := range clusterList {
				clusterMap := cluster.(map[string]interface{})
				clusterId, _ := clusterMap["ClusterId"].(string)

				// Apply cluster_id filter
				if filterClusterId != "" && filterClusterId != clusterId {
					continue
				}

				storageType := ""
				if instanceTypeList, ok := clusterMap["InstanceTypes"].(map[string]interface{})["InstanceType"].([]interface{}); ok {
					storageType = instanceTypeList[0].(map[string]interface{})["StorageType"].(string)
				}

				availableCapacity, _ := clusterMap["AvailableCapacity"].(int64)

				clusterData := map[string]interface{}{
					"cluster_id":         clusterId,
					"available_capacity": availableCapacity,
					"storage_type":       storageType,
				}
				clusters = append(clusters, clusterData)
			}

			// Skip zone if no clusters match filters
			if filterClusterId != "" && len(clusters) == 0 {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":       zoneId,
			"zone_id":  zoneId,
			"clusters": clusters,
		}

		ids = append(ids, zoneId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
