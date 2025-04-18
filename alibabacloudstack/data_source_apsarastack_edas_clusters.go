package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackEdasClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasClustersRead,

		Schema: map[string]*schema.Schema{
			"logical_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mem_used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	logicalRegionId := d.Get("logical_region_id").(string)
	request := edas.CreateListClusterRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.LogicalRegionId = logicalRegionId
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, id := range v.([]interface{}) {
			if id == nil {
				continue
			}
			idsMap[Trim(id.(string))] = Trim(id.(string))
		}
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListCluster(request)
	})

	response, ok := raw.(*edas.ListClusterResponse)
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_clusters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if response.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error(response.Message))
	}

	var filteredClusters []edas.Cluster
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
		for _, cluster := range response.ClusterList.Cluster {
			if r != nil && !r.MatchString(cluster.ClusterName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[cluster.ClusterId]; !ok {
					continue
				}
			}
			filteredClusters = append(filteredClusters, cluster)
		}
	} else {
		filteredClusters = response.ClusterList.Cluster
	}

	return edasClusterDescriptionAttributes(d, filteredClusters)
}

func edasClusterDescriptionAttributes(d *schema.ResourceData, clusters []edas.Cluster) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, cluster := range clusters {
		mapping := map[string]interface{}{
			"cluster_id":     cluster.ClusterId,
			"cluster_name":   cluster.ClusterName,
			"cluster_type":   cluster.ClusterType,
			"create_time":    cluster.CreateTime,
			"update_time":    cluster.UpdateTime,
			"cpu":            cluster.Cpu,
			"cpu_used":       cluster.CpuUsed,
			"mem":            cluster.Mem,
			"mem_used":       cluster.MemUsed,
			"network_mode":   cluster.NetworkMode,
			"node_num":       cluster.NodeNum,
			"vpc_id":         cluster.VpcId,
			"region_id":      cluster.RegionId,
		}
		ids = append(ids, cluster.ClusterId)
		s = append(s, mapping)
		names = append(names, cluster.ClusterName)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
