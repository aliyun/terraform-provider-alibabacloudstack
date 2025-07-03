package alibabacloudstack

import (
	"slices"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackNasZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasZonesRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
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
									"cluster_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_types": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"storage_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"protocol_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"protocols": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:  schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	response, err := client.DoTeaRequest("POST", "NAS", "2017-06-26", "DescribeZones", "", nil, nil, request)
	if err != nil {
		return err
	}

	resp, err := jsonpath.Get("$.Zones.Zone", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeZones", "$.Zones.Zone", response)
	}
	var zoneId, protocol string
	var zoneIds []string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
	}
	
	s := make([]map[string]interface{}, 0)
	for _, r := range resp.([]interface{}) {
		object := r.(map[string]interface{})
		if zoneId != "" && zoneId != object["ZoneId"].(string) {
			continue
		}
		var protocols []string
		for _, v := range object["Performance"].(map[string]interface{})["Protocol"].([]interface{}) {
			protocols = append(protocols, v.(string))
		}
		if protocol != "" && !slices.Contains(protocols, protocol) {
			continue
		}
		zoneIds = append(zoneIds, zoneId)
		
		clusters := []map[string]interface{}{}
		for _, ci := range object["Clusters"].(map[string]interface {})["Cluster"].([]interface{}) {
			c := ci.(map[string]interface{})
			cluster := map[string]interface{}{
				"cluster_id": c["ClusterId"],
				"cluster_type": c["ClusterType"],
				"cluster_version": c["ClusterVersion"],
			}
			instanceTypes := []map[string]interface{}{}
			for _, i := range c["InstanceTypes"].(map[string]interface{})["InstanceType"].([]interface{}){
				instanceType := i.(map[string]interface{})
				instanceTypes = append(instanceTypes, map[string]interface{}{
					"storage_type": instanceType["StorageType"],
					"protocol_type": instanceType["ProtocolType"],
				})
			}
			if len(instanceTypes) < 1 {
				continue
			}
			cluster["instance_types"] = instanceTypes
			clusters = append(clusters, cluster)
			}
		if len(clusters) < 1 {
			continue
		}
		mapping := map[string]interface{}{
			"zone_id": object["ZoneId"],
			"clusters": clusters,
			"protocols": object["Performance"].(map[string]interface{})["Protocol"].([]interface{}),
			
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zoneIds))

	if err := d.Set("zones", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
