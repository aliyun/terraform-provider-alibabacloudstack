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

func dataSourceAlibabacloudStackEvpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEvpcsRead,

		Schema: map[string]*schema.Schema{
			"evpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
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
			"evpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"evpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
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
						"resource_group": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
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

func dataSourceAlibabacloudStackEvpcsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{}
	if v, ok := d.GetOk("evpc_name"); ok {
		request["EvpcName"] = v.(string)
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListEvpc", "", nil, request, nil)
		addDebug("ListEvpc", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_evpc_evpcs", "ListEvpc", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ListEvpc", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	allEvpcs := response["EvpcList"].([]interface{})

	var filteredEvpcs []map[string]interface{}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}

	idsMap := getIdsStringFilter(d)

	for _, evpc := range allEvpcs {
		evpcMap := evpc.(map[string]interface{})

		if r != nil && !r.MatchString(evpcMap["EvpcName"].(string)) {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[evpcMap["EvpcId"].(string)]; !ok {
				continue
			}
		}

		if status, ok := d.GetOk("status"); ok && evpcMap["Status"].(string) != status.(string) {
			continue
		}

		filteredEvpcs = append(filteredEvpcs, evpcMap)
	}

	return evpcsDescriptionAttributes(d, filteredEvpcs)
}

func evpcsDescriptionAttributes(d *schema.ResourceData, evpcSetTypes []map[string]interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, evpc := range evpcSetTypes {
		mapping := map[string]interface{}{
			"evpc_id":             evpc["EvpcId"],
			"evpc_name":           evpc["EvpcName"],
			"status":              evpc["Status"],
			"description":         evpc["Description"],
			"cidr":                evpc["Cidr"],
			"tenant_id":           evpc["TenantId"],
			"department":          evpc["Department"],
			"department_name":     evpc["DepartmentName"],
			"region_id":           evpc["RegionId"],
			"resource_group":      evpc["ResourceGroup"],
			"resource_group_name": evpc["ResourceGroupName"],
			"cluster_id":          evpc["ClusterId"],
			"ascm_create_user":    evpc["AscmCreateUser"],
			"create_time":         evpc["CreateTime"],
			"update_time":         evpc["UpdateTime"],
		}

		ids = append(ids, evpc["EvpcId"].(string))
		names = append(names, evpc["EvpcName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("evpcs", s); err != nil {
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
