package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackGpdbInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_class": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'instance_class' is deprecated and will be removed in a future release. Please use new field 'db_instance_class' instead.",
				ConflictsWith: []string{"db_instance_class"},
			},
			"db_instance_class": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_class"},
			},
			"seg_node_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'instance_id' is deprecated and will be removed in a future release. Please use new field 'instance_id' instead.",
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Classic",
			},
			"instance_group_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid"}, false),
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_charge_type' is deprecated and will be removed in a future release. Please use new field 'payment_type' instead.",
				ConflictsWith: []string{"payment_type"},
			},
			"payment_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid"}, false),
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_charge_type"},
			},
			"description": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'description' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"description"},
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_inner_connection": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_inner_port": {
				Type:       schema.TypeString,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Field 'instance_inner_port' is deprecated and will be removed in a future release. Please use new field 'port' instead.",
			},
			"port": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"instance_vpc_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'instance_vpc_id' is deprecated and will be removed in a future release. Please use new field 'vpc_id' instead.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"gpdb"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"db_instance_storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"db_instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackGpdbInstanceCreate, resourceAlibabacloudStackGpdbInstanceRead, resourceAlibabacloudStackGpdbInstanceUpdate, resourceAlibabacloudStackGpdbInstanceDelete)
	return resource
}

func resourceAlibabacloudStackGpdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	instance, err := gpdbService.DescribeGpdbInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, instance.DBInstanceId, "db_instance_id", "instance_id")
	d.Set("region_id", instance.RegionId)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("status", instance.DBInstanceStatus)
	if instance.StorageType != "" {
		d.Set("db_instance_storage_type", instance.StorageType)
	}
	d.Set("instance_pay_type", instance.PayType)
	if instance.DBInstanceMode != "" {
		d.Set("db_instance_mode", instance.DBInstanceMode)
	}
	if instance.DBInstanceGroupCount != "" {
		d.Set("instance_group_count", instance.DBInstanceGroupCount)
	}
	if instance.DBInstanceClass != "" {
		d.Set("instance_class", instance.DBInstanceClass)
	}

	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "description")
	if instance.InstanceSpec != "" {
		connectivity.SetResourceData(d, instance.InstanceSpec, "db_instance_class", "instance_class")
	}
	connectivity.SetResourceData(d, instance.InstanceNetworkType, "network_type")
	security_ips, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("security_ip_list", security_ips)
	//d.Set("create_time", instance.CreationTime)
	connectivity.SetResourceData(d, instance.PayType, "payment_type", "instance_charge_type")
	d.Set("tags", gpdbService.tagsToMap(instance.Tags.Tag))
	if instance.ConnectionString != "" {
		d.Set("instance_inner_connection", instance.ConnectionString)
	}
	if instance.VpcId != "" {
		connectivity.SetResourceData(d, instance.VpcId, "vpc_id", "instance_vpc_id")
	}
	if instance.Port != "" {
		connectivity.SetResourceData(d, instance.Port, "port", "instance_inner_port")
	}
	return nil
}

func resourceAlibabacloudStackGpdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	request, err := buildGpdbCreateRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client.InitRpcRequest(*request.RpcRequest)

	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.CreateDBInstance(request)
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if raw == nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "API response is nil")
		}

		response, ok := raw.(*gpdb.CreateDBInstanceResponse)
		if !ok || response == nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, "Failed to cast API response")
		}

		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response, ok := raw.(*gpdb.CreateDBInstanceResponse)
	if !ok || response == nil {
		return errmsgs.Error("Failed to cast CreateDBInstance response")
	}

	d.SetId(response.DBInstanceId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackGpdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	// Begin Update
	d.Partial(true)

	// Update Instance Description
	if d.HasChanges("db_instance_description", "description") {
		request := gpdb.CreateModifyDBInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = connectivity.GetResourceData(d, "db_instance_description", "description").(string)
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			response, ok := raw.(*gpdb.ModifyDBInstanceDescriptionResponse)
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("db_instance_description")
	}

	// Update Security Ips
	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipStr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipStr == "" {
			ipStr = LOCAL_HOST_IP
		}
		if err := gpdbService.ModifyGpdbSecurityIps(d.Id(), ipStr); err != nil {
			return errmsgs.WrapError(err)
		}
		// d.SetPartial("security_ip_list")
	}

	if err := gpdbService.setInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}

	// Finish Update
	d.Partial(false)

	return nil
}

func resourceAlibabacloudStackGpdbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := gpdb.CreateDeleteDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.DeleteDBInstance(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			response, ok := raw.(*gpdb.DeleteDBInstanceResponse)
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return err
	}
	// because DeleteDBInstance is called synchronously, there is no wait or describe here.
	return nil
}

func buildGpdbCreateRequest(d *schema.ResourceData, meta interface{}) (*gpdb.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := gpdb.CreateCreateDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ZoneId = Trim(d.Get("availability_zone").(string))
	request.PayType = connectivity.GetResourceData(d, "payment_type", "instance_charge_type").(string)
	request.VSwitchId = Trim(d.Get("vswitch_id").(string))
	request.DBInstanceDescription = connectivity.GetResourceData(d, "db_instance_description", "description").(string)

	dbInstanceMode := Trim(d.Get("db_instance_mode").(string))
	if dbInstanceMode != "" {
		request.DBInstanceMode = dbInstanceMode
	}
	request.InstanceNetworkType = Trim(d.Get("network_type").(string))

	if request.DBInstanceMode == "StorageReserver" {
		storageType := Trim(d.Get("db_instance_storage_type").(string))
		if storageType == "" {
			return nil, errmsgs.WrapError(errmsgs.Error("storage_type is required when db_instance_mode is StorageReserver"))
		}
		request.StorageType = storageType

		segNodeNum := Trim(d.Get("seg_node_num").(string))
		if segNodeNum == "" {
			return nil, errmsgs.WrapError(errmsgs.Error("seg_node_num is required when db_instance_mode is StorageReserver"))
		}
		request.SegNodeNum = segNodeNum

		if dbInstanceClass := Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_class").(string)); dbInstanceClass != "" {
			request.InstanceSpec = dbInstanceClass
		} else {
			return nil, errmsgs.WrapError(errmsgs.Error("db_instance_class is required for StorageReserver mode"))
		}
	} else {
		dbInstanceClass := Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_class").(string))
		if dbInstanceClass == "" {
			return nil, errmsgs.WrapError(errmsgs.Error("db_instance_class is required when db_instance_mode is not StorageReserver"))
		}
		request.InstanceSpec = dbInstanceClass
		request.DBInstanceClass = dbInstanceClass
	}

	if d.Get("instance_group_count").(string) != "" {
		request.DBInstanceGroupCount = Trim(d.Get("instance_group_count").(string))
	}
	request.Engine = Trim(d.Get("engine").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))

	// Instance NetWorkType
	// request.InstanceNetworkType = string(Classic)
	if request.VSwitchId != "" {
		vpcService := VpcService{client}
		object, err := vpcService.DescribeVSwitch(request.VSwitchId)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the zone %s.", object.VSwitchId, request.ZoneId))
		}

		request.VPCId = object.VpcId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))
	}

	// Security Ips
	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	// ClientToken
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
