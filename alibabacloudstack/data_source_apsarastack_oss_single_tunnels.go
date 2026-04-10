package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOssSingleTunnels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOssSingleTunnelsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"tunnels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackOssSingleTunnelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}

	entries, err := ossService.listVpcipEntries()
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Create filters based on schema
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Filter results
	var filteredEntries []VpcipEntry
	for _, entry := range entries {
		id := fmt.Sprintf("%s:%s:%s", entry.Cluster, entry.VpcId, entry.Vip)

		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		if nameRegex != nil {
			if !nameRegex.MatchString(entry.Label) {
				continue
			}
		}

		filteredEntries = append(filteredEntries, entry)
	}

	// Prepare result data
	tunnels := make([]map[string]interface{}, 0, len(filteredEntries))
	ids := make([]string, 0, len(filteredEntries))

	for _, entry := range filteredEntries {
		id := fmt.Sprintf("%s:%s:%s", entry.Cluster, entry.VpcId, entry.Vip)
		tunnel := map[string]interface{}{
			"id":      id,
			"cluster": entry.Cluster,
			"label":   entry.Label,
			"vip":     entry.Vip,
			"vpc_id":  entry.VpcId,
			"shared":  entry.Shared,
		}
		tunnels = append(tunnels, tunnel)
		ids = append(ids, id)
	}

	// Set the data source attributes
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("tunnels", tunnels); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
