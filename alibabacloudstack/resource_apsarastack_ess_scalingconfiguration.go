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
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			"io_optimized": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_outdated": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_group_id": {
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
					},
				},
			},

			"substitute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceAlibabacloudStackEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	// Ensure instance_type is generation three
	client := meta.(*connectivity.AlibabacloudStackClient)
	request, err := buildAlibabacloudStackEssScalingConfigurationArgs(d, meta)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if err != nil {
		return WrapError(err)
	}

	request.IoOptimized = string(IOOptimized)
	if d.Get("is_outdated").(bool) == true {
		request.IoOptimized = string(NoneOptimized)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingConfiguration(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling, "IncorrectScalingGroupStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateScalingConfigurationResponse)
		d.SetId(response.ScalingConfigurationId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ess_scalingconfiguration", request.GetActionName(), AlibabacloudStackSdkGoERROR)
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

	if d.HasChange("active") {
		c, err := essService.DescribeEssScalingConfiguration(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}

		if d.Get("active").(bool) {
			if c.LifecycleState == string(Inactive) {

				err := essService.ActiveEssScalingConfiguration(c.ScalingGroupId, d.Id())
				if err != nil {
					return WrapError(err)
				}
			}
		} else {
			if c.LifecycleState == string(Active) {
				_, err := activeSubstituteScalingConfiguration(d, meta)
				if err != nil {
					return WrapError(err)
				}
			}
		}
		//d.SetPartial("active")
	}

	if err := enableEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}

	if err := modifyEssScalingConfiguration(d, meta); err != nil {
		return WrapError(err)
	}

	d.Partial(false)

	return resourceAlibabacloudStackEssScalingConfigurationRead(d, meta)
}

func modifyEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ess.CreateModifyScalingConfigurationRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ScalingConfigurationId = d.Id()

	if d.HasChange("override") {
		request.Override = requests.NewBoolean(d.Get("override").(bool))
		//d.SetPartial("override")
	}

	if d.HasChange("image_id") || d.Get("override").(bool) {
		request.ImageId = d.Get("image_id").(string)
		//d.SetPartial("image_id")
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

	hasChangeSecurityGroupId := d.HasChange("security_group_id")
	if hasChangeSecurityGroupId || d.Get("override").(bool) {
		securityGroupId := d.Get("security_group_id").(string)
		if securityGroupId == "" {
			return fmt.Errorf("securityGroupId must be assigned")

		}

		if securityGroupId != "" {
			request.SecurityGroupId = securityGroupId
		}
	}

	if d.HasChange("scaling_configuration_name") {
		request.ScalingConfigurationName = d.Get("scaling_configuration_name").(string)
		//d.SetPartial("scaling_configuration_name")
	}

	if d.HasChange("system_disk_category") {
		request.SystemDiskCategory = d.Get("system_disk_category").(string)
		//d.SetPartial("system_disk_category")
	}

	if d.HasChange("system_disk_size") {
		request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		//d.SetPartial("system_disk_size")
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
		//d.SetPartial("user_data")
	}

	if d.HasChange("role_name") {
		request.RamRoleName = d.Get("role_name").(string)
		//d.SetPartial("role_name")
	}

	if d.HasChange("key_name") {
		request.KeyPairName = d.Get("key_name").(string)
		//d.SetPartial("key_name")
	}

	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		//d.SetPartial("instance_name")
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := "{"
			for key, value := range v.(map[string]interface{}) {
				tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
			}
			request.Tags = strings.TrimSuffix(tags, ",") + "}"
		}
		//d.SetPartial("tags")
	}

	if d.HasChange("data_disk") {
		dds, ok := d.GetOk("data_disk")
		if ok {
			disks := dds.([]interface{})
			createDataDisks := make([]ess.ModifyScalingConfigurationDataDisk, 0, len(disks))
			for _, e := range disks {
				pack := e.(map[string]interface{})
				dataDisk := ess.ModifyScalingConfigurationDataDisk{
					Size:               strconv.Itoa(pack["size"].(int)),
					Category:           pack["category"].(string),
					SnapshotId:         pack["snapshot_id"].(string),
					Encrypted:          pack["encrypted"].(string),
					KMSKeyId:           pack["kms_key_id"].(string),
					DeleteWithInstance: strconv.FormatBool(pack["delete_with_instance"].(bool)),
				}
				createDataDisks = append(createDataDisks, dataDisk)
			}
			request.DataDisk = &createDataDisks
		}
		//d.SetPartial("data_disk")
	}
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingConfiguration(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
			return WrapError(err)
		}

		if d.Get("enable").(bool) {
			if group.LifecycleState == string(Inactive) {

				object, err := essService.DescribeEssScalingConfifurations(sgId)

				if err != nil {
					return WrapError(err)
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
					return WrapError(Error("Please active a scaling configuration before enabling scaling group %s. Its all scaling configuration are %s.",
						sgId, strings.Join(csIds, ",")))
				}

				request := ess.CreateEnableScalingGroupRequest()
				request.RegionId = client.RegionId
				if strings.ToLower(client.Config.Protocol) == "https" {
					request.Scheme = "https"
				} else {
					request.Scheme = "http"
				}
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

				request.ScalingGroupId = sgId
				request.ActiveScalingConfigurationId = activeConfig

				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Active, DefaultTimeout); err != nil {
					return WrapError(err)
				}

				//d.SetPartial("scaling_configuration_id")
			}
		} else {
			if group.LifecycleState == string(Active) {
				request := ess.CreateDisableScalingGroupRequest()
				request.RegionId = client.RegionId
				if strings.ToLower(client.Config.Protocol) == "https" {
					request.Scheme = "https"
				} else {
					request.Scheme = "http"
				}
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

				request.ScalingGroupId = sgId
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DisableScalingGroup(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err := essService.WaitForEssScalingGroup(sgId, Inactive, DefaultTimeout); err != nil {
					return WrapError(err)
				}
			}
		}
		//d.SetPartial("enable")
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
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("active", object.LifecycleState == string(Active))
	d.Set("image_id", object.ImageId)
	d.Set("scaling_configuration_name", object.ScalingConfigurationName)
	d.Set("internet_max_bandwidth_in", object.InternetMaxBandwidthIn)
	d.Set("system_disk_category", object.SystemDiskCategory)
	d.Set("system_disk_size", object.SystemDiskSize)
	d.Set("data_disk", essService.flattenDataDiskMappings(object.DataDisks.DataDisk))
	d.Set("role_name", object.RamRoleName)
	d.Set("key_name", object.KeyPairName)
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("tags", essTagsToMap(object.Tags.Tag))
	d.Set("instance_name", object.InstanceName)
	d.Set("override", d.Get("override").(bool))

	if sg, ok := d.GetOk("security_group_id"); ok && sg.(string) != "" {
		d.Set("security_group_id", object.SecurityGroupId)
	}
	if sgs, ok := d.GetOk("security_group_ids"); ok && len(sgs.([]interface{})) > 0 {
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
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ScalingGroupId = object.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return nil
	} else if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
		if d.Get("force_delete").(bool) {
			request := ess.CreateDeleteScalingGroupRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

			request.ScalingGroupId = object.ScalingGroupId
			request.ForceDelete = requests.NewBoolean(true)
			request.RegionId = client.RegionId
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DeleteScalingGroup(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
					return nil
				}
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
		}
		return WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), object.ScalingGroupId))
	}

	deleteScalingConfigurationRequest := ess.CreateDeleteScalingConfigurationRequest()
	deleteScalingConfigurationRequest.ScalingConfigurationId = d.Id()
	if strings.ToLower(client.Config.Protocol) == "https" {
		deleteScalingConfigurationRequest.Scheme = "https"
	} else {
		deleteScalingConfigurationRequest.Scheme = "http"
	}
	deleteScalingConfigurationRequest.Headers = map[string]string{"RegionId": client.RegionId}
	deleteScalingConfigurationRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	rawDeleteScalingConfiguration, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingConfiguration(deleteScalingConfigurationRequest)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound", "InvalidScalingConfigurationId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), rawDeleteScalingConfiguration, request.RpcRequest, request)

	return WrapError(essService.WaitForScalingConfiguration(d.Id(), Deleted, DefaultTimeout))
}

