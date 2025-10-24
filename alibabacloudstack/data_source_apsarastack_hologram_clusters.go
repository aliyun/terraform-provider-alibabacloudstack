package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackHologramClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackHologramClustersRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compute_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Standard",
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Follower"}, false),
			},
			"cpu": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "intel",
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"performance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_replica": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackHologramClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"x-acs-body": map[string]interface{}{
			"computeType": d.Get("compute_type"),
			"cpu":         d.Get("cpu"),
			"zoneId":      d.Get("zone_id"),
			"regionId":    client.RegionId,
		},
	}
	response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "ListClusters", "/api/v1/regions/listClusters", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instances", "ListClusters", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	clusters, err := jsonpath.Get("$.ClusterList", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instances", "ListClusters", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var ids []string
	var s []map[string]interface{}

	for _, item := range clusters.([]interface{}) {
		cluster := item.(map[string]interface{})
		if name_regex, ok := d.GetOk("name_regex"); ok && name_regex.(string) != "" {
			r := regexp.MustCompile(name_regex.(string))
			if !r.MatchString(cluster["Cluster"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[cluster["Cluster"].(string)]; !exist {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":              cluster["Cluster"],
			"type":            cluster["Type"],
			"cpu":             cluster["CpuBrand"],
			"cpu_arch":        cluster["CpuArch"],
			"performance":     cluster["Performance"],
			"support_replica": cluster["SupportReplica"],
			"cluster":         cluster["Cluster"],
		}

		ids = append(ids, cluster["Cluster"].(string))
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
