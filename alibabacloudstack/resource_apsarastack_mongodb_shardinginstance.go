package alibabacloudstack

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackMongoDBShardingInstance() *schema.Resource {
	resource := &schema.Resource{
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
			//			"instance_charge_type": {
			//				Type:         schema.TypeString,
			//				ValidateFunc: validation.StringInSlice([]string{string(PrePaid), string(PostPaid)}, false),
			//				Optional:     true,
			//				ForceNew:     true,
			//				Default:     string(PostPaid),
			//			},
			//			"period": {
			//				Type:             schema.TypeInt,
			//				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			//				Optional:         true,
			//				Computed:         true,
			//				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			//			},
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				Deprecated:    "Field 'name' is deprecated and will be removed in a future release. Please use new field 'db_instance_description' instead.",
				ConflictsWith: []string{"db_instance_description"},
			},
			"db_instance_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 256),
				ConflictsWith: []string{"name"},
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			//			"security_group_id": {
			//				Type:     schema.TypeString,
			//				Computed: true,
			//				Optional: true,
			//			},
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
				Type: schema.TypeString,
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
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"backup_period"},
			},
			"backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Computed:     true,
				Deprecated: "Field 'backup_time' is deprecated and will be removed in a future release. " +
					"Please use new field 'preferred_backup_time' instead.",
				ConflictsWith: []string{"preferred_backup_time"},
			},
			"preferred_backup_time": {
				Type:          schema.TypeString,
				ValidateFunc:  validation.StringInSlice(BACKUP_TIME, false),
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"backup_time"},
			},
			//Computed
			"retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_list": {
				Type: schema.TypeSet,
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
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateShardeNodeDescription(),
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					hashString := fmt.Sprintf("%s:%d:%s", m["node_class"].(string), m["node_storage"].(int), m["description"].(string))
					return hashcode.String(hashString)
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},

			"mongo_list": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_class": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateShardeNodeDescription(),
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					hashString := fmt.Sprintf("%s:%s", m["node_class"].(string), m["description"].(string))
					return hashcode.String(hashString)
				},
				Required: true,
				MinItems: 2,
				MaxItems: 32,
			},
			"configserver_list": {
				Type: schema.TypeSet,
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
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateShardeNodeDescription(),
						},
						//Computed
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					hashString := fmt.Sprintf("%s:%d:%s", m["node_class"].(string), m["node_storage"].(int), m["description"].(string))
					return hashcode.String(hashString)
				},
				Required: true,
				MinItems: 1,
				MaxItems: 1,
			},
			"audit_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disabled"}, false),
			},
			"audit_filter": auditFilterSchema([]string{"db", "mongos"}),
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackMongoDBShardingInstanceCreate,
		resourceAlibabacloudStackMongoDBShardingInstanceRead, resourceAlibabacloudStackMongoDBShardingInstanceUpdate, resourceAlibabacloudStackMongoDBShardingInstanceDelete)

	return resource
}

func resourceAlibabacloudStackMongoDBShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	reqQuery := map[string]interface{}{
		"EngineVersion":         d.Get("engine_version").(string),
		"Engine":                "MongoDB",
		"DBInstanceDescription": connectivity.GetResourceData(d, "db_instance_description", "name").(string),
		"ZoneId":                d.Get("zone_id").(string),
		"ChargeType":            string(PostPaid), //d.Get("instance_charge_type").(string),
	}

	reqQuery["AccountPassword"] = d.Get("account_password").(string)
	if reqQuery["AccountPassword"].(string) == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			reqQuery["AccountPassword"] = decryptResp.Plaintext
		}
	}

	shardList := []map[string]interface{}{}
	for _, d := range d.Get("shard_list").(*schema.Set).List() {
		item := d.(map[string]interface{})
		shardNode := map[string]interface{}{
			"Class":       item["node_class"],
			"Storage":     item["node_storage"],
			"Description": item["description"],
		}
		shardList = append(shardList, shardNode)
	}
	reqQuery["ReplicaSet"] = shardList

	mongoList := []map[string]interface{}{}
	for _, d := range d.Get("mongo_list").(*schema.Set).List() {
		item := d.(map[string]interface{})
		mongoNode := map[string]interface{}{
			"Class":       item["node_class"],
			"Description": item["description"],
		}
		mongoList = append(mongoList, mongoNode)
	}
	reqQuery["Mongos"] = mongoList

	csList := []map[string]interface{}{}
	for _, d := range d.Get("configserver_list").(*schema.Set).List() {
		item := d.(map[string]interface{})
		csNode := map[string]interface{}{
			"Class":       item["node_class"],
			"Storage":     item["node_storage"],
			"Description": item["description"],
		}
		csList = append(csList, csNode)
	}
	reqQuery["ConfigServer"] = csList

	vswitchId := Trim(d.Get("vswitch_id").(string))
	if vswitchId != "" {
		// check vswitchId in zone
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(vswitchId)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		zoneId := reqQuery["ZoneId"].(string)
		if zoneId == "" {
			zoneId = vsw.ZoneId
		} else if strings.Contains(zoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(zoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in multi the zone %s", vsw.VSwitchId, zoneId))
			}
		} else if zoneId != vsw.ZoneId {
			return errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the zone %s", vsw.VSwitchId, zoneId))
		}
		reqQuery["VSwitchId"] = vswitchId
		reqQuery["NetworkType"] = strings.ToUpper(string(Vpc))
		reqQuery["VpcId"] = vsw.VpcId
	} else {
		reqQuery["NetworkType"] = string(Classic)
	}

	//	if period, ok := d.GetOk("period"); ok && PayType(reqQuery["ChargeType"].(string)) == PrePaid {
	//		reqQuery["Period"] = period.(int)
	//	}

	if len(d.Get("security_ip_list").(*schema.Set).List()) > 0 {
		reqQuery["SecurityIPList"] = strings.Join(expandStringList(d.Get("security_ip_list").(*schema.Set).List()), COMMA_SEPARATED)
	} else {
		reqQuery["SecurityIPList"] = LOCAL_HOST_IP
	}
	response, err := client.DoTeaRequest("POST", "Dds", "2015-12-01", "CreateShardingDBInstance", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}

	d.SetId(response["DBInstanceId"].(string))

	stateConf := BuildStateConf([]string{"Creating"},
		[]string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ddsService.MongoDbInstanceStateRefreshFunc(d.Id(), []string{"failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func resourceAlibabacloudStackMongoDBShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	// 第一阶段：获取基础实例信息（串行）
	instance, err := ddsService.DescribeMongoDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, instance.DBInstanceDescription, "db_instance_description", "name")
	d.Set("engine_version", instance.EngineVersion)
	d.Set("storage_engine", instance.StorageEngine)
	d.Set("zone_id", instance.ZoneId)
	d.Set("vswitch_id", instance.VSwitchId)

	// 第二阶段：并发获取其他数据
	var wg sync.WaitGroup
	errChan := make(chan error, 5) // 根据实际任务数调整缓冲区大小
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 并发任务1：节点信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			nodes, err := ddsService.DdsDescribeShardingInstanceNodes(d.Id())
			if err != nil {
				errChan <- fmt.Errorf("DdsDescribeShardingInstanceNodes: %w", err)
				cancel()
				return
			}
			for key, info := range nodes {
				data := []map[string]interface{}{}
				for _, v := range info {
					data = append(data, v.(map[string]interface{}))
				}
				d.Set(key, data)
			}
		}
	}()

	// 并发任务2：TDE状态
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			tdeInfo, err := ddsService.DescribeMongoDBTDEInfo(d.Id())
			if err != nil {
				errChan <- fmt.Errorf("DescribeMongoDBTDEInfo: %w", err)
				cancel()
				return
			}
			if !(d.Get("tde_status") == "" && tdeInfo.TDEStatus == "disabled") {
				d.Set("tde_status", tdeInfo.TDEStatus)
			}
		}
	}()

	// 并发任务3：安全IP列表
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			ips, err := ddsService.DescribeMongoDBSecurityIps(d.Id())
			if err != nil {
				errChan <- fmt.Errorf("DescribeMongoDBSecurityIps: %w", err)
				cancel()
				return
			}
			d.Set("security_ip_list", ips)
		}
	}()

	// 并发任务4：审计策略
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			auditResponse, err := ddsService.DoDdsDescribeauditpolicyRequest(d.Id())
			if err != nil {
				errChan <- fmt.Errorf("DoDdsDescribeauditpolicyRequest: %w", err)
				cancel()
				return
			}
			d.Set("audit_status", auditResponse.LogAuditStatus)
		}
	}()

	// 并发任务5：审计过滤器
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			auditFilters, err := ddsService.GetAuditLogFilter(d.Id())
			if err != nil {
				errChan <- fmt.Errorf("GetAuditLogFilter: %w", err)
				cancel()
				return
			}
			d.Set("audit_filter", auditFilters)
		}
	}()

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 错误处理
	for e := range errChan {
		if e != nil {
			return errmsgs.WrapError(e)
		}
	}

	return nil
}

func resourceAlibabacloudStackMongoDBShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	if d.HasChanges("audit_status") && d.Get("audit_status").(string) != "" {
		request := client.NewCommonRequest("POST", "Dds", "2015-12-01", "ModifyAuditPolicy", "")

		request.QueryParams["DBInstanceId"] = d.Id()

		request.QueryParams["AuditStatus"] = d.Get("audit_status").(string)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_mongo_db_audit_policy", "ModifyAuditPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	if _, ok := d.GetOk("audit_filter"); ok && d.Get("audit_status").(string) == "Enable" && d.HasChange("audit_filter") {
		if err := ddsService.ModifyAuditLogFilter(d); err != nil {
			return err
		}
	}

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

	//	if d.HasChange("security_group_id") {
	//		request := dds.CreateModifySecurityGroupConfigurationRequest()
	//		client.InitRpcRequest(*request.RpcRequest)
	//		request.DBInstanceId = d.Id()
	//		request.SecurityGroupId = d.Get("security_group_id").(string)
	//
	//		raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
	//			return client.ModifySecurityGroupConfiguration(request)
	//		})
	//		if err != nil {
	//			errmsg := ""
	//			if raw != nil {
	//				response, ok := raw.(*dds.ModifySecurityGroupConfigurationResponse)
	//				if ok {
	//					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
	//				}
	//			}
	//			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	//		}
	//		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//		//d.SetPartial("security_group_id")
	//	}

	if d.IsNewResource() {
		return nil
	}

	for _, param := range []string{"shard_list", "mongo_list", "configserver_list"} {
		if d.HasChange(param) {
			err := ddsService.ModifyMongodbShardingInstanceNode(d, param)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	if d.HasChanges("db_instance_description", "name") {
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
	return nil
}

func resourceAlibabacloudStackMongoDBShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackMongoDBInstanceDelete(d, meta)
}
