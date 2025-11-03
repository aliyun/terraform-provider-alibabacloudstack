package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDatahubProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDatahubProjectsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDatahubProjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	listReq := client.NewCommonRequest("GET", "datahub", "2019-11-20", "ListProjects", "")
	listReq.QueryParams["PageNumber"] = "1"
	listReq.QueryParams["PageSize"] = "100"
	listResp, err := client.ProcessCommonRequest(listReq)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	listResult := &ListProjectResult{}
	err = json.Unmarshal(listResp.GetHttpContentBytes(), listResult)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	projectNames := []string{}
	var filteredProjects []map[string]interface{}
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	var r *regexp.Regexp
	if nameRegexOk && nameRegex.(string) != "" {
		r, err = regexp.Compile(nameRegex.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	for _, project := range listResult.List.Project {
		if nameRegexOk && nameRegex.(string) != "" {
			if !r.MatchString(project.ProjectName) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":               project.ProjectName,
			"name":             project.ProjectName,
			"comment":          project.Comment,
			"create_time":      strconv.FormatInt(project.CreateTime, 10),
			"last_modify_time": strconv.FormatInt(project.UpdateTime, 10),
		}
		filteredProjects = append(filteredProjects, mapping)
		projectNames = append(projectNames, project.ProjectName)
	}

	d.Set("ids", projectNames)
	d.SetId(dataResourceIdHash(projectNames))
	if err := d.Set("projects", filteredProjects); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
