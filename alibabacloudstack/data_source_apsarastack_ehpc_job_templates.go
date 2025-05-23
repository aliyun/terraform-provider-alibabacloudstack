package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackEhpcJobTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEhpcJobTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
						"array_request": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"package_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"re_runable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"runas_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stderr_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stdout_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"thread": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"variables": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEhpcJobTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "ListJobTemplates"
	request := make(map[string]interface{})

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		response, err := client.DoTeaRequest("GET", "ECS", "2018-04-12", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		resp, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Templates.JobTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"array_request":        object["ArrayRequest"],
			"clock_time":           object["ClockTime"],
			"command_line":         object["CommandLine"],
			"gpu":                  formatInt(object["Gpu"]),
			"id":                   fmt.Sprint(object["Id"]),
			"job_template_id":      fmt.Sprint(object["Id"]),
			"job_template_name":    object["Name"],
			"mem":                  object["Mem"],
			"node":                 formatInt(object["Node"]),
			"package_path":         object["PackagePath"],
			"priority":             formatInt(object["Priority"]),
			"queue":                object["Queue"],
			"re_runable":           object["ReRunable"],
			"runas_user":           object["RunasUser"],
			"stderr_redirect_path": object["StderrRedirectPath"],
			"stdout_redirect_path": object["StdoutRedirectPath"],
			"task":                 formatInt(object["Task"]),
			"thread":               formatInt(object["Thread"]),
			"variables":            object["Variables"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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
