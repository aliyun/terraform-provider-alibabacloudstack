package alibabacloudstack

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLaunchTemplate() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'launch_template_name' instead.",
				ConflictsWith: []string{"launch_template_name"},
			},
			"launch_template_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"host_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_owner_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "self", "others", "marketplace", ""}, false),
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
			},
			"internet_max_bandwidth_in": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 200),
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"io_optimized": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "optimized"}, false),
			},
			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"classic", "vpc"}, false),
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}, false),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_price_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk", "cloud_pperf", "cloud_sperf"}, false),
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(20, 500),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"userdata": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'userdata' is deprecated and will be removed in a future release. Please use new field 'user_data' instead.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(2, 256),
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(2, 128),
						},
						"primary_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(2, 256),
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(2, 128),
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(20, 500),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLaunchTemplateCreate, resourceAlibabacloudStackLaunchTemplateRead, resourceAlibabacloudStackLaunchTemplateUpdate, resourceAlibabacloudStackLaunchTemplateDelete)
	return resource
}

func resourceAlibabacloudStackLaunchTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateCreateLaunchTemplateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateName = connectivity.GetResourceData(d, "launch_template_name", "name").(string)
	request.Description = d.Get("description").(string)
	request.HostName = d.Get("host_name").(string)
	request.ImageId = d.Get("image_id").(string)
	request.ImageOwnerAlias = d.Get("image_owner_alias").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceType = d.Get("instance_type").(string)
	request.AutoReleaseTime = d.Get("auto_release_time").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	request.IoOptimized = d.Get("io_optimized").(string)
	request.KeyPairName = d.Get("key_pair_name").(string)
	request.NetworkType = d.Get("network_type").(string)

	request.RamRoleName = d.Get("ram_role_name").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.SecurityEnhancementStrategy = d.Get("security_enhancement_strategy").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.SpotPriceLimit = requests.NewFloat(d.Get("spot_price_limit").(float64))
	request.SpotStrategy = d.Get("spot_strategy").(string)
	request.SystemDiskDiskName = d.Get("system_disk_name").(string)
	request.SystemDiskCategory = d.Get("system_disk_category").(string)
	request.SystemDiskDescription = d.Get("system_disk_description").(string)
	request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	request.UserData = connectivity.GetResourceData(d, "user_data", "userdata").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
	netsRaw := d.Get("network_interfaces").([]interface{})
	if netsRaw != nil {
		var nets []ecs.CreateLaunchTemplateNetworkInterface
		for _, raw := range netsRaw {
			netRaw := raw.(map[string]interface{})
			net := ecs.CreateLaunchTemplateNetworkInterface{
				NetworkInterfaceName: netRaw["name"].(string),
				VSwitchId:            netRaw["vswitch_id"].(string),
				SecurityGroupId:      netRaw["security_group_id"].(string),
				Description:          netRaw["description"].(string),
				PrimaryIpAddress:     netRaw["primary_ip"].(string),
			}
			nets = append(nets, net)
		}
		request.NetworkInterface = &nets
	}

	disksRaw := d.Get("data_disks").([]interface{})
	if disksRaw != nil {
		var disks []ecs.CreateLaunchTemplateDataDisk
		for _, raw := range disksRaw {
			diskRaw := raw.(map[string]interface{})
			disk := ecs.CreateLaunchTemplateDataDisk{
				Size:               fmt.Sprintf("%d", diskRaw["size"].(int)),
				SnapshotId:         diskRaw["snapshot_id"].(string),
				Category:           diskRaw["category"].(string),
				Encrypted:          fmt.Sprintf("%v", diskRaw["encrypted"].(bool)),
				DiskName:           diskRaw["name"].(string),
				Description:        diskRaw["description"].(string),
				DeleteWithInstance: fmt.Sprintf("%v", diskRaw["delete_with_instance"].(bool)),
			}
			disks = append(disks, disk)
		}

		request.DataDisk = &disks
	}
	tagsRaw := d.Get("tags").(map[string]interface{})
	var tags []ecs.CreateLaunchTemplateTag
	for key, value := range tagsRaw {
		tags = append(tags, ecs.CreateLaunchTemplateTag{
			Key:   key,
			Value: value.(string),
		})
	}
	request.Tag = &tags

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateLaunchTemplate(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.CreateLaunchTemplateResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_launch_template", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateLaunchTemplateResponse)

	d.SetId(response.LaunchTemplateId)

	return nil
}

func resourceAlibabacloudStackLaunchTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeLaunchTemplate(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	latestVersion, err := ecsService.DescribeLaunchTemplateVersion(d.Id(), int(object.LatestVersionNumber))
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, latestVersion.LaunchTemplateName, "launch_template_name", "name")
	d.Set("description", latestVersion.LaunchTemplateData.Description)
	d.Set("host_name", latestVersion.LaunchTemplateData.HostName)
	d.Set("image_id", latestVersion.LaunchTemplateData.ImageId)
	d.Set("image_owner_alias", latestVersion.LaunchTemplateData.ImageOwnerAlias)
	d.Set("instance_charge_type", latestVersion.LaunchTemplateData.InstanceChargeType)
	d.Set("instance_name", latestVersion.LaunchTemplateData.InstanceName)
	d.Set("instance_type", latestVersion.LaunchTemplateData.InstanceType)
	d.Set("auto_release_time", latestVersion.LaunchTemplateData.AutoReleaseTime)
	d.Set("internet_charge_type", latestVersion.LaunchTemplateData.InternetChargeType)
	d.Set("internet_max_bandwidth_in", latestVersion.LaunchTemplateData.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", latestVersion.LaunchTemplateData.InternetMaxBandwidthOut)
	d.Set("io_optimized", latestVersion.LaunchTemplateData.IoOptimized)
	d.Set("key_pair_name", latestVersion.LaunchTemplateData.KeyPairName)
	d.Set("network_type", latestVersion.LaunchTemplateData.NetworkType)
	d.Set("ram_role_name", latestVersion.LaunchTemplateData.RamRoleName)
	d.Set("resource_group_id", latestVersion.LaunchTemplateData.ResourceGroupId)
	d.Set("security_enhancement_strategy", latestVersion.LaunchTemplateData.SecurityEnhancementStrategy)
	d.Set("security_group_id", latestVersion.LaunchTemplateData.SecurityGroupId)
	d.Set("spot_price_limit", latestVersion.LaunchTemplateData.SpotPriceLimit)
	d.Set("spot_strategy", latestVersion.LaunchTemplateData.SpotStrategy)
	d.Set("system_disk_name", latestVersion.LaunchTemplateData.SystemDiskDiskName)
	d.Set("system_disk_category", latestVersion.LaunchTemplateData.SystemDiskCategory)
	d.Set("system_disk_description", latestVersion.LaunchTemplateData.SystemDiskDescription)
	d.Set("system_disk_size", latestVersion.LaunchTemplateData.SystemDiskSize)
	d.Set("resource_group_id", latestVersion.LaunchTemplateData.ResourceGroupId)
	connectivity.SetResourceData(d, latestVersion.LaunchTemplateData.UserData, "user_data", "userdata")
	d.Set("vswitch_id", latestVersion.LaunchTemplateData.VSwitchId)
	d.Set("vpc_id", latestVersion.LaunchTemplateData.VpcId)
	d.Set("zone_id", latestVersion.LaunchTemplateData.ZoneId)
	var interfaces []map[string]interface{}
	for _, net := range latestVersion.LaunchTemplateData.NetworkInterfaces.NetworkInterface {
		ds := make(map[string]interface{})
		ds["vswitch_id"] = net.VSwitchId
		ds["security_group_id"] = net.SecurityGroupId
		ds["name"] = net.NetworkInterfaceName
		ds["description"] = net.Description
		ds["primary_ip"] = net.PrimaryIpAddress
		interfaces = append(interfaces, ds)
	}
	if err := d.Set("network_interfaces", interfaces); err != nil {
		return errmsgs.WrapError(err)
	}

	var disks []map[string]interface{}
	for _, disk := range latestVersion.LaunchTemplateData.DataDisks.DataDisk {
		ds := make(map[string]interface{})
		ds["size"] = disk.Size
		ds["snapshot_id"] = disk.SnapshotId
		ds["category"] = disk.Category
		ds["encrypted"] = (disk.Encrypted == "true")
		ds["name"] = disk.DiskName
		ds["description"] = disk.Description
		ds["delete_with_instance"] = disk.DeleteWithInstance
		disks = append(disks, ds)
	}
	if err := d.Set("data_disks", disks); err != nil {
		return errmsgs.WrapError(err)
	}

	tags := make(map[string]interface{})
	for _, tag := range latestVersion.LaunchTemplateData.Tags.InstanceTag {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceAlibabacloudStackLaunchTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	versions, err := getLaunchTemplateVersions(d.Id(), meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	// Remove one of the oldest and non-default version when the total number reach 30
	if len(versions) > 29 {
		var oldestVersion int64
		for _, version := range versions {
			if !version.DefaultVersion && (oldestVersion == 0 || version.VersionNumber < oldestVersion) {
				oldestVersion = version.VersionNumber
			}
		}

		err = deleteLaunchTemplateVersion(d.Id(), int(oldestVersion), meta)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil

}

func resourceAlibabacloudStackLaunchTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateDeleteLaunchTemplateRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteLaunchTemplate(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.DeleteLaunchTemplateResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	ecsService := EcsService{client}
	if err := ecsService.WaitForLaunchTemplate(d.Id(), Deleted, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func getLaunchTemplateVersions(id string, meta interface{}) ([]ecs.LaunchTemplateVersionSet, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateDescribeLaunchTemplateVersionsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = id
	request.PageSize = requests.NewInteger(50)

	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.DescribeLaunchTemplateVersions(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.DescribeLaunchTemplateVersionsResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplateVersionsResponse)
	if len(response.LaunchTemplateVersionSets.LaunchTemplateVersionSet) > 0 {
		return response.LaunchTemplateVersionSets.LaunchTemplateVersionSet, nil
	} else {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("LaunchTemplate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
}

func deleteLaunchTemplateVersion(id string, version int, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateDeleteLaunchTemplateVersionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = id
	request.DeleteVersion = &[]string{strconv.FormatInt(int64(version), 10)}

	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.DeleteLaunchTemplateVersion(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.DeleteLaunchTemplateVersionResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.ProviderERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func createLaunchTemplateVersion(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateCreateLaunchTemplateVersionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LaunchTemplateId = d.Id()
	request.Description = d.Get("description").(string)
	request.HostName = d.Get("host_name").(string)
	request.ImageId = d.Get("image_id").(string)
	request.ImageOwnerAlias = d.Get("image_owner_alias").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.LaunchTemplateName = connectivity.GetResourceData(d, "launch_template_name", "name").(string)
	request.InstanceType = d.Get("instance_type").(string)
	request.AutoReleaseTime = d.Get("auto_release_time").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	request.IoOptimized = d.Get("io_optimized").(string)
	request.KeyPairName = d.Get("key_pair_name").(string)
	request.NetworkType = d.Get("network_type").(string)

	request.RamRoleName = d.Get("ram_role_name").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.SecurityEnhancementStrategy = d.Get("security_enhancement_strategy").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.SpotPriceLimit = requests.NewFloat(d.Get("spot_price_limit").(float64))
	request.SpotStrategy = d.Get("spot_strategy").(string)
	request.SystemDiskDiskName = d.Get("system_disk_name").(string)
	request.SystemDiskCategory = d.Get("system_disk_category").(string)
	request.SystemDiskDescription = d.Get("system_disk_description").(string)
	request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	request.UserData = connectivity.GetResourceData(d, "user_data", "userdata").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
	netsRaw := d.Get("network_interfaces").([]interface{})
	if netsRaw != nil {
		var nets []ecs.CreateLaunchTemplateVersionNetworkInterface
		for _, raw := range netsRaw {
			netRaw := raw.(map[string]interface{})
			net := ecs.CreateLaunchTemplateVersionNetworkInterface{
				NetworkInterfaceName: netRaw["name"].(string),
				VSwitchId:            netRaw["vswitch_id"].(string),
				SecurityGroupId:      netRaw["security_group_id"].(string),
				Description:          netRaw["description"].(string),
				PrimaryIpAddress:     netRaw["primary_ip"].(string),
			}
			nets = append(nets, net)
		}
		request.NetworkInterface = &nets
	}

	disksRaw := d.Get("data_disks").([]interface{})
	if disksRaw != nil {
		var disks []ecs.CreateLaunchTemplateVersionDataDisk
		for _, raw := range disksRaw {
			diskRaw := raw.(map[string]interface{})
			disk := ecs.CreateLaunchTemplateVersionDataDisk{
				Size:               fmt.Sprintf("%d", diskRaw["size"].(int)),
				SnapshotId:         diskRaw["snapshot_id"].(string),
				Category:           diskRaw["category"].(string),
				Encrypted:          fmt.Sprintf("%v", diskRaw["encrypted"].(bool)),
				DiskName:           diskRaw["name"].(string),
				Description:        diskRaw["description"].(string),
				DeleteWithInstance: fmt.Sprintf("%v", diskRaw["delete_with_instance"].(bool)),
			}
			disks = append(disks, disk)
		}

		request.DataDisk = &disks
	}
	tagsRaw := d.Get("tags").(map[string]interface{})
	var tags []ecs.CreateLaunchTemplateVersionTag
	for key, value := range tagsRaw {
		tags = append(tags, ecs.CreateLaunchTemplateVersionTag{
			Key:   key,
			Value: value.(string),
		})
	}
	request.Tag = &tags

	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.CreateLaunchTemplateVersion(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*ecs.CreateLaunchTemplateVersionResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}