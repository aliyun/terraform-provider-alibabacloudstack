package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEssScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEssScalingGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_scaling_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cooldown_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"removal_policies": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"load_balancer_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"db_instance_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"lifecycle_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"active_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pending_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"removing_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEssScalingGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateDescribeScalingGroupsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var allScalingGroups []ess.ScalingGroup

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingGroups(request)
		})
		response, ok := raw.(*ess.DescribeScalingGroupsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scalinggroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(response.ScalingGroups.ScalingGroup) < 1 {
			break
		}

		allScalingGroups = append(allScalingGroups, response.ScalingGroups.ScalingGroup...)

		if len(response.ScalingGroups.ScalingGroup) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredScalingGroupsTemp = make([]ess.ScalingGroup, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}
	if okNameRegex || okIds {
		for _, group := range allScalingGroups {
			if okNameRegex && nameRegex != "" {
				var r = regexp.MustCompile(nameRegex.(string))
				if r != nil && !r.MatchString(group.ScalingGroupName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[group.ScalingGroupId]; !ok {
					continue
				}
			}
			filteredScalingGroupsTemp = append(filteredScalingGroupsTemp, group)
		}
	} else {
		filteredScalingGroupsTemp = allScalingGroups
	}
	return scalingGroupsDescriptionAttribute(d, filteredScalingGroupsTemp, meta)
}

func scalingGroupsDescriptionAttribute(d *schema.ResourceData, scalingGroups []ess.ScalingGroup, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, scalingGroup := range scalingGroups {
		mapping := map[string]interface{}{
			"id":                          scalingGroup.ScalingGroupId,
			"name":                        scalingGroup.ScalingGroupName,
			"active_scaling_configuration": scalingGroup.ActiveScalingConfigurationId,
			"region_id":                   scalingGroup.RegionId,
			"min_size":                    scalingGroup.MinSize,
			"max_size":                    scalingGroup.MaxSize,
			"cooldown_time":               scalingGroup.DefaultCooldown,
			"removal_policies":            scalingGroup.RemovalPolicies.RemovalPolicy,
			"load_balancer_ids":           scalingGroup.LoadBalancerIds.LoadBalancerId,
			"db_instance_ids":             scalingGroup.DBInstanceIds.DBInstanceId,
			"vswitch_ids":                 scalingGroup.VSwitchIds.VSwitchId,
			"lifecycle_state":             scalingGroup.LifecycleState,
			"total_capacity":              scalingGroup.TotalCapacity,
			"active_capacity":             scalingGroup.ActiveCapacity,
			"pending_capacity":            scalingGroup.PendingCapacity,
			"removing_capacity":           scalingGroup.RemovingCapacity,
			"creation_time":               scalingGroup.CreationTime,
		}
		ids = append(ids, scalingGroup.ScalingGroupId)
		names = append(names, scalingGroup.ScalingGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
