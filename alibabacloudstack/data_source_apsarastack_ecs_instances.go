package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				//must contain a valid status, expected Creating, Starting, Running, Stopping, Stopped
				ValidateFunc: validation.StringInSlice([]string{
					string(Running),
					string(Stopped),
					string(Creating),
					string(Starting),
					string(Stopping),
				}, false),
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateDescribeInstancesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request.InstanceIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeInstancesTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	var allInstances []ecs.Instance
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		bresponse, ok := raw.(*ecs.DescribeInstancesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if len(bresponse.Instances.Instance) < 1 {
			break
		}

		allInstances = append(allInstances, bresponse.Instances.Instance...)

		if len(bresponse.Instances.Instance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredInstancesTemp []ecs.Instance

	nameRegex, ok := d.GetOk("name_regex")
	imageId, okImg := d.GetOk("image_id")
	if (ok && nameRegex.(string) != "") || (okImg && imageId.(string) != "") {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, inst := range allInstances {
			if r != nil && !r.MatchString(inst.InstanceName) {
				continue
			}
			if imageId.(string) != "" && inst.ImageId != imageId.(string) {
				continue
			}
			filteredInstancesTemp = append(filteredInstancesTemp, inst)
		}
	} else {
		filteredInstancesTemp = allInstances
	}
	// Filter by ram role name and fetch the instance role name
	instanceIds := make([]string, 0)
	for _, inst := range filteredInstancesTemp {
		if inst.InstanceNetworkType == "classic" {
			continue
		}
		instanceIds = append(instanceIds, inst.InstanceId)
	}
	instanceRoleNameMap := make(map[string]string)
	for index := 0; index < len(instanceIds); index += 100 {
		// DescribeInstanceRamRole parameter InstanceIds supports at most 100 items once
		request := ecs.CreateDescribeInstanceRamRoleRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceIds = convertListToJsonString(convertListStringToListInterface(instanceIds[index:IntMin(index+100, len(instanceIds))]))
		request.RamRoleName = d.Get("ram_role_name").(string)
		request.PageSize = requests.NewInteger(PageSizeLarge)
		request.PageNumber = requests.NewInteger(1)
		for {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeInstanceRamRole(request)
			})
			bresponse, ok := raw.(*ecs.DescribeInstanceRamRoleResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			if len(bresponse.InstanceRamRoleSets.InstanceRamRoleSet) < 1 {
				break
			}
			for _, role := range bresponse.InstanceRamRoleSets.InstanceRamRoleSet {
				instanceRoleNameMap[role.InstanceId] = role.RamRoleName
			}

			if len(bresponse.InstanceRamRoleSets.InstanceRamRoleSet) < PageSizeLarge {
				break
			}

			if page, err := getNextpageNumber(request.PageNumber); err != nil {
				return errmsgs.WrapError(err)
			} else {
				request.PageNumber = page
			}
		}
	}
	instanceDiskMappings, err := getInstanceDisksMappings(instanceRoleNameMap, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return instancessDescriptionAttributes(d, filteredInstancesTemp, instanceRoleNameMap, instanceDiskMappings, meta)
}

// populate the numerous fields that the instance description returns.
func instancessDescriptionAttributes(d *schema.ResourceData, instances []ecs.Instance, instanceRoleNameMap map[string]string, instanceDisksMap map[string][]map[string]interface{}, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, inst := range instances {
		// if instance can not in instanceRoleNameMap, it should be removed.
		if _, ok := instanceRoleNameMap[inst.InstanceId]; !ok {
			continue
		}
		mapping := map[string]interface{}{
			"id":                         inst.InstanceId,
			"region_id":                  inst.RegionId,
			"availability_zone":          inst.ZoneId,
			"status":                     inst.Status,
			"name":                       inst.InstanceName,
			"instance_type":              inst.InstanceType,
			"vpc_id":                     inst.VpcAttributes.VpcId,
			"vswitch_id":                 inst.VpcAttributes.VSwitchId,
			"image_id":                   inst.ImageId,
			"description":                inst.Description,
			"security_groups":            inst.SecurityGroupIds.SecurityGroupId,
			"eip":                        inst.EipAddress.IpAddress,
			"key_name":                   inst.KeyPairName,
			"ram_role_name":              instanceRoleNameMap[inst.InstanceId],
			"creation_time":              inst.CreationTime,
			"instance_charge_type":       inst.InstanceChargeType,
			"internet_max_bandwidth_out": inst.InternetMaxBandwidthOut,
			"disk_device_mappings":       instanceDisksMap[inst.InstanceId],
			"tags":                       ecsTagsToMap(inst.Tags.Tag),
		}
		if len(inst.InnerIpAddress.IpAddress) > 0 {
			mapping["private_ip"] = inst.InnerIpAddress.IpAddress[0]
		} else {
			mapping["private_ip"] = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}

		ids = append(ids, inst.InstanceId)
		names = append(names, inst.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("names", names)
	if err := d.Set("instances", s); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

// Returns a mapping of instance disks
func getInstanceDisksMappings(instanceMap map[string]string, meta interface{}) (map[string][]map[string]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateDescribeDisksRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeXLarge)
	request.PageNumber = requests.NewInteger(1)
	instanceDisks := make(map[string][]map[string]interface{})
	var allDisks []ecs.Disk
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		bresponse, ok := raw.(*ecs.DescribeDisksResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return instanceDisks, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_instances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if bresponse == nil || len(bresponse.Disks.Disk) < 1 {
			break
		}

		allDisks = append(allDisks, bresponse.Disks.Disk...)

		if len(bresponse.Disks.Disk) < PageSizeXLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return instanceDisks, errmsgs.WrapError(err)
		}
		request.PageNumber = page
	}
	for _, disk := range allDisks {
		if _, ok := instanceMap[disk.InstanceId]; !ok {
			continue
		}
		mapping := map[string]interface{}{
			"device":   disk.Device,
			"size":     disk.Size,
			"category": disk.Category,
			"type":     disk.Type,
		}
		instanceDisks[disk.InstanceId] = append(instanceDisks[disk.InstanceId], mapping)
	}

	return instanceDisks, nil
}
