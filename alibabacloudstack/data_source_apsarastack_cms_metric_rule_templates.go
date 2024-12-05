package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCmsMetricRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCmsMetricRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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
			"is_default":{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"history":{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rest_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
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

func dataSourceAlibabacloudStackCmsMetricRuleTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var objects []cms.Template
	var templateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		templateNameRegex = r
	}
	
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	
	request := requests.NewCommonRequest()
	request.Product = "cms"
	request.Version = "2019-01-01"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "DescribeMetricRuleTemplateList"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		
		
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "cms",
		"RegionId":        client.RegionId,
		"Action":          "DescribeMetricRuleTemplateList",
		"Version":         "2019-01-01",
		"pageSize":        "10",
		//"roleType":        roleType,
	}
	
	
	if v, ok := d.GetOk("keyword"); ok {
		request.QueryParams["Keyword"] = v.(string)
	}
	if v, ok := d.GetOk("template_id"); ok {
		request.QueryParams["TemplateId"] =fmt.Sprint(v.(int))
	}
	request.QueryParams["IsDefault"] =fmt.Sprint(d.Get("is_default").(bool))
	request.QueryParams["History"] =fmt.Sprint(d.Get("history").(bool))
	var resp *cms.DescribeMetricRuleTemplateListResponse
	for {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "cms_metric_rule_templates", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*responses.CommonResponse)
		err = json.Unmarshal(response.GetHttpContentBytes(), &resp)

		for _, item := range resp.Templates.Template {
			if templateNameRegex != nil {
				if !templateNameRegex.MatchString(item.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item.TemplateId)]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(resp.Templates.Template) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":               object.Description,
			"name":                      object.Name,
			"rest_version":              object.RestVersion,
			"id":                        object.TemplateId,
		}
		ids = append(ids, fmt.Sprint(object.TemplateId))
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	
	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}