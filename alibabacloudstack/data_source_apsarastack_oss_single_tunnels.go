package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
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

	request := client.NewCommonRequest("GET", "OneRouter", "2018-12-12", "DoOpenApi", "")
	request.QueryParams["OpenApiAction"] = "ListVpcip"
	request.QueryParams["ProductName"] = "oss"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("ListVpcip", bresponse, request, request.QueryParams)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	vpcipList, err := jsonpath.Get("$.Data.ListVpcipResult.Vpcip", response)
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
	var filteredVpcips []interface{}
	for _, item := range vpcipList.([]interface{}) {
		vpcip := item.(map[string]interface{})

		// For this resource, we'll use the VIP as an identifier for filtering by IDs
		// The ID format appears to be cluster:vpc_id:vip based on the resource code
		id := fmt.Sprintf("%s:%s:%s", vpcip["Cluster"], vpcip["VpcId"], vpcip["Vip"])

		// Check if ID matches the provided IDs filter
		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		// Check if label matches the name regex filter (using Label as name)
		if nameRegex != nil {
			if label, exists := vpcip["Label"]; exists {
				if !nameRegex.MatchString(label.(string)) {
					continue
				}
			} else {
				continue
			}
		}

		filteredVpcips = append(filteredVpcips, item)
	}

	// Prepare result data
	tunnels := make([]map[string]interface{}, 0, len(filteredVpcips))
	ids := make([]string, 0, len(filteredVpcips))

	for _, item := range filteredVpcips {
		vpcip := item.(map[string]interface{})

		tunnel := make(map[string]interface{})
		tunnel["id"] = fmt.Sprintf("%s:%s:%s", vpcip["Cluster"], vpcip["VpcId"], vpcip["Vip"])
		// Map fields according to schema
		if val, ok := vpcip["Cluster"]; ok {
			tunnel["cluster"] = val
		}
		if val, ok := vpcip["Label"]; ok {
			tunnel["label"] = val
		}
		if val, ok := vpcip["Vip"]; ok {
			tunnel["vip"] = val
		}
		if val, ok := vpcip["VpcId"]; ok {
			tunnel["vpc_id"] = val
		}
		if val, ok := vpcip["shared"]; ok {
			tunnel["shared"] = val
		}
		tunnels = append(tunnels, tunnel)

		// Add ID to the list of IDs
		id := fmt.Sprintf("%s:%s:%s", vpcip["Cluster"], vpcip["VpcId"], vpcip["Vip"])
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
