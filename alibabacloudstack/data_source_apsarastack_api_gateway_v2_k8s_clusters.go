package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2K8sClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2K8sClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of k8s cluster IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by k8s cluster name.",
			},
			"k8s_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the k8s cluster used as a filter.",
			},
			"cs_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The container service cluster ID as a filter.",
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the k8s cluster.",
						},
						"k8s_cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the k8s cluster.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the k8s cluster.",
						},
						"cs_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The container service cluster ID.",
						},
						"cs_cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The container service cluster name.",
						},
						"slb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SLB type of the k8s cluster.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID associated with the k8s cluster.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2K8sClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	requestParams := make(map[string]interface{})
	if v, ok := d.GetOk("k8s_cluster_name"); ok {
		requestParams["clusterName"] = v.(string)
	}
	requestParams["current"] = 1
	requestParams["size"] = 100

	// Call ListClusters API
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListClusters", "/k8s/listClusters", nil, nil, requestParams)
	if err != nil {
		return err
	}

	// Parse response data
	data, ok := resp["data"]
	if !ok || data == nil {
		d.SetId("")
		return nil
	}

	recordsData, ok := data.(map[string]interface{})["records"]
	if !ok || recordsData == nil {
		d.SetId("")
		return nil
	}

	records, ok := recordsData.([]interface{})
	if !ok {
		d.SetId("")
		return nil
	}

	// Filter by ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Compile name_regex
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var csClusterId string
	if v, ok := d.GetOk("cs_cluster_id"); ok && v.(string) != "" {
		csClusterId = v.(string)
	}

	var filteredRecords []interface{}
	for _, record := range records {
		r, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		k8sClusterCode, ok := r["k8sClusterCode"].(string)
		if !ok {
			continue
		}

		k8sClusterName, ok := r["k8sClusterName"].(string)
		if !ok {
			continue
		}

		// Apply name_regex filter
		if nameRegex != nil && !nameRegex.MatchString(k8sClusterName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[k8sClusterCode]; !exists {
				continue
			}
		}

		filteredRecords = append(filteredRecords, r)
	}

	// Prepare result
	clusters := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, record := range filteredRecords {
		r := record.(map[string]interface{})

		cluster := make(map[string]interface{})
		cluster["id"] = r["k8sClusterCode"]
		cluster["k8s_cluster_name"] = r["k8sClusterName"]
		cluster["cluster_type"] = r["k8sClusterType"]

		// Extract k8sClusterAttribute
		if attrData, ok := r["k8sClusterAttribute"].(map[string]interface{}); ok {
			if csClusterId != "" && csClusterId != attrData["csClusterId"] {
				continue
			}
			cluster["cs_cluster_id"] = attrData["csClusterId"]
			cluster["cs_cluster_name"] = attrData["csClusterName"]
			cluster["slb_type"] = attrData["slbType"]
			if v, exist := attrData["vpcId"]; exist {
				cluster["vpc_id"] = v
			}
		}

		clusters = append(clusters, cluster)
		ids = append(ids, r["k8sClusterCode"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", clusters); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}
