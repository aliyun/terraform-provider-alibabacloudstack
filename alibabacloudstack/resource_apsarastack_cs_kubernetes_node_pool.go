package alibabacloudstack

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	roacs "github.com/alibabacloud-go/cs-20151215/v5/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ResourceName = "resource_alibabacloudstack_cs_kubernetes_permissions"

func resourceAlibabacloudStackCSKubernetesNodePool() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instances"},
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
			},
			"instance_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 10,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key_name", "kms_encrypted_password"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "kms_encrypted_password"},
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"system_disk_performance_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
				DiffSuppressFunc: csNodepoolDiskPerformanceLevelDiffSuppressFunc,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AliyunLinux", "Windows", "CentOS", "WindowsCore"}, false),
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12, 24, 36, 48, 60}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Month"}, false),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csNodepoolInstancePostPaidDiffSuppressFunc,
			},
			"install_cloud_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"unschedulable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"data_disks": {
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
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk", "cloud_pperf", "cloud_sperf"}, false),
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"labels": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"taints": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scaling_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"max_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"cpu", "gpu", "gpushare", "spot"}, false),
						},
						"is_bond_eip": {
							Type:          schema.TypeBool,
							Optional:      true,
							ConflictsWith: []string{"internet_charge_type"},
						},
						"eip_internet_charge_type": {
							Type:          schema.TypeString,
							Optional:      true,
							ValidateFunc:  validation.StringInSlice([]string{"PayByBandwidth", "PayByTraffic"}, false),
							ConflictsWith: []string{"internet_charge_type"},
						},
						"eip_bandwidth": {
							Type:          schema.TypeInt,
							Optional:      true,
							ValidateFunc:  validation.IntBetween(1, 500),
							ConflictsWith: []string{"internet_charge_type"},
						},
					},
				},
				ConflictsWith: []string{"instances"},
			},
			"scaling_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"release", "recycle"}, false),
				DiffSuppressFunc: csNodepoolScalingPolicyDiffSuppressFunc,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayByTraffic", "PayByBandwidth"}, false),
			},
			"internet_max_bandwidth_out": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SpotWithPriceLimit"}, false),
			},
			"spot_price_limit": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"price_limit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csNodepoolSpotInstanceSettingDiffSuppressFunc,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:      100,
				ConflictsWith: []string{"node_count", "scaling_config"},
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"format_disk": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCSKubernetesNodePoolCreate, resourceAlibabacloudStackCSNodePoolRead, resourceAlibabacloudStackCSKubernetesNodePoolUpdate, resourceAlibabacloudStackCSNodePoolDelete)
	return resource
}

type NodePoolCommonResponse struct {
	Response
	NodePoolID string `json:"nodepool_id,omitempty"`

	TaskId string `json:"task_id,omitempty"`
}
type nodePoolDataDisk struct {
	Category string `json:"category"`

	Encrypted string `json:"encrypted"` // true|false

	Size int `json:"size"`
}
type autoScaling struct {
	Enable       bool   `json:"enable"`
	MaxInstances int64  `json:"max_instances"`
	MinInstances int64  `json:"min_instances"`
	Type         string `json:"type"`
}
type kubernetesConfig struct {
	Taints []cs.Taint `json:"taints"`
	Labels []cs.Label `json:"labels"`

	UserData string `json:"user_data"`

	Runtime        string `json:"runtime,omitempty"`
	RuntimeVersion string `json:"runtime_version"`
	CmsEnabled     bool   `json:"cms_enabled"`

	Unschedulable bool `json:"unschedulable"`
}
type tEEConfig struct {
	TEEEnable bool `json:"tee_enable"`
}
type CreateClusterNodePoolRequest struct {
	Count            int64            `json:"count"`
	NodePoolInfo     NodePoolInfo     `json:"node_pool_info"`
	ScalingGroup     scalingGroup     `json:"scaling_group"`
	KubernetesConfig kubernetesConfig `json:"kubernetes_config"`
	AutoScaling      autoScaling      `json:"auto_scaling"`
}

