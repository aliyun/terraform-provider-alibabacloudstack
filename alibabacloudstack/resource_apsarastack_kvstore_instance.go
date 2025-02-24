package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKVStoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKVStoreInstanceCreate,
		Read:   resourceAlibabacloudStackKVStoreInstanceRead,
		Update: resourceAlibabacloudStackKVStoreInstanceUpdate,
		Delete: resourceAlibabacloudStackKVStoreInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tair_instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"instance_name"},
			},
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'instance_name' is deprecated and will be removed in a future release. Please use new field 'tair_instance_name' instead.",
				ConflictsWith: []string{"tair_instance_name"},
			},
			"password": {
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
			"instance_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  KVStore5Dot0,
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "Field 'availability_zone' is deprecated and will be removed in a future release. Please use new field 'zone_id' instead.",
				ConflictsWith: []string{"zone_id"},
			},
			// 			"payment_type": {
			// 				Type:          schema.TypeString,
			// 				ValidateFunc:  validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			// 				Optional:      true,
			// 				Computed:      true,
			// 				ConflictsWith: []string{"instance_charge_type"},
			// 			},
			// 			"instance_charge_type": {
			// 				Type:          schema.TypeString,
			// 				ValidateFunc:  validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			// 				Optional:      true,
			// 				Computed:      true,
			// 				Deprecated:    "Field 'instance_charge_type' is deprecated and will be removed in a future release. Please use new field 'payment_type' instead.",
			// 				ConflictsWith: []string{"payment_type"},
			// 			},
			// 			"period": {
			// 				Type:             schema.TypeInt,
			// 				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			// 				Optional:         true,
			// 				Default:          1,
			// 				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			// 			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(KVStoreRedis),
				ValidateFunc: validation.StringInSlice([]string{
					string(KVStoreMemcache),
					string(KVStoreRedis),
				}, false),
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"connection_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
				// Optional: true, ASCM不支持设定
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"vpc_auth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Close"}, false),
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
				Set: func(v interface{}) int {
					return hashcode.String(
						v.(map[string]interface{})["name"].(string) + "|" + v.(map[string]interface{})["value"].(string))
				},
				Optional: true,
				Computed: true,
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
			"cpu_type": {
				Type:     schema.TypeString,
				Optional: true,
				Deprecated:    "Field 'cpu_type' is deprecated and will be removed in a future release.",
			},
			"node_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE", "double"}, false),
				DiffSuppressFunc: NodeTypeDiffSuppressFunc,
			},
			"architecture_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"cluster", "rwsplit", "standard"}, false),
				DiffSuppressFunc: ArchitectureTypeDiffSuppressFunc,
			},
			"series": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"community", "enterprise"}, false),
			},
			"tde_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "TDE does not support the simultaneous use of `encryption_key` and `role_arn`.",
				ConflictsWith: []string{"role_arn"},
			},
			"role_arn": {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "TDE does not support the simultaneous use of `encryption_key` and `role_arn`.",
				ConflictsWith: []string{"encryption_key"},
			},
		},
	}
}

func resourceAlibabacloudStackKVStoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	kvstoreService := KvstoreService{client}
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})

	request["ClientToken"] = buildClientToken("CreateInstance")
	if v, ok := connectivity.GetResourceDataOk(d, "tair_instance_name", "instance_name"); ok {
		request["InstanceName"] = v.(string)
	}
		if v, ok := d.GetOk("cpu_type"); ok {
			request["CpuType"] = v.(string)
		}
	if v, ok := d.GetOk("node_type"); ok {
		request["NodeType"] = v.(string)
	}
	if v, ok := d.GetOk("architecture_type"); ok && v.(string) != "" {
		request["ArchitectureType"] = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v.(string)
	}
	if v, ok := d.GetOk("engine_version"); ok {
		request["EngineVersion"] = v.(string)
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v.(string)
	}
	// 	if v, ok := connectivity.GetResourceDataOk(d, "payment_type", "instance_charge_type"); ok {
	// 		request["ChargeType"] = v.(string)
	// 	} else {
	// 		request["ChargeType"] = string(PostPaid)
	// 	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v.(string)
	}
	if v, ok := d.GetOk("series"); ok {
		request["Series"] = v.(string)
	}

	if request["Password"] == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request["Password"] = decryptResp.Plaintext
		}
	}
	if v, ok := d.GetOk("backup_id"); ok {
		request["BackupId"] = v.(string)
	}

	// 	if request["ChargeType"] == PrePaid {
	// 		request["Period"] = strconv.Itoa(d.Get("period").(int))
	// 	}

	if v, ok := connectivity.GetResourceDataOk(d, "zone_id", "availability_zone"); ok && Trim(v.(string)) != "" {
		request["ZoneId"] = Trim(v.(string))
	}
	request["NetworkType"] = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request["VSwitchId"] = vswitchId.(string)
		request["NetworkType"] = strings.ToUpper(string(Vpc))
		// request["PrivateIpAddress"] = Trim(d.Get("private_ip").(string))

		// check vswitchId in zone
		object, err := vpcService.DescribeVSwitch(vswitchId.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}

		if request["ZoneId"] == "" {
			request["ZoneId"] = object.ZoneId
		}

		request["VpcId"] = object.VpcId
	}

	response, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", action, "", nil, nil, request)
	log.Printf(" create kvstroe instances Finished !! response: %v", response)
	if err != nil {
		return err
	}
	if value, exist := response["asapiSuccess"]; !exist || !value.(bool) {
		err = errmsgs.Error("create kvstroe instances Failed !!")
		return errmsgs.WrapErrorf(err, "create kvstroe instances Failed !! %s", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))
	log.Printf("begin describe kvstroe instances !!")
	// wait instance status change from Creating to Normal
	stateConf := BuildStateConfByTimes([]string{"Creating"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}), 200)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapError(err)
	}

	log.Printf("begin update kvstroe instances !!")
	if tde, ok := d.GetOk("tde_status"); ok && tde.(string) == "Enabled" {
		tde_req := make(map[string]interface{})
		tde_req["InstanceId"] = d.Id()
		tde_req["TDEStatus"] = tde.(string)
		if encryption_key, ok := d.GetOk("encryption_key"); ok && encryption_key != "" {
			tde_req["EncryptionKey"] = encryption_key
		} else {
			if role_arn, ok := d.GetOk("role_arn"); ok && role_arn.(string) != "" {
				tde_req["RoleArn"] = d.Get("role_arn").(string)
			} else if client.Config.RamRoleArn != "" {
				tde_req["RoleArn"] = client.Config.RamRoleArn
			}
		}

		tde_response, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", "ModifyInstanceTDE", "", nil, nil, tde_req)

		addDebug("ModifyInstanceTDE", tde_response, tde_req)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "apsarastack_kvstroe_instance", "ModifyInstanceTDE", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if value, exist := tde_response["asapiSuccess"]; !exist || !value.(bool) {
			err = errmsgs.Error("kvstroe ModifyInstanceTDE Failed !!")
			return errmsgs.WrapErrorf(err, "kvstroe ModifyInstanceTDE Failed !! %s", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}

		log.Print("enabled TDE")
	}
	return resourceAlibabacloudStackKVStoreInstanceUpdate(d, meta)
}

func resourceAlibabacloudStackKVStoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	d.Partial(true)
	stateConf := BuildStateConf([]string{"DBInstanceClassChanging", "DBInstanceNetTypeChanging", "Changing"}, []string{"Normal"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if d.HasChange("parameters") {
		config := make(map[string]interface{})
		documented := d.Get("parameters").(*schema.Set).List()
		if len(documented) > 0 {
			for _, i := range documented {
				key := i.(map[string]interface{})["name"].(string)
				value := i.(map[string]interface{})["value"]
				config[key] = value
			}
			cfg, _ := json.Marshal(config)
			if err := kvstoreService.ModifyInstanceConfig(d.Id(), string(cfg)); err != nil {
				return errmsgs.WrapError(err)
			}
		}

		//d.SetPartial("parameters")
	}

	if d.HasChange("security_ips") {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		request := r_kvstore.CreateModifySecurityIpsRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.SecurityIpGroupName = "default"
		request.InstanceId = d.Id()
		if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
			request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		} else {
			return errmsgs.WrapError(errmsgs.Error("Security ips cannot be empty"))
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifySecurityIps(request)
		})
		if err != nil {
			return err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("security_ips")
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	if d.HasChange("vpc_auth_mode") {
		if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
			// vpc_auth_mode works only if the network type is VPC
			instanceType := d.Get("instance_type").(string)
			if string(KVStoreRedis) == instanceType {
				// wait instance status is Normal before modifying
				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapError(err)
				}

				request := r_kvstore.CreateModifyInstanceVpcAuthModeRequest()
				client.InitRpcRequest(*request.RpcRequest)
				request.InstanceId = d.Id()
				request.VpcAuthMode = d.Get("vpc_auth_mode").(string)

				raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
					return rkvClient.ModifyInstanceVpcAuthMode(request)
				})
				if err != nil {
					errmsg := ""
					if raw != nil {
						baseResponse := &responses.BaseResponse{}
						err = json.Unmarshal([]byte(raw.(*r_kvstore.ModifyInstanceVpcAuthModeResponse).GetHttpContentString()), baseResponse)
						if err == nil {
							errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
						}
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				//d.SetPartial("vpc_auth_mode")

				// The auth mode take some time to be effective, so wait to ensure the state !
				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapError(err)
				}
			}
		}
	}
	// 	configPayType := PayType(connectivity.GetResourceData(d, "payment_type", "instance_charge_type").(string))
	// 	if !d.IsNewResource() && (d.HasChanges("payment_type", "instance_charge_type")) && configPayType == PrePaid {
	// 		// for now we just support charge change from PostPaid to PrePaid
	// 		prePaidRequest := r_kvstore.CreateTransformToPrePaidRequest()
	// 		client.InitRpcRequest(*prePaidRequest.RpcRequest)
	// 		prePaidRequest.InstanceId = d.Id()
	// 		prePaidRequest.Period = requests.Integer(strconv.Itoa(d.Get("period").(int)))
	// 		prePaidRequest.AutoPay = requests.NewBoolean(true)
	// 		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
	// 			return rkvClient.TransformToPrePaid(prePaidRequest)
	// 		})
	// 		if err != nil {
	// 			errmsg := ""
	// 			if raw != nil {
	// 				baseResponse := &responses.BaseResponse{}
	// 				err = json.Unmarshal([]byte(raw.(*r_kvstore.TransformToPrePaidResponse).GetHttpContentString()), baseResponse)
	// 				if err == nil {
	// 					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
	// 				}
	// 			}
	// 			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), prePaidRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	// 		}
	// 		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
	// 		// wait instance status is Normal after modifying
	// 		if _, err := stateConf.WaitForState(); err != nil {
	// 			return errmsgs.WrapError(err)
	// 		}
	// 		//d.SetPartial("instance_charge_type")
	// 		//d.SetPartial("period")
	// 	}

	if d.HasChanges("maintain_start_time", "maintain_end_time") {
		request := r_kvstore.CreateModifyInstanceMaintainTimeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.ModifyInstanceMaintainTime(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				baseResponse := &responses.BaseResponse{}
				err = json.Unmarshal([]byte(raw.(*r_kvstore.ModifyInstanceMaintainTimeResponse).GetHttpContentString()), baseResponse)
				if err == nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("maintain_start_time")
		//d.SetPartial("maintain_end_time")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlibabacloudStackKVStoreInstanceRead(d, meta)
	}

	if d.HasChange("instance_class") {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}

		request := r_kvstore.CreateModifyInstanceSpecRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		request.EffectiveTime = "Immediately"

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.ModifyInstanceSpec(request)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{"MissingRedisUsedmemoryUnsupportPerfItem"}) {
					time.Sleep(time.Duration(5) * time.Second)
					return resource.RetryableError(err)
				}
				errmsg := ""
				if raw != nil {
					baseResponse := &responses.BaseResponse{}
					err = json.Unmarshal([]byte(raw.(*r_kvstore.ModifyInstanceSpecResponse).GetHttpContentString()), baseResponse)
					if err == nil {
						errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
					}
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return err
		}
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		// There needs more time to sync instance class update
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if object.InstanceClass != request.InstanceClass {
				return resource.RetryableError(errmsgs.Error("Waitting for instance class is changed timeout. Expect instance class %s, got %s.",
					object.InstanceClass, request.InstanceClass))
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapError(err)
		}

		//d.SetPartial("instance_class")
	}

	request := r_kvstore.CreateModifyInstanceAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()
	update := false
	if d.HasChanges("tair_instance_name", "instance_name") {
		request.InstanceName = connectivity.GetResourceData(d, "tair_instance_name", "instance_name").(string)
		update = true
	}

	if d.HasChanges("password", "kms_encrypted_password") {
		if v := d.Get("password").(string); v != "" {
			//d.SetPartial("password")
			request.NewPassword = v
			update = true
		}
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request.NewPassword = decryptResp.Plaintext
			//d.SetPartial("kms_encrypted_password")
			//d.SetPartial("kms_encryption_context")
			update = true
		}
	}

	if update {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyInstanceAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				baseResponse := &responses.BaseResponse{}
				err = json.Unmarshal([]byte(raw.(*r_kvstore.ModifyInstanceAttributeResponse).GetHttpContentString()), baseResponse)
				if err == nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("instance_name")
		//d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlibabacloudStackKVStoreInstanceRead(d, meta)
}

func resourceAlibabacloudStackKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.InstanceName, "tair_instance_name", "instance_name")
	d.Set("instance_class", object.InstanceClass)
	// 目前查询接口没有返回是否为企业版本，只能从类型判断
	if strings.HasPrefix(object.InstanceClass, strings.ToLower(object.InstanceType)+".amber.") {
		d.Set("series", "enterprise")
	} else {
		d.Set("series", "community")
	}
	connectivity.SetResourceData(d, object.ZoneId, "zone_id", "availability_zone")
	// 	connectivity.SetResourceData(d, object.ChargeType, "payment_type", "instance_charge_type")
	d.Set("instance_type", object.InstanceType)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("connection_domain", object.ConnectionDomain)
	d.Set("engine_version", object.EngineVersion)
	d.Set("private_ip", object.PrivateIp)
	d.Set("security_ips", strings.Split(object.SecurityIPList, COMMA_SEPARATED))
	d.Set("vpc_auth_mode", object.VpcAuthMode)
	d.Set("maintain_start_time", object.MaintainStartTime)
	d.Set("maintain_end_time", object.MaintainEndTime)

	if object.ChargeType == string(PrePaid) {
		request := r_kvstore.CreateDescribeInstanceAutoRenewalAttributeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				baseResponse := &responses.BaseResponse{}
				err = json.Unmarshal([]byte(raw.(*r_kvstore.DescribeInstanceAutoRenewalAttributeResponse).GetHttpContentString()), baseResponse)
				if err == nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		period, err := computePeriodByUnit(object.CreateTime, object.EndTime, d.Get("period").(int), "Month")
		if err != nil {
			return errmsgs.WrapError(err)
		}
		d.Set("period", period)
	}
	//refresh parameters
	if err = refreshParameters(d, meta); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackKVStoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if PayType(object.ChargeType) == PrePaid {
		return errmsgs.WrapError(errmsgs.Error("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically"))
	}
	request := r_kvstore.CreateDeleteInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteInstance(request)
	})

	if err != nil {
		if !errmsgs.IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			errmsg := ""
			if raw != nil {
				baseResponse := &responses.BaseResponse{}
				err = json.Unmarshal([]byte(raw.(*r_kvstore.DeleteInstanceResponse).GetHttpContentString()), baseResponse)
				if err == nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{"Creating", "Active", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return errmsgs.WrapError(err)
}

func refreshParameters(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	var param []map[string]interface{}
	documented, ok := d.GetOk("parameters")
	if !ok {
		d.Set("parameters", param)
		return nil
	}
	object, err := kvstoreService.DescribeParameters(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.Parameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, i := range object.ConfigParameters.Parameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}

	d.Set("parameters", param)
	return nil
}
