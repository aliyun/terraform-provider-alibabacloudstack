package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDnsRecursorAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsRecursorAclsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of Recursor ACL IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by Recursor ACL name.",
			},
			"recursor_acls": {
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
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsRecursorAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}

	response, err := client.DoTeaRequest("POST", "CloudDns", "2022-06-24", "DescribeRecursorAcls", "", nil, nil, request)
	if err != nil {
		return err
	}

	// Prepare result
	ids := make([]string, 0)
	recursorAcls := make([]map[string]interface{}, 0)

	for _, item := range response["Data"].([]interface{}) {
		acl := item.(map[string]interface{})

		id := acl["Id"].(string)
		name := acl["Name"].(string)
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		mapping := map[string]interface{}{
			"id":               id,
			"name":             name,
			"policy":           acl["Policy"],
			"line_ids":         acl["LineIds"],
			"remark":           acl["Remark"],
			"create_timestamp": acl["CreateTimestamp"],
			"update_timestamp": acl["UpdateTimestamp"],
		}
		ids = append(ids, id)
		recursorAcls = append(recursorAcls, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	if err := d.Set("recursor_acls", recursorAcls); err != nil {
		return err
	}

	return nil
}
