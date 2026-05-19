package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackLogClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackLogClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackLogClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("GET", "Sls", "2020-03-31", "ListClusters", "/sls/v1/cluster/listClusters")

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_log_clusters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	var result ListClustersResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &result)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var ids []string
	var s []map[string]interface{}

	for _, cluster := range result.Clusters {
		if len(idsMap) > 0 {
			if _, ok := idsMap[cluster.Name]; !ok {
				continue
			}
		}

		if nameRegex != nil && !nameRegex.MatchString(cluster.Name) {
			continue
		}

		ids = append(ids, cluster.Name)

		mapping := map[string]interface{}{
			"name":        cluster.Name,
			"data_server": cluster.DataServer,
			"status":      cluster.Status,
			"zone":        cluster.Zone,
			"description": cluster.Description,
		}
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

type ListClustersResponse struct {
	Count    int              `json:"count"`
	Total    int              `json:"total"`
	Clusters []SlsClusterInfo `json:"clusters"`
}

type SlsClusterInfo struct {
	Name        string `json:"name"`
	RegionId    string `json:"regionId"`
	DataServer  string `json:"dataServer"`
	Status      string `json:"status"`
	Zone        string `json:"zone"`
	Description string `json:"description"`
}
