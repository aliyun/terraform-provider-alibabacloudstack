package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCmsMetricRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCmsMetricRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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

	var objects []Template
	var templateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		templateNameRegex = r
	}

	idsMap := getIdsStringFilter(d)
	pageNumber := 1

	request := client.NewCommonRequest("GET", "Cms", "2019-01-01", "DescribeMetricRuleTemplateList", "")
	request.QueryParams["PageSize"] = "10"
	request.QueryParams["IsDefault"] = fmt.Sprint(d.Get("is_default").(bool))
	request.QueryParams["History"] = "false"

	var templateId int
	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(int)
		request.QueryParams["TemplateId"] = strconv.Itoa(templateId)
	}

	var resp *DescribeMetricRuleTemplateListResponse
	for {
		request.QueryParams["PageNumber"] = strconv.Itoa(pageNumber)
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw DescribeMetricRuleTemplateList : %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_metric_rule_templates", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &resp)

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
			
			if templateId != 0 && templateId != int(item.TemplateId) {
				continue
			}
			
			objects = append(objects, item)
		}
		if len(resp.Templates.Template) < PageSizeLarge {
			break
		}

		pageNumber += 1
	}

	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":  object.Description,
			"name":         object.Name,
			"rest_version": object.RestVersion,
			"id":           object.TemplateId,
		}
		ids = append(ids, fmt.Sprint(object.TemplateId))
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
