package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDBInstanceCreate,
		Read:   resourceAlibabacloudStackDBInstanceRead,
		Update: resourceAlibabacloudStackDBInstanceUpdate,
		Delete: resourceAlibabacloudStackDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
			"tags": tagsSchema(),
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
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Deprecated:    "Field 'force_restart' is deprecated and will be removed in a future release, and not for any use now.",
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceAlibabacloudStackDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	vpcService := VpcService{client}

	var VSwitchId, InstanceNetworkType, ZoneIdSlave1, ZoneIdSlave2, ZoneId, VPCId, arnrole string
	var encryption bool
	EncryptionKey := d.Get("encryption_key").(string)
	encryption = d.Get("encryption").(bool)
	log.Print("Encryption key input")
	if EncryptionKey != "" && encryption == true {
		log.Print("Encryption key condition passed")
		req := client.NewCommonRequest("POST", "Rds", "2014-08-15", "CheckCloudResourceAuthorized", "")
		req.QueryParams["TargetRegionId"] = client.RegionId
		ram, err := client.WithRdsClient(func(RdsClient *rds.Client) (interface{}, error) {
			return RdsClient.ProcessCommonRequest(req)
		})
		resparn, ok := ram.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(resparn.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_db_instance", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		var arnresp RoleARN
		addDebug(req.GetActionName(), ram, req)
		log.Printf("raw response %v", resparn)
		err = json.Unmarshal(resparn.GetHttpContentBytes(), &arnresp)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CheckCloudResourceAuthorized", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
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

	request := client.NewCommonRequest("POST", "Rds", "2014-08-15", "CreateDBInstance", "")
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

	log.Printf("request245 %v", request.QueryParams)
	raw, err := client.WithRdsClient(func(RdsClient *rds.Client) (interface{}, error) {
		return RdsClient.ProcessCommonRequest(request)
	})
	response, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_db_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if arnrole != "" {
		log.Print("arnrole has been added")
	} else {
		log.Print("arnrole has not been added")
	}
	var resp CreateDBInstanceResponse
	addDebug(request.GetActionName(), raw, request)
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_db_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	log.Printf("response25 %v", response)
	d.SetId(resp.DBInstanceId)
	d.Set("connection_string", resp.ConnectionString)

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	if tde := d.Get("tde_status"); tde == true {
		client := meta.(*connectivity.AlibabacloudStackClient)
		rdsService = RdsService{client}
		tde_req := rds.CreateModifyDBInstanceTDERequest()
		client.InitRpcRequest(*tde_req.RpcRequest)
		tde_req.QueryParams["RoleARN"] = arnrole
		tde_req.DBInstanceId = d.Id()
		if EncryptionKey != "" {
			tde_req.EncryptionKey = EncryptionKey
		}
		tde_req.TDEStatus = "Enabled"
		tderaw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(tde_req)
		})
		if err != nil {
			errmsg := ""
			if tderaw != nil {
				response, ok := tderaw.(*rds.ModifyDBInstanceTDEResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}

		log.Print("enabled TDE")
		addDebug(request.GetActionName(), tderaw, request)
	}
	if ssl := d.Get("enable_ssl"); ssl == true {
		ssl_req := rds.CreateModifyDBInstanceSSLRequest()
		client.InitRpcRequest(*ssl_req.RpcRequest)
		ssl_req.QueryParams["Forwardedregionid"] = client.RegionId
		ssl_req.DBInstanceId = d.Id()
		ssl_req.SSLEnabled = "1"
		ssl_req.ConnectionString = d.Get("connection_string").(string)
		sslraw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceSSL(ssl_req)
		})
		if err != nil {
			errmsg := ""
			if sslraw != nil {
				response, ok := sslraw.(*rds.ModifyDBInstanceSSLResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		log.Print("enabled SSL")
		addDebug(ssl_req.GetActionName(), sslraw, ssl_req)
	}
	return resourceAlibabacloudStackDBInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if d.HasChange("parameters") {
		if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if err := rdsService.setInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}

	payType := Postpaid
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok && Trim(v.(string)) != "" {
		payType = PayType(Trim(v.(string)))
	}
	if !d.IsNewResource() && d.HasChanges("instance_charge_type", "payment_type") && payType == Prepaid {
		prePaidRequest := rds.CreateModifyDBInstancePayTypeRequest()
		client.InitRpcRequest(*prePaidRequest.RpcRequest)
		prePaidRequest.DBInstanceId = d.Id()
		prePaidRequest.PayType = string(payType)
		prePaidRequest.AutoPay = "true"
		period := d.Get("period").(int)
		prePaidRequest.UsedTime = requests.Integer(strconv.Itoa(period))
		prePaidRequest.Period = string(Month)
		if period > 9 {
			prePaidRequest.UsedTime = requests.Integer(strconv.Itoa(period / 12))
			prePaidRequest.Period = string(Year)
		}

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstancePayType(prePaidRequest)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyDBInstancePayTypeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), prePaidRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest, prePaidRequest.QueryParams)
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	if payType == Prepaid && d.HasChanges("auto_renew", "auto_renew_period") {
		request := rds.CreateModifyInstanceAutoRenewalAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		auto_renew := d.Get("auto_renew").(bool)
		if auto_renew {
			request.AutoRenew = "True"
		} else {
			request.AutoRenew = "False"
		}
		request.Duration = strconv.Itoa(d.Get("auto_renew_period").(int))

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyInstanceAutoRenewalAttributeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(request.GetActionName(), raw, request, request.QueryParams)
	}

	if d.HasChange("monitoring_period") {
		period := d.Get("monitoring_period").(int)
		request := rds.CreateModifyDBInstanceMonitorRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.Period = strconv.Itoa(period)

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMonitor(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyDBInstanceMonitorResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	}

	if d.HasChange("maintain_time") {
		request := rds.CreateModifyDBInstanceMaintainTimeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyDBInstanceMaintainTimeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackDBInstanceRead(d, meta)
	}

	if d.HasChanges("instance_name", "db_instance_description") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = connectivity.GetResourceData(d, "db_instance_description", "instance_name").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*rds.ModifyDBInstanceDescriptionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}
		if err := rdsService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	update := false
	request := rds.CreateModifyDBInstanceSpecRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok {
		request.PayType = v.(string)
	} else {
		request.PayType = string(Postpaid)
	}

	if d.HasChanges("instance_type", "db_instance_class") {
		request.DBInstanceClass = connectivity.GetResourceData(d, "db_instance_class", "instance_type").(string)
		if err := errmsgs.CheckEmpty(request.DBInstanceClass, schema.TypeString, "db_instance_class", "instance_type"); err != nil {
			return errmsgs.WrapError(err)
		}
		update = true
	}

	if d.HasChanges("instance_storage", "db_instance_storage") {
		request.DBInstanceStorage = requests.NewInteger(connectivity.GetResourceData(d, "db_instance_storage", "instance_storage").(int))
		if err := errmsgs.CheckEmpty(request.DBInstanceStorage, schema.TypeString, "db_instance_storage", "instance_storage"); err != nil {
			return errmsgs.WrapError(err)
		}
		update = true
	}
	if update {
		// wait instance status is running before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
					return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
				}
				errmsg := ""
				if raw != nil {
					response, ok := raw.(*rds.ModifyDBInstanceSpecResponse)
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
					}
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errmsgs.WrapError(err)
		}

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	d.Partial(false)
	if d.HasChange("tde_status") {
		tde_req := client.NewCommonRequest("POST", "Rds", "2014-08-15", "ModifyDBInstanceTDE", "")
		tde_req.QueryParams["DBInstanceId"] = d.Id()
		tde_req.QueryParams["TDEStatus"] = "Enabled"
		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ProcessCommonRequest(tde_req)
		})
		response, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), tde_req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		log.Print("Updated TDE")
		addDebug(tde_req.GetActionName(), raw, tde_req)
	}

	if d.HasChange("enable_ssl") {
		ssl := d.Get("enable_ssl").(bool)
		ssl_req := client.NewCommonRequest("POST", "Rds", "2014-08-15", "ModifyDBInstanceSSL", "")
		ssl_req.QueryParams["DBInstanceId"] = d.Id()
		ssl_req.QueryParams["ConnectionString"] = d.Get("connection_string").(string)
		if ssl == true {
			ssl_req.QueryParams["SSLEnabled"] = "1"
		} else {
			ssl_req.QueryParams["SSLEnabled"] = "0"
		}
		sslraw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ProcessCommonRequest(ssl_req)
		})
		response, ok := sslraw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), ssl_req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		if ssl == true {
			log.Print("Updated SSL to true")
		} else {
			log.Print("Updated SSL to false")
		}
		addDebug(ssl_req.GetActionName(), sslraw, ssl_req)
	}
	return resourceAlibabacloudStackDBInstanceRead(d, meta)
}

func resourceAlibabacloudStackDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
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

	ips, err := rdsService.GetSecurityIps(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	tags, err := rdsService.describeTags(d)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", rdsService.tagsToMap(tags))
	}

	monitoringPeriod, err := rdsService.DescribeDbInstanceMonitor(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("monitoring_period", monitoringPeriod)
	d.Set("security_ips", ips)
	d.Set("security_ip_mode", instance.SecurityIPMode)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	connectivity.SetResourceData(d, instance.DBInstanceClass, "db_instance_class", "instance_type")
	d.Set("port", instance.Port)
	connectivity.SetResourceData(d, instance.DBInstanceStorage, "db_instance_storage", "instance_storage")
	d.Set("zone_id", instance.ZoneId)
	if instance.PayType != "" {
		// 专有云场景下不会返回pay type
		connectivity.SetResourceData(d, instance.PayType, "payment_type", "instance_charge_type")
	}
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "instance_name")
	d.Set("maintain_time", instance.MaintainTime)
	connectivity.SetResourceData(d, instance.DBInstanceStorageType, "db_instance_storage_type", "storage_type")

	if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
		return errmsgs.WrapError(err)
	}

	if instance.PayType == string(Prepaid) {
		request := client.NewCommonRequest("POST", "Rds", "2014-08-15", "DescribeInstanceAutoRenewalAttribute", "")
		request.QueryParams["DBInstanceId"] = d.Id()
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ProcessCommonRequest(request)
		})
		response, ok := raw.(*rds.DescribeInstanceAutoRenewalAttributeResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		if response != nil && len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			d.Set("auto_renew", renew.AutoRenew == "True")
			d.Set("auto_renew_period", renew.Duration)
		}
		period, err := computePeriodByUnit(instance.CreationTime, instance.ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("period", period)
	}

	return nil
}

func resourceAlibabacloudStackDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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

	request := client.NewCommonRequest("POST", "Rds", "2014-08-15", "DeleteDBInstance", "")
	request.QueryParams["DBInstanceId"] = d.Id()

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ProcessCommonRequest(request)
		})

		if err != nil && !errmsgs.NotFoundError(err) {
			if errmsgs.IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			response, ok := raw.(*responses.CommonResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

type CreateDBInstanceResponse struct {
	*responses.BaseResponse
	RequestId        string `json:"RequestId" xml:"RequestId"`
	DBInstanceId     string `json:"DBInstanceId" xml:"DBInstanceId"`
	OrderId          string `json:"OrderId" xml:"OrderId"`
	ConnectionString string `json:"ConnectionString" xml:"ConnectionString"`
	Port             string `json:"Port" xml:"Port"`
}

type RoleARN struct {
	AuthorizationState int    `json:"AuthorizationState"`
	EagleEyeTraceID    string `json:"eagleEyeTraceId"`
	AsapiSuccess       bool   `json:"asapiSuccess"`
	AsapiRequestID     string `json:"asapiRequestId"`
	RequestID          string `json:"RequestId"`
	RoleArn            string `json:"RoleArn"`
}
