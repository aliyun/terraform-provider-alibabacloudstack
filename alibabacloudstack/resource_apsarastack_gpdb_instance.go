package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackGpdbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAlibabacloudStackGpdbInstanceRead,
		Create: resourceAlibabacloudStackGpdbInstanceCreate,
		Update: resourceAlibabacloudStackGpdbInstanceUpdate,
		Delete: resourceAlibabacloudStackGpdbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Type:     schema.TypeString,
				Required: true,
				Deprecated:   "Field 'instance_class' is deprecated and will be removed in a future release. Please use new field 'db_instance_class' instead.",
				ConflictsWith: []string{"db_instance_class"},
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				ConflictsWith: []string{"instance_class"},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
				Deprecated:   "Field 'instance_id' is deprecated and will be removed in a future release. Please use new field 'instance_id' instead.",
				ConflictsWith: []string{"db_instance_id"},
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
				ConflictsWith: []string{"instance_id"},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_network_type": {
				Type:     schema.TypeString,
				Computed: true,
				Deprecated:   "Field 'instance_network_type' is deprecated and will be removed in a future release. Please use new field 'network_type' instead.",
				ConflictsWith: []string{"network_type"},
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
				ConflictsWith: []string{"instance_network_type"},
			},
			"instance_group_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid"}, false),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Deprecated:   "Field 'instance_charge_type' is deprecated and will be removed in a future release. Please use new field 'payment_type' instead.",
				ConflictsWith: []string{"payment_type"},
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid"}, false),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ConflictsWith: []string{"instance_charge_type"},
			},
			"description": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Optional:     true,
				Deprecated:   "Field 'description' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Optional:     true,
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
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Deprecated:   "Field 'instance_inner_port' is deprecated and will be removed in a future release. Please use new field 'port' instead.",
				ConflictsWith: []string{"port"},
			},
			"port": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"instance_inner_port"},
			},
			"instance_vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Deprecated:   "Field 'instance_vpc_id' is deprecated and will be removed in a future release. Please use new field 'vpc_id' instead.",
				ConflictsWith: []string{"vpc_id"},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"instance_vpc_id"},
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
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackGpdbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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
	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "description")
	connectivity.SetResourceData(d, instance.DBInstanceClass, "db_instance_class", "instance_class")
	connectivity.SetResourceData(d, instance.InstanceNetworkType, "network_type", "instance_network_type")
	d.Set("instance_group_count", instance.DBInstanceGroupCount)
	security_ips, err := gpdbService.DescribeGpdbSecurityIps(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("security_ip_list", security_ips)
	//d.Set("create_time", instance.CreationTime)
	connectivity.SetResourceData(d, instance.PayType, "payment_type", "instance_charge_type")
	d.Set("tags", gpdbService.tagsToMap(instance.Tags.Tag))
	d.Set("instance_inner_connection", instance.ConnectionString)
	connectivity.SetResourceData(d, instance.Port, "port", "instance_inner_port")
	connectivity.SetResourceData(d, instance.VpcId, "vpc_id", "instance_vpc_id")
	return nil
}

func resourceAlibabacloudStackGpdbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	request, err := buildGpdbCreateRequest(d, meta)
	client.InitRpcRequest(*request.RpcRequest)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.CreateDBInstance(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"SYSTEM.CONCURRENT_OPERATE"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, ok := raw.(*gpdb.CreateDBInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_gpdb_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(response.DBInstanceId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Minute, gpdbService.GpdbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return resourceAlibabacloudStackGpdbInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackGpdbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	// Begin Update
	d.Partial(true)

	// Update Instance Description
	if d.HasChanges("db_instance_description", "description"){
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
		//d.SetPartial("security_ip_list")
	}

	if err := gpdbService.setInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}

	// Finish Update
	d.Partial(false)

	return resourceAlibabacloudStackGpdbInstanceRead(d, meta)
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
	request.DBInstanceClass = Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_class").(string))
	request.DBInstanceGroupCount = Trim(d.Get("instance_group_count").(string))
	request.Engine = Trim(d.Get("engine").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))

	// Instance NetWorkType
	request.InstanceNetworkType = string(Classic)
	if request.VSwitchId != "" {
		// check vswitchId in zone
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
