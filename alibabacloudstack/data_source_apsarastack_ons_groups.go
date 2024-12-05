package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id_regex": {
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
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	namespaceid := d.Get("instance_id").(string)

	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleGroupList", "")
	request.QueryParams["InstanceId"] = namespaceid
	request.QueryParams["OnsRegionId"] = client.RegionId
	request.QueryParams["PreventCache"] = ""

	response := OnsGroup{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleGroupList : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_groups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Code == 200 || len(response.Data) < 1 {
			break
		}
	}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("group_id_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, ons := range response.Data {
		if r != nil && !r.MatchString(ons.GroupID) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                  ons.ID,
			"remark":              ons.Remark,
			"instance_id":         ons.NamespaceID,
			"group_id":            ons.GroupID,
			"owner":               ons.Owner,
			"independent_naming":  ons.IndependentNaming,
			"create_time":         ons.CreateTime,
		}
		ids = append(ids, string(rune(ons.ID)))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = writeToFile(output.(string), s)
	}
	return nil
}
