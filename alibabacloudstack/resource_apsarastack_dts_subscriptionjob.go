package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDtsSubscriptionJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDtsSubscriptionJobCreate,
		Read:   resourceAlibabacloudStackDtsSubscriptionJobRead,
		Update: resourceAlibabacloudStackDtsSubscriptionJobUpdate,
		Delete: resourceAlibabacloudStackDtsSubscriptionJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"checkpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"compute_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"db_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delay_notice": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delay_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delay_rule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADS", "DB2", "DRDS", "DataHub", "Greenplum", "MSSQL", "MySQL", "PolarDB", "PostgreSQL", "Redis", "Tablestore", "as400", "clickhouse", "kafka", "mongodb", "odps", "oracle", "polardb_o", "polardb_pg", "tidb"}, false),
			},
			"destination_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"error_notice": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"error_phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"large", "medium", "micro", "small", "xlarge", "xxlarge"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"payment_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"payment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"reserve": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "Oracle"}, false),
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CEN", "DRDS", "ECS", "Express", "LocalInstance", "PolarDB", "RDS", "dg"}, false),
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Abnormal", "Downgrade", "Locked", "Normal", "NotStarted", "NotStarted", "PreCheckPass", "PrecheckFailed", "Prechecking", "Retrying", "Starting", "Upgrade"}, false),
			},
			"subscription_data_type_ddl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"subscription_data_type_dml": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"subscription_instance_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"classic", "vpc"}, false),
			},
			"subscription_instance_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscription_instance_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"bidirectional", "oneway"}, false),
			},
			"synchronization_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlibabacloudStackDtsSubscriptionJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "CreateDtsInstance"
	request := client.NewCommonRequest("POST", "Dts", "2020-01-01", action, "")
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-regionid"] = client.RegionId
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	request.QueryParams["AutoPay"] = "false"
	request.QueryParams["AutoStart"] = "true"
	if v, ok := d.GetOk("compute_unit"); ok {
		request.QueryParams["ComputeUnit"] = v.(string)
	}
	if v, ok := d.GetOk("database_count"); ok {
		request.QueryParams["DatabaseCount"] = fmt.Sprintf("%v", v.(int))
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request.QueryParams["DestinationEndpointEngineName"] = v.(string)
	}
	if v, ok := d.GetOk("destination_region"); ok {
		request.QueryParams["DestinationRegion"] = v.(string)
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request.QueryParams["InstanceClass"] = v.(string)
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request.QueryParams["PayType"] = v.(string)
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request.QueryParams["Period"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request.QueryParams["SourceEndpointEngineName"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request.QueryParams["SourceRegion"] = v.(string)
	}
	if v, ok := d.GetOk("sync_architecture"); ok {
		request.QueryParams["SyncArchitecture"] = v.(string)
	}
	request.QueryParams["Type"] = "SUBSCRIBE"
	if v, ok := d.GetOk("payment_duration"); ok {
		request.QueryParams["UsedTime"] = fmt.Sprintf("%v", v.(int))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	response := make(map[string]interface{})
	request.Domain = client.Config.Endpoints[connectivity.DTSCode]
	var err error
	var dtsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dtsClient.Domain = client.Config.Endpoints[connectivity.DTSCode]
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dtsClient.ProcessCommonRequest(request)
		addDebug(action, raw, request, request.QueryParams)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if err != nil {
			errmsg := ""
			if raw != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.Set("dts_instance_id", response["InstanceId"])
	configureSubscriptionReq := client.NewCommonRequest("POST", "Dts", "2020-01-01", "ConfigureSubscription", "")
	configureSubscriptionReq.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	configureSubscriptionReq.QueryParams["DbList"] = d.Get("db_list").(string)
	configureSubscriptionReq.QueryParams["SubscriptionInstanceNetworkType"] = d.Get("subscription_instance_network_type").(string)
	configureSubscriptionReq.QueryParams["DtsInstanceId"] = d.Get("dts_instance_id").(string)
	configureSubscriptionReq.QueryParams["DtsJobName"] = d.Get("dts_job_name").(string)
	configureSubscriptionReq.Headers["x-acs-regionid"] = client.RegionId
	configureSubscriptionReq.Headers["x-acs-content-type"] = "application/json"
	configureSubscriptionReq.Headers["Content-type"] = "application/json"
	if v, ok := d.GetOk("checkpoint"); ok {
		configureSubscriptionReq.QueryParams["Checkpoint"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointDatabaseName"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointEngineName"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointIP"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointInstanceID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_instance_type"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointInstanceType"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointOracleSID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointOwnerID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointPassword"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointPort"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointRegion"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointRole"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointUserName"] = v.(string)
	}
	if v, ok := d.GetOkExists("subscription_data_type_ddl"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionDataTypeDDL"] = v.(string)
	}
	if v, ok := d.GetOkExists("subscription_data_type_dml"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionDataTypeDML"] = v.(string)
	}
	if v, ok := d.GetOk("subscription_instance_vpc_id"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionInstanceVPCId"] = v.(string)
	}
	if v, ok := d.GetOk("subscription_instance_vswitch_id"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionInstanceVSwitchId"] = v.(string)
	}
	if v, ok := d.GetOkExists("delay_notice"); ok {
		configureSubscriptionReq.QueryParams["DelayNotice"] = fmt.Sprintf("%t", v.(bool))
	}
	if v, ok := d.GetOk("delay_phone"); ok {
		configureSubscriptionReq.QueryParams["DelayPhone"] = v.(string)
	}
	if v, ok := d.GetOk("delay_rule_time"); ok {
		configureSubscriptionReq.QueryParams["DelayRuleTime"] = v.(string)
	}
	if v, ok := d.GetOkExists("error_notice"); ok {
		configureSubscriptionReq.QueryParams["ErrorNotice"] = fmt.Sprintf("%t", v.(bool))
	}
	if v, ok := d.GetOk("error_phone"); ok {
		configureSubscriptionReq.QueryParams["ErrorPhone"] = v.(string)
	}
	if v, ok := d.GetOk("reserve"); ok {
		configureSubscriptionReq.QueryParams["Reserve"] = v.(string)
	}
	configureSubscriptionRsp := make(map[string]interface{})
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dtsClient.ProcessCommonRequest(configureSubscriptionReq)
		addDebug(action, raw, configureSubscriptionReq, configureSubscriptionReq.QueryParams)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if err != nil {
			errmsg := ""
			if raw != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		err = json.Unmarshal(raw.GetHttpContentBytes(), &configureSubscriptionRsp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if fmt.Sprint(configureSubscriptionRsp["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "ConfigureSubscription", configureSubscriptionRsp))
	}
	fmt.Printf("==================================  %s", fmt.Sprint(configureSubscriptionRsp["DtsJobId"]))
	d.SetId(fmt.Sprint(configureSubscriptionRsp["DtsJobId"]))
	return resourceAlibabacloudStackDtsSubscriptionJobUpdate(d, meta)
}

func resourceAlibabacloudStackDtsSubscriptionJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSubscriptionJob(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dts_subscription_job dtsService.DescribeDtsSubscriptionJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("checkpoint", fmt.Sprint(formatInt(object["Checkpoint"])))
	d.Set("db_list", object["DbObject"])
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	d.Set("payment_type", convertDtsPaymentTypeResponse(object["PayType"]))
	d.Set("source_endpoint_database_name", object["SourceEndpoint"].(map[string]interface{})["DatabaseName"])
	d.Set("source_endpoint_engine_name", object["SourceEndpoint"].(map[string]interface{})["EngineName"])
	d.Set("source_endpoint_ip", object["SourceEndpoint"].(map[string]interface{})["Ip"])
	d.Set("source_endpoint_instance_id", object["SourceEndpoint"].(map[string]interface{})["InstanceID"])
	d.Set("source_endpoint_instance_type", object["SourceEndpoint"].(map[string]interface{})["InstanceType"])
	d.Set("source_endpoint_oracle_sid", object["SourceEndpoint"].(map[string]interface{})["OracleSID"])
	d.Set("source_endpoint_owner_id", object["SourceEndpoint"].(map[string]interface{})["AliyunUid"])
	d.Set("source_endpoint_port", object["SourceEndpoint"].(map[string]interface{})["Port"])
	d.Set("source_endpoint_region", object["SourceEndpoint"].(map[string]interface{})["Region"])
	d.Set("source_endpoint_role", object["SourceEndpoint"].(map[string]interface{})["RoleName"])
	d.Set("source_endpoint_user_name", object["SourceEndpoint"].(map[string]interface{})["UserName"])
	d.Set("status", object["Status"])
	d.Set("subscription_data_type_ddl", object["SubscriptionDataType"].(map[string]interface{})["Ddl"])
	d.Set("subscription_data_type_dml", object["SubscriptionDataType"].(map[string]interface{})["Dml"])

	var jsonData map[string]interface{}
	json.Unmarshal([]byte(object["Reserved"].(string)), &jsonData)
	if jsonData["netType"] != nil {
		d.Set("subscription_instance_network_type", strings.ToLower(jsonData["netType"].(string)))
	}
	d.Set("subscription_instance_vpc_id", jsonData["vpcId"])
	d.Set("subscription_instance_vswitch_id", jsonData["vswitchId"])
	listTagResourcesObject, err := dtsService.ListTagResources(object["DtsInstanceID"].(string), "ALIYUN::DTS::INSTANCE")
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAlibabacloudStackDtsSubscriptionJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	d.Partial(true)
	if d.HasChange("tags") {
		if err := dtsService.SetResourceTags(d, "ALIYUN::DTS::INSTANCE"); err != nil {
			return errmsgs.WrapError(err)
		}
		//d.SetPartial("tags")
	}
	var err error
	var dtsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dtsClient.Domain = client.Config.Endpoints[connectivity.DTSCode]
	if err != nil {
		return errmsgs.WrapError(err)
	}
	update := false
	request := client.NewCommonRequest("POST", "Dts", "2020-01-01", "", "")
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-regionid"] = client.RegionId
	request.Headers["x-acs-content-type"] = "application/json"
	request.QueryParams["DtsJobId"] = d.Id()
	if d.HasChange("dts_job_name") {
		update = true
		if v, ok := d.GetOk("dts_job_name"); ok {
			request.QueryParams["DtsJobName"] = v.(string)
		}
	}
	request.Domain = client.Config.Endpoints[connectivity.DTSCode]
	if update {
		action := "ModifyDtsJobName"
		response := make(map[string]interface{})
		request.ApiName = action
		request.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(request)
			addDebug(action, raw, request, request.QueryParams)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if err != nil {
				errmsg := ""
				if raw != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		//d.SetPartial("dts_job_name")
	}
	update = false
	modifyDtsJobPasswordReq := client.NewCommonRequest("POST", "Dts", "2020-01-01", "ModifyDtsJobPassword", "")
	modifyDtsJobPasswordReq.QueryParams["DtsJobId"] = d.Id()
	modifyDtsJobPasswordReq.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	modifyDtsJobPasswordReq.Headers["x-acs-regionid"] = client.RegionId
	modifyDtsJobPasswordReq.Headers["x-acs-content-type"] = "application/json"
	modifyDtsJobPasswordReq.Headers["Content-type"] = "application/json"
	modifyDtsJobPasswordReq.QueryParams["Endpoint"] = "src"
	if !d.IsNewResource() && d.HasChange("source_endpoint_password") {
		update = true
		if v, ok := d.GetOk("source_endpoint_password"); ok {
			modifyDtsJobPasswordReq.QueryParams["Password"] = v.(string)
		}
		if v, ok := d.GetOk("source_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq.QueryParams["UserName"] = v.(string)
		}
	}
	if update {
		action := "ModifyDtsJobPassword"
		response := make(map[string]interface{})
		modifyDtsJobPasswordReq.ApiName = action
		modifyDtsJobPasswordReq.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(modifyDtsJobPasswordReq)
			addDebug(action, raw, modifyDtsJobPasswordReq, modifyDtsJobPasswordReq.QueryParams)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if err != nil {
				errmsg := ""
				if raw != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		target := d.Get("status").(string)
		err = resourceAlibabacloudStackDtsSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, d.Get("status")))
		}
	}

	if !d.IsNewResource() && d.HasChange("status") {
		target := d.Get("status").(string)
		err := resourceAlibabacloudStackDtsSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, d.Get("status")))
		}
	}

	update = false
	configureSubscriptionReq := client.NewCommonRequest("POST", "Dts", "2020-01-01", "ConfigureSubscription", "")
	configureSubscriptionReq.QueryParams["DtsJobId"] = d.Id()
	configureSubscriptionReq.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	configureSubscriptionReq.Headers["x-acs-regionid"] = client.RegionId
	configureSubscriptionReq.Headers["x-acs-content-type"] = "application/json"
	configureSubscriptionReq.Headers["Content-type"] = "application/json"
	if d.IsNewResource() {
		update = true
	}
	if d.HasChange("db_list") {
		update = true
	}
	if v, ok := d.GetOk("db_list"); ok {
		configureSubscriptionReq.QueryParams["DbList"] = v.(string)
	}
	if v, ok := d.GetOk("dts_job_name"); ok {
		configureSubscriptionReq.QueryParams["DtsJobName"] = v.(string)
	}
	if d.HasChange("subscription_instance_network_type") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_network_type"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionInstanceNetworkType"] = v.(string)
	}
	if d.HasChange("checkpoint") {
		update = true
	}
	if v, ok := d.GetOk("checkpoint"); ok {
		configureSubscriptionReq.QueryParams["Checkpoint"] = v.(string)
	}
	if d.HasChange("source_endpoint_database_name") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointDatabaseName"] = v.(string)
	}
	if !d.IsNewResource() && d.HasChange("source_endpoint_engine_name") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointEngineName"] = v.(string)
	}
	if d.HasChange("source_endpoint_ip") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointIP"] = v.(string)
	}
	if d.HasChange("source_endpoint_instance_id") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointInstanceID"] = v.(string)
	}
	if d.HasChange("source_endpoint_instance_type") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_instance_type"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointInstanceType"] = v.(string)
	}
	if d.HasChange("source_endpoint_oracle_sid") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointOracleSID"] = v.(string)
	}
	if d.HasChange("source_endpoint_owner_id") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointOwnerID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointPassword"] = v.(string)
	}
	if d.HasChange("source_endpoint_port") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointPort"] = v.(string)
	}
	if d.HasChange("source_endpoint_region") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_region"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointRegion"] = v.(string)
	}
	if d.HasChange("source_endpoint_role") {
		update = true
	}
	if v, ok := d.GetOk("source_endpoint_role"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointRole"] = v.(string)
	}

	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		configureSubscriptionReq.QueryParams["SourceEndpointUserName"] = v.(string)
	}
	if d.HasChange("subscription_data_type_ddl") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("subscription_data_type_ddl"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionDataTypeDDL"] = v.(string)
	}
	if d.HasChange("subscription_data_type_dml") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("subscription_data_type_dml"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionDataTypeDML"] = v.(string)
	}
	if d.HasChange("subscription_instance_vpc_id") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_vpc_id"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionInstanceVPCId"] = v.(string)
	}
	if d.HasChange("subscription_instance_vswitch_id") {
		update = true
	}
	if v, ok := d.GetOk("subscription_instance_vswitch_id"); ok {
		configureSubscriptionReq.QueryParams["SubscriptionInstanceVSwitchId"] = v.(string)
	}
	if update {

		if v, ok := d.GetOkExists("delay_notice"); ok {
			configureSubscriptionReq.QueryParams["DelayNotice"] = fmt.Sprintf("%t", v.(bool))
		}
		if v, ok := d.GetOk("delay_phone"); ok {
			configureSubscriptionReq.QueryParams["DelayPhone"] = v.(string)
		}
		if v, ok := d.GetOk("delay_rule_time"); ok {
			configureSubscriptionReq.QueryParams["DelayRuleTime"] = v.(string)
		}
		if v, ok := d.GetOkExists("error_notice"); ok {
			configureSubscriptionReq.QueryParams["ErrorNotice"] = fmt.Sprintf("%t", v.(bool))
		}
		if v, ok := d.GetOk("error_phone"); ok {
			configureSubscriptionReq.QueryParams["ErrorPhone"] = v.(string)
		}
		if v, ok := d.GetOk("reserve"); ok {
			configureSubscriptionReq.QueryParams["Reserve"] = v.(string)
		}
		action := "ConfigureSubscription"
		response := make(map[string]interface{})
		configureSubscriptionReq.ApiName = action
		configureSubscriptionReq.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(configureSubscriptionReq)
			addDebug(action, raw, configureSubscriptionReq, configureSubscriptionReq.QueryParams)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if err != nil {
				errmsg := ""
				if raw != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		//d.SetPartial("db_list")
		//d.SetPartial("dts_job_name")
		//d.SetPartial("subscription_instance_network_type")
		//d.SetPartial("checkpoint")
		//d.SetPartial("source_endpoint_database_name")
		//d.SetPartial("source_endpoint_engine_name")
		//d.SetPartial("source_endpoint_ip")
		//d.SetPartial("source_endpoint_instance_id")
		//d.SetPartial("source_endpoint_instance_type")
		//d.SetPartial("source_endpoint_oracle_sid")
		//d.SetPartial("source_endpoint_owner_id")
		//d.SetPartial("source_endpoint_password")
		//d.SetPartial("source_endpoint_port")
		//d.SetPartial("source_endpoint_region")
		//d.SetPartial("source_endpoint_role")
		//d.SetPartial("source_endpoint_user_name")
		//d.SetPartial("subscription_data_type_ddl")
		//d.SetPartial("subscription_data_type_dml")
		//d.SetPartial("subscription_instance_vpc_id")
		//d.SetPartial("subscription_instance_vswitch_id")

		target := d.Get("status").(string)
		err = resourceAlibabacloudStackDtsSubscriptionJobStatusFlow(d, meta, target)
		if err != nil {
			return errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, d.Get("status")))
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackDtsSubscriptionJobRead(d, meta)
}

