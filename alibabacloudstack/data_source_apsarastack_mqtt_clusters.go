package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMqttClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMqttClustersRead,

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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMqttClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "Ons-inner", "2018-02-05", "ConsoleMqttClusterList", "")
	request.QueryParams["OnsRegionId"] = string(client.Region)
	request.QueryParams["PreventCache"] = fmt.Sprintf("%v", time.Now().Unix()*1000)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "ConsoleClusterList", errmsgs.AlibabacloudStackOssGoSdk, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	data, ok := result["Data"]
	if !ok || len(data.([]interface{})) == 0 {
		return errmsgs.Error(fmt.Sprintf("ConsoleClusterList Failed! region: %s \n %#v", client.RegionId, bresponse.GetHttpContentString()))
	}
	idsMap := getIdsStringFilter(d)
	ids := make([]string, 0)
	clusters := make([]map[string]interface{}, 0)
	for _, v := range data.([]interface{}) {
		object := v.(string)
		if len(idsMap) > 0 {
			if _, ok := idsMap[object]; !ok {
				continue
			}
		}

		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(object) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":   object,
			"name": object,
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
