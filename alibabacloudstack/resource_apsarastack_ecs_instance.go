package alibabacloudstack

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'availability_zone' is deprecated and will be removed in a future release. Please use new field 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			"zone_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},

			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},

			"security_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ECS-Instance",
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
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
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password") == ""
				},
				Elem: schema.TypeString,
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Default:      DiskCloudEfficiency,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_efficiency", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40,
			},

			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"system_disk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 128),
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_efficiency", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
							Default:      DiskCloudEfficiency,
							ForceNew:     true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(2, 256),
						},
					},
				},
			},

			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vswitch_id"},
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hpc_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: vpcTypeResourceDiffSuppressFunc,
			},

			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"storage_set_partition_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
			},
			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}, false),
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_address_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 10),
			},
			"ipv6_address_list": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"tags":             tagsSchema(),
			"system_disk_tags": tagsSchema(),
			"data_disk_tags":   tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackInstanceCreate, resourceAlibabacloudStackInstanceRead, resourceAlibabacloudStackInstanceUpdate, resourceAlibabacloudStackInstanceDelete)
	return resource
}

func resourceAlibabacloudStackInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	request, err := buildAlibabacloudStackInstanceArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request.IoOptimized = string(IOOptimized)

	if d.Get("is_outdated").(bool) {
		request.IoOptimized = string(NoneOptimized)
	}
	client.InitRpcRequest(*request.RpcRequest)

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.RunInstances(request)
		})
		if err != nil {
			if errmsgs.IsThrottling(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.RunInstancesResponse)
		for _, k := range response.InstanceIdSets.InstanceIdSet {
			d.SetId(k)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopped"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{"Stopping"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	if v, ok := d.GetOk("security_groups"); ok {
		sgs := expandStringList(v.(*schema.Set).List())
		if len(sgs) > 1 {
			err := ecsService.JoinSecurityGroups(d.Id(), sgs[1:])
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}
	if d.Get("enable_ipv6").(bool) && d.Get("ipv6_address_count").(int) > 0 {
		_, err := AssignIpv6AddressesFunc(d.Id(), d.Get("ipv6_address_count").(int), d.Get("ipv6_address_list").([]string), meta)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func resourceAlibabacloudStackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	instance, err := ecsService.DescribeInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	log.Printf("[ECS Creation]: Getting Instance Details Successfully: %s", instance.Status)
	system_disks, err := ecsService.DescribeInstanceDisksByType(d.Id(), client.ResourceGroup, "system")
	if err != nil {
		return errmsgs.WrapError(err)
	}
	system_disk_tags := getOnlySystemTags(d, system_disks[0].Tags.Tag)

	d.Set("system_disk_category", system_disks[0].Category)
	d.Set("system_disk_size", system_disks[0].Size)
	d.Set("system_disk_name", system_disks[0].DiskName)
	d.Set("system_disk_description", system_disks[0].Description)
	d.Set("system_disk_id", system_disks[0].DiskId)
	d.Set("system_disk_tags", ecsService.tagsToMap(system_disk_tags))
	d.Set("instance_name", instance.InstanceName)
	d.Set("description", instance.Description)
	d.Set("status", instance.Status)
	connectivity.SetResourceData(d, instance.ZoneId, "zone_id", "availability_zone")
	d.Set("host_name", instance.HostName)
	d.Set("image_id", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	d.Set("password", d.Get("password").(string))
	d.Set("internet_max_bandwidth_out", instance.InternetMaxBandwidthOut)
	d.Set("internet_max_bandwidth_in", instance.InternetMaxBandwidthIn)
	d.Set("key_name", instance.KeyPairName)

	d.Set("hpc_cluster_id", instance.HpcClusterId)
	d.Set("tags", ecsService.tagsToMap(instance.Tags.Tag))

	d.Set("vswitch_id", instance.VpcAttributes.VSwitchId)

	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		d.Set("private_ip", instance.VpcAttributes.PrivateIpAddress.IpAddress[0])
	} else {
		d.Set("private_ip", strings.Join(instance.InnerIpAddress.IpAddress, ","))
	}

	sgs := make([]string, 0, len(instance.SecurityGroupIds.SecurityGroupId))
	for _, sg := range instance.SecurityGroupIds.SecurityGroupId {
		sgs = append(sgs, sg)
	}
	if err := d.Set("security_groups", sgs); err != nil {
		return errmsgs.WrapError(err)
	}

	if !d.IsNewResource() || d.HasChange("user_data") {
		dataRequest := ecs.CreateDescribeUserDataRequest()
		client.InitRpcRequest(*dataRequest.RpcRequest)
		dataRequest.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeUserData(dataRequest)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*ecs.DescribeUserDataResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), dataRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(dataRequest.GetActionName(), raw, dataRequest.RpcRequest, dataRequest)
		response, _ := raw.(*ecs.DescribeUserDataResponse)
		old_s := base64.StdEncoding.EncodeToString([]byte(response.UserData))
		d.Set("user_data", d.Get("user_data").(string))
		log.Printf("data : %s", old_s)
	}

	if len(instance.VpcAttributes.VSwitchId) > 0 && (!d.IsNewResource() || d.HasChange("role_name")) {
		request := ecs.CreateDescribeInstanceRamRoleRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceIds = convertListToJsonString([]interface{}{d.Id()})
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*ecs.DescribeInstanceRamRoleResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
		log.Printf("[ECS Creation]: Getting Instance RamRole Details: %s ", response.InstanceRamRoleSets.InstanceRamRoleSet)
		if len(response.InstanceRamRoleSets.InstanceRamRoleSet) >= 1 {
			d.Set("role_name", response.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
		}
	}

	if instance.InstanceChargeType == string(PrePaid) {
		request := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceAutoRenewAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*ecs.DescribeInstanceAutoRenewAttributeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	// 查询 ipv6地址集
	// request := ecs.CreateDescribeNetworkInterfacesRequest()
	// client.InitRpcRequest(*request.RpcRequest)
	// request.InstanceId = d.Id()
	// raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	// 	return ecsClient.DescribeNetworkInterfaces(request)
	// })
	// bresponse, ok := raw.(*ecs.DescribeNetworkInterfacesResponse)
	// if err != nil {
	// 	errmsg := ""
	// 	if ok {
	// 		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
	// 	}
	// 	return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	// }
	// ipv6_address_list := make([]string, 0)
	// if len(bresponse.NetworkInterfaceSets.NetworkInterfaceSet) > 0 && len(bresponse.NetworkInterfaceSets.NetworkInterfaceSet[0].Ipv6Sets.Ipv6Set) > 0 {
	// 	for _, ipv6 := range bresponse.NetworkInterfaceSets.NetworkInterfaceSet[0].Ipv6Sets.Ipv6Set {
	// 		ipv6_address_list = append(ipv6_address_list, ipv6.Ipv6Address)
	// 	}
	// 	d.Set("ipv6_address_list", ipv6_address_list)
	// }

	return nil
}

func resourceAlibabacloudStackInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	d.Partial(true)

	if !d.IsNewResource() {
		err := setTags(client, TagResourceInstance, d)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	if d.HasChange("security_groups") {
		if !d.IsNewResource() || d.Get("vswitch_id").(string) == "" {
			o, n := d.GetChange("security_groups")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)

			rl := expandStringList(os.Difference(ns).List())
			al := expandStringList(ns.Difference(os).List())

			if len(al) > 0 {
				err := ecsService.JoinSecurityGroups(d.Id(), al)
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
			if len(rl) > 0 {
				err := ecsService.LeaveSecurityGroups(d.Id(), rl)
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}
	}

	run := false
	imageUpdate, err := modifyInstanceImage(d, meta, run)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	vpcUpdate, err := modifyVpcAttribute(d, meta, run)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	passwordUpdate, err := modifyInstanceAttribute(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	typeUpdate, err := modifyInstanceType(d, meta, run)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if imageUpdate || vpcUpdate || passwordUpdate || typeUpdate {
		run = true
		instance, errDesc := ecsService.DescribeInstance(d.Id())
		if errDesc != nil {
			return errmsgs.WrapError(errDesc)
		}
		if instance.Status == string(Running) {
			stopRequest := ecs.CreateStopInstanceRequest()
			client.InitRpcRequest(*stopRequest.RpcRequest)
			stopRequest.InstanceId = d.Id()
			stopRequest.ForceStop = requests.NewBoolean(false)
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.StopInstance(stopRequest)
			})
			if err != nil {
				errmsg := ""
				if bresponse, ok := raw.(*ecs.StopInstanceResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), stopRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(stopRequest.GetActionName(), raw)
		}

		stateConf := BuildStateConf([]string{"Pending", "Running", "Stopping"}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

		if _, err = stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if _, err := modifyInstanceImage(d, meta, run); err != nil {
			return errmsgs.WrapError(err)
		}

		if _, err := modifyVpcAttribute(d, meta, run); err != nil {
			return errmsgs.WrapError(err)
		}

		if _, err := modifyInstanceType(d, meta, run); err != nil {
			return errmsgs.WrapError(err)
		}

		startRequest := ecs.CreateStartInstanceRequest()
		client.InitRpcRequest(*startRequest.RpcRequest)
		startRequest.InstanceId = d.Id()

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.StartInstance(startRequest)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
					time.Sleep(time.Second)
					return resource.RetryableError(err)
				}
				errmsg := ""
				if bresponse, ok := raw.(*ecs.StartInstanceResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), startRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				return resource.NonRetryableError(err)
			}
			addDebug(startRequest.GetActionName(), raw)
			return nil
		})

		if err != nil {
			return err
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"Pending", "Starting", "Stopped"},
			Target:     []string{"Running"},
			Refresh:    ecsService.InstanceStateRefreshFunc(d.Id(), []string{}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		if _, err = stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if err := modifyInstanceNetworkSpec(d, meta); err != nil {
		return errmsgs.WrapError(err)
	}

	if d.HasChanges("system_disk_tags", "system_disk_id", "image_id", "tags") {
		var oraw, nraw map[string]interface{}
		system_disks, err := ecsService.DescribeInstanceDisksByType(d.Id(), client.ResourceGroup, "system")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		system_disk_tags := getOnlySystemTags(d, system_disks[0].Tags.Tag)
		oraw = make(map[string]interface{})
		for k, v := range ecsService.tagsToMap(system_disk_tags) {
			oraw[k] = v
		}

		if v, ok := d.GetOk("system_disk_tags"); ok {
			nraw = v.(map[string]interface{})
		}

		err = updateTags(client, []string{system_disks[0].DiskId}, "disk", oraw, nraw)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if data_disk_tags, ok := d.GetOk("data_disk_tags"); ok {
		disks, err := ecsService.DescribeInstanceDisksByType(d.Id(), client.ResourceGroup, "data")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if len(disks) > 0 {
			oraw := make(map[string]interface{})
			diskids := make([]string, 0, len(disks))
			datadisk_tags := ecsMergeTags(d, data_disk_tags.(map[string]interface{}))
			for _, disk := range disks {
				diskids = append(diskids, disk.DiskId)
			}
			err = updateTags(client, diskids, "disk", oraw, datadisk_tags)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}

	}

	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	stopRequest := ecs.CreateStopInstanceRequest()
	client.InitRpcRequest(*stopRequest.RpcRequest)
	stopRequest.InstanceId = d.Id()
	stopRequest.ForceStop = requests.NewBoolean(true)

	deleteRequest := ecs.CreateDeleteInstanceRequest()
	client.InitRpcRequest(*deleteRequest.RpcRequest)
	deleteRequest.InstanceId = d.Id()
	deleteRequest.Force = requests.NewBoolean(true)

	wait := incrementalWait(1*time.Second, 1*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(deleteRequest)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectInstanceStatus", "DependencyViolation.RouteEntry", "IncorrectInstanceStatus.Initializing"}) {
				return resource.RetryableError(err)
			}
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, "LastTokenProcessing"}) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*ecs.StartInstanceResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "StartInstance", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
		}
		addDebug(deleteRequest.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, errmsgs.EcsNotFound) {
			return nil
		}
		return err
	}

	stateConf := BuildStateConf([]string{"Pending", "Running", "Stopped", "Stopping"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecsService.InstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func buildAlibabacloudStackInstanceArgs(d *schema.ResourceData, meta interface{}) (*ecs.RunInstancesRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ecs.CreateRunInstancesRequest()
	client.InitRpcRequest(*request.RpcRequest)

	if v := d.Get("enable_ipv6").(bool); v {
		request.QueryParams["Ipv6CidrBlock"] = d.Get("ipv6_cidr_block").(string)
		if d.Get("ipv6_address_count").(int) <= 0 {
			return nil, errmsgs.WrapError(errmsgs.Error("if enable_ipv6 = true, ipv6_address_count must be greater than 0"))
		}
		request.QueryParams["Ipv6AddressCount"] = "0"
	}
	request.InstanceType = d.Get("instance_type").(string)
	request.SystemDiskDiskName = d.Get("system_disk_name").(string)
	imageID := d.Get("image_id").(string)

	request.ImageId = imageID
	if v := d.Get("system_disk_description").(string); v != "" {
		request.SystemDiskDescription = v
	}
	systemDiskCategory := DiskCategory(d.Get("system_disk_category").(string))

	if v, ok := connectivity.GetResourceDataOk(d, "zone_id", "availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}

	request.SystemDiskCategory = string(systemDiskCategory)
	request.SystemDiskSize = strconv.Itoa(d.Get("system_disk_size").(int))

	if v, ok := d.GetOk("security_groups"); ok {
		sgs := expandStringList(v.(*schema.Set).List())
		request.SecurityGroupId = sgs[0]
	}

	if v := d.Get("instance_name").(string); v != "" {
		request.InstanceName = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	if v, ok := d.GetOk("internet_max_bandwidth_in"); ok {
		request.InternetMaxBandwidthIn = requests.NewInteger(v.(int))
	}

	if v := d.Get("host_name").(string); v != "" {
		request.HostName = v
	}

	if v := d.Get("password").(string); v != "" {
		request.Password = v
	}

	if v := d.Get("kms_encrypted_password").(string); v != "" {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return request, errmsgs.WrapError(err)
		}
		request.Password = decryptResp.Plaintext
	}

	vswitchValue := d.Get("subnet_id").(string)
	if vswitchValue == "" {
		vswitchValue = d.Get("vswitch_id").(string)
	}
	if vswitchValue != "" {
		request.VSwitchId = vswitchValue
		if v, ok := d.GetOk("private_ip"); ok && v.(string) != "" {
			request.PrivateIpAddress = v.(string)
		}
	}

	if v := d.Get("user_data").(string); v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v)
		if base64DecodeError == nil {
			request.UserData = v
		} else {
			request.UserData = base64.StdEncoding.EncodeToString([]byte(v))
		}
	}

	if v := d.Get("role_name").(string); v != "" {
		request.RamRoleName = v
	}

	if v := d.Get("key_name").(string); v != "" {
		request.KeyPairName = v
	}
	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	i := d.Get("storage_set_partition_number").(int)
	if v := d.Get("storage_set_id").(string); v != "" {
		request.StorageSetId = v
		if i >= 1 {
			request.StorageSetPartitionNumber = requests.NewInteger(d.Get("storage_set_partition_number").(int))
		} else {
			return nil, fmt.Errorf("can't empty storage_set_partition_number when you set storage_set_id and >=2 ")
		}
	}

	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		request.SecurityEnhancementStrategy = v.(string)
	}

	v, ok := d.GetOk("tags")
	if ok && len(v.(map[string]interface{})) > 0 {
		tags := make([]ecs.RunInstancesTag, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.RunInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	if v, ok := d.GetOk("data_disks"); ok {
		disks := v.([]interface{})
		var dataDiskRequests []ecs.RunInstancesDataDisk
		for i := range disks {
			disk := disks[i].(map[string]interface{})

			dataDiskRequest := ecs.RunInstancesDataDisk{
				Category:           disk["category"].(string),
				DeleteWithInstance: strconv.FormatBool(disk["delete_with_instance"].(bool)),
				Encrypted:          strconv.FormatBool(disk["encrypted"].(bool)),
			}
			if enc, ok := disk["encrypted"]; ok {
				if enc.(bool) {
					if j, ok := disk["kms_key_id"]; ok {
						dataDiskRequest.KMSKeyId = j.(string)
					}
					if dataDiskRequest.KMSKeyId == "" {
						return nil, errmsgs.WrapError(errors.New("KmsKeyId can not be empty if encrypted is set to \"true\""))
					}
				}
			}
			if kms, ok := disk["kms_key_id"]; ok {
				dataDiskRequest.KMSKeyId = kms.(string)
			}
			if name, ok := disk["name"]; ok {
				dataDiskRequest.DiskName = name.(string)
			}
			if snapshotId, ok := disk["snapshot_id"]; ok {
				dataDiskRequest.SnapshotId = snapshotId.(string)
			}
			if description, ok := disk["description"]; ok {
				dataDiskRequest.Description = description.(string)
			}

			dataDiskRequest.Size = fmt.Sprintf("%d", disk["size"].(int))
			dataDiskRequest.Category = disk["category"].(string)
			if dataDiskRequest.Category == string(DiskEphemeralSSD) {
				dataDiskRequest.DeleteWithInstance = ""
			}

			dataDiskRequests = append(dataDiskRequests, dataDiskRequest)
		}
		request.DataDisk = &dataDiskRequests
	}
	return request, nil
}

func modifyInstanceImage(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChanges("image_id", "system_disk_size") {
		update = true
		if !run {
			return update, nil
		}
		instance, err := ecsService.DescribeInstance(d.Id())
		if err != nil {
			return update, errmsgs.WrapError(err)
		}
		keyPairName := instance.KeyPairName
		request := ecs.CreateReplaceSystemDiskRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.ImageId = d.Get("image_id").(string)
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		request.ClientToken = buildClientToken(request.GetActionName())
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ReplaceSystemDisk(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*ecs.ReplaceSystemDiskResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return update, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// Ensure instance's image has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, errDesc := ecsService.DescribeInstance(d.Id())
			if errDesc != nil {
				return update, errmsgs.WrapError(errDesc)
			}
			var disks []ecs.Disk
			err := resource.Retry(2*time.Minute, func() *resource.RetryError {
				disks, err = ecsService.DescribeInstanceDisksByType(d.Id(), client.ResourceGroup, "system")
				if err != nil {
					if errmsgs.NotFoundError(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return update, errmsgs.WrapError(err)
			}

			if instance.ImageId == d.Get("image_id") && disks[0].Size == d.Get("system_disk_size").(int) {
				break
			}
			time.Sleep(DefaultIntervalShort * time.Second)

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, errmsgs.WrapError(errmsgs.GetTimeErrorFromString(fmt.Sprintf("Replacing instance %s system disk timeout.", d.Id())))
			}
		}

		// After updating image, it need to re-attach key pair
		if keyPairName != "" {
			if err := ecsService.AttachKeyPair(keyPairName, []interface{}{d.Id()}); err != nil {
				return update, errmsgs.WrapError(err)
			}
		}
	}
	return update, nil
}

func modifyInstanceAttribute(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	reboot := false
	request := ecs.CreateModifyInstanceAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()

	if d.HasChange("instance_name") {
		//d.SetPartial("instance_name")
		request.InstanceName = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("description") {
		//d.SetPartial("description")
		request.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("user_data") {
		//d.SetPartial("user_data")
		old, new := d.GetChange("user_data")
		old_s := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(old)))
		if fmt.Sprint(new) != old_s {

			if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
				_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
				if base64DecodeError == nil {
					request.UserData = v.(string)
				} else {
					request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
				}
			}
			update = true
			reboot = true
		}
	}

	if d.HasChange("host_name") {
		//d.SetPartial("host_name")
		request.HostName = d.Get("host_name").(string)
		update = true
		reboot = true
	}

	if d.HasChanges("password", "kms_encrypted_password") {
		if v := d.Get("password").(string); v != "" {
			//d.SetPartial("password")
			request.Password = v
			update = true
			reboot = true
		}
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return reboot, errmsgs.WrapError(err)
			}
			request.Password = decryptResp.Plaintext
			//d.SetPartial("kms_encrypted_password")
			//d.SetPartial("kms_encryption_context")
			update = true
			reboot = true
		}
	}

	if update {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"InvalidChargeType.ValueNotSupported"}) {
					time.Sleep(time.Minute)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return reboot, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	return reboot, nil
}

func modifyVpcAttribute(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}

	update := false
	request := ecs.CreateModifyInstanceVpcAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()
	request.VSwitchId = d.Get("vswitch_id").(string)

	if d.HasChange("vswitch_id") {
		update = true
		if d.Get("vswitch_id").(string) == "" {
			return update, errmsgs.WrapError(errmsgs.Error("Field 'vswitch_id' is required when modifying the instance VPC attribute."))
		}
		//d.SetPartial("vswitch_id")
	}

	if d.HasChange("subnet_id") {
		update = true
		if d.Get("subnet_id").(string) == "" {
			return update, errmsgs.WrapError(errmsgs.Error("Field 'subnet_id' is required when modifying the instance VPC attribute."))
		}
		request.VSwitchId = d.Get("subnet_id").(string)
		//d.SetPartial("subnet_id")
	}

	if request.VSwitchId != "" && d.HasChange("private_ip") {
		request.PrivateIpAddress = d.Get("private_ip").(string)
		update = true
		//d.SetPartial("private_ip")
	}

	if !run {
		return update, nil
	}

	if update {
		client := meta.(*connectivity.AlibabacloudStackClient)
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceVpcAttribute(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"OperationConflict"}) {
					time.Sleep(1 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return update, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		ecsService := EcsService{client}
		if err := ecsService.WaitForVpcAttributesChanged(d.Id(), request.VSwitchId, request.PrivateIpAddress); err != nil {
			return update, errmsgs.WrapError(err)
		}
	}
	return update, nil
}

func modifyInstanceType(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		d.Partial(false)
		return false, nil
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("instance_type") {
		update = true
		if !run {
			return update, nil
		}

		//An instance that was successfully modified once cannot be modified again within 5 minutes.
		request := ecs.CreateModifyInstanceSpecRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.InstanceType = d.Get("instance_type").(string)
		request.ClientToken = buildClientToken(request.GetActionName())

		err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			args := *request
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceSpec(&args)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return update, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}

		// Ensure instance's type has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, err := ecsService.DescribeInstance(d.Id())

			if err != nil {
				return update, errmsgs.WrapError(err)
			}

			if instance.InstanceType == d.Get("instance_type").(string) {
				break
			}

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, instance.InstanceType, d.Get("instance_type"), errmsgs.ProviderERROR)
			}

			time.Sleep(DefaultIntervalShort * time.Second)
		}
		//d.SetPartial("instance_type")
	}
	return update, nil
}

