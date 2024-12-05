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

func dataSourceAlibabacloudstackCmsAlarmContacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudstackCmsAlarmContactsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"chanel_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"chanel_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_contact_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudstackCmsAlarmContactsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("GET", "Cms", "2019-01-01", "ListCmsContacts", "")

	response := CmsContact{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListCmsContacts : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_alarm_contacts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	var s []map[string]interface{}
	for _, c := range response.Data {
		if r != nil && !r.MatchString(c.Name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                 c.Cid,
			"alarm_contact_name": c.Name,
		}
		ids = append(ids, c.Cid)

		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("contacts", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
