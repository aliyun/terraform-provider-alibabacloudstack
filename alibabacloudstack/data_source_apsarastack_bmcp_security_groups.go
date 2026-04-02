package alibabacloudstack

import (
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func dataSourceAlibabacloudStackBmcpSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sg_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"department": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"department_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ascm_create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v.(string)
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroup", "", nil, request, nil)
		addDebug("ListSecurityGroup", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bmcp_security_groups", "ListSecurityGroup", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ListSecurityGroup", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	allSecurityGroups := response["data"].([]interface{})

	var filteredSecurityGroups []map[string]interface{}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	idsMap := getIdsStringFilter(d)

	for _, securityGroup := range allSecurityGroups {
		securityGroupMap := securityGroup.(map[string]interface{})

		if r != nil && !r.MatchString(securityGroupMap["name"].(string)) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[securityGroupMap["sgId"].(string)]; !ok {
				continue
			}
		}

		filteredSecurityGroups = append(filteredSecurityGroups, securityGroupMap)
	}

	return securityGroupsDescriptionAttributes(d, filteredSecurityGroups)
}

func securityGroupsDescriptionAttributes(d *schema.ResourceData, securityGroupSetTypes []map[string]interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, securityGroup := range securityGroupSetTypes {
		mapping := map[string]interface{}{
			"sg_id":             securityGroup["sgId"],
			"name":              securityGroup["name"],
			"description":       securityGroup["description"],
			"vpc_id":            securityGroup["vpcId"],
			"resource_group":    securityGroup["ResourceGroup"],
			"resource_group_name": securityGroup["ResourceGroupName"],
			"department":        securityGroup["Department"],
			"department_name":   securityGroup["DepartmentName"],
			"region_id":         securityGroup["RegionId"],
			"ascm_create_user":  securityGroup["AscmCreateUser"],
			"create_time":       securityGroup["createTime"],
			"update_time":       securityGroup["updateTime"],
		}

		ids = append(ids, securityGroup["sgId"].(string))
		names = append(names, securityGroup["name"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("security_groups", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
