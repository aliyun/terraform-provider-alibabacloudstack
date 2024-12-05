package alibabacloudstack

import (
	"log"
	"strings"
	"time"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBReadonlyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDBReadonlyInstanceCreate,
		Read:   resourceAlibabacloudStackDBReadonlyInstanceRead,
		Update: resourceAlibabacloudStackDBReadonlyInstanceUpdate,
		Delete: resourceAlibabacloudStackDBReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"master_db_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				Deprecated:   "Field 'master_db_instance_id' is deprecated and will be removed in a future release. Please use 'master_instance_id' instead.",
			},
			"master_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Computed:     true,
				Deprecated:   "Field 'instance_name' is deprecated and will be removed in a future release. Please use 'db_instance_description' instead.",
			},
			"db_instance_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Computed:     true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				Deprecated:   "Field 'instance_type' is deprecated and will be removed in a future release. Please use 'db_instance_class' instead.",
			},
			"db_instance_class": {
				Type:         schema.TypeString,
				Required:     true,
			},

			"instance_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				Deprecated:   "Field 'instance_storage' is deprecated and will be removed in a future release. Please use 'db_instance_storage' instead.",
			},
			"db_instance_storage": {
				Type:         schema.TypeInt,
				Required:     true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"db_instance_storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3", "cloud_pperf", "cloud_sperf"}, false),
			},

			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},

			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackDBReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	request, err := buildDBReadonlyCreateRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	log.Print("wait for instance to be ready")

	if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}
	log.Print("instance is ready")
	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.CreateReadOnlyDBInstance(request)
	})

	bresponse, ok := raw.(*rds.CreateReadOnlyDBInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp := bresponse
	d.SetId(resp.DBInstanceId)

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackDBReadonlyInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackDBReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	d.Partial(true)

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackDBInstanceRead(d, meta)
	}

	if d.HasChange("db_instance_description") || d.HasChange("instance_name") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "db_instance_description", "instance_name"); err == nil {
			request.DBInstanceDescription = v.(string)
		} else {
			return err
		}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceDescription(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
					return resource.RetryableError(err)
				}
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*rds.ModifyDBInstanceDescriptionResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

				return resource.NonRetryableError(err)
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			//d.SetPartial("instance_name")
			return nil
		})

		if err != nil {
			return err
		}

	}

	update := false
	request := rds.CreateModifyDBInstanceSpecRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.PayType = string(Postpaid)

	if d.HasChange("db_instance_storage_type") {
		request.DBInstanceStorageType = d.Get("db_instance_storage_type").(string)
		update = true
	}
	if d.HasChange("db_instance_class") || d.HasChange("instance_type") {
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "db_instance_class", "instance_type"); err == nil {
			request.DBInstanceClass = v.(string)
		} else {
			return err
		}
		update = true
	}

	if d.HasChange("db_instance_storage") || d.HasChange("instance_storage") {
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(0), "db_instance_storage", "instance_storage"); err == nil {
			request.DBInstanceStorage = requests.NewInteger(v.(int))
		} else {
			return err
		}
		update = true
	}

	if update {
		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		_, err := stateConf.WaitForState()
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
					return resource.RetryableError(err)
				}
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*rds.ModifyDBInstanceSpecResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			//d.SetPartial("instance_type")
			//d.SetPartial("instance_storage")
			//d.SetPartial("db_instance_storage_type")
			return nil
		})

		if err != nil {
			return err
		}

		// wait instance status is running after modifying
		_, err = stateConf.WaitForState()
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAlibabacloudStackDBReadonlyInstanceRead(d, meta)
}

func resourceAlibabacloudStackDBReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("engine", instance.Engine)
	connectivity.SetResourceData(d, instance.MasterInstanceId, "master_instance_id", "master_db_instance_id")
	d.Set("engine_version", instance.EngineVersion)
	connectivity.SetResourceData(d, instance.DBInstanceClass, "db_instance_class", "instance_type")
	d.Set("port", instance.Port)
	connectivity.SetResourceData(d, instance.DBInstanceStorage, "db_instance_storage", "instance_storage")
	d.Set("zone_id", instance.ZoneId)
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("db_instance_storage_type", instance.DBInstanceStorageType)

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return err
	}

	tags, err := rdsService.describeTags(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", rdsService.tagsToMap(tags))
	}

	return nil
}

func resourceAlibabacloudStackDBReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if PayType(instance.PayType) == Prepaid {
		return errmsgs.WrapError(errmsgs.Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := rds.CreateDeleteDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"RwSplitNetType.Exist", "OperationDenied.DBInstanceStatus", "OperationDenied.MasterDBInstanceState"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.DeleteDBInstanceResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return err
	}
	waitSecondsIfWithTest(600)
	return nil
}

func buildDBReadonlyCreateRequest(d *schema.ResourceData, meta interface{}) (*rds.CreateReadOnlyDBInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := rds.CreateCreateReadOnlyDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "master_instance_id", "master_db_instance_id"); err == nil {
		request.DBInstanceId = Trim(v.(string))
	} else {
		return nil, err
	}
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(0), "db_instance_storage", "instance_storage"); err == nil {
		request.DBInstanceStorage = requests.NewInteger(v.(int))
	} else {
		return nil, err
	}
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "db_instance_class", "instance_type"); err == nil {
		request.DBInstanceClass = Trim(v.(string))
	} else {
		return nil, err
	}
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "db_instance_description", "instance_name"); err == nil {
		request.DBInstanceDescription = v.(string)
	} else {
		return nil, err
	}
	request.DBInstanceStorageType = d.Get("db_instance_storage_type").(string)

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchId := Trim(d.Get("vswitch_id").(string))

	request.InstanceNetworkType = string(Classic)

	if vswitchId != "" {
		request.VSwitchId = vswitchId
		request.InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
		}

		request.VPCId = vsw.VpcId
	}

	request.PayType = string(Postpaid)
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
