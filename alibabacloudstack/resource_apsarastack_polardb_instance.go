package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
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
			"param_group_id": {
				Type:     schema.TypeString,
				Optional: true,
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
					if v, ok := d.GetOk("tde_status"); ok && v.(bool) {
						return old == new
					}
					return true
				},
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"storage_type": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'storage_type' is deprecated and will be removed in a future release. Please use new field 'db_instance_storage_type' instead.",
				ConflictsWith: []string{"db_instance_storage_type"},
				AtLeastOneOf:  []string{"db_instance_storage_type"},
			},
			"db_instance_storage_type": {
				Type:          schema.TypeString,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_type"},
				AtLeastOneOf:  []string{"storage_type"},
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Field 'encryption' is deprecated and will be removed in a future release. Please use new field 'tde_status' instead.",
			},
			"instance_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_type' is deprecated and will be removed in a future release. Please use new field 'db_instance_class' instead.",
				ConflictsWith: []string{"db_instance_class"},
				AtLeastOneOf:  []string{"db_instance_class"},
			},
			"db_instance_class": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_type"},
				AtLeastOneOf:  []string{"instance_type"},
			},
			"instance_storage": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'instance_storage' is deprecated and will be removed in a future release. Please use new field 'db_instance_storage' instead.",
				ConflictsWith: []string{"db_instance_storage"},
				AtLeastOneOf:  []string{"db_instance_storage"},
			},
			"db_instance_storage": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_storage"},
				AtLeastOneOf:  []string{"instance_storage"},
			},
			"instance_charge_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:      true,
				Computed:      true,
				Deprecated:    "The field `instance_charge_type` has been deprecated and is scheduled for removal in version 3.21.0.",
				ConflictsWith: []string{"payment_type"},
				DiffSuppressFunc: DeprecatedDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
			"payment_type": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_charge_type"},
				Deprecated:    "The field `payment_type` has been deprecated and is scheduled for removal in version 3.21.0.",
				DiffSuppressFunc: DeprecatedDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: DeprecatedDiffSuppressFunc,
				DiffSuppressOnRefresh: true,
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
				DiffSuppressFunc: DeprecatedDiffSuppressFunc,
				Deprecated:    "The field `auto_renew` has been deprecated and is scheduled for removal in version 3.21.0.",
				DiffSuppressOnRefresh: true,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 12),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: DeprecatedDiffSuppressFunc,
				Deprecated:    "The field `auto_renew_period` has been deprecated and is scheduled for removal in version 3.21.0.",
				DiffSuppressOnRefresh: true,
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
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"acl": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"prefer", "require", "verify-ca", "verify-full", "cert"}, false),
				Optional:     true,
				Computed:     true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								oldAll, newAll := d.GetChange("parameters")
								oldSet := oldAll.(*schema.Set)
								newSet := newAll.(*schema.Set)

								removed := make(map[string]string)
								for _, i := range oldSet.Difference(newSet).List() {
									item := i.(map[string]interface{})
									if item["name"].(string) == "" {
										continue
									}
									removed[item["name"].(string)] = item["value"].(string)
								}
								changekey := old
								if _, ok := removed[changekey]; ok {
									return true
								}
								return old == new
							},
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								oldAll, newAll := d.GetChange("parameters")
								oldSet := oldAll.(*schema.Set)
								newSet := newAll.(*schema.Set)

								removed := make(map[string]string)
								for _, i := range oldSet.Difference(newSet).List() {
									item := i.(map[string]interface{})
									if item["name"].(string) == "" {
										continue
									}
									removed[item["name"].(string)] = item["value"].(string)
								}
								parts := strings.Split(k, ".")
								lastIndex := len(parts) - 1
								parts[lastIndex] = "name"
								name_k := strings.Join(parts, ".")
								changekey := d.Get(name_k).(string)
								if _, ok := removed[changekey]; ok {
									return true
								}
								return old == new
							},
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return hashcode.String(m["name"].(string))
				},
				Optional: true,
				Computed: true,
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": caseInsensitiveTagsSchema(),
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
	var err error
	var encryption bool
	EncryptionKey := d.Get("encryption_key").(string)
	encryption = d.Get("tde_status").(bool)
	EncryptAlgorithm := d.Get("encrypt_algorithm").(string)
	log.Print("Encryption key input")
	if EncryptionKey != "" && encryption {
		log.Print("Encryption key condition passed")
		arnrole, err = PolardbService.CheckCloudResourceAuthorized()
		if err != nil {
			return errmsgs.WrapErrorf(err, "CheckCloudResourceAuthorized", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	} else if EncryptionKey == "" && encryption {
		return errmsgs.WrapErrorf(nil, "Add EncryptionKey or Set encryption to false", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
	} else if EncryptionKey != "" && !encryption {
		return errmsgs.WrapErrorf(nil, "Set encryption to true", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
	} else {
		log.Print("Encryption key condition failed")
	}
	log.Printf("encryptionbool %v", d.Get("encryption").(bool))
	log.Printf("check arnrole %v", arnrole)
	enginever := Trim(d.Get("engine_version").(string))
	engine := Trim(d.Get("engine").(string))
	DBInstanceStorage := connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int)
	DBInstanceClass := Trim(connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string))
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
	DBInstanceStorageType := connectivity.GetResourceData(d, "db_instance_storage_type", "storage_type").(string)
	ZoneIdSlave1 = d.Get("zone_id_slave1").(string)
	ZoneIdSlave2 = d.Get("zone_id_slave2").(string)
	SecurityIPList := LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		SecurityIPList = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CreateDBInstance", "")
	PolardbCreatedbinstanceResponse := PolardbCreatedbinstanceResponse{}

	mergeMaps(request.QueryParams, map[string]string{
		"EngineVersion":         enginever,
		"Engine":                engine,
		"DBInstanceStorage":     strconv.Itoa(DBInstanceStorage),
		"DBInstanceClass":       DBInstanceClass,
		"DBInstanceNetType":     DBInstanceNetType,
		"DBInstanceDescription": DBInstanceDescription,
		"InstanceNetworkType":   InstanceNetworkType,
		"VSwitchId":             VSwitchId,
		"PayType":               payType,
		"DBInstanceStorageType": DBInstanceStorageType,
		"SecurityIPList":        SecurityIPList,
		"ZoneIdSlave1":          ZoneIdSlave1,
		"ZoneIdSlave2":          ZoneIdSlave2,
		"ZoneId":                ZoneId,
		"VPCId":                 VPCId,
	})
	if v, ok := d.GetOk("cpu_type"); ok && v.(string) != "" {
		request.QueryParams["CpuType"] = v.(string)
	}
	if v, ok := d.GetOk("param_group_id"); ok && v.(string) != "" {
		request.QueryParams["DBParamGroupId"] = v.(string)
	}
	if tde := d.Get("tde_status"); tde.(bool) && engine != "MySQL" {
		request.QueryParams["TdeStatus"] = "1"
		request.QueryParams["EncryptAlgorithm"] = EncryptAlgorithm
		request.QueryParams["Encryption"] = strconv.FormatBool(encryption)
		request.QueryParams["RoleARN"] = arnrole
		if EncryptionKey != "" {
			request.QueryParams["EncryptionKey"] = EncryptionKey
		}
	}
	log.Printf("request245 %v", request.QueryParams)
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug("CreateDBInstance", bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "CreateDBInstance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCreatedbinstanceResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "CreateDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(PolardbCreatedbinstanceResponse.DBInstanceId)
	d.Set("connection_string", PolardbCreatedbinstanceResponse.ConnectionString)

	stateConf := BuildStateConfByTimes([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}), 100)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	tde := d.Get("tde_status").(bool)
	log.Printf(" ============================================= tde:%t, engine:%s", tde, engine)
	if tde && engine == "MySQL" {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceTDE", "")
		PolardbModifydbinstancetdeResponse := PolardbModifydbinstancetdeResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["TDEStatus"] = "Enabled"
		request.QueryParams["RoleArn"] = arnrole

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
		stateConf := BuildStateConf([]string{"Disabled"}, []string{"Enabled"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, PolardbService.PolardbDBInstanceTdeStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}

		log.Print("enabled TDE")
	}
	if ssl := d.Get("enable_ssl"); ssl.(bool) {
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
		// var target, process string
		// if engine == "MySQL" {
		// 	target = "Yes"
		// 	process = "No"
		// } else {
		// 	target = "on"
		// 	process = "off"
		// }
		stateConf := BuildStateConf([]string{"SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{}))
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
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
	engine := Trim(d.Get("engine").(string))

	if err := PolardbService.ModifyParameters(d, client); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := PolardbService.SetInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}
	payType := Postpaid
	if !d.IsNewResource() && d.HasChanges("instance_type", "db_instance_class", "instance_storage", "db_instance_storage") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSpec", "")
		PolardbModifydbinstancespecResponse := PolardbModifydbinstancespecResponse{}
		request.QueryParams["DBInstanceId"] = d.Id()
		request.QueryParams["PayType"] = string(payType)

		request.QueryParams["DBInstanceClass"] = connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string)
		request.QueryParams["DBInstanceStorage"] = strconv.Itoa(connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int))
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

	if d.HasChanges("instance_type", "db_instance_class", "monitoring_period") && d.Get("monitoring_period").(int) != 0 {
		// XXX: It needs to be reset `monitoring_period` after the change type
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

	if d.HasChange("enable_ssl") {
		ssl := d.Get("enable_ssl").(bool)
		ssl_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSSL", "")
		ssl_req.QueryParams["DBInstanceId"] = d.Id()
		ssl_req.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		if ssl {
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
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_dbinstance", ssl_req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	log.Printf("================================== acl:%s ,acl_change:%t, engine:%s", d.Get("acl").(string), d.HasChange("acl"), engine)
	if d.HasChange("acl") && d.Get("acl").(string) != "" && engine != "MySQL" {
		time.Sleep(5 * time.Second)
		ssl_object, err := PolardbService.DescribeDBInstanceSSL(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if ssl_object["SSLEnabled"].(string) != "on" {
			return errmsgs.WrapError(errmsgs.Error("To modify the ACL, the instance must first enable enable_ssl."))
		}
		acl_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceSSL", "")
		acl_req.QueryParams["DBInstanceId"] = d.Id()
		acl_req.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		acl_req.QueryParams["ACL"] = d.Get("acl").(string)
		acl_response, err := client.ProcessCommonRequest(acl_req)
		addDebug("ModifyDBInstanceSSL-ACL", acl_response, acl_req, acl_req.QueryParams)
		if err != nil {
			if acl_response == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(acl_response.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_dbinstance", acl_req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"SSL_MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 2*time.Minute, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if d.IsNewResource() {
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

	if d.HasChange("tde_status") && d.Get("tde_status").(bool) && engine == "MySQL" {
		arnrole, err := PolardbService.CheckCloudResourceAuthorized()
		if err != nil {
			return errmsgs.WrapErrorf(err, "CheckCloudResourceAuthorized", "CheckCloudResourceAuthorized", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		tde_req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyDBInstanceTDE", "")
		tde_req.QueryParams["DBInstanceId"] = d.Id()
		tde_req.QueryParams["TDEStatus"] = "Enabled"
		tde_req.QueryParams["RoleArn"] = arnrole

		if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
			tde_req.QueryParams["EncryptionKey"] = v.(string)
		}

		bresponse, err := client.ProcessCommonRequest(tde_req)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_dbinstance", "ModifyDBInstanceTDE", tde_req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"TDE_MODIFYING", "CONFIG_ENCRYPTING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		log.Print("Updated TDE")
	}

	return nil
}

func resourceAlibabacloudStackPolardbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	PolardbService := PolardbService{client}

	instance, err := PolardbService.DoPolardbDescribedbinstanceattributeRequest(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	ips, err := PolardbService.GetSecurityIps(d, client)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	tags, err := PolardbService.describeTags(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tags", PolardbService.tagsToMap(tags))

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
	d.Set("cpu_type", instance.Items.DBInstanceAttribute[0].CpuType)
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
	d.Set("network_type", instance.Items.DBInstanceAttribute[0].InstanceNetworkType)
	d.Set("connection_string", instance.Items.DBInstanceAttribute[0].ConnectionString)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("maintain_time", instance.Items.DBInstanceAttribute[0].MaintainTime)
	connectivity.SetResourceData(d, instance.Items.DBInstanceAttribute[0].DBInstanceStorageType, "db_instance_storage_type", "storage_type")
	engine := Trim(d.Get("engine").(string))
	ssl_object, err := PolardbService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	ssl := false
	if (engine == "MySQL" && ssl_object["SSLEnabled"].(string) == "Yes") || (engine != "MySQL" && ssl_object["SSLEnabled"].(string) == "on") {
		ssl = true
	}
	d.Set("enable_ssl", ssl)
	acl, ok := ssl_object["ACL"]
	if ok && acl.(string) != "" {
		d.Set("acl", acl.(string))
	}
	tde_object, err := PolardbService.DescribeDBInstanceTDE(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tde_status", tde_object["TDEStatus"].(string) == "Enabled")
	d.Set("encrypt_algorithm", tde_object["EncryptAlgorithm"].(string))
	encryptionKey := PolardbService.DescribeDBInstanceEncryptionKey(d.Id())
	if encryptionKey != "" {
		d.Set("encryption_key", encryptionKey)
	}
	if arnrole, err := PolardbService.CheckCloudResourceAuthorized(); err == nil && arnrole != "" {
		d.Set("role_arn", arnrole)
	}
	if err = PolardbService.RefreshParameters(d, client); err != nil {
		return errmsgs.WrapError(err)
	}

	if instance.Items.DBInstanceAttribute[0].PayType == string(Prepaid) {
		response, err := PolardbService.DoPolardbDescribeinstanceautorenewalattributeRequest(d, client)
		if err != nil {
			return errmsgs.WrapError(err)
		}
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
	stateConf := BuildStateConf([]string{"SSL_MODIFYING", "TDE_MODIFYING", "DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
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
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_dbinstance", "DeleteAccount", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	stateConf = BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, PolardbService.PolardbDBInstanceStateRefreshFunc(d, client, d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		if !errmsgs.NotFoundError(err) {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
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
