package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackMongoDBInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(10, 2000),
				Required:     true,
			},
			"replication_factor": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{3, 5, 7}),
				Optional:     true,
				Computed:     true,
			},
			"storage_engine": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"WiredTiger", "RocksDB"}, false),
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
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
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				ConflictsWith: []string{"name"},
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Deprecated: "Field 'backup_period' is deprecated and will be removed in a future release. Please use new field 'preferred_backup_period' instead.",
				ConflictsWith: []string{"preferred_backup_period"},
			},
			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				ConflictsWith: []string{"backup_period"},
			},
			"backup_time": {
				Type:     schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
				Deprecated: "Field 'backup_time' is deprecated and will be removed in a future release. Please use new field 'preferred_backup_time' instead.",
				ConflictsWith: []string{"preferred_backup_time"},
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
				ConflictsWith: []string{"backup_time"},
			},
			"ssl_action": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close", "Update"}, false),
				Optional:     true,
				Computed:     true,
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replica_set_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tde_status": {
				Type: schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "" && new == "disabled" || old == "enabled"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
				Optional:     true,
			},
			"maintain_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackMongoDBInstanceCreate, resourceAlibabacloudStackMongoDBInstanceRead, resourceAlibabacloudStackMongoDBInstanceUpdate, resourceAlibabacloudStackMongoDBInstanceDelete)
	return resource
}

func buildMongoDBCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateDBInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := dds.CreateCreateDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
	request.DBInstanceStorage = requests.NewInteger(d.Get("db_instance_storage").(int))
	request.DBInstanceClass = Trim(d.Get("db_instance_class").(string))
	request.DBInstanceDescription = connectivity.GetResourceData(d, "db_instance_description", "name").(string)

	request.AccountPassword = d.Get("account_password").(string)
	if request.AccountPassword == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, errmsgs.WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext
		}
	}

	request.ZoneId = d.Get("zone_id").(string)
	request.StorageEngine = d.Get("storage_engine").(string)

	if replication_factor, ok := d.GetOk("replication_factor"); ok {
		request.ReplicationFactor = strconv.Itoa(replication_factor.(int))
	}

	request.NetworkType = string(Classic)
	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, errmsgs.WrapError(fmt.Errorf("The specified vswitch %s isn't in the multi zone %s.", vsw.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, errmsgs.WrapError(fmt.Errorf("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if PayType(request.ChargeType) == PrePaid && ok {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List())[:], COMMA_SEPARATED)
	}

	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceAlibabacloudStackMongoDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBCreateRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateDBInstance(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*dds.CreateDBInstanceResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	response, _ := raw.(*dds.CreateDBInstanceResponse)

	d.SetId(response.DBInstanceId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackMongoDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	instance, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	backupPolicy, err := ddsService.DescribeMongoDBBackupPolicy(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, backupPolicy.PreferredBackupTime, "preferred_backup_time", "backup_time")
	connectivity.SetResourceData(d, strings.Split(backupPolicy.PreferredBackupPeriod, ","), "preferred_backup_period", "backup_period")
	retention_period, _ := strconv.Atoi(backupPolicy.BackupRetentionPeriod)
	d.Set("retention_period", retention_period)

	ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("security_ip_list", ips)

	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "name")
	d.Set("engine_version", instance.EngineVersion)
	d.Set("db_instance_class", instance.DBInstanceClass)
	d.Set("db_instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("instance_charge_type", instance.ChargeType)
	if instance.ChargeType == "PrePaid" {
		period, err := computePeriodByUnit(instance.CreationTime, instance.ExpireTime, d.Get("period").(int), "Month")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("period", period)
	}
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("storage_engine", instance.StorageEngine)
	d.Set("maintain_start_time", instance.MaintainStartTime)
	d.Set("maintain_end_time", instance.MaintainEndTime)
	d.Set("replica_set_name", instance.ReplicaSetName)

	sslAction, err := ddsService.DescribeDBInstanceSSL(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("ssl_status", sslAction.SSLStatus)

	if replication_factor, err := strconv.Atoi(instance.ReplicationFactor); err == nil {
		d.Set("replication_factor", replication_factor)
	}
	tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if !(d.Get("tde_status") == "" && tdeInfo.TDEStatus == "disabled") {
		d.Set("tde_status", tdeInfo.TDEStatus)
	}

	d.Set("tags", ddsService.tagsInAttributeToMap(instance.Tags.Tag))
	return nil
}

func resourceAlibabacloudStackMongoDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	d.Partial(true)

	if !d.IsNewResource() && (d.HasChange("instance_charge_type") && d.Get("instance_charge_type").(string) == "PrePaid") {
		prePaidRequest := dds.CreateTransformToPrePaidRequest()
		client.InitRpcRequest(*prePaidRequest.RpcRequest)
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.AutoPay = requests.NewBoolean(true)
		prePaidRequest.Period = requests.NewInteger(d.Get("period").(int))

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.TransformToPrePaid(prePaidRequest)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.TransformToPrePaidResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), prePaidRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is running after modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 0, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("instance_charge_type")
		//d.SetPartial("period")
	}

	if d.HasChanges("preferred_backup_time", "preferred_backup_period", "backup_time", "backup_period"){
		if err := ddsService.MotifyMongoDBBackupPolicy(d); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("preferred_backup_time")
		//d.SetPartial("preferred_backup_period")
	}

	if d.HasChange("tde_status") {
		request := dds.CreateModifyDBInstanceTDERequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.TDEStatus = d.Get("tde_status").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceTDE(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifyDBInstanceTDEResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("tde_status")
	}

	if d.HasChanges("maintain_start_time", "maintain_end_time") {
		request := dds.CreateModifyDBInstanceMaintainTimeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifyDBInstanceMaintainTime(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifyDBInstanceMaintainTimeResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("maintain_start_time")
		//d.SetPartial("maintain_end_time")
	}

	if d.HasChange("security_group_id") {
		request := dds.CreateModifySecurityGroupConfigurationRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.SecurityGroupId = d.Get("security_group_id").(string)

		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.ModifySecurityGroupConfiguration(request)
		})
		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifySecurityGroupConfigurationResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("security_group_id")
	}

	if err := ddsService.setInstanceTags(d); err != nil {
		return errmsgs.WrapError(err)
	}

	if d.IsNewResource() {
		d.Partial(false)
		return nil
	}

	if d.HasChanges("db_instance_description", "name"){
		request := dds.CreateModifyDBInstanceDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.DBInstanceDescription = connectivity.GetResourceData(d, "db_instance_description", "name").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceDescription(request)
		})

		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifyDBInstanceDescriptionResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("db_instance_description")
	}

	if d.HasChange("security_ip_list") {
		ipList := expandStringList(d.Get("security_ip_list").(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := ddsService.ModifyMongoDBSecurityIps(d.Id(), ipstr); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("security_ip_list")
	}

	if d.HasChanges("account_password", "kms_encrypted_password") {
		var accountPassword string
		if accountPassword = d.Get("account_password").(string); accountPassword != "" {
			//d.SetPartial("account_password")
		} else if kmsPassword := d.Get("kms_encrypted_password").(string); kmsPassword != "" {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			accountPassword = decryptResp.Plaintext
			//d.SetPartial("kms_encrypted_password")
			//d.SetPartial("kms_encryption_context")
		}

		err := ddsService.ResetAccountPassword(d, accountPassword)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("ssl_action") {
		request := dds.CreateModifyDBInstanceSSLRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.SSLAction = d.Get("ssl_action").(string)

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSSL(request)
		})

		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifyDBInstanceSSLResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("ssl_action")
	}

	if d.HasChange("db_instance_storage") ||
		d.HasChange("db_instance_class") ||
		d.HasChange("replication_factor") {

		request := dds.CreateModifyDBInstanceSpecRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()

		request.DBInstanceClass = d.Get("db_instance_class").(string)
		request.DBInstanceStorage = strconv.Itoa(d.Get("db_instance_storage").(int))
		request.ReplicationFactor = strconv.Itoa(d.Get("replication_factor").(int))

		// wait instance status is running before modifying
		stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}

		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.ModifyDBInstanceSpec(request)
		})

		if err != nil {
			errmsg := ""
			if bresponse, ok := raw.(*dds.ModifyDBInstanceSpecResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("db_instance_class")
		//d.SetPartial("db_instance_storage")
		//d.SetPartial("replication_factor")

		// wait instance status is running after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackMongoDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDeleteDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})
		bresponse, ok := raw.(*dds.DeleteDBInstanceResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)

			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
	}
	stateConf := BuildStateConf([]string{"Creating", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return errmsgs.WrapError(err)
}