func resourceAlibabacloudStackDtsSubscriptionJobDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok {
		if v.(string) == "Subscription" {
			return nil
		}
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteDtsJob"
	var err error
	var dtsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dtsClient.Domain = client.Config.Endpoints[connectivity.DTSCode]
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := client.NewCommonRequest("POST", "Dts", "2020-01-01", action, "")
	request.QueryParams["DtsJobId"] = d.Id()
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-regionid"] = client.RegionId
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	if v, ok := d.GetOk("dts_instance_id"); ok {
		request.QueryParams["DtsInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("synchronization_direction"); ok {
		request.QueryParams["SynchronizationDirection"] = v.(string)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	response := make(map[string]interface{})
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dtsClient.ProcessCommonRequest(request)
		addDebug(action, raw, request, request.QueryParams)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}

func resourceAlibabacloudStackDtsSubscriptionJobStatusFlow(d *schema.ResourceData, meta interface{}, target string) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSubscriptionJob(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if object["Status"].(string) != target {
		var err error
		var dtsClient *sts.Client
		if client.Config.SecurityToken == "" {
			dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
		} else {
			dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
		}
		dtsClient.Domain = client.Config.Endpoints[connectivity.DTSCode]
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if target == "NotConfigured" {
			action := "ResetDtsJob"
			request := client.NewCommonRequest("POST","Dts","2020-01-01","ResetDtsJob","")
			request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
			request.Headers["x-acs-content-type"] = "application/json"
			request.Headers["Content-type"] = "application/json"
			request.QueryParams["DtsJobId"] = d.Id()
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request.QueryParams["SynchronizationDirection"] = v.(string)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			response := make(map[string]interface{})
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				raw, err := dtsClient.ProcessCommonRequest(request)
				addDebug(action, raw, request, request.QueryParams)
				if err != nil {
					if errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			stateConf := BuildStateConf([]string{}, []string{"NotConfigured"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.DtsSubscriptionJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		if target == "Normal" || (target == "Abnormal" && object["Status"].(string) == "NotStarted") {
			action := "StartDtsJob"
			request := client.NewCommonRequest("POST","Dts","2020-01-01","StartDtsJob","")
			request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
			request.Headers["x-acs-content-type"] = "application/json"
			request.Headers["Content-type"] = "application/json"
			request.QueryParams["DtsJobId"] =       d.Id()
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request.QueryParams["SynchronizationDirection"] = v.(string)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			response := make(map[string]interface{})
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				raw, err := dtsClient.ProcessCommonRequest(request)
				addDebug(action, raw, request, request.QueryParams)
				if err != nil {
					if errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			stateConf := BuildStateConf([]string{}, []string{"Starting", "Normal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, dtsService.DtsSubscriptionJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		if target == "Abnormal" {
			action := "SuspendDtsJob"
			request := client.NewCommonRequest("POST","Dts","2020-01-01",action,"")
			request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
			request.Headers["x-acs-content-type"] = "application/json"
			request.Headers["Content-type"] = "application/json"
			request.QueryParams["DtsJobId"] = d.Id()
			if v, ok := d.GetOk("synchronization_direction"); ok {
				request.QueryParams["SynchronizationDirection"] = v.(string)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			response := make(map[string]interface{})
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				raw, err := dtsClient.ProcessCommonRequest(request)
				addDebug(action, raw, request, request.QueryParams)
				if err != nil {
					if errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			stateConf := BuildStateConf([]string{}, []string{"Abnormal"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, dtsService.DtsSubscriptionJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		//d.SetPartial("status")
	}

	return nil
}

func convertDtsPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
func convertDtsPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
