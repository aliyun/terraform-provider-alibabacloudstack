package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbReadonlyInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
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
			"tde_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encrypt_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "aes-256",
				ValidateFunc: validation.StringInSlice([]string{"sm4-128", "aes-256"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("engine").(string) == "MySQL" {
						return true
					}
					if v, ok := d.GetOk("tde_status"); ok && v.(bool) {
						return old == new
					}
					return true
				},
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Optional: true,
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
	setResourceFunc(resource, resourceAlibabacloudStackPolardbReadonlyInstanceCreate, resourceAlibabacloudStackPolardbReadonlyInstanceRead, resourceAlibabacloudStackPolardbReadonlyInstanceUpdate, resourceAlibabacloudStackPolardbReadonlyInstanceDelete)
	return resource
}

func resourceAlibabacloudStackPolardbReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CreateReadOnlyDBInstance", "")
	PolardbCreatereadonlydbinstanceResponse := PolardbCreatereadonlydbinstanceResponse{}
	master_instance_id := connectivity.GetResourceData(d, "master_instance_id", "master_db_instance_id").(string)
	if err := PolardbService.WaitForDBInstance(master_instance_id, Running, DefaultLongTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	request.QueryParams["DBInstanceId"] = master_instance_id
	if err := errmsgs.CheckEmpty(request.QueryParams["DBInstanceId"], schema.TypeString, "master_instance_id", "master_db_instance_id"); err != nil {
		return errmsgs.WrapError(err)
	}
	engine := d.Get("engine").(string)
	request.QueryParams["Engine"] = engine
	request.QueryParams["EngineVersion"] = Trim(d.Get("engine_version").(string))
	// To be tested
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
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
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
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackPolardbReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	if d.HasChange("parameters") {
		if err := PolardbService.ModifyParameters(d, client); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.IsNewResource() {
		return resourceAlibabacloudStackPolardbReadonlyInstanceRead(d, meta)
	}

	if d.HasChanges("db_instance_description", "instance_name") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceDescription", "")
		PolardbModifydbinstancedescriptionResponse := PolardbModifydbinstancedescriptionResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["DBInstanceDescription"] = connectivity.GetResourceData(d, "db_instance_description", "instance_name").(string)

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
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
		request.QueryParams["PayType"] = string(Postpaid)
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
		_, err := stateConf.WaitForState()
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
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
		// wait instance status change from Creating to running
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if d.HasChange("enable_ssl") {
		ssl := d.Get("enable_ssl").(bool)
		if d.IsNewResource() && ssl == false {
			// New resource defaults to false
			return nil
		}
		ssl_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSSL", "")
		ssl_req.QueryParams["DBInstanceId"] = d.Id()
		ssl_req.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		if ssl == true {
			ssl_req.QueryParams["SSLEnabled"] = "1"
		} else {
			ssl_req.QueryParams["SSLEnabled"] = "0"
		}
		bresponse, err := client.ProcessCommonRequest(ssl_req)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if ssl == true {
			log.Print("Updated SSL to true")
		} else {
			log.Print("Updated SSL to false")
		}
	}
	return nil
}

func resourceAlibabacloudStackPolardbReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instance, err := PolardbService.DoPolardbDescribedbinstanceattributeRequest(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	engine := Trim(instance.Items.DBInstanceAttribute[0].Engine)
	d.Set("engine", engine)
	d.Set("engine_version", instance.Items.DBInstanceAttribute[0].EngineVersion)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceClass, "db_instance_class", "instance_type")
	d.Set("port", instance.Items.DBInstanceAttribute[0].Port)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceStorage, "db_instance_storage", "instance_storage")
	d.Set("zone_id", instance.Items.DBInstanceAttribute[0].ZoneId)
	d.Set("vswitch_id", instance.Items.DBInstanceAttribute[0].VSwitchId)
	d.Set("connection_string", instance.Items.DBInstanceAttribute[0].ConnectionString)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("db_instance_storage_type", instance.Items.DBInstanceAttribute[0].DBInstanceStorageType)
	d.Set("master_db_instance_id", instance.Items.DBInstanceAttribute[0].MasterInstanceId)
	if err = PolardbService.RefreshParameters(d, client); err != nil {
		return errmsgs.WrapError(err)
	}
	ssl_object, err := PolardbService.DescribeDBInstanceSSL(d.Id())
	ssl := false
	if (engine == "MySQL" && ssl_object["SSLEnabled"].(string) == "Yes") || (engine != "MySQL" && ssl_object["SSLEnabled"].(string) == "on") {
		ssl = true
	}
	d.Set("enable_ssl", ssl)
	return nil
}

func resourceAlibabacloudStackPolardbReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instance, err := PolardbService.DoPolardbDescribedbinstanceattributeRequest(d.Id())
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
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	return PolardbService.WaitForDBInstance(d.Id(), Deleted, DefaultLongTimeout)
}

type PolardbCreatereadonlydbinstanceResponse struct {
	RequestId        string `json:"RequestId"`
	DBInstanceId     string `json:"DBInstanceId"`
	OrderId          string `json:"OrderId"`
	ConnectionString string `json:"ConnectionString"`
	Port             string `json:"Port"`
}
