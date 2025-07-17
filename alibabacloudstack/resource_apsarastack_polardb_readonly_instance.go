package alibabacloudstack

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbReadonlyInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackPolardbReadonlyInstanceCreate,
		Read:   resourceAlibabacloudStackPolardbReadonlyInstanceRead,
		Update: resourceAlibabacloudStackPolardbReadonlyInstanceUpdate,
		Delete: resourceAlibabacloudStackPolardbReadonlyInstanceDelete,
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'master_db_instance_id' is deprecated and will be removed in a future release. Please use new field 'master_instance_id' instead.",
				ConflictsWith: []string{"master_instance_id"},
			},
			"master_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"master_db_instance_id"},
			},

			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Computed:      true,
				Deprecated:    "Field 'instance_name' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Computed:      true,
				ConflictsWith: []string{"instance_name"},
			},

			"instance_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_type' is deprecated and will be removed in a future release. Please use new field 'db_instance_class' instead.",
				ConflictsWith: []string{"db_instance_class"},
			},
			"db_instance_class": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_type"},
			},

			"instance_storage": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_storage' is deprecated and will be removed in a future release. Please use new field 'db_instance_storage' instead.",
				ConflictsWith: []string{"db_instance_storage"},
			},
			"db_instance_storage": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_storage"},
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
		},
	}
}

func resourceAlibabacloudStackPolardbReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CreateReadOnlyDBInstance", "")
	PolardbCreatereadonlydbinstanceResponse := PolardbCreatereadonlydbinstanceResponse{}

	if err := PolardbService.WaitForDBInstance(d, client, Running, DefaultLongTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	request.QueryParams["DBInstanceId"] = connectivity.GetResourceData(d, "master_instance_id", "master_db_instance_id").(string)
	if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceId"], schema.TypeString, "master_instance_id", "master_db_instance_id"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["EngineVersion"] = Trim(d.Get("engine_version").(string))
	//待测
	request.QueryParams["DBInstanceStorage"] = strconv.Itoa(connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int))
	if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceStorage"], schema.TypeInt, "db_instance_storage", "instance_storage"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["DBInstanceClass"] = Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string))
	if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceClass"], schema.TypeString, "db_instance_class", "instance_type"); err != nil {
		return errmsgs.WrapError(err)
	}
	request.QueryParams["DBInstanceDescription"] = connectivity.GetResourceData(d, "db_instance_description", "instance_name").(string)
	request.QueryParams["DBInstanceStorageType"] = d.Get("db_instance_storage_type").(string)
	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.QueryParams["ZoneId"] = Trim(zone.(string))
	}
	vswitchId := Trim(d.Get("vswitch_id").(string))
	vpcService := VpcService{client}
	request.QueryParams["InstanceNetworkType"] = string(Classic)
	if vswitchId != "" {
		request.QueryParams["VSwitchId"] = vswitchId
		request.QueryParams["InstanceNetworkType"] = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if request.QueryParams["ZoneId"] == "" {
			request.QueryParams["ZoneId"] = vsw.ZoneId
		} else if strings.Contains(request.QueryParams["ZoneId"], MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.QueryParams["ZoneId"], "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.QueryParams["ZoneId"]))
			}
		} else if request.QueryParams["ZoneId"] != vsw.ZoneId {
			return errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.QueryParams["ZoneId"]))
		}

		request.QueryParams["VPCId"] = vsw.VpcId
	}
	request.QueryParams["PayType"] = string(Postpaid)
	request.QueryParams["ClientToken"] = buildClientToken(request.GetActionName())

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCreatereadonlydbinstanceResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.SetId(PolardbCreatereadonlydbinstanceResponse.DBInstanceId)

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return resourceAlibabacloudStackPolardbReadonlyInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackPolardbReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	d.Partial(true)
	if d.HasChange("parameters") {
		if err := PolardbService.ModifyParameters(d, client, "parameters"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackPolardbReadonlyInstanceRead(d, meta)
	}

	if d.HasChanges("db_instance_description", "instance_name") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceDescription", "")
		PolardbModifydbinstancedescriptionResponse := PolardbModifydbinstancedescriptionResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["DBInstanceDescription"] = connectivity.GetResourceData(d, "db_instance_description", "instance_name").(string)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceDescription", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancedescriptionResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceDescription", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}

	update := false
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSpec", "")
	PolardbModifydbinstancespecResponse := PolardbModifydbinstancespecResponse{}
	request.QueryParams["DBInstanceId"] = d.Id()
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok {
		request.QueryParams["PayType"] = v.(string)
	} else {
		request.QueryParams["PayType"] = string(Postpaid)
	}

	if d.HasChanges("instance_type", "db_instance_class") {
		request.QueryParams["DBInstanceClass"] = connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string)
		if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceClass"], schema.TypeString, "db_instance_class", "instance_type"); err != nil {
			return errmsgs.WrapError(err)
		}
		update = true
	}

	if d.HasChanges("instance_storage", "db_instance_storage") {
		request.QueryParams["DBInstanceStorage"] = strconv.Itoa(connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int))
		if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceStorage"], schema.TypeString, "db_instance_storage", "instance_storage"); err != nil {
			return errmsgs.WrapError(err)
		}
		update = true
	}
	if update {

		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
		_, err := stateConf.WaitForState()
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceSpec", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancespecResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceSpec", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAlibabacloudStackPolardbReadonlyInstanceRead(d, meta)
}

func resourceAlibabacloudStackPolardbReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instance, err := PolardbService.DoPolardbDescribedbinstanceattributeRequest(d.Id(), client)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("engine", instance.Items.DBInstanceAttribute[0].Engine)
	d.Set("engine_version", instance.Items.DBInstanceAttribute[0].EngineVersion)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceClass, "db_instance_class", "instance_type")
	d.Set("port", instance.Items.DBInstanceAttribute[0].Port)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceStorage, "db_instance_storage", "instance_storage")
	d.Set("zone_id", instance.Items.DBInstanceAttribute[0].ZoneId)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].PayType, "payment_type", "instance_charge_type")
	d.Set("vswitch_id", instance.Items.DBInstanceAttribute[0].VSwitchId)
	d.Set("connection_string", instance.Items.DBInstanceAttribute[0].ConnectionString)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("db_instance_storage_type", instance.Items.DBInstanceAttribute[0].DBInstanceStorageType)

	if err = PolardbService.RefreshParameters(d, client, "parameters"); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackPolardbReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instance, err := PolardbService.DoPolardbDescribedbinstanceattributeRequest(d.Id(), client)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if PayType(instance.Items.DBInstanceAttribute[0].PayType) == Prepaid {
		return errmsgs.WrapError(errmsgs.Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DeleteDBInstance", "")
	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	waitSecondsIfWithTest(600)
	return nil
}

type PolardbCreatereadonlydbinstanceResponse struct {
	RequestId        string `json:"RequestId"`
	DBInstanceId     string `json:"DBInstanceId"`
	OrderId          string `json:"OrderId"`
	ConnectionString string `json:"ConnectionString"`
	Port             string `json:"Port"`
}
