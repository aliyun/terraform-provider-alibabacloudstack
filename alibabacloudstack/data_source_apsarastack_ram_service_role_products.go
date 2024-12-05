package alibabacloudstack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
)

func dataSourceAlibabacloudstackRamServiceRoleProducts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudstackRamServiceRoleProductsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chinese_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"ascii_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudstackRamServiceRoleProductsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoleProducts", "")

	response := RoleProducts{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRAMServiceRoleProducts : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ram_service_role_products", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Success == true || len(response.Data) < 1 {
			break
		}

	}

	var r *regexp.Regexp
	if rt, ok := d.GetOk("key"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.Key) {
			continue
		}
		mapping := map[string]interface{}{
			"chinese_name": rg.ChineseName,
			"ascii_name":   rg.ASCIIName,
			"key":          rg.Key,
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("products", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

type RoleProducts struct {
	Data []struct {
		ChineseName string `json:"chineseName"`
		ASCIIName   string `json:"asciiName"`
		Key         string `json:"key"`
	} `json:"data"`
	Message         string `json:"message"`
	ServerRole      string `json:"serverRole"`
	AsapiRequestID  string `json:"asapiRequestId"`
	Success         bool   `json:"success"`
	Domain          string `json:"domain"`
	PureListData    bool   `json:"pureListData"`
	API             string `json:"api"`
	AsapiErrorCode  string `json:"asapiErrorCode"`
}
