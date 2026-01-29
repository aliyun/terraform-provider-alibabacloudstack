package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOtsClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOtsClustersRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_replica": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOtsClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Create common request for Tablestore API
	request := client.NewCommonRequest("POST", "Tablestore", "2020-12-09", "ListClusterInfo", "/v2/openapi/listclusterinfo")

	body, err := json.Marshal(map[string]interface{}{
		"OnlyPublicCluster": "true",
	})
	if err != nil {
		return err
	}
	// Set OnlyPublicCluster parameter if provided
	request.SetContent(body)
	request.SetContentType("application/json")

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("ListClusterInfo", bresponse, request, request.FormParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ListClusterInfo", "Tablestore", errmsg)
	}

	// Parse response
	var resp map[string]interface{}
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &resp); err != nil {
		return errmsgs.WrapErrorf(err, "Failed to unmarshal ListClusterInfo response: %s", bresponse.GetHttpContentString())
	}

	// Extract cluster list
	clustersRaw, ok := resp["ClusterList"]
	if !ok {
		return errmsgs.Error(fmt.Sprintf("ListClusterInfo response missing 'Clusters' field: %#v", resp))
	}

	clustersList, ok := clustersRaw.([]interface{})
	if !ok || len(clustersList) == 0 {
		clustersList = []interface{}{}
	}

	// Prepare filters
	namesMap := getStringListFilters(d, "names")
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	// Process clusters
	names := make([]string, 0)
	clusters := make([]map[string]interface{}, 0)

	for _, clusterRaw := range clustersList {
		cluster, ok := clusterRaw.(map[string]interface{})
		if !ok {
			continue
		}

		clusterName, ok := cluster["ClusterName"].(string)
		if !ok {
			continue
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(clusterName) {
			continue
		}

		// Apply IDs filter
		if len(namesMap) > 0 {
			if _, found := namesMap[clusterName]; !found {
				continue
			}
		}

		// Build cluster mapping
		mapping := map[string]interface{}{
			"cluster_name":    clusterName,
			"cluster_type":    cluster["ClusterType"],
			"alias_name":      cluster["AliasName"],
			"support_replica": cluster["SupportReplica"],
		}

		names = append(names, clusterName)
		clusters = append(clusters, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("clusters", clusters); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
