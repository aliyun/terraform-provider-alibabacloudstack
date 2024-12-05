package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"strconv"

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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				Deprecated:   "Field 'instance_name' is deprecated and will be removed in a future release. Please use 'tair_instance_name' instead.",
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Deprecated:   "Field 'availability_zone' is deprecated and will be removed in a future release. Please use 'zone_id' instead.",
			},
			"payment_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
				Deprecated:   "Field 'instance_charge_type' is deprecated and will be removed in a future release. Please use 'payment_type' instead.",
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},
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
				Optional: true,
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
				Required: true,
			},
			"node_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"MASTER_SLAVE", "STAND_ALONE"}, false),
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
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"community", "enterprise"}, false),
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
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "tair_instance_name", "instance_name"); err == nil {
		request["InstanceName"] = v.(string)
	} else {
		return err
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		request["CpuType"] = v.(string)
	}
	if v, ok := d.GetOk("node_type"); ok {
		request["NodeType"] = v.(string)
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
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "payment_type", "instance_charge_type"); err == nil {
		request["ChargeType"] = v.(string)
	} else {
		return err
	}
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

	if request["ChargeType"] == PrePaid {
		request["Period"] = strconv.Itoa(d.Get("period").(int))
	}

	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "zone_id", "availability_zone"); err == nil && Trim(v.(string)) != "" {
		request["ZoneId"] = Trim(v.(string))
	} else if err != nil {
		return err
	}
	request["NetworkType"] = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request["VSwitchId"] = vswitchId.(string)
		request["NetworkType"] = strings.ToUpper(string(Vpc))
		request["PrivateIpAddress"] = Trim(d.Get("private_ip").(string))

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

	response, err := client.DoTeaRequest("POST", "R-kvstore", "2015-01-01", action, "", nil, request)
	log.Printf(" create kvstroe instances Finished !! response: %v", response)
	if err != nil {
		return err
	}
	if !response["asapiSuccess"].(bool) {
		err = errmsgs.Error("create kvstroe instances Failed !!")
		return errmsgs.WrapErrorf(err, " create kvstroe instances Failed !! %s", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))
	log.Printf("begin describe kvstroe instances !!")
	// wait instance status change from Creating to Normal
	stateConf := BuildStateConfByTimes([]string{"Creating"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}), 200)
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapError(err)
	}

	log.Printf("begin update kvstroe instances !!")

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
	configPayType := PayType("")
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "payment_type", "instance_charge_type"); err == nil {
		configPayType = PayType(v.(string))
	} else {
		return err
	}
	if !d.IsNewResource() && (d.HasChange("payment_type")||d.HasChange("instance_charge_type")) && configPayType == PrePaid {
		// for now we just support charge change from PostPaid to PrePaid
		prePaidRequest := r_kvstore.CreateTransformToPrePaidRequest()
		client.InitRpcRequest(*prePaidRequest.RpcRequest)
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.Period = requests.Integer(strconv.Itoa(d.Get("period").(int)))
		prePaidRequest.AutoPay = requests.NewBoolean(true)
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.TransformToPrePaid(prePaidRequest)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				baseResponse := &responses.BaseResponse{}
				err = json.Unmarshal([]byte(raw.(*r_kvstore.TransformToPrePaidResponse).GetHttpContentString()), baseResponse)
				if err == nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(baseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), prePaidRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("instance_charge_type")
		//d.SetPartial("period")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
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
	if d.HasChange("tair_instance_name") || d.HasChange("instance_name") {
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "tair_instance_name", "instance_name"); err == nil {
			request.InstanceName = v.(string)
			update = true
		} else {
			return err
		}
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
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
	connectivity.SetResourceData(d, object.ZoneId, "zone_id", "availability_zone")
	connectivity.SetResourceData(d, object.ChargeType, "payment_type", "instance_charge_type")
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

func buildKVStoreCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	request := r_kvstore.CreateCreateInstanceRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "tair_instance_name", "instance_name"); err == nil {
		request.InstanceName = Trim(v.(string))
	} else {
		return request, err
	}
	request.InstanceType = Trim(d.Get("instance_type").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "payment_type", "instance_charge_type"); err == nil {
		request.ChargeType = Trim(v.(string))
	} else {
		return request, err
	}
	request.Password = Trim(d.Get("password").(string))

	if request.Password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, errmsgs.WrapError(err)
			}
			request.Password = decryptResp.Plaintext
		}
	}

	request.BackupId = Trim(d.Get("backup_id").(string))

	if PayType(request.ChargeType) == PrePaid {
		request.Period = strconv.Itoa(d.Get("period").(int))
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	request.NetworkType = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request.VSwitchId = vswitchId.(string)
		request.NetworkType = strings.ToUpper(string(Vpc))
		request.PrivateIpAddress = Trim(d.Get("private_ip").(string))

		// check vswitchId in zone
		object, err := vpcService.DescribeVSwitch(vswitchId.(string))
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the multi zone %s", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, errmsgs.WrapError(errmsgs.Error("The specified vswitch %s isn't in the zone %s", object.VSwitchId, request.ZoneId))
		}

		request.VpcId = object.VpcId
	}

	request.Token = buildClientToken(request.GetActionName())

	return request, nil
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