func modifyInstanceNetworkSpec(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	allocate := false
	update := false
	request := ecs.CreateModifyInstanceNetworkSpecRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())

	if d.HasChange("internet_max_bandwidth_out") {
		o, n := d.GetChange("internet_max_bandwidth_out")
		if o.(int) <= 0 && n.(int) > 0 {
			allocate = true
		}
		request.InternetMaxBandwidthOut = requests.NewInteger(n.(int))
		update = true
		//d.SetPartial("internet_max_bandwidth_out")
	}

	if d.HasChange("internet_max_bandwidth_in") {
		request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
		update = true
		//d.SetPartial("internet_max_bandwidth_in")
	}

	//An instance that was successfully modified once cannot be modified again within 5 minutes.
	wait := incrementalWait(2*time.Second, 2*time.Second)

	if update {
		if err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceNetworkSpec(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, "LastOrderProcessing", "LastRequestProcessing", "LastTokenProcessing"}) {
					wait()
					return resource.RetryableError(err)
				}
				if errmsgs.IsExpectedErrors(err, []string{"InternalError"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		ecsService := EcsService{client: client}

		deadline := time.Now().Add(DefaultTimeout * time.Second)
		for {
			instance, err := ecsService.DescribeInstance(d.Id())
			if err != nil {
				return errmsgs.WrapError(err)
			}

			if instance.InternetMaxBandwidthOut == d.Get("internet_max_bandwidth_out").(int) &&
				instance.InternetMaxBandwidthIn == d.Get("internet_max_bandwidth_in").(int) {
				break
			}

			if time.Now().After(deadline) {
				return errmsgs.WrapError(errmsgs.Error(`wait for internet update timeout! expect internet_charge_type value %s, get %s
					expect internet_max_bandwidth_out value %d, get %d, expect internet_max_bandwidth_out value %d, get %d,`,
					//d.Get("internet_charge_type").(string),
					"default",
					instance.InternetChargeType, d.Get("internet_max_bandwidth_out").(int),
					instance.InternetMaxBandwidthOut, d.Get("internet_max_bandwidth_in").(int), instance.InternetMaxBandwidthIn))
			}
			time.Sleep(1 * time.Second)
		}

		if allocate {
			request := ecs.CreateAllocatePublicIpAddressRequest()
			client.InitRpcRequest(*request.RpcRequest)
			request.InstanceId = d.Id()
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.AllocatePublicIpAddress(request)
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
	}
	return nil
}

func AssignIpv6AddressesFunc(id string, ipv6_addresses_count int, ipv6_addresses []string, meta interface{}) ([]string, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}
	var ipv6s []string
	request := ecs.CreateAssignIpv6AddressesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	instance, err := ecsService.DescribeInstance(id)

	if err != nil {
		return ipv6s, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "DescribeInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	request.NetworkInterfaceId = instance.NetworkInterfaces.NetworkInterface[0].NetworkInterfaceId
	request.Ipv6AddressCount = requests.NewInteger(ipv6_addresses_count)
	request.Ipv6Address = &ipv6_addresses
	raw, ipv6_err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.AssignIpv6Addresses(request)
	})
	if ipv6_err != nil {
		return ipv6s, errmsgs.WrapErrorf(ipv6_err, errmsgs.DefaultErrorMsg, "AssignIpv6Addresses", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("AssignIpv6Addresses", raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.AssignIpv6AddressesResponse)
	ipv6s = response.Ipv6Sets.Ipv6Address
	return ipv6s, nil
}

func ecsMergeTags(d *schema.ResourceData, disktags map[string]interface{}) map[string]interface{} {
	if intance_tags, ok := d.GetOk("tags"); ok && len(intance_tags.(map[string]interface{})) > 0 {
		mergedMap := make(map[string]interface{})
		for k, v := range intance_tags.(map[string]interface{}) {
			mergedMap[k] = v
		}
		for k, v := range disktags {
			mergedMap[k] = v
		}
		return mergedMap
	}

	return disktags
}

func getOnlySystemTags(d *schema.ResourceData, tags []ecs.Tag) []ecs.Tag {
	var only_system_tags []ecs.Tag
	old_s_tags := d.Get("system_disk_tags").(map[string]interface{})
	ecs_tags := d.Get("tags").(map[string]interface{})
	only_ecs_tags := make([]string, 0)
	// 获取只属于ecs的tags 的key列表
	for k, _ := range ecs_tags {
		if _, ok := old_s_tags[k]; !ok {
			only_ecs_tags = append(only_ecs_tags, k)
		}
	}
	// 剔除只属于ecs的tags
	for _, tag := range tags {
		in_only_ecs_tags := false
		for _, only_ecs_tag := range only_ecs_tags {
			if tag.TagKey == only_ecs_tag {
				in_only_ecs_tags = true
				break
			}
		}
		if !in_only_ecs_tags {
			only_system_tags = append(only_system_tags, tag)
		}
	}
	return only_system_tags
}
