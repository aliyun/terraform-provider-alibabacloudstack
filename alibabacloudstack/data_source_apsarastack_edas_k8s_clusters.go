package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackEdasK8sClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasK8SClusterRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
						"cs_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_region_id": {
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
						"network_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_import_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasK8SClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var clusters []EdasK8sCluster

	request := client.NewCommonRequest("POST", "Edas", "2017-08-01", "GetK8sCluster", "/pop/v5/k8s_clusters")
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_k8s_cluster", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response := EdasGetK8sClusterResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if response.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error("list clusters failed for " + response.Message))
	}

	clusters = response.ClusterPage.ClusterList

	// 如果设置了name_regex过滤器，则过滤结果
	var filteredClusters []EdasK8sCluster
	var r *regexp.Regexp
	if nameRegex, hasNameRegex := d.GetOk("name_regex"); hasNameRegex {
		nameRegex, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		r = nameRegex
	}
	idsMap := getIdsStringFilter(d)
	clusterIds := make([]string, 0, len(filteredClusters))
	clusterList := make([]map[string]interface{}, 0, len(filteredClusters))

	for _, cluster := range clusters {
		if r != nil && !r.MatchString(cluster.ClusterName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, exist := idsMap[cluster.ClusterId]; !exist {
				continue
			}
		}
		if v, ok := d.GetOk("logical_region_id"); ok && v != cluster.RegionId {
			continue
		}
		mapping := map[string]interface{}{
			"id":                    cluster.ClusterId,
			"cluster_id":            cluster.ClusterId,
			"cs_cluster_id":         cluster.CsClusterId,
			"logical_region_id":     cluster.RegionId,
			"cluster_name":          cluster.ClusterName,
			"cluster_type":          cluster.ClusterType,
			"network_mode":          cluster.NetworkMode,
			"vpc_id":                cluster.VpcId,
			"cluster_import_status": cluster.ClusterStatus,
		}
		clusterList = append(clusterList, mapping)
		clusterIds = append(clusterIds, cluster.ClusterId)
	}

	d.SetId(dataResourceIdHash(clusterIds))
	if err := d.Set("clusters", clusterList); err != nil {
		return err
	}
	if err := d.Set("ids", clusterIds); err != nil {
		return err
	}

	return nil
}
