package alibabacloudstack

import (
	"fmt"
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
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeInt,
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
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "The field `network_type` no longer requires manual input; it will be automatically populated based on whether the `vswitch_id` field is configured. This field will be deprecated in version 3.21.0.",
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Classic", "StorageReserver"}, false),
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
	d.Set("cpu_type", instance.CpuType)
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
	var instnaceClass string
	if instance.DBInstanceClass != "" {
		instnaceClass = instance.DBInstanceClass
	} else if instance.InstanceSpec != "" {
		instnaceClass = instance.InstanceSpec
	} else if instance.CpuCores != 0 && instance.MemorySize != 0 {
		instnaceClass = fmt.Sprintf("%dC%dG", instance.CpuCores, instance.MemorySize)
	}
	connectivity.SetResourceData(d, instnaceClass, "db_instance_class", "instance_class")

	d.Set("seg_node_num", instance.SegNodeNum)

	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "description")
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

	reqQuery, err := buildGpdbCreateRequest(d, meta)
	if err != nil {
		return err
	}

	response := map[string]interface{}{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "gpdb", "2016-05-03", "CreateDBInstance", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_instance", "CreateDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(response["DBInstanceId"].(string))

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

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

func buildGpdbCreateRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{}
	if v, ok := d.GetOk("availability_zone"); ok {
		reqQuery["ZoneId"] = v
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		reqQuery["CpuType"] = v
	}
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok {
		reqQuery["PayType"] = v
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		reqQuery["VSwitchId"] = v

		vpcService := VpcService{client}
		object, err := vpcService.DescribeVSwitch(v.(string))
		if err != nil {
			return reqQuery, errmsgs.WrapError(err)
		}

		if zoneId, existed := reqQuery["ZoneId"]; !existed || zoneId.(string) == "" {
			reqQuery["ZoneId"] = object.ZoneId
		} else if strings.Contains(zoneId.(string), MULTI_IZ_SYMBOL) {
			zoneStr := strings.Split(strings.SplitAfter(zoneId.(string), "(")[1], ")")[0]
			if !strings.Contains(zoneStr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return reqQuery, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", object.VSwitchId, zoneId))
			}
		} else if zoneId.(string) != object.ZoneId {
			return reqQuery, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the zone %s.", object.VSwitchId, zoneId))
		}

		reqQuery["VPCId"] = object.VpcId
		reqQuery["InstanceNetworkType"] = "VPC"
	} else {
		reqQuery["InstanceNetworkType"] = "Classic"
	}

	if v, ok := connectivity.GetResourceDataOk(d, "db_instance_description", "description"); ok {
		reqQuery["DBInstanceDescription"] = v
	}

	var dbMode string
	if v, ok := d.GetOk("db_instance_mode"); ok {
		reqQuery["DBInstanceMode"] = v
		dbMode = v.(string)
	}

	if dbMode == "StorageReserver" {
		if v, ok := d.GetOk("db_instance_storage_type"); ok && Trim(v.(string)) != "" {
			reqQuery["StorageType"] = Trim(v.(string))
		} else {

			return reqQuery, errmsgs.WrapError(errmsgs.Error("storage_type is required when db_instance_mode is StorageReserver"))
		}
		if v, ok := d.GetOk("seg_node_num"); ok && v.(int) > 0 {
			reqQuery["SegNodeNum"] = v.(int)
		} else {

			return reqQuery, errmsgs.WrapError(errmsgs.Error("seg_node_num is required when db_instance_mode is StorageReserver"))
		}
	}

	if v, ok := connectivity.GetResourceDataOk(d, "db_instance_class", "instance_class"); ok && Trim(v.(string)) != "" {
		reqQuery["InstanceSpec"] = Trim(v.(string))
		reqQuery["DBInstanceClass"] = Trim(v.(string))
	} else {
		return reqQuery, errmsgs.WrapError(errmsgs.Error("db_instance_class or instance_class is necessory"))
	}

	if v, ok := d.GetOk("instance_group_count"); ok {
		reqQuery["DBInstanceGroupCount"] = v
	}
	if v, ok := d.GetOk("engine"); ok {
		reqQuery["Engine"] = v
	}
	if v, ok := d.GetOk("engine_version"); ok {
		reqQuery["EngineVersion"] = v
	}

	if v, ok := d.GetOk("security_ip_list"); ok && len(v.(*schema.Set).List()) > 0 {

		reqQuery["SecurityIPList"] = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	} else {
		reqQuery["SecurityIPList"] = LOCAL_HOST_IP
	}

	return reqQuery, nil
}
