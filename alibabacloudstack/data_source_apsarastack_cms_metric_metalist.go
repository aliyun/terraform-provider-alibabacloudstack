package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"regexp"
)

func dataSourceAlibabacloudstackCmsMetricMetalist() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudstackCmsMetricMetalistRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"periods": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"statistics": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"namespace": {
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

func dataSourceAlibabacloudstackCmsMetricMetalistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	Namespace := d.Get("namespace").(string)

	request := client.NewCommonRequest("GET", "cms", "2019-01-01", "DescribeMetricMetaList", "")
	request.QueryParams["Namespace"] = Namespace
	request.QueryParams["ProductName"] = "cms"

	response := MetaList{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw DescribeMetricMetaList : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_metric_metalist", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if len(response.Resources.Resource) < 1 || response.Success == true {
			break
		}
	}
	var r *regexp.Regexp
	if rt, ok := d.GetOk("namespace"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Resources.Resource {
		if r != nil && !r.MatchString(Namespace) {
			continue
		}
		mapping := map[string]interface{}{
			"description":  rg.Description,
			"namespace":    rg.Namespace,
			"labels":       rg.Labels,
			"metric_name":  rg.MetricName,
			"dimensions":   rg.Dimensions,
			"periods":      rg.Periods,
			"statistics":   rg.Statistics,
			"unit":         rg.Unit,
		}
		ids = append(ids, fmt.Sprintf(rg.Namespace))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("resources", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