type NodePoolInfo struct {
	NodePoolId      string    `json:"nodepool_id"`
	RegionId        string    `json:"region_id"`
	Name            string    `json:"name"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	IsDefault       bool      `json:"is_default"`
	NodePoolType    string    `json:"type"`
	ResourceGroupId string    `json:"resource_group_id"`
}
type scalingGroup struct {
	VswitchIds    []string `json:"vswitch_ids"`
	InstanceTypes []string `json:"instance_types"`
	LoginPassword string   `json:"login_password"`

	SystemDiskCategory string `json:"system_disk_category"`
	SystemDiskSize     int64  `json:"system_disk_size"`

	DataDisks []nodePoolDataDisk `json:"data_disks"` //支持多个数据盘
	Tags      []cs.Tag           `json:"tags"`
	ImageId   string             `json:"image_id"`
	Platform  string             `json:"platform"`
	// 支持包年包月
	InstanceChargeType string `json:"instance_charge_type"`

	ScalingPolicy string `json:"scaling_policy"`

	// 公网ip
	InternetChargeType      string `json:"internet_charge_type"`
	InternetMaxBandwidthOut int    `json:"internet_max_bandwidth_out"`
}
type CreateNodePoolRequest struct {
	ClusterID        string           `json:"ClusterId"`
	NodepoolID       string           `json:"NodepoolId"`
	UpdateNodes      bool             `json:"update_nodes"`
	ScalingGroup     scalingGroup     `json:"scaling_group"`
	KubernetesConfig kubernetesConfig `json:"kubernetes_config"`
	TEEConfig        tEEConfig        `json:"tee_config"`
	AutoScaling      autoScaling      `json:"auto_scaling"`
}

func resourceAlibabacloudStackCSKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	// prepare args and set default value
	request, err := buildNodePoolArgs(d, meta)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes_node_pool", "PrepareKubernetesNodePoolArgs", err)
	}
	request.Headers["x-acs-asapi-gateway-version"] = "3.0"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cs_kubernetes_node_pool", "CreateClusterNodePool", response, errmsg)
	}
	nodepoolresponse := NodePoolCommonResponse{}
	if err := json.Unmarshal(response.GetHttpContentBytes(), &nodepoolresponse); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_cs_kubernetes_node_pool", "NodePoolCommonResponse", response)
	}

	d.SetId(nodepoolresponse.NodePoolID)

	// reset interval to 10s
	stateConf := BuildStateConf([]string{"initial", "scaling"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), d.Get("cluster_id").(string), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, "ResourceID:%s , TaskID:%s ", d.Id(), nodepoolresponse.TaskId)
	}

	// attach existing node
	if v, ok := d.GetOk("instances"); ok && v != nil {
		attachExistingInstance(d, meta)
	}

	return nil
}

func resourceAlibabacloudStackCSKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	clusterId := d.Get("cluster_id").(string)
	d.Partial(true)
	update := false

	args := &CreateNodePoolRequest{
		ClusterID:        clusterId,
		NodepoolID:       d.Id(),
		UpdateNodes:      true,
		ScalingGroup:     scalingGroup{},
		KubernetesConfig: kubernetesConfig{},
		TEEConfig:        tEEConfig{},
		AutoScaling:      autoScaling{},
	}
	if d.HasChange("node_count") {
		oldV, newV := d.GetChange("node_count")

		oldValue, ok := oldV.(int)
		if !ok {
			return errmsgs.WrapErrorf(fmt.Errorf("node_count old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := newV.(int)
		if !ok {
			return errmsgs.WrapErrorf(fmt.Errorf("node_count new value can not be parsed"), "parseError %d", newValue)
		}
		log.Printf("UUUUUUUUUUUUUUUUUUUUUUUUUUUU %d , %d", newValue, oldValue)
		if newValue < oldValue {
			err := RemoveNodePoolNodes(d, meta, clusterId, d.Id(), nil, nil)
			if err != nil {
				return err
			}

			// The removal of a node is logically independent.
			// The removal of a node should not involve parameter changes.
			return nil

		}
		//update = true
		if newValue > oldValue {
			err := ScaleClusterNodePool(d, meta, clusterId, d.Id(), oldValue, newValue)
			if err != nil {
				return err
			}

			// The removal of a node is logically independent.
			// The removal of a node should not involve parameter changes.
			return nil

		}
	}

	if d.HasChange("vswitch_ids") {
		update = true

		args.ScalingGroup.VswitchIds = expandStringList(d.Get("vswitch_ids").([]interface{}))
	}

	if d.HasChange("install_cloud_monitor") {
		update = true
		args.KubernetesConfig.CmsEnabled = d.Get("install_cloud_monitor").(bool)
	}

	if d.HasChange("unschedulable") {
		update = true
		args.KubernetesConfig.Unschedulable = d.Get("unschedulable").(bool)
	}

	if d.HasChange("instance_types") {
		update = true
		args.ScalingGroup.InstanceTypes = expandStringList(d.Get("instance_types").([]interface{}))
	}

	// password is required by update method
	args.ScalingGroup.LoginPassword = d.Get("password").(string)
	if d.HasChange("password") {
		update = true
		args.ScalingGroup.LoginPassword = d.Get("password").(string)
	}

	if d.HasChange("system_disk_category") {
		update = true
		args.ScalingGroup.SystemDiskCategory = d.Get("system_disk_category").(string)
	}

	if d.HasChange("system_disk_size") {
		update = true
		args.ScalingGroup.SystemDiskSize = int64(d.Get("system_disk_size").(int))
	}

	if d.HasChange("image_id") {
		update = true
		args.ScalingGroup.ImageId = d.Get("image_id").(string)
	}

	if d.HasChange("data_disks") {
		update = true
		setNodePoolDataDisks(&args.ScalingGroup, d)
	}

	if d.HasChange("tags") {
		update = true
		setNodePoolTags(&args.ScalingGroup, d)
	}

	if d.HasChange("labels") {
		update = true
		setNodePoolLabels(&args.KubernetesConfig, d)
	}

	if d.HasChange("taints") {
		update = true
		setNodePoolTaints(&args.KubernetesConfig, d)
	}

	if d.HasChange("user_data") {
		update = true
		if v := d.Get("user_data").(string); v != "" {
			_, base64DecodeError := base64.StdEncoding.DecodeString(v)
			if base64DecodeError == nil {
				args.KubernetesConfig.UserData = v
			} else {
				args.KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v))
			}
		}
	}

	if d.HasChange("scaling_config") {
		update = true
		if v, ok := d.GetOk("scaling_config"); ok {
			args.AutoScaling = setAutoScalingConfig(v.([]interface{}))
		}
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		update = true
		args.ScalingGroup.InternetChargeType = v.(string)
	}

	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		update = true
		args.ScalingGroup.InternetMaxBandwidthOut = v.(int)
	}

	if v, ok := d.GetOk("platform"); ok {
		update = true
		args.ScalingGroup.Platform = v.(string)
	}

	if d.HasChange("scaling_policy") {
		update = true
		args.ScalingGroup.ScalingPolicy = d.Get("scaling_policy").(string)
	}

	if update {
		//begin
		request := client.NewCommonRequest("POST", "CS", "2015-12-15", "ModifyClusterNodePool", fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, d.Id()))
		request.QueryParams["ClusterId"] = clusterId
		request.QueryParams["SignatureVersion"] = "1.0"
		request.Headers["x-acs-asapi-gateway-version"] = "3.0"
		jsonData, err := json.Marshal(args)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
		}
		request.SetContentType(requests.Json)
		request.SetContent(jsonData)
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			if response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateKubernetesNodePool", response, errmsg)
		}

		stateConf := BuildStateConf([]string{"scaling", "updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), clusterId, []string{"deleting", "failed"}))

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	// attach or remove existing node
	if d.HasChange("instances") {
		rawOldValue, rawNewValue := d.GetChange("instances")
		oldValue, ok := rawOldValue.([]interface{})
		if !ok {
			return errmsgs.WrapErrorf(fmt.Errorf("instances old value can not be parsed"), "parseError %d", oldValue)
		}
		newValue, ok := rawNewValue.([]interface{})
		if !ok {
			return errmsgs.WrapErrorf(fmt.Errorf("instances new value can not be parsed"), "parseError %d", oldValue)
		}

		if len(newValue) > len(oldValue) {
			attachExistingInstance(d, meta)
		} else {
			err := RemoveNodePoolNodes(d, meta, clusterId, d.Id(), oldValue, newValue)
			if err != nil {
				return err
			}
		}
	}

	update = false
	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackCSNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Get("cluster_id").(string)
	csService := CsService{client}

	object, err := csService.DescribeCsKubernetesNodePool(d.Id(), clusterId)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("node_count", object.Status.TotalNodes)
	d.Set("name", object.NodepoolInfo.Name)
	//d.Set("vpc_id", object.VpcId)
	d.Set("vswitch_ids", object.ScalingGroup.VswitchIds)
	d.Set("instance_types", object.ScalingGroup.InstanceTypes)
	d.Set("key_name", object.ScalingGroup.KeyPair)
	d.Set("security_group_id", object.ScalingGroup.SecurityGroupID)
	d.Set("system_disk_category", object.ScalingGroup.SystemDiskCategory)
	d.Set("system_disk_size", object.ScalingGroup.SystemDiskSize)

	d.Set("image_id", object.ScalingGroup.ImageID)
	d.Set("platform", object.ScalingGroup.Platform)
	d.Set("scaling_policy", object.ScalingGroup.ScalingPolicy)
	d.Set("node_name_mode", object.KubernetesConfig.NodeNameMode)
	d.Set("user_data", object.KubernetesConfig.UserData)
	d.Set("scaling_group_id", object.ScalingGroup.ScalingGroupID)
	d.Set("unschedulable", object.KubernetesConfig.Unschedulable)
	d.Set("instance_charge_type", object.ScalingGroup.InstanceChargeType)
	d.Set("resource_group_id", object.NodepoolInfo.ResourceGroupID)
	d.Set("spot_strategy", object.ScalingGroup.SpotStrategy)
	//d.Set("internet_charge_type", object.ScalingGroup.InternetChargeType)
	//d.Set("internet_max_bandwidth_out", object.ScalingGroup.InternetMaxBandwidthOut)
	d.Set("install_cloud_monitor", object.KubernetesConfig.CmsEnabled)
	if object.ScalingGroup.InstanceChargeType == "PrePaid" {
		d.Set("period", object.ScalingGroup.Period)
		d.Set("period_unit", object.ScalingGroup.PeriodUnit)
		d.Set("auto_renew", object.ScalingGroup.AutoRenew)
		d.Set("auto_renew_period", object.ScalingGroup.AutoRenewPeriod)
	}

	if passwd, ok := d.GetOk("password"); ok && passwd.(string) != "" {
		d.Set("password", passwd)
	}

	// if parts, err := ParseResourceId(d.Id(), 2); err != nil {
	// 	return errmsgs.WrapError(err)
	// } else {
	// 	d.Set("cluster_id", string(parts[0]))
	// }

	if err := d.Set("data_disks", flattenNodeDataDisksConfig(object.ScalingGroup.DataDisks)); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("taints", flattenTaintsConfig(object.KubernetesConfig.Taints)); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("labels", flattenLabelsConfig(object.KubernetesConfig.Labels)); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("tags", flattenTagsConfig(object.ScalingGroup.Tags)); err != nil {
		return errmsgs.WrapError(err)
	}

	// if object.Management.Enable == true {
	// 	if err := d.Set("management", flattenManagementNodepoolConfig(&object.Management)); err != nil {
	// 		return errmsgs.WrapError(err)
	// 	}
	// }

	if object.AutoScaling.Enable == true {
		if err := d.Set("scaling_config", flattenAutoScalingConfig(&object.AutoScaling)); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if err := d.Set("spot_price_limit", flattenSpotPriceLimit(object.ScalingGroup.SpotPriceLimit)); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackCSNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}
	clusterId := d.Get("cluster_id").(string)
	var raw interface{}
	// delete all nodes
	err := RemoveNodePoolNodes(d, meta, clusterId, d.Id(), nil, nil)
	if err != nil {
		return err
	}

	req := client.NewCommonRequest("DELETE", "CS", "2015-12-15", "DeleteClusterNodepool", fmt.Sprintf("/clusters/%s/nodepools/%s", clusterId, d.Id()))
	req.QueryParams["ClusterId"] = clusterId
	req.QueryParams["NodepoolId"] = d.Id()
	req.Headers["x-acs-asapi-gateway-version"] = "3.0"

	response, err := client.ProcessCommonRequest(req)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteClusterNodePool", raw, errmsg)
	}

	stateConf := BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), clusterId, []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func buildNodePoolArgs(d *schema.ResourceData, meta interface{}) (*requests.CommonRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	password := d.Get("password").(string)
	if password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt2(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return nil, errmsgs.WrapError(err)
			}
			password = decryptResp
		} else if v := d.Get("key_name").(string); v == "" {
				return nil, errmsgs.WrapError(fmt.Errorf("password is require while kms_encrypted_password or key_name not set"))
		}
	}
	request := client.NewCommonRequest("POST", "CS", "2015-12-15", "CreateClusterNodePool", fmt.Sprintf("/clusters/%s/nodepools", d.Get("cluster_id").(string)))

	NodePoolInfo := NodePoolInfo{
		Name:         d.Get("name").(string),
		NodePoolType: "ess", // hard code the type
	}
	ScalingGroup := scalingGroup{

		VswitchIds:    expandStringList(d.Get("vswitch_ids").([]interface{})),
		InstanceTypes: expandStringList(d.Get("instance_types").([]interface{})),
		LoginPassword: password,

		SystemDiskCategory: d.Get("system_disk_category").(string),
		SystemDiskSize:     int64(d.Get("system_disk_size").(int)),

		ImageId:  d.Get("image_id").(string),
		Platform: d.Get("platform").(string),
	}
	KubernetesConfig := kubernetesConfig{}
	AutoScaling := autoScaling{}

	setNodePoolDataDisks(&ScalingGroup, d)
	setNodePoolTags(&ScalingGroup, d)
	setNodePoolTaints(&KubernetesConfig, d)
	setNodePoolLabels(&KubernetesConfig, d)

	if v, ok := d.GetOk("instance_charge_type"); ok {
		ScalingGroup.InstanceChargeType = v.(string)

	}

	if v, ok := d.GetOk("password"); ok {
		ScalingGroup.LoginPassword = v.(string)

	}
	if v, ok := d.GetOk("install_cloud_monitor"); ok {
		KubernetesConfig.CmsEnabled = v.(bool)
	}

	if v, ok := d.GetOk("unschedulable"); ok {
		KubernetesConfig.Unschedulable = v.(bool)
	}

	if v, ok := d.GetOk("user_data"); ok && v != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			KubernetesConfig.UserData = v.(string)
		} else {
			KubernetesConfig.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	// set auto scaling config
	if v, ok := d.GetOk("scaling_policy"); ok {
		ScalingGroup.ScalingPolicy = v.(string)
	}

	if v, ok := d.GetOk("scaling_config"); ok {
		if sc, ok := v.([]interface{}); len(sc) > 0 && ok {
			AutoScaling = setAutoScalingConfig(sc)
		}
	}

	// set manage nodepool params

	// if v, ok := d.GetOk("resource_group_id"); ok {
	// 	ScalingGroup.ResourceGroupId = v.(string)
	// }

	// setting spot instance

	if v, ok := d.GetOk("internet_charge_type"); ok {
		ScalingGroup.InternetChargeType = v.(string)
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		ScalingGroup.InternetMaxBandwidthOut = v.(int)
	}
	request.QueryParams["ClusterId"] = d.Get("cluster_id").(string)
	request.QueryParams["Password"] = d.Get("password").(string)
	request.QueryParams["SignatureVersion"] = "1.0"
	body := CreateClusterNodePoolRequest{
		Count:            int64(d.Get("node_count").(int)),
		NodePoolInfo:     NodePoolInfo,
		ScalingGroup:     ScalingGroup,
		KubernetesConfig: KubernetesConfig,
		AutoScaling:      AutoScaling,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	return request, nil
}

func ConvertCsTags(d *schema.ResourceData) ([]cs.Tag, error) {
	tags := make([]cs.Tag, 0)
	tagsMap, ok := d.Get("tags").(map[string]interface{})
	if ok {
		for key, value := range tagsMap {
			if value != nil {
				if v, ok := value.(string); ok {
					tags = append(tags, cs.Tag{
						Key:   key,
						Value: v,
					})
				}
			}
		}
	}

	return tags, nil
}

func setNodePoolTags(scalingGroup *scalingGroup, d *schema.ResourceData) error {
	if _, ok := d.GetOk("tags"); ok {
		if tags, err := ConvertCsTags(d); err == nil {
			scalingGroup.Tags = tags
		}
	}

	return nil
}

func setNodePoolLabels(config *kubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("labels"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		labels := make([]cs.Label, 0)
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				labels = append(labels, cs.Label{
					Key:   m["key"].(string),
					Value: m["value"].(string),
				})
			}

		}
		config.Labels = labels
	}

	return nil
}

func setNodePoolDataDisks(scalingGroup *scalingGroup, d *schema.ResourceData) error {
	if dds, ok := d.GetOk("data_disks"); ok {
		disks := dds.([]interface{})
		createDataDisks := make([]nodePoolDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := nodePoolDataDisk{
				Size: pack["size"].(int),

				Category: pack["category"].(string),

				Encrypted: pack["encrypted"].(string),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		scalingGroup.DataDisks = createDataDisks
	}

	return nil
}

func setNodePoolTaints(config *kubernetesConfig, d *schema.ResourceData) error {
	if v, ok := d.GetOk("taints"); ok && len(v.([]interface{})) > 0 {
		vl := v.([]interface{})
		taints := make([]cs.Taint, 0)
		for _, i := range vl {
			if m, ok := i.(map[string]interface{}); ok {
				taints = append(taints, cs.Taint{
					Key:    m["key"].(string),
					Value:  m["value"].(string),
					Effect: cs.Effect(m["effect"].(string)),
				})
			}

		}
		config.Taints = taints
	}

	return nil
}

func setManagedNodepoolConfig(l []interface{}) (config cs.Management) {
	if len(l) == 0 || l[0] == nil {
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "management" is set, we think of it as creating a managed node pool
	config.Enable = true

	if v, ok := m["auto_repair"].(bool); ok {
		config.AutoRepair = v
	}
	if v, ok := m["auto_upgrade"].(bool); ok {
		config.UpgradeConf.AutoUpgrade = v
	}
	if v, ok := m["surge"].(int); ok {
		config.UpgradeConf.Surge = int64(v)
	}
	if v, ok := m["surge_percentage"].(int); ok {
		config.UpgradeConf.SurgePercentage = int64(v)
	}
	if v, ok := m["max_unavailable"].(int); ok {
		config.UpgradeConf.MaxUnavailable = int64(v)
	}

	return config
}

func setAutoScalingConfig(l []interface{}) (config autoScaling) {
	if len(l) == 0 || l[0] == nil {
		return config
	}

	m := l[0].(map[string]interface{})

	// Once "scaling_config" is set, we think of it as creating a auto scaling node pool
	config.Enable = true

	if v, ok := m["min_size"].(int); ok {
		config.MinInstances = int64(v)
	}
	if v, ok := m["max_size"].(int); ok {
		config.MaxInstances = int64(v)
	}
	if v, ok := m["type"].(string); ok {
		config.Type = v
	}

	return config
}

func setSpotPriceLimit(l []interface{}) (config []cs.SpotPrice) {
	if len(l) == 0 || l[0] == nil {
		return config
	}
	for _, v := range l {
		if m, ok := v.(map[string]interface{}); ok {
			config = append(config, cs.SpotPrice{
				InstanceType: m["instance_type"].(string),
				PriceLimit:   m["price_limit"].(string),
			})
		}
	}

	return
}

func flattenSpotPriceLimit(config []cs.SpotPrice) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, spotInfo := range config {
		m = append(m, map[string]interface{}{
			"instance_type": spotInfo.InstanceType,
			"price_limit":   spotInfo.PriceLimit,
		})
	}

	return m
}

func flattenAutoScalingConfig(config *cs.AutoScaling) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"min_size":                 config.MinInstances,
		"max_size":                 config.MaxInstances,
		"type":                     config.Type,
		"is_bond_eip":              config.IsBindEip,
		"eip_internet_charge_type": config.EipInternetChargeType,
		"eip_bandwidth":            config.EipBandWidth,
	})

	return
}

func flattenManagementNodepoolConfig(config *cs.Management) (m []map[string]interface{}) {
	if config == nil {
		return
	}
	m = append(m, map[string]interface{}{
		"auto_repair":      config.AutoRepair,
		"auto_upgrade":     config.UpgradeConf.AutoUpgrade,
		"surge":            config.UpgradeConf.Surge,
		"surge_percentage": config.UpgradeConf.SurgePercentage,
		"max_unavailable":  config.UpgradeConf.MaxUnavailable,
	})

	return
}

func flattenNodeDataDisksConfig(config []cs.NodePoolDataDisk) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, disks := range config {
		m = append(m, map[string]interface{}{
			"size":              disks.Size,
			"category":          disks.Category,
			"encrypted":         disks.Encrypted,
			"performance_level": disks.PerformanceLevel,
		})
	}

	return m
}

func flattenTaintsConfig(config []cs.Taint) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, taint := range config {
		m = append(m, map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		})
	}

	return m
}

func flattenLabelsConfig(config []cs.Label) (m []map[string]interface{}) {
	if config == nil {
		return []map[string]interface{}{}
	}

	for _, label := range config {
		m = append(m, map[string]interface{}{
			"key":   label.Key,
			"value": label.Value,
		})
	}

	return m
}

func flattenTagsConfig(config []cs.Tag) map[string]string {
	m := make(map[string]string, len(config))
	if len(config) < 0 {
		return m
	}

	for _, tag := range config {
		if tag.Key != DefaultClusterTag {
			m[tag.Key] = tag.Value
		}
	}

	return m
}
func RemoveNodePoolNodes(d *schema.ResourceData, meta interface{}, clusterid, nodepoolid string, oldNodes []interface{}, newNodes []interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	// list all nodes of the nodepool
	object, err := csService.DescribeClusterNodes(clusterid, d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	// fetch the NodeName of all nodes
	var allNodeName []string
	for _, value := range object.Nodes {
		allNodeName = append(allNodeName, value.NodeName)
	}

	removeNodesName := allNodeName

	// remove automatically created nodes
	if d.HasChange("node_count") {
		o, n := d.GetChange("node_count")
		count := o.(int) - n.(int)
		removeNodesName = allNodeName[:count]
	}

	// remove manually added nodes
	if d.HasChange("instances") {
		var removeInstanceList []string
		var attachNodeList []string
		if oldNodes != nil && newNodes != nil {
			attachNodeList = difference(expandStringList(oldNodes), expandStringList(newNodes))
		}
		if len(newNodes) == 0 {
			attachNodeList = expandStringList(oldNodes)
		}
		for _, v := range object.Nodes {
			for _, name := range attachNodeList {
				if name == v.InstanceID {
					removeInstanceList = append(removeInstanceList, v.NodeName)
				}
			}
		}
		removeNodesName = removeInstanceList
	}
	if len(removeNodesName) > 0 {
		req := csService.client.NewCommonRequest("POST", "CS", "2015-12-15", "RemoveClusterNodes", fmt.Sprintf("/api/v2/clusters/%s/nodes/remove", clusterid))
		req.QueryParams["SignatureVersion"] = "1.0"
		req.Headers["x-acs-asapi-gateway-version"] = "3.0"
		body := map[string]interface{}{
			"release_node": true,
			"drain_node":   true,
			"nodes":        removeNodesName,
			"ClusterId":    clusterid,
		}
		jsonData, err := json.Marshal(body)
		if err != nil {
			return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
		}
		req.SetContentType(requests.Json)
		req.SetContent(jsonData)
		resp, err := csService.client.ProcessCommonRequest(req)
		if err != nil {
			if resp == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteKubernetesClusterNodes", errmsgs.DenverdinoAliyungo)
		}
		stateConf := BuildStateConf([]string{"removing"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), clusterid, []string{"deleting", "failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}
func ScaleClusterNodePool(d *schema.ResourceData, meta interface{}, clusterid, nodepoolid string, oldValue, newValue int) error {
	var raw interface{}
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	// list all nodes of the nodepool
	req := csService.client.NewCommonRequest("POST", "CS", "2015-12-15", "ScaleClusterNodePool", fmt.Sprintf("/clusters/%s/nodepools/%s)", clusterid, nodepoolid))
	req.QueryParams["SignatureVersion"] = "1.0"
	body := map[string]interface{}{
		"ClusterId":  clusterid,
		"NodepoolId": nodepoolid,
		"count":      int64(newValue) - int64(oldValue),
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	req.SetContentType(requests.Json)
	req.SetContent(jsonData)

	req.Headers["x-acs-asapi-gateway-version"] = "3.0"
	response, err := csService.client.ProcessCommonRequest(req)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ScaleClusterNodePool", raw)
	}

	stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), clusterid, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
func attachExistingInstance(d *schema.ResourceData, meta interface{}) error {
	csService := CsService{meta.(*connectivity.AlibabacloudStackClient)}
	client, err := meta.(*connectivity.AlibabacloudStackClient).NewRoaCsClient()
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, ResourceName, "InitializeClient", err)
	}
	clusterId := d.Get("cluster_id").(string)

	args := &roacs.AttachInstancesRequest{
		NodepoolId:       tea.String(d.Id()),
		FormatDisk:       tea.Bool(false),
		KeepInstanceName: tea.Bool(true),
	}

	if v, ok := d.GetOk("password"); ok {
		args.Password = tea.String(v.(string))
	}

	if v, ok := d.GetOk("key_name"); ok {
		args.KeyPair = tea.String(v.(string))
	}

	if v, ok := d.GetOk("format_disk"); ok {
		args.FormatDisk = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("keep_instance_name"); ok {
		args.KeepInstanceName = tea.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_id"); ok {
		args.ImageId = tea.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		args.Instances = tea.StringSlice(expandStringList(v.([]interface{})))
	}

	_, err = client.AttachInstances(tea.String(clusterId), args)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, ResourceName, "AttachInstances", errmsgs.AliyunTablestoreGoSdk)
	}

	stateConf := BuildStateConf([]string{"scaling"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, csService.CsKubernetesNodePoolStateRefreshFunc(d.Id(), clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
