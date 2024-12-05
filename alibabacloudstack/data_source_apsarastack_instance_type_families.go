package alibabacloudstack

import (
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackInstanceTypeFamilies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackInstanceTypeFamiliesRead,

		Schema: map[string]*schema.Schema{
			"generation": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ecs-1", "ecs-2", "ecs-3", "ecs-4"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"families": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackInstanceTypeFamiliesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	request := ecs.CreateDescribeInstanceTypeFamiliesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	if v, ok := d.GetOk("generation"); ok {
		request.Generation = v.(string)
	}

	zones, err := ecsService.DescribeZones(d)
	if err != nil {
		return WrapErrorf(err, "DescribeZones", "alibabacloudstack_instance_type_families", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	families := make(map[string]map[string]string)
	for _, zone := range zones {
		for _, infos := range zone.AvailableResources.ResourcesInfo {
			for _, family := range infos.InstanceTypeFamilies.SupportedInstanceTypeFamily {
				if _, ok := families[family]; !ok {
					families[family] = make(map[string]string)
				}
				families[family][zone.ZoneId] = ""
			}
		}
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceTypeFamilies(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_instance_type_families", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	var instanceTypeFamilies []ecs.InstanceTypeFamily
	response, _ := raw.(*ecs.DescribeInstanceTypeFamiliesResponse)
	if response != nil {
		for _, family := range response.InstanceTypeFamilies.InstanceTypeFamily {
			if _, ok := families[family.InstanceTypeFamilyId]; !ok {
				continue
			}
			instanceTypeFamilies = append(instanceTypeFamilies, family)
		}
	}
	return instanceTypeFamiliesDescriptionAttributes(d, instanceTypeFamilies, families)
}

func instanceTypeFamiliesDescriptionAttributes(d *schema.ResourceData, typeFamilies []ecs.InstanceTypeFamily, families map[string]map[string]string) error {
	var ids []string
	var s []map[string]interface{}
	for _, f := range typeFamilies {

		mapping := map[string]interface{}{
			"id":         f.InstanceTypeFamilyId,
			"generation": f.Generation,
		}
		var zoneIds []string
		for zoneId := range families[f.InstanceTypeFamilyId] {
			zoneIds = append(zoneIds, zoneId)
		}
		sort.Strings(zoneIds)
		mapping["zone_ids"] = zoneIds

		ids = append(ids, f.InstanceTypeFamilyId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("families", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
