package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackFlinkNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackFlinkNamespacesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"owner_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackFlinkNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	response, err := client.DoTeaRequest("GET", "ververica", "2020-05-01", "DescribeNamespaces", "/flink/namespace/list", nil, nil, nil)
	if err != nil {
		return err
	}
	ns, err := jsonpath.Get("$.data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}

	var names []string
	var s []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	for _, o := range ns.([]interface{}) {
		object := o.(map[string]interface{})
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(object["Name"].(string)) {
				continue
			}
		}

		if ower_id, ok := d.GetOk("owner_uid"); ok {
			if _, exist := object["Uid"]; !exist {
				continue
			}
			if ower_id.(string) != object["Uid"].(string) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[object["Name"].(string)]; !exist {
				continue
			}
		}

		cu, err := strconv.Atoi(object["GuaranteeQuota"].(string))
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"name":      object["Name"].(string),
			"cu":        cu,
			"cpu_type":  object["CpuBrand"].(string),
			"owner_uid": object["Uid"].(string),
		}

		names = append(names, object["Name"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("namespaces", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
