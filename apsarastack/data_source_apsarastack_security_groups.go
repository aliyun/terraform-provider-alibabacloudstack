package apsarastack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SecurityGroup struct {
	Attributes        ecs.DescribeSecurityGroupAttributeResponse
	CreationTime      string
	SecurityGroupType string
	ResourceGroupId   string
	Tags              ecs.TagsInDescribeSecurityGroups
}

func dataSourceApsaraStackSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"groups": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	request := ecs.CreateDescribeSecurityGroupsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	var sg []SecurityGroup
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeSecurityGroupsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeSecurityGroupsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	// ids
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "security_groups", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if len(response.SecurityGroups.SecurityGroup) < 1 {
			break
		}

		for _, item := range response.SecurityGroups.SecurityGroup {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecurityGroupName) {
					continue
				}
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[item.SecurityGroupId]; !ok {
					continue
				}
			}

			attr, err := ecsService.DescribeSecurityGroup(item.SecurityGroupId)
			if err != nil {
				return WrapError(err)
			}

			sg = append(sg,
				SecurityGroup{
					Attributes:        attr,
					CreationTime:      item.CreationTime,
					SecurityGroupType: item.SecurityGroupType,
					ResourceGroupId:   item.ResourceGroupId,
					Tags:              item.Tags,
				},
			)
		}

		if len(response.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	return securityGroupsDescription(d, sg, meta)
}

func securityGroupsDescription(d *schema.ResourceData, sg []SecurityGroup, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, item := range sg {
		mapping := map[string]interface{}{
			"id":            item.Attributes.SecurityGroupId,
			"name":          item.Attributes.SecurityGroupName,
			"description":   item.Attributes.Description,
			"vpc_id":        item.Attributes.VpcId,
			"creation_time": item.CreationTime,
			"tags":          ecsService.tagsToMap(item.Tags.Tag),
		}

		ids = append(ids, string(item.Attributes.SecurityGroupId))
		names = append(names, item.Attributes.SecurityGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
