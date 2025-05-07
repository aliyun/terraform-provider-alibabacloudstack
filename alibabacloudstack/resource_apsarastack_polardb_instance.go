package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPolardbInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PolarDB_PPAS", "PolarDB_PG"}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"zone_id_slave1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id_slave2": {
				Type:     schema.TypeString,
				Optional: true,
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
					if v,ok:= d.GetOk("tde_status"); ok && v.(bool) {
						return old == new
					}
					return true
				},

			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"storage_type": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'storage_type' is deprecated and will be removed in a future release. Please use new field 'db_instance_storage_type' instead.",
				ConflictsWith: []string{"db_instance_storage_type"},
			},
			"db_instance_storage_type": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_type"},
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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
			"instance_charge_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_charge_type' is deprecated and will be removed in a future release. Please use new field 'payment_type' instead.",
				ConflictsWith: []string{"payment_type"},
			},
			"payment_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_charge_type"},
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"monitoring_period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{5, 60, 300}),
				Optional:     true,
				Computed:     true,
			},
			"auto_renew": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidAndRenewDiffSuppressFunc,
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
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Deprecated:    "Field 'instance_name' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				ConflictsWith: []string{"instance_name"},
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_ip_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{NormalMode, SafetyMode}, false),
				Optional:     true,
				Default:      NormalMode,
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
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbInstanceCreate, resourceAlibabacloudStackPolardbInstanceRead, resourceAlibabacloudStackPolardbInstanceUpdate, resourceAlibabacloudStackPolardbInstanceDelete)
	return resource
}

func resourceAlibabacloudStackPolardbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	vpcService := VpcService{client}

	var VSwitchId, InstanceNetworkType, ZoneIdSlave1, ZoneIdSlave2, ZoneId, VPCId, arnrole string

	var encryption bool
	EncryptionKey := d.Get("encryption_key").(string)
	encryption = d.Get("encryption").(bool)
	EncryptAlgorithm := d.Get("encrypt_algorithm").(string)
	log.Print("Encryption key input")
	if EncryptionKey != "" && encryption == true {
		log.Print("Encryption key condition passed")
		req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CheckCloudResourceAuthorized", "")
		req.QueryParams["TargetRegionId"] = client.RegionId
		var arnresp RoleARN
		bresponse, err := client.ProcessCommonRequest(req)
		addDebug(req.GetActionName(), bresponse, req, req.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "CheckCloudResourceAuthorized", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &arnresp)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		arnrole = arnresp.RoleArn
		d.Set("role_arn", arnrole)
		log.Printf("check arnrole %v", arnrole)
	} else if EncryptionKey == "" && encryption == true {
		return errmsgs.WrapErrorf(nil, "Add EncryptionKey or Set encryption to false", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
	} else if EncryptionKey != "" && encryption == false {
		return errmsgs.WrapErrorf(nil, "Set encryption to true", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
	} else {
		log.Print("Encryption key condition failed")
	}
	d.Set("encryption", encryption)
	log.Printf("encryptionbool %v", d.Get("encryption").(bool))

	enginever := Trim(d.Get("engine_version").(string))
	engine := Trim(d.Get("engine").(string))
	DBInstanceStorage := connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int)
	if err := errmsgs.CheckEmpty(DBInstanceStorage, schema.TypeString, "db_instance_storage", "instance_storage"); err != nil {
		return errmsgs.WrapError(err)
	}
	DBInstanceClass := Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string))
	if err := errmsgs.CheckEmpty(DBInstanceClass, schema.TypeString, "db_instance_class", "instance_type"); err != nil {
		return errmsgs.WrapError(err)
	}
	DBInstanceNetType := string(Intranet)
	DBInstanceDescription := connectivity.GetResourceData(d, "db_instance_description", "instance_name").(string)
	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		ZoneId = Trim(zone.(string))
	}
	vswitchId := Trim(d.Get("vswitch_id").(string))

	InstanceNetworkType = string(Classic)
	if vswitchId != "" {
		VSwitchId = vswitchId
		InstanceNetworkType = strings.ToUpper(string(Vpc))

		// check vswitchId in zone
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil
		}

		if ZoneId == "" {
			ZoneId = vsw.ZoneId
		}

		VPCId = vsw.VpcId
	}
	payType := string(Postpaid)
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok && Trim(v.(string)) != "" {
		payType = Trim(v.(string))
	}
	DBInstanceStorageType := connectivity.GetResourceData(d, "db_instance_storage_type", "storage_type").(string)
	if err := errmsgs.CheckEmpty(DBInstanceStorageType, schema.TypeString, "db_instance_storage_type", "storage_type"); err != nil {
		return errmsgs.WrapError(err)
	}
	ZoneIdSlave1 = d.Get("zone_id_slave1").(string)
	ZoneIdSlave2 = d.Get("zone_id_slave2").(string)
	SecurityIPList := LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		SecurityIPList = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	ClientToken := fmt.Sprintf("Terraform-AlibabacloudStack-%d-%s", time.Now().Unix(), uuid)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CreateDBInstance", "")
	PolardbCreatedbinstanceResponse := PolardbCreatedbinstanceResponse{}

	mergeMaps(request.QueryParams, map[string]string{
		"EngineVersion":         enginever,
		"Engine":                engine,
		"Encryption":            strconv.FormatBool(encryption),
		"DBInstanceStorage":     strconv.Itoa(DBInstanceStorage),
		"DBInstanceClass":       DBInstanceClass,
		"DBInstanceNetType":     DBInstanceNetType,
		"DBInstanceDescription": DBInstanceDescription,
		"InstanceNetworkType":   InstanceNetworkType,
		"VSwitchId":             VSwitchId,
		"PayType":               payType,
		"DBInstanceStorageType": DBInstanceStorageType,
		"SecurityIPList":        SecurityIPList,
		"ClientToken":           ClientToken,
		"ZoneIdSlave1":          ZoneIdSlave1,
		"ZoneIdSlave2":          ZoneIdSlave2,
		"EncryptionKey":         EncryptionKey,
		"ZoneId":                ZoneId,
		"VPCId":                 VPCId,
		"RoleARN":               arnrole,
	})
	if tde := d.Get("tde_status"); tde == true && engine != "MySQL" {
		request.QueryParams["TdeStatus"] = "1"
		request.QueryParams["EncryptAlgorithm"] = EncryptAlgorithm
		request.QueryParams["RoleARN"] = arnrole
		if EncryptionKey != "" {
			request.QueryParams["EncryptionKey"] = EncryptionKey
		}
	}
	log.Printf("request245 %v", request.QueryParams)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "CreateDBInstance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("CreateDBInstance", bresponse, request, request.QueryParams)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCreatedbinstanceResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "CreateDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(PolardbCreatedbinstanceResponse.DBInstanceId)
	d.Set("connection_string", PolardbCreatedbinstanceResponse.ConnectionString)

	stateConf := BuildStateConfByTimes([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}), 100)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	tde := d.Get("tde_status").(bool)
	log.Printf(" ============================================= tde:%t, engine:%s", tde, engine)
	if tde == true && engine == "MySQL" {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceTDE", "")
		PolardbModifydbinstancetdeResponse := PolardbModifydbinstancetdeResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["TDEStatus"] = "Enabled"
		request.QueryParams["RoleARN"] = arnrole

		if EncryptionKey != "" {
			request.QueryParams["EncryptionKey"] = EncryptionKey
		}
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("ModifyDBInstanceTDE", bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceTDE", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancetdeResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceTDE", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"Disabled"}, []string{"Enabled"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, PolardbService.PolardbDBInstanceTdeStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		log.Print("enabled TDE")
	}
	if ssl := d.Get("enable_ssl"); ssl == true {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSSL", "")
		PolardbModifydbinstancesslResponse := PolardbModifydbinstancesslResponse{}

		//request.QueryParams["Forwardedregionid"] = client.RegionId
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["SSLEnabled"] = "1"
		request.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("ModifyDBInstanceSSL", bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceSSL", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancesslResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceSSL", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		var target, process string
		if engine == "MySQL" {
			target = "Yes"
			process = "No"
		} else {
			target = "on"
			process = "off"
		}
		stateConf := BuildStateConf([]string{process}, []string{target}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, PolardbService.PolardbDBInstanceSslStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		log.Print("enabled SSL")
	}
	return nil
}

func resourceAlibabacloudStackPolardbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))

	if d.HasChange("parameters") {
		if err := PolardbService.ModifyParameters(d, client, "parameters"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	payType := Postpaid
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok && Trim(v.(string)) != "" {
		payType = PayType(Trim(v.(string)))
	}

	if !d.IsNewResource() && d.HasChanges("instance_charge_type", "payment_type") && payType == Prepaid {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstancePayType", "")
		PolardbModifydbinstancepaytypeResponse := PolardbModifydbinstancepaytypeResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["PayType"] = string(payType)
		request.QueryParams["AutoPay"] = "true"
		period := d.Get("period").(int)
		request.QueryParams["UsedTime"] = strconv.Itoa(period)
		request.QueryParams["Period"] = string(Month)
		if period > 9 {
			request.QueryParams["UsedTime"] = strconv.Itoa(period / 12)
			request.QueryParams["Period"] = string(Year)
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstancePayType", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancepaytypeResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstancePayType", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if payType == Prepaid && d.HasChanges("auto_renew", "auto_renew_period") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyInstanceAutoRenewalAttribute", "")
		PolardbModifyinstanceautorenewalattributeResponse := PolardbModifyinstanceautorenewalattributeResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		auto_renew := d.Get("auto_renew").(bool)
		if auto_renew {
			request.QueryParams["AutoRenew"] = "True"
		} else {
			request.QueryParams["AutoRenew"] = "False"
		}
		request.QueryParams["Duration"] = strconv.Itoa(d.Get("auto_renew_period").(int))

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyInstanceAutoRenewalAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifyinstanceautorenewalattributeResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyInstanceAutoRenewalAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}

	if d.HasChange("monitoring_period") {
		period := d.Get("monitoring_period").(int)
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceMonitor", "")
		PolardbModifydbinstancemonitorResponse := PolardbModifydbinstancemonitorResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["Period"] = strconv.Itoa(period)
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceMonitor", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancemonitorResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceMonitor", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}

	if d.HasChange("maintain_time") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceMaintainTime", "")
		PolardbModifydbinstancemaintaintimeResponse := PolardbModifydbinstancemaintaintimeResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["MaintainTime"] = d.Get("maintain_time").(string)
		request.QueryParams["ClientToken"] = buildClientToken(request.GetActionName())

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceMaintainTime", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifydbinstancemaintaintimeResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifyDBInstanceMaintainTime", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	if d.HasChanges("instance_name", "db_instance_description") {
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

	if d.HasChange("security_ips") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifySecurityIps", "")
		PolardbModifysecurityipsResponse := PolardbModifysecurityipsResponse{}

		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}
		if err := PolardbService.ModifyDBSecurityIps(d, client, ipstr); err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["SecurityIps"] = ipstr
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifysecurityipsResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", errmsgs.AlibabacloudStackSdkGoERROR)
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
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		// wait instance status is running before modifying
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
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	d.Partial(false)
	engine := Trim(d.Get("engine").(string))
	if d.HasChange("tde_status") && d.Get("tde_status").(bool) && engine == "MySQL" {
		tde_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceTDE", "")
		tde_req.QueryParams["DBInstanceId"] = d.Id()
		tde_req.QueryParams["TDEStatus"] = "Enabled"

		bresponse, err := client.ProcessCommonRequest(tde_req)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"Disabled"}, []string{"Enabled"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, PolardbService.PolardbDBInstanceTdeStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		log.Print("Updated TDE")
	}

	if d.HasChange("enable_ssl") {
		ssl := d.Get("enable_ssl").(bool)
		ssl_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSSL", "")
		ssl_req.QueryParams["DBInstanceId"] = d.Id()
		ssl_req.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		var target, process string
		engine := Trim(d.Get("engine").(string))
		if ssl == true {
			ssl_req.QueryParams["SSLEnabled"] = "1"
			if engine == "MySQL" {
				target = "Yes"
				process = "No"
			} else {
				target = "on"
				process = "off"
			}

		} else {
			ssl_req.QueryParams["SSLEnabled"] = "0"
			if engine == "MySQL" {
				target = "off"
				process = "on"
			} else {
				target = "No"
				process = "Yes"
			}
		}
		bresponse, err := client.ProcessCommonRequest(ssl_req)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_account", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{process}, []string{target}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, PolardbService.PolardbDBInstanceSslStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		if err := PolardbService.WaitForDBInstance(d, client, Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
		if ssl == true {
			log.Print("Updated SSL to true")
		} else {
			log.Print("Updated SSL to false")
		}
	}
	return nil
}

func resourceAlibabacloudStackPolardbInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	ips, err := PolardbService.GetSecurityIps(d, client)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// 未完成
	// tags, err := rdsService.describeTags(d)
	// if err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	// if len(tags) > 0 {
	// 	d.Set("tags", rdsService.tagsToMap(tags))
	// }

	monitoringPeriod, err := PolardbService.DoPolardbDescribedbinstancemonitorRequest(d, client)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	period, err := strconv.Atoi(monitoringPeriod.Period)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("monitoring_period", period)
	d.Set("security_ips", ips)
	d.Set("security_ip_mode", instance.Items.DBInstanceAttribute[0].SecurityIPMode)
	d.Set("engine", instance.Items.DBInstanceAttribute[0].Engine)
	d.Set("engine_version", instance.Items.DBInstanceAttribute[0].EngineVersion)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceClass, "db_instance_class", "instance_type")
	d.Set("port", instance.Items.DBInstanceAttribute[0].Port)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceStorage, "db_instance_storage", "instance_storage")
	d.Set("zone_id", instance.Items.DBInstanceAttribute[0].ZoneId)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].PayType, "payment_type", "instance_charge_type")
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.Items.DBInstanceAttribute[0].VSwitchId)
	d.Set("connection_string", instance.Items.DBInstanceAttribute[0].ConnectionString)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("maintain_time", instance.Items.DBInstanceAttribute[0].MaintainTime)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceStorageType, "db_instance_storage_type", "storage_type")
	engine := Trim(d.Get("engine").(string))
	ssl_object, err := PolardbService.DescribeDBInstanceSSL(d.Id())
	ssl := false
	if (engine == "MySQL" && ssl_object["SSLEnabled"].(string) == "Yes") || (engine != "MySQL" && ssl_object["SSLEnabled"].(string) == "on") {
		ssl = true
	}
	d.Set("enable_ssl", ssl)
	tde_object, err := PolardbService.DescribeDBInstanceTDE(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tde_status", tde_object["TDEStatus"].(string) == "Enabled")
	d.Set("encrypt_algorithm", tde_object["EncryptAlgorithm"].(string))
	if err = PolardbService.RefreshParameters(d, client, "parameters"); err != nil {
		return errmsgs.WrapError(err)
	}

	if instance.Items.DBInstanceAttribute[0].PayType == string(Prepaid) {
		response, err := PolardbService.DoPolardbDescribeinstanceautorenewalattributeRequest(d, client)
		if response != nil && len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			d.Set("auto_renew", renew.AutoRenew == "True")
			d.Set("auto_renew_period", renew.Duration)
		}
		period, err := computePeriodByUnit(instance.Items.DBInstanceAttribute[0].CreationTime, instance.Items.DBInstanceAttribute[0].ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("period", period)
	}

	return nil
}

func resourceAlibabacloudStackPolardbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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

	return PolardbService.WaitForDBInstance(d, client, Deleted, DefaultLongTimeout)
}

type PolardbCreatedbinstanceResponse struct {
	RequestId        string `json:"RequestId"`
	DBInstanceId     string `json:"DBInstanceId"`
	OrderId          string `json:"OrderId"`
	ConnectionString string `json:"ConnectionString"`
	Port             string `json:"Port"`
}

type PolardbModifydbinstancedescriptionResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbinstancemaintaintimeResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbinstancemonitorResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbinstancepaytypeResponse struct {
	RequestId    string `json:"RequestId"`
	DBInstanceId string `json:"DBInstanceId"`
	OrderId      int    `json:"OrderId"`
}
type PolardbModifydbinstancesslResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbinstancespecResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifydbinstancetdeResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifyinstanceautorenewalattributeResponse struct {
	RequestId string `json:"RequestId"`
}

type PolardbDeletedbinstanceResponse struct {
	RequestId string `json:"RequestId"`
}
type PolardbModifysecurityipsResponse struct {
	RequestId string      `json:"RequestId"`
	TaskId    interface{} `json:"TaskId"`
}
