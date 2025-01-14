package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasDeployGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasDeployGroupsRead,

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateFunc: validation.StringIsValidRegExp,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasDeployGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	appId := d.Get("app_id").(string)
	request := edas.CreateListDeployGroupRequest()
	client.InitRoaRequest(*request.RoaRequest)
	request.AppId = appId
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListDeployGroup(request)
	})

	response, ok := raw.(*edas.ListDeployGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_edas_deploy_groups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	if response.Code != 200 {
		return errmsgs.WrapError(errmsgs.Error(response.Message))
	}

	var filteredGroups []edas.DeployGroup
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r, err = regexp.Compile(nameRegex.(string))
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
		for _, group := range response.DeployGroupList.DeployGroup {
			if r != nil && !r.MatchString(group.GroupName) {
				continue
			}

			filteredGroups = append(filteredGroups, group)
		}
	} else {
		filteredGroups = response.DeployGroupList.DeployGroup
	}

	return edasDeployGroupAttributes(d, filteredGroups)
}

func edasDeployGroupAttributes(d *schema.ResourceData, groups []edas.DeployGroup) error {
	var ids []string
	var s []map[string]interface{}
	var names []string

	for _, group := range groups {
		mapping := map[string]interface{}{
			"group_id":           group.GroupId,
			"group_name":         group.GroupName,
			"group_type":         group.GroupType,
			"create_time":        group.CreateTime,
			"update_time":        group.UpdateTime,
			"app_id":             group.AppId,
			"cluster_id":         group.ClusterId,
			"package_version_id": group.PackageVersionId,
			"app_version_id":     group.AppVersionId,
		}
		ids = append(ids, group.GroupId)
		s = append(s, mapping)
		names = append(names, group.GroupName)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
