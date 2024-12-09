package alibabacloudstack

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEssScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackEssScalingConfigurationCreate,
		Read:   resourceAlibabacloudStackEssScalingConfigurationRead,
		Update: resourceAlibabacloudStackEssScalingConfigurationUpdate,
		Delete: resourceAlibabacloudStackEssScalingConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Deprecated: "Field 'status' is deprecated and will be removed in a future release. Please use new field 'active' instead.",
				ConflictsWith: []string{"active"},
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"status"},
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^ecs\..*`), "prefix must be 'ecs.'"),
			},
			"instance_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				MaxItems: int(MaxScalingConfigurationInstanceTypes),
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"deployment_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scaling_configuration_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(2, 40),
			},
			"internet_max_bandwidth_in": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "ephemeral_ssd", "cloud_ssd", "cloud_efficiency", "cloud_pperf", "cloud_sperf"}, false),
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(40, 500),
			},
			"data_disk": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "cloud",
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_efficiency", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"substitute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk_auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Deprecated: "Field 'ram_role_name' is deprecated and will be removed in a future release. Please use new field 'role_name' instead.",
				ConflictsWith: []string{"role_name"},
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"ram_role_name"},
			},
			"key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				Deprecated: "Field 'key_pair_name' is deprecated and will be removed in a future release. Please use new field 'key_name' instead.",
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"key_pair_name"},
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ESS-Instance",
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"override": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request, err := buildAlibabacloudStackEssScalingConfigurationArgs(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	client.InitRpcRequest(*request.RpcRequest)
	if d.Get("is_outdated").(bool) == true {
		request.IoOptimized = string(NoneOptimized)
	} else {
		request.IoOptimized = string(IOOptimized)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingConfiguration(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling, "IncorrectScalingGroupStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if bresponse, ok := raw.(*ess.CreateScalingConfigurationResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ess_scalingconfiguration", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateScalingConfigurationResponse)
		d.SetId(response.ScalingConfigurationId)
		return nil
	}); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ess_scalingconfiguration", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return resourceAlibabacloudStackEssScalingConfigurationUpdate(d, meta)
}

func resourceAlibabacloudStackEssScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	d.Partial(true)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	if d.HasChanges("status", "active") {
		c, err := essService.DescribeEssScalingConfiguration(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return errmsgs.WrapError(err)
		}

		if connectivity.GetResourceData(d, "active", "status").(bool) {
			if c.LifecycleState == string(Inactive) {
				err := essService.ActiveEssScalingConfiguration(c.ScalingGroupId, d.Id())
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
		} else {
			if c.LifecycleState == string(Active) {
				_, err := activeSubstituteScalingConfiguration(d, meta)
				if err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}
	}

	if err := enableEssScalingConfiguration(d, meta); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := modifyEssScalingConfiguration(d, meta); err != nil {
		return errmsgs.WrapError(err)
	}

	d.Partial(false)

	return resourceAlibabacloudStackEssScalingConfigurationRead(d, meta)
}

func modifyEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyScalingConfigurationRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScalingConfigurationId = d.Id()

	if d.HasChange("override") {
		request.Override = requests.NewBoolean(d.Get("override").(bool))
	}

	if d.HasChange("image_id") || d.Get("override").(bool) {
		request.ImageId = d.Get("image_id").(string)
	}

	hasChangeInstanceType := d.HasChange("instance_type")
	hasChangeInstanceTypes := d.HasChange("instance_types")
	if hasChangeInstanceType || hasChangeInstanceTypes || d.Get("override").(bool) {
		instanceType := d.Get("instance_type").(string)
		instanceTypes := d.Get("instance_types").([]interface{})
		if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
			return fmt.Errorf("instance_type must be assigned")
		}
		types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
		if instanceTypes != nil && len(instanceTypes) > 0 {
			types = expandStringList(instanceTypes)
		}
		if instanceType != "" {
			types = append(types, instanceType)
		}
		request.InstanceTypes = &types
	}

	hasChangeSecurityGroupId := d.HasChange("security_group_ids")
	if hasChangeSecurityGroupId || d.Get("override").(bool) {
		securityGroupIds := d.Get("security_group_ids").([]interface{})
		if len(securityGroupIds) <= 0 {
			return fmt.Errorf("securityGroupIds must be assigned")
		}
		sgs := make([]string, 0, len(securityGroupIds))
		for _, sg := range securityGroupIds {
			sgs = append(sgs, sg.(string))
		}
		request.SecurityGroupIds = &sgs
	}

	if d.HasChange("scaling_configuration_name") {
		request.ScalingConfigurationName = d.Get("scaling_configuration_name").(string)
	}

	if d.HasChange("system_disk_category") {
		request.SystemDiskCategory = d.Get("system_disk_category").(string)
	}

	if d.HasChange("system_disk_size") {
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	}

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
			if base64DecodeError == nil {
				request.UserData = v.(string)
			} else {
				request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
			}
		}
	}

	if d.HasChanges("role_name", "ram_role_name") {
		request.RamRoleName = connectivity.GetResourceData(d, "role_name", "ram_role_name").(string)
	}

	if d.HasChanges("key_name", "key_pair_name") {
		request.KeyPairName = connectivity.GetResourceData(d, "key_name", "key_pair_name").(string)
	}

	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
	}
	if d.HasChange("system_disk_auto_snapshot_policy_id") {
		request.SystemDiskAutoSnapshotPolicyId = d.Get("system_disk_auto_snapshot_policy_id").(string)
	}
	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := "{"
			for key, value := range v.(map[string]interface{}) {
				tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
			}
			request.Tags = strings.TrimSuffix(tags, ",") + "}"
		}
	}
	if d.HasChange("host_name") {
		request.HostName = d.Get("host_name").(string)
	}
	if d.HasChange("data_disk") {
		dds, ok := d.GetOk("data_disk")
		if ok {
			disks := dds.([]interface{})
			createDataDisks := make([]ess.ModifyScalingConfigurationDataDisk, 0, len(disks))
			for _, e := range disks {
				pack := e.(map[string]interface{})
				dataDisk := ess.ModifyScalingConfigurationDataDisk{
					Size:                 strconv.Itoa(pack["size"].(int)),
					Category:             pack["category"].(string),
					SnapshotId:           pack["snapshot_id"].(string),
					Encrypted:            pack["encrypted"].(string),
					KMSKeyId:             pack["kms_key_id"].(string),
					DeleteWithInstance:   strconv.FormatBool(pack["delete_with_instance"].(bool)),
					Device:               pack["device"].(string),
					Description:          pack["description"].(string),
					AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
					DiskName:             pack["name"].(string),
				}
				createDataDisks = append(createDataDisks, dataDisk)
			}
			request.DataDisk = &createDataDisks
		}
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingConfiguration(request)
	})
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*ess.ModifyScalingConfigurationResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func enableEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	if d.HasChange("enable") {
		sgId := d.Get("scaling_group_id").(string)
		group, err := essService.DescribeEssScalingGroup(sgId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if d.Get("enable").(bool) {
			if group.LifecycleState == string(Inactive) {
				object, err := essService.DescribeEssScalingConfifurations(sgId)
				if err != nil {
					return errmsgs.WrapError(err)
				}
				activeConfig := ""
				var csIds []string
				for _, c := range object {
					csIds = append(csIds, c.ScalingConfigurationId)
					if c.LifecycleState == string(Active) {
						activeConfig = c.ScalingConfigurationId
					}
				}

				if activeConfig == "" {
					return errmsgs.WrapError(errmsgs.Error("Please active a scaling configuration before enabling scaling group %s. Its all scaling configuration are %s.",
						sgId, strings.Join(csIds, ",")))
				}

				request := ess.CreateEnableScalingGroupRequest()
				client.InitRpcRequest(*request.RpcRequest)

				request.ScalingGroupId = sgId
				request.ActiveScalingConfigurationId = activeConfig

				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(request)
				})
				if err != nil {
					errmsg := ""
					if bresponse, ok := raw.(*ess.EnableScalingGroupResponse); ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Active, DefaultTimeout); err != nil {
					return errmsgs.WrapError(err)
				}
			}
		} else {
			if group.LifecycleState == string(Active) {
				request := ess.CreateDisableScalingGroupRequest()
				client.InitRpcRequest(*request.RpcRequest)

				request.ScalingGroupId = sgId
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DisableScalingGroup(request)
				})
				if err != nil {
					errmsg := ""
					if bresponse, ok := raw.(*ess.DisableScalingGroupResponse); ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Inactive, DefaultTimeout); err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackEssScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}
	object, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	connectivity.SetResourceData(d, object.LifecycleState == string(Active), "active", "status")
	d.Set("image_id", object.ImageId)
	d.Set("scaling_configuration_name", object.ScalingConfigurationName)
	d.Set("internet_max_bandwidth_in", object.InternetMaxBandwidthIn)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("system_disk_auto_snapshot_policy_id", object.SystemDiskAutoSnapshotPolicyId)
	d.Set("data_disk", essService.flattenDataDiskMappings(object.DataDisks.DataDisk))
	connectivity.SetResourceData(d, object.RamRoleName, "role_name", "ram_role_name")
	connectivity.SetResourceData(d, object.KeyPairName, "key_name", "key_pair_name")
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("tags", essTagsToMap(object.Tags.Tag))
	d.Set("instance_name", object.InstanceName)
	d.Set("override", d.Get("override").(bool))
	d.Set("host_name", object.HostName)
	if sg, ok := d.GetOk("security_group_ids"); ok && len(sg.([]interface{})) >= 0 {
		d.Set("security_group_ids", object.SecurityGroupIds.SecurityGroupId)
	}

	if instanceType, ok := d.GetOk("instance_type"); ok && instanceType.(string) != "" {
		d.Set("instance_type", object.InstanceType)
	}
	if instanceTypes, ok := d.GetOk("instance_types"); ok && len(instanceTypes.([]interface{})) > 0 {
		d.Set("instance_types", object.InstanceTypes.InstanceType)
	}
	userData := d.Get("user_data")
	if userData.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(userData.(string))
		if base64DecodeError == nil {
			d.Set("user_data", object.UserData)
		} else {
			d.Set("user_data", userDataHashSum(object.UserData))
		}
	} else {
		d.Set("user_data", userDataHashSum(object.UserData))
	}
	return nil
}

func resourceAlibabacloudStackEssScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	object, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScalingGroupId = object.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*ess.DescribeScalingConfigurationsResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return nil
	} else if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
		if d.Get("force_delete").(bool) {
			request := ess.CreateDeleteScalingGroupRequest()
			client.InitRpcRequest(*request.RpcRequest)

			request.ScalingGroupId = object.ScalingGroupId
			request.ForceDelete = requests.NewBoolean(true)

			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DeleteScalingGroup(request)
			})

			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
					return nil
				}
				errmsg := ""
				if bresponse, ok := raw.(*ess.DeleteScalingGroupResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return errmsgs.WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
		}
		return errmsgs.WrapError(errmsgs.Error("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), object.ScalingGroupId))
	}

	deleteScalingConfigurationRequest := ess.CreateDeleteScalingConfigurationRequest()
	client.InitRpcRequest(*deleteScalingConfigurationRequest.RpcRequest)

	deleteScalingConfigurationRequest.ScalingConfigurationId = d.Id()

	rawDeleteScalingConfiguration, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingConfiguration(deleteScalingConfigurationRequest)
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound", "InvalidScalingConfigurationId.NotFound"}) {
			return nil
		}
		errmsg := ""
		if bresponse, ok := rawDeleteScalingConfiguration.(*ess.DeleteScalingConfigurationResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), rawDeleteScalingConfiguration, request.RpcRequest, request)

	return errmsgs.WrapError(essService.WaitForScalingConfiguration(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssScalingConfigurationArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingConfigurationRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateCreateScalingConfigurationRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	security_group_ids := d.Get("security_group_ids").([]interface{})
	if len(security_group_ids) <= 0 {
		return nil, errmsgs.WrapError(errmsgs.Error("security_group_ids must be assigned"))
	}
	sgs := make([]string, 0, len(security_group_ids))
	for _, v := range security_group_ids {
		sgs = append(sgs, v.(string))
	}

	request.SecurityGroupIds = &sgs
	request.DeploymentSetId = d.Get("deployment_set_id").(string)
	request.InstanceName = d.Get("instance_name").(string)

	types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
	instanceType := d.Get("instance_type").(string)
	instanceTypes := d.Get("instance_types").([]interface{})
	if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
		return nil, errmsgs.WrapError(errmsgs.Error("instance_type or instance_types must be assigned"))
	}

	if instanceTypes != nil && len(instanceTypes) > 0 {
		types = expandStringList(instanceTypes)
	}

	if instanceType != "" {
		types = append(types, instanceType)
	}

	request.InstanceTypes = &types

	if v := d.Get("scaling_configuration_name").(string); v != "" {
		request.ScalingConfigurationName = v
	}
	if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
		request.InternetMaxBandwidthIn = requests.NewInteger(v)
	}

	if v := d.Get("system_disk_category").(string); v != "" {
		request.SystemDiskCategory = v
	}
	if v := d.Get("system_disk_auto_snapshot_policy_id").(string); v != "" {
		request.SystemDiskAutoSnapshotPolicyId = v
	}
	if v := d.Get("system_disk_size").(int); v != 0 {
		request.SystemDiskSize = requests.NewInteger(v)
	}

	dds, ok := d.GetOk("data_disk")
	if ok {
		disks := dds.([]interface{})
		createDataDisks := make([]ess.CreateScalingConfigurationDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := ess.CreateScalingConfigurationDataDisk{
				Size:                 strconv.Itoa(pack["size"].(int)),
				Category:             pack["category"].(string),
				SnapshotId:           pack["snapshot_id"].(string),
				DeleteWithInstance:   strconv.FormatBool(pack["delete_with_instance"].(bool)),
				Device:               pack["device"].(string),
				Encrypted:            pack["encrypted"].(string),
				KMSKeyId:             pack["kms_key_id"].(string),
				DiskName:             pack["name"].(string),
				Description:          pack["description"].(string),
				AutoSnapshotPolicyId: pack["auto_snapshot_policy_id"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		request.DataDisk = &createDataDisks
	}

	if v, ok :=  connectivity.GetResourceDataOk(d, "role_name", "ram_role_name"); ok && v.(string) != "" {
		request.RamRoleName = v.(string)
	}
	
	if v, ok := connectivity.GetResourceDataOk(d, "key_name", "key_pair_name"); ok && v.(string) != "" {
		request.KeyPairName = v.(string)
	}
	
	if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			request.UserData = v.(string)
		} else {
			request.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := "{"
		for key, value := range v.(map[string]interface{}) {
			tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
		}
		request.Tags = strings.TrimSuffix(tags, ",") + "}"
	}

	if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
		request.InstanceName = v.(string)
	}
	if v, ok := d.GetOk("host_name"); ok && v.(string) != "" {
		request.HostName = v.(string)
	}
	return request, nil
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) (configures []ess.ScalingConfigurationInDescribeScalingConfigurations, err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	substituteId, ok := d.GetOk("substitute")

	c, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*ess.DescribeScalingConfigurationsResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return
	}

	if !ok || substituteId.(string) == "" {
		if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
			return configures, errmsgs.WrapError(errmsgs.Error("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId))
		}

		var configs []string
		for _, cc := range response.ScalingConfigurations.ScalingConfiguration {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return configures, errmsgs.WrapError(errmsgs.Error("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ",")))
	}

	err = essService.ActiveEssScalingConfiguration(c.ScalingGroupId, substituteId.(string))
	if err != nil {
		return configures, errmsgs.WrapError(errmsgs.Error("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substituteId.(string)))
	}

	return response.ScalingConfigurations.ScalingConfiguration, nil
}
