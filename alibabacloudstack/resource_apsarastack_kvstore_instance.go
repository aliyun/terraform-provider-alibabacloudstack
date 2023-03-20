package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
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
				Default:  KVStore2Dot8,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Optional:     true,
				Default:      PostPaid,
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
			"architecture_type": {
				Type:     schema.TypeString,
				Required: true,
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
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request["Product"] = "R-kvstore"
	request["product"] = "R-kvstore"
	request["OrganizationId"] = client.Department
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateInstance")
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v.(string)
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		request["CpuType"] = v.(string)
	}
	if v, ok := d.GetOk("architecture_type"); ok {
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

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request["ChargeType"] = v.(string)
	}
	if v, ok := d.GetOk("password"); ok {
		request["Password"] = v.(string)
	}

	if request["Password"] == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
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

	if zone, ok := d.GetOk("availability_zone"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}

	request["NetworkType"] = strings.ToUpper(string(Classic))
	if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
		request["VSwitchId"] = vswitchId.(string)
		request["NetworkType"] = strings.ToUpper(string(Vpc))
		request["PrivateIpAddress"] = Trim(d.Get("private_ip").(string))

		// check vswitchId in zone
		object, err := vpcService.DescribeVSwitch(vswitchId.(string))
		if err != nil {
			return WrapError(err)
		}

		if request["ZoneId"] == "" {
			request["ZoneId"] = object.ZoneId
		}

		request["VpcId"] = object.VpcId
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	//addDebug(action, response, request)
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_data_works_folder", action, ApsaraStackSdkGoERROR)
	//}

	d.SetId(fmt.Sprint(response["InstanceId"]))
	//client := meta.(*connectivity.ApsaraStackClient)

	//request, err := buildKVStoreCreateRequest(d, meta)
	//if err != nil {
	//	return WrapError(err)
	//}
	//
	//raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
	//	return rkvClient.CreateInstance(request)
	//})
	//
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kvstore_instance", request.GetActionName(), ApsaraStackSdkGoERROR)
	//}
	//addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//response, _ := raw.(*r_kvstore.CreateInstanceResponse)
	//d.SetId(response.InstanceId)

	// wait instance status change from Creating to Normal
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Normal"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
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
				return WrapError(err)
			}
		}

		//d.SetPartial("parameters")
	}

	if d.HasChange("security_ips") {
		// wait instance status is Normal before modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateModifySecurityIpsRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.SecurityIpGroupName = "default"
		request.InstanceId = d.Id()
		if len(d.Get("security_ips").(*schema.Set).List()) > 0 {
			request.SecurityIps = strings.Join(expandStringList(d.Get("security_ips").(*schema.Set).List())[:], COMMA_SEPARATED)
		} else {
			return WrapError(Error("Security ips cannot be empty"))
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifySecurityIps(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("security_ips")
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("vpc_auth_mode") {
		if vswitchId, ok := d.GetOk("vswitch_id"); ok && vswitchId.(string) != "" {
			// vpc_auth_mode works only if the network type is VPC
			instanceType := d.Get("instance_type").(string)
			if string(KVStoreRedis) == instanceType {
				// wait instance status is Normal before modifying
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapError(err)
				}

				request := r_kvstore.CreateModifyInstanceVpcAuthModeRequest()
				request.RegionId = client.RegionId
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
				request.InstanceId = d.Id()
				request.VpcAuthMode = d.Get("vpc_auth_mode").(string)

				raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
					return rkvClient.ModifyInstanceVpcAuthMode(request)
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				//d.SetPartial("vpc_auth_mode")

				// The auth mode take some time to be effective, so wait to ensure the state !
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapError(err)
				}
			}
		}
	}
	configPayType := PayType(d.Get("instance_charge_type").(string))
	if !d.IsNewResource() && d.HasChange("instance_charge_type") && configPayType == PrePaid {
		// for now we just support charge change from PostPaid to PrePaid
		prePaidRequest := r_kvstore.CreateTransformToPrePaidRequest()
		prePaidRequest.RegionId = client.RegionId
		prePaidRequest.Headers = map[string]string{"RegionId": client.RegionId}
		prePaidRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		prePaidRequest.InstanceId = d.Id()
		prePaidRequest.Period = requests.Integer(strconv.Itoa(d.Get("period").(int)))

		prePaidRequest.AutoPay = requests.NewBoolean(true)
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.TransformToPrePaid(prePaidRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), prePaidRequest.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(prePaidRequest.GetActionName(), raw, prePaidRequest.RpcRequest, prePaidRequest)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		//d.SetPartial("instance_charge_type")
		//d.SetPartial("period")
	}

	if d.HasChange("maintain_start_time") || d.HasChange("maintain_end_time") {
		request := r_kvstore.CreateModifyInstanceMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.InstanceId = d.Id()
		request.MaintainStartTime = d.Get("maintain_start_time").(string)
		request.MaintainEndTime = d.Get("maintain_end_time").(string)

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.ModifyInstanceMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
			return WrapError(err)
		}

		request := r_kvstore.CreateModifyInstanceSpecRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		request.EffectiveTime = "Immediately"

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.ModifyInstanceSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"MissingRedisUsedmemoryUnsupportPerfItem"}) {
					time.Sleep(time.Duration(5) * time.Second)
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
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		// There needs more time to sync instance class update
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if object.InstanceClass != request.InstanceClass {
				return resource.RetryableError(Error("Waitting for instance class is changed timeout. Expect instance class %s, got %s.",
					object.InstanceClass, request.InstanceClass))
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}

		//d.SetPartial("instance_class")
	}

	request := r_kvstore.CreateModifyInstanceAttributeRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = d.Id()
	update := false
	if d.HasChange("instance_name") {
		request.InstanceName = d.Get("instance_name").(string)
		update = true
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
				return WrapError(err)
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
			return WrapError(err)
		}
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance status is Normal after modifying
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapError(err)
		}
		//d.SetPartial("instance_name")
		//d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlibabacloudStackKVStoreInstanceRead(d, meta)
}

func resourceAlibabacloudStackKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_name", object.InstanceName)
	d.Set("instance_class", object.InstanceClass)
	d.Set("availability_zone", object.ZoneId)
	d.Set("instance_charge_type", object.ChargeType)
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
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = d.Id()

		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.DescribeInstanceAutoRenewalAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		period, err := computePeriodByUnit(object.CreateTime, object.EndTime, d.Get("period").(int), "Month")
		if err != nil {
			return WrapError(err)
		}
		d.Set("period", period)
	}
	//refresh parameters
	if err = refreshParameters(d, meta); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlibabacloudStackKVStoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	object, err := kvstoreService.DescribeKVstoreInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if PayType(object.ChargeType) == PrePaid {
		return WrapError(Error("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically"))
	}
	request := r_kvstore.CreateDeleteInstanceRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = d.Id()

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteInstance(request)
	})

	if err != nil {
		if !IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{"Creating", "Active", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, kvstoreService.RdsKvstoreInstanceStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return WrapError(err)
}

func buildKVStoreCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	request := r_kvstore.CreateCreateInstanceRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceName = Trim(d.Get("instance_name").(string))

	request.InstanceType = Trim(d.Get("instance_type").(string))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	request.ChargeType = Trim(d.Get("instance_charge_type").(string))

	request.Password = d.Get("password").(string)
	if request.Password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return request, WrapError(err)
			}
			request.Password = decryptResp.Plaintext
		}
	}

	request.BackupId = Trim(d.Get("backup_id").(string))

	if PayType(request.ChargeType) == PrePaid {
		request.Period = strconv.Itoa(d.Get("period").(int))
	}

	if zone, ok := d.GetOk("availability_zone"); ok && Trim(zone.(string)) != "" {
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
			return nil, WrapError(err)
		}

		if request.ZoneId == "" {
			request.ZoneId = object.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(object.ZoneId)[len(object.ZoneId)-1])) {
				return nil, WrapError(Error("The specified vswitch %s isn't in the multi zone %s", object.VSwitchId, request.ZoneId))
			}
		} else if request.ZoneId != object.ZoneId {
			return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s", object.VSwitchId, request.ZoneId))
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
		return WrapError(err)
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
