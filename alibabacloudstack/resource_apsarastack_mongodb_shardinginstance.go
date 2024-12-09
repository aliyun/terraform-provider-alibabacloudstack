package alibabacloudstack

import (
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

func resourceAlibabacloudStackMongoDBShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackMongoDBShardingInstanceCreate,
		Read:   resourceAlibabacloudStackMongoDBShardingInstanceRead,
		Update: resourceAlibabacloudStackMongoDBShardingInstanceUpdate,
		Delete: resourceAlibabacloudStackMongoDBShardingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
				ForceNew:     true,
				Computed:     true,
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
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
				Deprecated:   "Field 'name' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:         schema.TypeString,
				Optional:     true,
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
			"tde_status": {
				Type:             schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "" && new == "disabled" || old == "enabled"
				},
				ValidateFunc: validation.StringInSlice([]string{"enabled", "disabled"}, false),
				Optional:     true,
			},
			"backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Deprecated: "Field 'backup_period' is deprecated and will be removed in a future release. " +
					"Please use new field 'preferred_backup_period' instead.",
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
				Optional: true,
				Computed: true,
				Deprecated: "Field 'backup_time' is deprecated and will be removed in a future release. " +
					"Please use new field 'preferred_backup_time' instead.",
				ConflictsWith: []string{"preferred_backup_time"},
			},
			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
				ConflictsWith: []string{"backup_time"},
			},
			//Computed
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						"node_storage": {
							Type:     schema.TypeInt,
							Required: true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},

			"mongo_list": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},
		},
	}
}

func buildMongoDBShardingCreateRequest(d *schema.ResourceData, meta interface{}) (*dds.CreateShardingDBInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := dds.CreateCreateShardingDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.Engine = "MongoDB"
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

	shardList, ok := d.GetOk("shard_list")
	if ok {
		replicaSets := []dds.CreateShardingDBInstanceReplicaSet{}
		for _, rew := range shardList.([]interface{}) {
			item := rew.(map[string]interface{})
			class := item["node_class"].(string)
			nodeStorage := item["node_storage"].(int)
			var csr dds.CreateShardingDBInstanceReplicaSet
			csr.Class = class
			csr.Storage = strconv.Itoa(nodeStorage)
			replicaSets = append(replicaSets, csr)
		}
		request.ReplicaSet = &replicaSets
	}

	mongoList, ok := d.GetOk("mongo_list")
	if ok {
		mongos := []dds.CreateShardingDBInstanceMongos{}
		for _, rew := range mongoList.([]interface{}) {
			item := rew.(map[string]interface{})
			class := item["node_class"].(string)
			mongos = append(mongos, dds.CreateShardingDBInstanceMongos{class})
		}
		request.Mongos = &mongos
	}

	request.ConfigServer = &[]dds.CreateShardingDBInstanceConfigServer{{"20", "dds.cs.mid"}}

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
				return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch " + vsw.VSwitchId + " isn't in multi the zone " + request.ZoneId))
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch " + vsw.VSwitchId + " isn't in the zone " + request.ZoneId))
		}
		request.VSwitchId = vswitchId
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.VpcId = vsw.VpcId
	}

	request.ChargeType = d.Get("instance_charge_type").(string)
	period, ok := d.GetOk("period")
	if ok && PayType(request.ChargeType) == PrePaid {
		request.Period = requests.NewInteger(period.(int))
	}

	request.SecurityIPList = LOCAL_HOST_IP
	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		request.SecurityIPList = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List()), COMMA_SEPARATED)
	}

	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func resourceAlibabacloudStackMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	request, err := buildMongoDBShardingCreateRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.CreateShardingDBInstance(request)
	})

	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*dds.CreateShardingDBInstanceResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_mongodb_sharding_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	response, _ := raw.(*dds.CreateShardingDBInstanceResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(response.DBInstanceId)

	if err := ddsService.WaitForMongoDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackMongoDBShardingInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
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
	connectivity.SetResourceData(d, backupPolicy.PreferredBackupPeriod, "preferred_backup_period", "backup_period")
	retention_period, _ := strconv.Atoi(backupPolicy.BackupRetentionPeriod)
	d.Set("retention_period", retention_period)

	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "name")
	d.Set("engine_version", instance.EngineVersion)
	d.Set("storage_engine", instance.StorageEngine)
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

	mongosList := []map[string]interface{}{}
	for _, item := range instance.MongosList.MongosAttribute {
		mongo := map[string]interface{}{
			"node_class":      item.NodeClass,
			"node_id":         item.NodeId,
			"port":            item.Port,
			"connect_string":  item.ConnectSting,
		}
		mongosList = append(mongosList, mongo)
	}
	err = d.Set("mongo_list", mongosList)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	shardList := []map[string]interface{}{}
	for _, item := range instance.ShardList.ShardAttribute {
		shard := map[string]interface{}{
			"node_id":       item.NodeId,
			"node_storage":  item.NodeStorage,
			"node_class":    item.NodeClass,
		}
		shardList = append(shardList, shard)
	}
	err = d.Set("shard_list", shardList)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	if !(d.Get("tde_status") == "" && tdeInfo.TDEStatus == "disabled") {
		d.Set("tde_status", tdeInfo.TDEStatus)
	}

	ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("security_ip_list", ips)
	// 混合云不支持
	//	groupIp, err := ddsService.DescribeMongoDBSecurityGroupId(d.Id())
	//	if err != nil {
	//		return errmsgs.WrapError(err)
	//	}
	//	if len(groupIp.Items.RdsEcsSecurityGroupRel) > 0 {
	//		d.Set("security_group_id", groupIp.Items.RdsEcsSecurityRel[0].SecurityGroupId)
	//	}

	return nil
}

func resourceAlibabacloudStackMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}
	d.Partial(true)

	if d.HasChanges("preferred_backup_time", "preferred_backup_period", "backup_time", "backup_period") {
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
			if raw != nil {
				response, ok := raw.(*dds.ModifyDBInstanceTDEResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("tde_status")
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
			if raw != nil {
				response, ok := raw.(*dds.ModifySecurityGroupConfigurationResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("security_group_id")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackMongoDBShardingInstanceRead(d, meta)
	}

	if d.HasChange("shard_list") {
		state, diff := d.GetChange("shard_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d.Id(), MongoDBShardingNodeShard, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("shard_list")
	}

	if d.HasChange("mongo_list") {
		state, diff := d.GetChange("mongo_list")
		err := ddsService.ModifyMongodbShardingInstanceNode(d.Id(), MongoDBShardingNodeMongos, state.([]interface{}), diff.([]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("mongo_list")
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
			if raw != nil {
				response, ok := raw.(*dds.ModifyDBInstanceDescriptionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("db_instance_description")
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
		//d.SetPartial("account_password")
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
	d.Partial(false)
	return resourceAlibabacloudStackMongoDBShardingInstanceRead(d, meta)
}

func resourceAlibabacloudStackMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	request := dds.CreateDeleteDBInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()

	err := resource.Retry(10*5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*dds.DeleteDBInstanceResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return errmsgs.WrapError(ddsService.WaitForMongoDBInstance(d.Id(), Deleted, DefaultTimeout))
}
