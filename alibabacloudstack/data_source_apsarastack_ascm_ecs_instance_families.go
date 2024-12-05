package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func dataSourceAlibabacloudStackEcsInstanceFamilies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEcsInstanceFamiliesRead,
		Schema: map[string]*schema.Schema{

			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"families": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type_family_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEcsInstanceFamiliesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "ascm", "2019-05-10", "DescribeInstanceTypeFamilies", "")

	response := EcsInstanceFamily{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw DescribeInstanceTypeFamilies : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ecs_instance_families", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == 200 || len(response.Data.InstanceTypeFamilies) < 1 {
			break
		}

	}

	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data.InstanceTypeFamilies {
		mapping := map[string]interface{}{
			"instance_type_family_id": rg.InstanceTypeFamilyID,
			"generation":              rg.Generation,
		}
		ids = append(ids, rg.InstanceTypeFamilyID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("families", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