func buildAlibabacloudStackEssScalingConfigurationArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingConfigurationRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := ess.CreateCreateScalingConfigurationRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	request.ImageId = d.Get("image_id").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)

	securityGroupId := d.Get("security_group_id").(string)

	if securityGroupId == "" {
		return nil, WrapError(Error("security_group_id must be assigned"))
	}

	if securityGroupId != "" {
		request.SecurityGroupId = securityGroupId
	}

	types := make([]string, 0, int(MaxScalingConfigurationInstanceTypes))
	instanceType := d.Get("instance_type").(string)
	instanceTypes := d.Get("instance_types").([]interface{})
	if instanceType == "" && (instanceTypes == nil || len(instanceTypes) == 0) {
		return nil, WrapError(Error("instance_type or instance_types must be assigned"))
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
				Size:               strconv.Itoa(pack["size"].(int)),
				Category:           pack["category"].(string),
				SnapshotId:         pack["snapshot_id"].(string),
				Encrypted:          pack["encrypted"].(string),
				KMSKeyId:           pack["kms_key_id"].(string),
				DeleteWithInstance: strconv.FormatBool(pack["delete_with_instance"].(bool)),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		request.DataDisk = &createDataDisks
	}

	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		request.RamRoleName = v.(string)
	}

	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
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

	return request, nil
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) (configures []ess.ScalingConfiguration, err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	substituteId, ok := d.GetOk("substitute")

	c, err := essService.DescribeEssScalingConfiguration(d.Id())
	if err != nil {
		err = WrapError(err)
		return
	}

	request := ess.CreateDescribeScalingConfigurationsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if len(response.ScalingConfigurations.ScalingConfiguration) < 1 {
		return
	}

	if !ok || substituteId.(string) == "" {

		if len(response.ScalingConfigurations.ScalingConfiguration) == 1 {
			return configures, WrapError(Error("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId))
		}

		var configs []string
		for _, cc := range response.ScalingConfigurations.ScalingConfiguration {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return configures, WrapError(Error("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ",")))

	}

	err = essService.ActiveEssScalingConfiguration(c.ScalingGroupId, substituteId.(string))
	if err != nil {
		return configures, WrapError(Error("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substituteId.(string)))
	}

	return response.ScalingConfigurations.ScalingConfiguration, nil
}
