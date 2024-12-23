package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOnsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOnsInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"topic_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tps_receive_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tps_send_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOnsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleInstanceList", "")
	request.Headers["OnsRegionId"] = client.RegionId
	request.QueryParams["OnsRegionId"] = client.RegionId
	request.QueryParams["PreventCache"] = ""

	response := OInstance{}

	for {
		raw, err := client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleInstanceList : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == "200" || len(response.Data) < 1 {
			break
		}
	}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range response.Data {
		if r != nil && !r.MatchString(item.InstanceName) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                 item.InstanceID,
			"instance_id":        item.InstanceID,
			"instance_status":    item.InstanceStatus,
			"create_time":        item.CreateTime,
			"instance_type":      item.InstanceType,
			"instance_name":      item.InstanceName,
			"cluster":            item.Cluster,
			"tps_receive_max":    item.TpsReceiveMax,
			"tps_send_max":       item.TpsMax,
			"topic_capacity":     item.TopicCapacity,
			"independent_naming": item.IndependentNaming,
		}

		names = append(names, item.InstanceName)
		ids = append(ids, item.InstanceID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
