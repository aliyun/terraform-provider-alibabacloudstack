package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackHsmClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackHsmClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"vsm_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"evsm", "gvsm", "svsm"}, false),
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
						"vsm_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_white_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_master": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"abnormal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_no": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hsm_cluster_items": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"is_master": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"zone_no": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vsm_type": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlibabacloudStackHsmClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	currentPage := 1
	request["PageSize"] = 100
	request["CurrentPage"] = currentPage
	if v, ok := d.GetOk("vsm_type"); ok {
		request["VsmType"] = v.(string)
	}
	var clusters []interface{}
	for true {
		resp, err := client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeClusters", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		// Parse response
		clustersRaw, ok := resp["HsmClusters"].([]interface{})
		totalCount, _ := resp["TotalCount"].(json.Number).Int64()
		if ok {
			clusters = append(clusters, clustersRaw...)
		}
		if int(totalCount) <= currentPage*100 {
			break
		}
		currentPage++
	}

	// Handle ids filter
	idsMap := getIdsStringFilter(d)

	// Handle name_regex filter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var clusterList []map[string]interface{}
	var ids []string
	for _, v := range clusters {
		cluster := v.(map[string]interface{})
		clusterId, _ := cluster["ClusterId"].(string)
		// Apply id filter
		if len(idsMap) > 0 {
			clusterId, _ := cluster["ClusterId"].(string)
			if _, exists := idsMap[clusterId]; !exists {
				continue
			}
		}

		// Apply name_regex filter
		if nameRegex != nil {
			clusterName, _ := cluster["ClusterName"].(string)
			if !nameRegex.MatchString(clusterName) {
				continue
			}
		}

		mapping := make(map[string]interface{})
		mapping["id"] = clusterId
		mapping["cluster_id"] = clusterId
		mapping["vsm_type"], _ = cluster["VsmType"].(string)
		mapping["cluster_name"], _ = cluster["ClusterName"].(string)
		mapping["ip_white_list"], _ = cluster["IpWhiteList"].(string)
		mapping["vpc_id"], _ = cluster["VpcId"].(string)
		mapping["cluster_size"], _ = cluster["ClusterSize"].(string)
		mapping["cluster_master"], _ = cluster["ClusterMaster"].(string)
		if gmtCreateStr, ok := cluster["GmtCreate"].(string); ok {
			if gmtCreate, err := strconv.ParseInt(gmtCreateStr, 10, 64); err == nil {
				mapping["gmt_create"] = int(gmtCreate)
			}
		}
		mapping["master_ip"], _ = cluster["MasterIp"].(string)
		mapping["abnormal_type"], _ = cluster["AbnormalType"].(string)
		mapping["vendor_code"], _ = cluster["VendorCode"].(string)
		mapping["product_code"], _ = cluster["ProductCode"].(string)

		// Process cluster zones
		if clusterZonesRaw, ok := cluster["ClusterZones"].([]interface{}); ok {
			var clusterZones []map[string]interface{}
			for _, zoneRaw := range clusterZonesRaw {
				if zone, ok := zoneRaw.(map[string]interface{}); ok {
					zoneMapping := make(map[string]interface{})
					zoneMapping["vswitch_id"], _ = zone["VSwitchId"].(string)
					zoneMapping["zone_no"], _ = zone["ZoneNo"].(string)
					clusterZones = append(clusterZones, zoneMapping)
				}
			}
			mapping["cluster_zones"] = clusterZones
		}

		// Process HSM cluster items
		if hsmClusterItemsRaw, ok := cluster["HsmClusterItems"].([]interface{}); ok {
			var hsmClusterItems []map[string]interface{}
			for _, itemRaw := range hsmClusterItemsRaw {
				if item, ok := itemRaw.(map[string]interface{}); ok {
					itemMapping := make(map[string]interface{})
					if idStr, ok := item["Id"].(string); ok {
						if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
							itemMapping["id"] = int(id)
						}
					}
					itemMapping["instance_id"], _ = item["InstanceId"].(string)
					if statusStr, ok := item["Status"].(string); ok {
						if status, err := strconv.Atoi(statusStr); err == nil {
							itemMapping["status"] = status
						}
					}
					if isMasterStr, ok := item["IsMaster"].(string); ok {
						if isMaster, err := strconv.Atoi(isMasterStr); err == nil {
							itemMapping["is_master"] = isMaster
						}
					}
					itemMapping["zone_no"], _ = item["ZoneNo"].(string)
					itemMapping["ip"], _ = item["Ip"].(string)
					itemMapping["vsm_type"], _ = item["VsmType"].(string)
					hsmClusterItems = append(hsmClusterItems, itemMapping)
				}
			}
			mapping["hsm_cluster_items"] = hsmClusterItems
		}
		ids = append(ids, clusterId)
		clusterList = append(clusterList, mapping)
	}

	// Set results
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", clusterList); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
