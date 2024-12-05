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

func dataSourceAlibabacloudStackOnsTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOnsTopicsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"relation": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"relation_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"independent_naming": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOnsTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	namespaceid := d.Get("instance_id").(string)

	request := client.NewCommonRequest("POST", "Ons-inner", "2018-02-05", "ConsoleTopicList", "")
	request.QueryParams["namespaceId"] = namespaceid
	request.QueryParams["OnsRegionId"] = client.RegionId
	request.QueryParams["PreventCache"] = ""

	response := Topic{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ConsoleTopicList : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ons_topics", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, ons := range response.Data {
		if r != nil && !r.MatchString(ons.Topic) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                 ons.ID,
			"topic":              ons.Topic,
			"remark":             ons.Remark,
			"instance_id":        ons.NamespaceID,
			"owner":              ons.Owner,
			"relation":           ons.Relation,
			"relation_name":      ons.RelationName,
			"message_type":       ons.OrderType,
			"independent_naming": ons.IndependentNaming,
			"create_time":        ons.CreateTime,
		}
		ids = append(ids, string(rune(ons.ID)))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("topics", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		_ = writeToFile(output.(string), s)
	}
	return nil
}
