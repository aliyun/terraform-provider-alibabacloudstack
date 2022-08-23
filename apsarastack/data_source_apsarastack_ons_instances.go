package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
)

func dataSourceApsaraStackOnsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackOnsInstancesRead,

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

func dataSourceApsaraStackOnsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ons-inner"
	request.Version = "2018-02-05"
	request.Scheme = "http"
	request.RegionId = client.RegionId
	request.ApiName = "ConsoleInstanceList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ons-inner",
		"RegionId":        client.RegionId,
		"Action":          "ConsoleInstanceList",
		"Version":         "2018-02-05",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"OnsRegionId":     client.RegionId,
		"PreventCache":    "",
	}
	response := OInstance{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleInstanceList : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_ons_instances", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		bresponse, _ := raw.(*responses.CommonResponse)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if err != nil {
			return WrapError(err)
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
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
