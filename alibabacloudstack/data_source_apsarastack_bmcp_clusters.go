package alibabacloudstack

import (
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"
)

func dataSourceAlibabacloudStackBmcpClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpClustersRead,

		Schema: map[string]*schema.Schema{
			"cluster_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_id_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"clusters": {
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
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gpu": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gpu_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"gpu_model": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"flops_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"video_memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_arch_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"switch_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   PageSizeLarge,
	}

	resp, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListBmcpCluster", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListBmcpCluster", "GET", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response
	clusterList, err := jmespath.Search("ClusterList", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListBmcpCluster", "ClusterList", resp)
	}

	clusterListSlice, ok := clusterList.([]interface{})
	if !ok {
		clusterListSlice = []interface{}{}
	}

	var filteredClusters []interface{}
	for _, item := range clusterListSlice {
		cluster := item.(map[string]interface{})

		clusterName, _ := cluster["ClusterName"].(string)
		clusterId, _ := cluster["ClusterId"].(string)
		status, _ := cluster["Status"].(string)
		regionId, _ := cluster["RegionId"].(string)

		// Apply cluster_name_regex filter
		if clusterNameRegex, ok := d.GetOk("cluster_name_regex"); ok {
			if !strings.Contains(clusterName, clusterNameRegex.(string)) {
				continue
			}
		}

		// Apply cluster_id_regex filter
		if clusterIdRegex, ok := d.GetOk("cluster_id_regex"); ok {
			if !strings.Contains(clusterId, clusterIdRegex.(string)) {
				continue
			}
		}

		// Apply status_regex filter
		if statusRegex, ok := d.GetOk("status_regex"); ok {
			if !strings.Contains(status, statusRegex.(string)) {
				continue
			}
		}

		// Apply region_id_regex filter
		if regionIdRegex, ok := d.GetOk("region_id_regex"); ok {
			if !strings.Contains(regionId, regionIdRegex.(string)) {
				continue
			}
		}

		filteredClusters = append(filteredClusters, cluster)
	}

	return bmcpClustersAttributes(d, filteredClusters)
}

func bmcpClustersAttributes(d *schema.ResourceData, clusters []interface{}) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range clusters {
		cluster := item.(map[string]interface{})

		clusterId := formatAnyToString(cluster["ClusterId"])
		clusterName := formatAnyToString(cluster["ClusterName"])
		status := formatAnyToString(cluster["Status"])
		regionId := formatAnyToString(cluster["RegionId"])
		zone := formatAnyToString(cluster["Zone"])
		vpcId := formatAnyToString(cluster["VpcId"])
		evpcId := formatAnyToString(cluster["EvpcId"])
		organization := formatAnyToString(cluster["Organization"])
		nodeCount := formatAnyToInt(cluster["NodeCount"])
		cpuCount := formatAnyToInt(cluster["CpuCount"])
		memCount := formatAnyToInt(cluster["MemCount"])
		flopsCount := formatAnyToFloat(cluster["FlopsCount"])
		videoMemory := formatAnyToInt(cluster["VideoMemory"])
		clusterArchType := formatAnyToString(cluster["ClusterArchType"])
		switchMethod := formatAnyToString(cluster["SwitchMethod"])
		enableIpv6, _ := cluster["EnableIpv6"].(bool)
		createTime := formatAnyToString(cluster["CreateTime"])
		updateTime := formatAnyToString(cluster["UpdateTime"])

		// Build GPU list
		var gpuList []map[string]interface{}
		if gpuListRaw, ok := cluster["Gpu"].([]interface{}); ok {
			for _, gpu := range gpuListRaw {
				gpuMap := gpu.(map[string]interface{})
				gpuNum := formatAnyToInt(gpuMap["GpuNum"])
				gpuModel := formatAnyToString(gpuMap["GpuModel"])
				gpuList = append(gpuList, map[string]interface{}{
					"gpu_num":   gpuNum,
					"gpu_model": gpuModel,
				})
			}
		}

		mapping := map[string]interface{}{
			"id":                clusterId,
			"cluster_id":        clusterId,
			"cluster_name":      clusterName,
			"status":            status,
			"region_id":         regionId,
			"zone":              zone,
			"vpc_id":            vpcId,
			"evpc_id":           evpcId,
			"organization":      organization,
			"node_count":        nodeCount,
			"cpu_count":         cpuCount,
			"mem_count":         memCount,
			"gpu":               gpuList,
			"flops_count":       flopsCount,
			"video_memory":      videoMemory,
			"cluster_arch_type": clusterArchType,
			"switch_method":     switchMethod,
			"enable_ipv6":       enableIpv6,
			"create_time":       createTime,
			"update_time":       updateTime,
		}

		ids = append(ids, clusterId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
