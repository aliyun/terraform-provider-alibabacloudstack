package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOssClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOssClustersRead,

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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_apsara_stack": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"api_zonelocal_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_zonelocal_public_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_public_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_ha_enable_single_cluster_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_cs_public_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_unique_domain": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_master_zone": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOssClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoApi", "")
	// request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
	// request.QueryParams["AppAction"] = "GetOssEndpointList"
	// request.QueryParams["AppName"] = "one-console-app-oss"
	// region := client.RegionId
	// request.QueryParams["Params"] = fmt.Sprintf("{\"region\":\"%s\", \"params\":{\"region\":\"%s\"}}", region, region)
	// bresponse, err := client.ProcessCommonRequest(request)
	// addDebug("GetOssEndpointList", bresponse, request, request.QueryParams)
	// if err != nil {
	// 	if bresponse == nil {
	// 		return errmsgs.WrapErrorf(err, "Process Common Request Failed")
	// 	}
	// 	errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
	// 	return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "GetOssEndpointList", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	// }
	// result := make(map[string]interface{})
	// _ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	// data, ok := result["Data"]
	// if !ok || len(data.([]interface{})) == 0 {
	// 	return errmsgs.Error(fmt.Sprintf("GetOssEndpointList Failed! region: %s \n %#v", client.RegionId, bresponse.GetHttpContentString()))
	// }
	ossService := OssService{client}
	data, err := ossService.GetOssEndpointListForIot()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap := getIdsStringFilter(d)
	ids := make([]string, 0)
	clusters := make([]map[string]interface{}, 0)
	for _, v := range data {
		object := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[object["cluster"].(string)]; !ok {
				continue
			}
		}

		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(object["clusterName"].(string)) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                                  object["cluster"],
			"cluster":                             object["cluster"],
			"ha_apsara_stack":                     object["haApsaraStack"],
			"api_zonelocal_endpoint":              object["api-zonelocal-endpoint"],
			"api_zonelocal_public_endpoint":       object["api-zonelocal-public-endpoint"],
			"oss_public_endpoint":                 object["oss-public-endpoint"],
			"real_zone":                           object["real-zone"],
			"oss_ha_enable_single_cluster_access": object["ossHaEnableSingleClusterAccess"],
			"oss_cs_public_endpoint":              object["oss-cs-public-endpoint"],
			"oss_unique_domain":                   object["oss-unique-domain"],
			"cluster_name":                        object["clusterName"],
			"is_master_zone":                      object["isMasterZone"],
			"location":                            object["location"],
			"oss_endpoint":                        object["oss-endpoint"],
			"oss_suffix":                          object["oss-suffix"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		clusters = append(clusters, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("clusters", clusters); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
