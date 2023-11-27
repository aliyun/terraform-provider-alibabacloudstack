package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/go-uuid"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_ssd", "cloud_pperf", "cloud_sperf"}, false),
				Required:     true,
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
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(Postpaid), string(Prepaid)}, false),
				Optional:     true,
				Default:      Postpaid,
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
			//"multi_zone": {
			//	Type:             schema.TypeBool,
			//	Optional:         true,
			//	Default:          false,
			//},
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
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
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	//request, err := buildDBCreateRequest(d, meta)
	//if err != nil {
	//	return WrapError(err)
	//}
	//client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	//request := rds.CreateCreateDBInstanceRequest()
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Rds"
	request.Domain = client.Domain
	request.Version = "2014-08-15"
	request.Scheme = "http"
	request.ApiName = "CreateDBInstance"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.RegionId = string(client.Region)
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	//request.Headers = map[string]string{"RegionId": string(client.RegionId)}
	var VSwitchId, InstanceNetworkType, ZoneIdSlave1, ZoneIdSlave2, ZoneId, VPCId, arnrole string
	var encryption bool
	EncryptionKey := d.Get("encryption_key").(string)
	encryption = d.Get("encryption").(bool)
	log.Print("Encryption key input")
	if EncryptionKey != "" && encryption == true {
		log.Print("Encryption key condition passed")
		req := requests.NewCommonRequest()
		req.Method = "POST"
		req.Product = "Rds"
		req.Domain = client.Domain
		req.Version = "2014-08-15"
		req.Scheme = "http"
		req.ApiName = "CheckCloudResourceAuthorized"
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.RegionId = string(client.Region)
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Product":         "Rds",
			"RegionId":        client.RegionId,
			"TargetRegionId":  client.RegionId,
		}
		ram, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
			return crClient.ProcessCommonRequest(req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		var arnresp RoleARN
		addDebug(request.GetActionName(), ram, req)
		resparn, _ := ram.(*responses.CommonResponse)
		log.Printf("raw response %v", resparn)
		err = json.Unmarshal(resparn.GetHttpContentBytes(), &arnresp)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "CheckCloudResourceAuthorized", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		arnrole = arnresp.RoleArn
		d.Set("role_arn", arnrole)
		log.Printf("check arnrole %v", arnrole)
	} else if EncryptionKey == "" && encryption == true {
		return WrapErrorf(nil, "Add EncryptionKey or Set encryption to false", "CheckCloudResourceAuthorized", request.GetActionName())
	} else if EncryptionKey != "" && encryption == false {
		return WrapErrorf(nil, "Set encryption to true", "CheckCloudResourceAuthorized", request.GetActionName())
	} else {
		log.Print("Encryption key condition failed")
	}
	d.Set("encryption", encryption)
	log.Printf("encryptionbool %v", d.Get("encryption").(bool))

	enginever := Trim(d.Get("engine_version").(string))
	engine := Trim(d.Get("engine").(string))
	DBInstanceStorage := requests.NewInteger(d.Get("instance_storage").(int))
	DBInstanceClass := Trim(d.Get("instance_type").(string))
	DBInstanceNetType := string(Intranet)
	DBInstanceDescription := d.Get("instance_name").(string)
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
	PayType := Trim(d.Get("instance_charge_type").(string))
	DBInstanceStorageType := d.Get("storage_type").(string)
	ZoneIdSlave1 = d.Get("zone_id_slave1").(string)
	ZoneIdSlave2 = d.Get("zone_id_slave2").(string)
	//if d.Get("multi_zone").(bool)==true{
	//	if ZoneIdSlave1=d.Get("zone_id_slave1").(string);ZoneIdSlave1==""{
	//		return WrapErrorf(nil, "ZoneIdSlave1 should be set if Multi Zone is true", d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	//	}
	//}
	SecurityIPList := LOCAL_HOST_IP
	if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
		SecurityIPList = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	ClientToken := fmt.Sprintf("Terraform-AlibabacloudStack-%d-%s", time.Now().Unix(), uuid)
	request.QueryParams = map[string]string{
		"AccessKeySecret":       client.SecretKey,
		"Product":               "rds",
		"Department":            client.Department,
		"ResourceGroup":         client.ResourceGroup,
		"EngineVersion":         enginever,
		"Engine":                engine,
		"Encryption":            strconv.FormatBool(encryption),
		"DBInstanceStorage":     string(DBInstanceStorage),
		"DBInstanceClass":       DBInstanceClass,
		"DBInstanceNetType":     DBInstanceNetType,
		"DBInstanceDescription": DBInstanceDescription,
		//"MultiZone":MultiZone,
		"InstanceNetworkType":   InstanceNetworkType,
		"VSwitchId":             VSwitchId,
		"PayType":               PayType,
		"DBInstanceStorageType": DBInstanceStorageType,
		"SecurityIPList":        SecurityIPList,
		"ClientToken":           ClientToken,
		"ZoneIdSlave1":          ZoneIdSlave1,
		"ZoneIdSlave2":          ZoneIdSlave2,
		"EncryptionKey":         EncryptionKey,
		"ZoneId":                ZoneId,
		"VPCId":                 VPCId,
		"RoleARN":               arnrole,
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	//request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	log.Printf("request245 %v", request.QueryParams)
	//log.Printf("request245 %v",request)
	//raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
	//	return rdsClient.CreateDBInstance(request)
	//})
	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	if arnrole != "" {
		log.Print("arnrole has been added")
	} else {
		log.Print("arnrole has not been added")
	}
	var resp CreateDBInstanceResponse
	addDebug(request.GetActionName(), raw, request)
	response, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(response.GetHttpContentBytes(), &resp)
	//log.Printf("response for create %v", &resp)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_db_instance", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	log.Printf("response25 %v", response)
	d.SetId(resp.DBInstanceId)
	d.Set("connection_string", resp.ConnectionString)

	// wait instance status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if tde := d.Get("tde_status"); tde == true {
		client := meta.(*connectivity.AlibabacloudStackClient)
		rdsService = RdsService{client}
		tde_req := rds.CreateModifyDBInstanceTDERequest()
		tde_req.RegionId = client.RegionId
		tde_req.Headers = map[string]string{"RegionId": client.RegionId}
		tde_req.DBInstanceId = d.Id()
		if EncryptionKey != "" {
			tde_req.EncryptionKey = EncryptionKey
		}
		tde_req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "RoleARN": arnrole}

		tde_req.TDEStatus = "Enabled"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		tderaw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(tde_req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_db_instance", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}

		log.Print("enabled TDE")
		addDebug(request.GetActionName(), tderaw, request)
	}
	if ssl := d.Get("enable_ssl"); ssl == true {
		ssl_req := rds.CreateModifyDBInstanceSSLRequest()
		ssl_req.RegionId = client.RegionId
		ssl_req.Headers = map[string]string{"RegionId": client.RegionId}
		ssl_req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Forwardedregionid": client.RegionId}
		ssl_req.DBInstanceId = d.Id()
		ssl_req.SSLEnabled = "1"
		ssl_req.ConnectionString = d.Get("connection_string").(string)
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		sslraw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceSSL(ssl_req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		log.Print("enabled SSL")
		addDebug(request.GetActionName(), sslraw, request)
	}
	return resourceAlibabacloudStackDBInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	//const tde_set = false
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Minute, rdsService.RdsDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	//if d.HasChange("parameters") {
	//	if err := rdsService.ModifyParameters(d, "parameters"); err != nil {
	//		return WrapError(err)
	//	}
	//}

	if err := rdsService.setInstanceTags(d); err != nil {
		return WrapError(err)
	}

	payType := PayType(d.Get("instance_charge_type").(string))
	if !d.IsNewResource() && d.HasChange("instance_charge_type") && payType == Prepaid {
		prePaidRequest := rds.CreateModifyDBInstancePayTypeRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			prePaidRequest.Scheme = "https"
		} else {
			prePaidRequest.Scheme = "http"
		}
		prePaidRequest.RegionId = client.RegionId
		prePaidRequest.Headers = map[string]string{"RegionId": client.RegionId}
		prePaidRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		//d.SetPartial("instance_charge_type")
		//d.SetPartial("period")

	}

	if payType == Prepaid && (d.HasChange("auto_renew") || d.HasChange("auto_renew_period")) {
		request := rds.CreateModifyInstanceAutoRenewalAttributeRequest()
		request.DBInstanceId = d.Id()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": string(client.RegionId)}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		//d.SetPartial("auto_renew")
		//d.SetPartial("auto_renew_period")
	}

	if d.HasChange("monitoring_period") {
		period := d.Get("monitoring_period").(int)
		request := rds.CreateModifyDBInstanceMonitorRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": string(client.RegionId)}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()
		request.Period = strconv.Itoa(period)

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMonitor(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("maintain_time") {
		request := rds.CreateModifyDBInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": string(client.RegionId)}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)
		request.ClientToken = buildClientToken(request.GetActionName())

		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("maintain_time")
	}

	//if d.HasChange("security_ip_mode") && d.Get("security_ip_mode").(string) == SafetyMode {
	//	request := rds.CreateMigrateSecurityIPModeRequest()
	//	request.RegionId = client.RegionId
	//	if strings.ToLower(client.Config.Protocol) == "https" {
	//		request.Scheme = "https"
	//	} else {
	//		request.Scheme = "http"
	//	}
	//	request.Headers = map[string]string{"RegionId": string(client.RegionId)}
	//	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//	request.DBInstanceId = d.Id()
	//	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
	//		return rdsClient.MigrateSecurityIPMode(request)
	//	})
	//	if err != nil {
	//		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	//	}
	//	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//	//d.SetPartial("security_ip_mode")
	//}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackDBInstanceRead(d, meta)
	}

	if d.HasChange("instance_name") {
		request := rds.CreateModifyDBInstanceDescriptionRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": string(client.RegionId)}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = d.Get("instance_name").(string)

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyDBInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("instance_name")
	}

	if d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := rdsService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		//d.SetPartial("security_ips")
	}

	update := false
	request := rds.CreateModifyDBInstanceSpecRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": string(client.RegionId)}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = d.Id()
	request.PayType = d.Get("instance_charge_type").(string)

	if d.HasChange("instance_type") {
		request.DBInstanceClass = d.Get("instance_type").(string)
		update = true
	}

	if d.HasChange("instance_storage") {
		request.DBInstanceStorage = requests.NewInteger(d.Get("instance_storage").(int))
		update = true
	}
	if update {
		// wait instance status is running before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidOrderTask.NotSupport"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			//d.SetPartial("instance_type")
			//d.SetPartial("instance_storage")
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	if d.HasChange("tde_status") {
		//if tde:=d.Get("tde_status");tde==true{
		client := meta.(*connectivity.AlibabacloudStackClient)
		rdsService = RdsService{client}
		tde_req := rds.CreateModifyDBInstanceTDERequest()
		tde_req.RegionId = client.RegionId
		tde_req.Headers = map[string]string{"RegionId": client.RegionId}
		tde_req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		tde_req.DBInstanceId = d.Id()
		tde_req.TDEStatus = "Enabled"
		//tde_req.RoleArn=d.Get("role_arn").(string)
		//tde_req.EncryptionKey=d.Get("encryption_key").(string)

		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		raw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(tde_req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		log.Print("Updated TDE")
		addDebug(request.GetActionName(), raw, request)
		//}

	}
	if d.HasChange("enable_ssl") {
		ssl := d.Get("enable_ssl").(bool)
		if ssl == true {
			ssl_req := rds.CreateModifyDBInstanceSSLRequest()
			ssl_req.RegionId = client.RegionId
			ssl_req.Headers = map[string]string{"RegionId": client.RegionId}
			ssl_req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Forwardedregionid": client.RegionId}
			ssl_req.DBInstanceId = d.Id()
			ssl_req.SSLEnabled = "1"
			ssl_req.ConnectionString = d.Get("connection_string").(string)
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			sslraw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
				return client.ModifyDBInstanceSSL(ssl_req)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			log.Print("Updated SSL to true")
			addDebug(request.GetActionName(), sslraw, request)
		} else {
			ssl_req := rds.CreateModifyDBInstanceSSLRequest()
			ssl_req.RegionId = client.RegionId
			ssl_req.Headers = map[string]string{"RegionId": client.RegionId}
			ssl_req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Forwardedregionid": client.RegionId}
			ssl_req.DBInstanceId = d.Id()
			ssl_req.SSLEnabled = "0"
			ssl_req.ConnectionString = d.Get("connection_string").(string)
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			sslraw, err := client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
				return client.ModifyDBInstanceSSL(ssl_req)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
			}
			if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			log.Print("Updated SSL to false")
			addDebug(request.GetActionName(), sslraw, request)
		}
	}
	return resourceAlibabacloudStackDBInstanceRead(d, meta)
}

func resourceAlibabacloudStackDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	ips, err := rdsService.GetSecurityIps(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tags, err := rdsService.describeTags(d)
	if err != nil {
		return WrapError(err)
	}
	if len(tags) > 0 {
		d.Set("tags", rdsService.tagsToMap(tags))
	}

	monitoringPeriod, err := rdsService.DescribeDbInstanceMonitor(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("monitoring_period", monitoringPeriod)
	d.Set("security_ips", ips)
	d.Set("security_ip_mode", instance.SecurityIPMode)
	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("instance_type", instance.DBInstanceClass)
	d.Set("port", instance.Port)
	d.Set("instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.PayType)
	d.Set("period", d.Get("period"))
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("connection_string", instance.ConnectionString)
	d.Set("instance_name", instance.DBInstanceDescription)
	d.Set("maintain_time", instance.MaintainTime)
	d.Set("storage_type", instance.DBInstanceStorageType)

	//if err = rdsService.RefreshParameters(d, "parameters"); err != nil {
	//	return WrapError(err)
	//}

	if instance.PayType == string(Prepaid) {
		request := rds.CreateDescribeInstanceAutoRenewalAttributeRequest()
		request.RegionId = client.RegionId
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": string(client.RegionId)}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*rds.DescribeInstanceAutoRenewalAttributeResponse)
		if response != nil && len(response.Items.Item) > 0 {
			renew := response.Items.Item[0]
			d.Set("auto_renew", renew.AutoRenew == "True")
			d.Set("auto_renew_period", renew.Duration)
		}
		period, err := computePeriodByUnit(instance.CreationTime, instance.ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(instance.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": string(client.RegionId)}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = d.Id()

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(request)
		})

		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus", "OperationDenied.ReadDBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	//stateConf := BuildStateConf([]string{"Processing", "Pending", "NoStart", "Failed", "Default"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, rdsService.RdsTaskStateRefreshFunc(d.Id(), "DeleteDBInstance"))
	//if _, err = stateConf.WaitForState(); err != nil {
	//	return WrapErrorf(err, IdMsg, d.Id())
	//}
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
