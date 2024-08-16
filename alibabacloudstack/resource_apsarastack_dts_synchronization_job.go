package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDtsSynchronizationJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDtsSynchronizationJobCreate,
		Read:   resourceAlibabacloudStackDtsSynchronizationJobRead,
		Update: resourceAlibabacloudStackDtsSynchronizationJobUpdate,
		Delete: resourceAlibabacloudStackDtsSynchronizationJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dts_job_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dts_job_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"checkpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"xxlarge", "xlarge", "large", "medium", "small"}, false),
			},
			"data_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"data_synchronization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"structure_initialization": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"synchronization_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Forward", "Reverse"}, false),
			},
			"db_list": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserve": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CEN", "DG", "DISTRIBUTED_DMSLOGICDB", "ECS", "EXPRESS", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS"}, false),
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AS400", "DB2", "DMSPOLARDB", "HBASE", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "POSTGRESQL", "TERADATA"}, false),
			},
			"source_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_endpoint_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_endpoint_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ads", "CEN", "DATAHUB", "DG", "ECS", "EXPRESS", "GREENPLUM", "MONGODB", "OTHER", "PolarDB", "POLARDBX20", "RDS"}, false),
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ADB20", "ADB30", "AS400", "DATAHUB", "DB2", "GREENPLUM", "KAFKA", "MONGODB", "MSSQL", "MySQL", "ORACLE", "PolarDB", "POLARDBX20", "POLARDB_O", "PostgreSQL"}, false),
			},
			"destination_endpoint_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination_endpoint_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_endpoint_oracle_sid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delay_notice": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delay_phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delay_rule_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"error_notice": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"error_phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Synchronizing", "Suspending"}, false),
			},
		},
	}
}

func resourceAlibabacloudStackDtsSynchronizationJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "ConfigureDtsJob"
	request := requests.NewCommonRequest()
	request.ApiName = action
	request.Version = "2020-01-01"
	request.Method = "POST"
	request.Product = "Dts"
	request.RegionId = client.RegionId
	request.Domain = client.Domain
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-regionid"] = client.RegionId
	request.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{
		"RegionId":                        client.RegionId,
		"Product":                         "Dts",
		"Version":                         "2020-01-01",
		"Action":                          "ConfigureDtsJob",
		"OrganizationId":                  client.Department,
		"DbList":                          d.Get("db_list").(string),
		"DtsJobName":                      d.Get("dts_job_name").(string),
		"DataInitialization":              fmt.Sprintf("%t", d.Get("data_initialization").(bool)),
		"DataSynchronization":             fmt.Sprintf("%t", d.Get("data_synchronization").(bool)),
		"StructureInitialization":         fmt.Sprintf("%t", d.Get("structure_initialization").(bool)),
		"SynchronizationDirection":        d.Get("synchronization_direction").(string),
		"DestinationEndpointInstanceType": d.Get("destination_endpoint_instance_type").(string),
		"SourceEndpointInstanceType":      d.Get("source_endpoint_instance_type").(string),
		"JobType":                         "SYNC",
	}
	if v, ok := d.GetOk("dts_job_id"); ok {
		request.QueryParams["DtsJobId"] = v.(string)
	}
	if v, ok := d.GetOk("dts_instance_id"); ok {
		request.QueryParams["DtsInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("checkpoint"); ok {
		request.QueryParams["Checkpoint"] = v.(string)
	}
	if v, ok := d.GetOkExists("delay_notice"); ok {
		request.QueryParams["DelayNotice"] = fmt.Sprintf("%t", v.(bool))
	}
	if v, ok := d.GetOk("delay_phone"); ok {
		request.QueryParams["DelayPhone"] = v.(string)
	}
	if v, ok := d.GetOk("delay_rule_time"); ok {
		request.QueryParams["DelayRuleTime"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_database_name"); ok {
		request.QueryParams["DestinationEndpointDataBaseName"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request.QueryParams["DestinationEndpointEngineName"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_ip"); ok {
		request.QueryParams["DestinationEndpointIP"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_instance_id"); ok {
		request.QueryParams["DestinationEndpointInstanceID"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_oracle_sid"); ok {
		request.QueryParams["DestinationEndpointOracleSID"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_password"); ok {
		request.QueryParams["DestinationEndpointPassword"] = v.(string)
	}
	if v, ok := d.GetOk("destination_endpoint_port"); ok {
		request.QueryParams["DestinationEndpointPort"] = v.(string)
	}

	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request.QueryParams["DestinationEndpointRegion"] = v.(string)
	}

	if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
		request.QueryParams["DestinationEndpointUserName"] = v.(string)
	}
	if v, ok := d.GetOkExists("error_notice"); ok {
		request.QueryParams["ErrorNotice"] = fmt.Sprintf("%t", v.(bool))
	}
	if v, ok := d.GetOk("error_phone"); ok {
		request.QueryParams["ErrorPhone"] = v.(string)
	}
	if v, ok := d.GetOk("reserve"); ok {
		request.QueryParams["Reserve"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_database_name"); ok {
		request.QueryParams["SourceEndpointDatabaseName"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request.QueryParams["SourceEndpointEngineName"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_ip"); ok {
		request.QueryParams["SourceEndpointIP"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_instance_id"); ok {
		request.QueryParams["SourceEndpointInstanceID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_oracle_sid"); ok {
		request.QueryParams["SourceEndpointOracleSID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_owner_id"); ok {
		request.QueryParams["SourceEndpointOwnerID"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_password"); ok {
		request.QueryParams["SourceEndpointPassword"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_port"); ok {
		request.QueryParams["SourceEndpointPort"] = v.(string)
	}

	if v, ok := d.GetOk("source_endpoint_region"); ok {
		request.QueryParams["SourceEndpointRegion"] = v.(string)
	}

	if v, ok := d.GetOk("source_endpoint_role"); ok {
		request.QueryParams["SourceEndpointRole"] = v.(string)
	}
	if v, ok := d.GetOk("source_endpoint_user_name"); ok {
		request.QueryParams["SourceEndpointUserName"] = v.(string)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	request.Domain = client.Config.AscmEndpoint
	var err error
	var dtsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dtsClient.Domain = client.Config.AscmEndpoint
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dtsClient.ProcessCommonRequest(request)
		addDebug(action, raw, request, request.QueryParams)
		if err != nil {
			if NeedRetry(err) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_synchronization_job", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DtsJobId"]))
	d.Set("dts_instance_id", response["DtsInstanceId"])
	dtsService := DtsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Synchronizing"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{"InitializeFailed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlibabacloudStackDtsSynchronizationJobUpdate(d, meta)
}
func resourceAlibabacloudStackDtsSynchronizationJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSynchronizationJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dts_synchronization_job dtsService.DescribeDtsSynchronizationJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	migrationModeObj := object["MigrationMode"].(map[string]interface{})
	destinationEndpointObj := object["DestinationEndpoint"].(map[string]interface{})
	sourceEndpointObj := object["SourceEndpoint"].(map[string]interface{})
	d.Set("checkpoint", fmt.Sprint(formatInt(object["Checkpoint"])))
	d.Set("data_initialization", migrationModeObj["DataInitialization"])
	d.Set("data_synchronization", migrationModeObj["DataSynchronization"])
	d.Set("db_list", object["DbObject"])
	d.Set("destination_endpoint_database_name", destinationEndpointObj["DatabaseName"])
	d.Set("destination_endpoint_engine_name", destinationEndpointObj["EngineName"])
	d.Set("destination_endpoint_ip", destinationEndpointObj["Ip"])
	d.Set("destination_endpoint_instance_id", destinationEndpointObj["InstanceID"])
	d.Set("destination_endpoint_instance_type", destinationEndpointObj["InstanceType"])
	d.Set("destination_endpoint_oracle_sid", destinationEndpointObj["OracleSID"])
	d.Set("destination_endpoint_port", destinationEndpointObj["Port"])
	d.Set("destination_endpoint_region", destinationEndpointObj["Region"])
	d.Set("destination_endpoint_user_name", destinationEndpointObj["UserName"])
	d.Set("dts_instance_id", object["DtsInstanceID"])
	d.Set("dts_job_name", object["DtsJobName"])
	d.Set("source_endpoint_database_name", sourceEndpointObj["DatabaseName"])
	d.Set("source_endpoint_engine_name", sourceEndpointObj["EngineName"])
	d.Set("source_endpoint_ip", sourceEndpointObj["Ip"])
	d.Set("source_endpoint_instance_id", sourceEndpointObj["InstanceID"])
	d.Set("source_endpoint_instance_type", sourceEndpointObj["InstanceType"])
	d.Set("source_endpoint_oracle_sid", sourceEndpointObj["OracleSID"])
	d.Set("source_endpoint_owner_id", sourceEndpointObj["AliyunUid"])
	d.Set("source_endpoint_port", sourceEndpointObj["Port"])
	d.Set("source_endpoint_region", sourceEndpointObj["Region"])
	d.Set("source_endpoint_role", sourceEndpointObj["RoleName"])
	d.Set("source_endpoint_user_name", sourceEndpointObj["UserName"])
	d.Set("status", object["Status"])
	d.Set("structure_initialization", migrationModeObj["StructureInitialization"])
	d.Set("synchronization_direction", object["SynchronizationDirection"])

	return nil
}
func resourceAlibabacloudStackDtsSynchronizationJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	d.Partial(true)
	var err error
	var dtsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dtsClient.Domain = client.Config.AscmEndpoint
	if err != nil {
		return WrapError(err)
	}
	request := requests.NewCommonRequest()
	request.Version = "2020-01-01"
	request.Method = "POST"
	request.Product = "Dts"
	request.RegionId = client.RegionId
	request.Domain = client.Domain
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-regionid"] = client.RegionId
	request.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{
		"RegionId":       client.RegionId,
		"product":        "Dts",
		"OrganizationId": client.Department,
		"DtsJobId":       d.Id(),
	}
	update := false
	if !d.IsNewResource() && d.HasChange("dts_job_name") {
		update = true
		request.QueryParams["DtsJobName"] = d.Get("dts_job_name").(string)
	}
	if update {
		action := "ModifyDtsJobName"
		request.ApiName = action
		response := make(map[string]interface{})
		request.ApiName = action
		request.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(request)
			addDebug(action, raw, request, request.QueryParams)
			if err != nil {
				if NeedRetry(err) {
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
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_synchronization_job", action, AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		// d.SetPartial("dts_job_name")
	}
	modifyDtsJobPasswordReq := requests.NewCommonRequest()
	modifyDtsJobPasswordReq.Version = "2020-01-01"
	modifyDtsJobPasswordReq.Method = "POST"
	modifyDtsJobPasswordReq.Product = "Dts"
	modifyDtsJobPasswordReq.RegionId = client.RegionId
	modifyDtsJobPasswordReq.Domain = client.Domain
	modifyDtsJobPasswordReq.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	modifyDtsJobPasswordReq.Headers["x-acs-regionid"] = client.RegionId
	modifyDtsJobPasswordReq.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
	modifyDtsJobPasswordReq.Headers["x-acs-organizationid"] = client.Department
	modifyDtsJobPasswordReq.Headers["x-acs-content-type"] = "application/json"
	modifyDtsJobPasswordReq.Headers["Content-type"] = "application/json"
	if strings.ToLower(client.Config.Protocol) == "https" {
		modifyDtsJobPasswordReq.Scheme = "https"
	} else {
		modifyDtsJobPasswordReq.Scheme = "http"
	}
	modifyDtsJobPasswordReq.QueryParams = map[string]string{
		"DtsJobId":       d.Id(),
		"RegionId":       client.RegionId,
		"Product":        "Dts",
		"Version":        "2020-01-01",
		"OrganizationId": client.Department,
		"ResourceId":     client.ResourceGroup,
	}
	if !d.IsNewResource() && d.HasChange("source_endpoint_password") {

		modifyDtsJobPasswordReq.QueryParams["Endpoint"] = "src"
		if v, ok := d.GetOk("source_endpoint_password"); ok {
			modifyDtsJobPasswordReq.QueryParams["Password"] = v.(string)
		}
		if v, ok := d.GetOk("source_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq.QueryParams["UserName"] = v.(string)
		}

		action := "ModifyDtsJobPassword"
		modifyDtsJobPasswordReq.ApiName = action
		response := make(map[string]interface{})
		modifyDtsJobPasswordReq.ApiName = action
		modifyDtsJobPasswordReq.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(modifyDtsJobPasswordReq)
			addDebug(action, raw, modifyDtsJobPasswordReq, modifyDtsJobPasswordReq.QueryParams)
			if err != nil {
				if NeedRetry(err) {
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
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_synchronization_job", action, AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		// d.SetPartial("source_endpoint_password")
		// d.SetPartial("source_endpoint_user_name")

		target := d.Get("status").(string)
		err = resourceAlibabacloudStackDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	if !d.IsNewResource() && d.HasChange("destination_endpoint_password") {

		modifyDtsJobPasswordReq.QueryParams["Endpoint"] = "src"
		if v, ok := d.GetOk("destination_endpoint_password"); ok {
			modifyDtsJobPasswordReq.QueryParams["Password"] = v.(string)
		}
		if v, ok := d.GetOk("destination_endpoint_user_name"); ok {
			modifyDtsJobPasswordReq.QueryParams["UserName"] = v.(string)
		}

		action := "ModifyDtsJobPassword"
		modifyDtsJobPasswordReq.ApiName = action
		response := make(map[string]interface{})
		modifyDtsJobPasswordReq.ApiName = action
		modifyDtsJobPasswordReq.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(modifyDtsJobPasswordReq)
			addDebug(action, raw, modifyDtsJobPasswordReq, modifyDtsJobPasswordReq.QueryParams)
			if err != nil {
				if NeedRetry(err) {
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
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_synchronization_job", action, AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		// d.SetPartial("source_endpoint_password")
		// d.SetPartial("source_endpoint_user_name")

		target := d.Get("status").(string)
		err = resourceAlibabacloudStackDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	update = false
	transferInstanceClassReq := requests.NewCommonRequest()
	transferInstanceClassReq.Version = "2020-01-01"
	transferInstanceClassReq.Method = "POST"
	transferInstanceClassReq.Product = "Dts"
	transferInstanceClassReq.RegionId = client.RegionId
	transferInstanceClassReq.Domain = client.Domain
	transferInstanceClassReq.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	transferInstanceClassReq.Headers["x-acs-regionid"] = client.RegionId
	transferInstanceClassReq.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
	transferInstanceClassReq.Headers["x-acs-organizationid"] = client.Department
	transferInstanceClassReq.Headers["x-acs-content-type"] = "application/json"
	transferInstanceClassReq.Headers["Content-type"] = "application/json"
	if strings.ToLower(client.Config.Protocol) == "https" {
		transferInstanceClassReq.Scheme = "https"
	} else {
		transferInstanceClassReq.Scheme = "http"
	}
	transferInstanceClassReq.QueryParams = map[string]string{
		"DtsJobId":       d.Id(),
		"RegionId":       client.RegionId,
		"Product":        "Dts",
		"Version":        "2020-01-01",
		"OrganizationId": client.Department,
		"ResourceId":     client.ResourceGroup,
		"OrderType":      "UPGRADE",
	}
	if !d.IsNewResource() && d.HasChange("instance_class") {
		if v, ok := d.GetOk("instance_class"); ok {
			transferInstanceClassReq.QueryParams["InstanceClass"] = v.(string)
		}
		update = true
	}

	if update {
		action := "TransferInstanceClass"
		transferInstanceClassReq.ApiName = action
		response := make(map[string]interface{})
		transferInstanceClassReq.QueryParams["Action"] = action
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := dtsClient.ProcessCommonRequest(transferInstanceClassReq)
			addDebug(action, raw, transferInstanceClassReq, transferInstanceClassReq.QueryParams)
			if err != nil {
				if NeedRetry(err) {
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
			return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_synchronization_job", action, AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	if !d.IsNewResource() && d.HasChange("status") {
		target := d.Get("status").(string)
		err := resourceAlibabacloudStackDtsSynchronizationJobStatusFlow(d, meta, target)
		if err != nil {
			return WrapError(Error(FailedToReachTargetStatus, d.Get("status")))
		}
	}

	d.Partial(false)
	return resourceAlibabacloudStackDtsSynchronizationJobRead(d, meta)
}
func resourceAlibabacloudStackDtsSynchronizationJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "ResetDtsJob"
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}

	if v, ok := d.GetOk("dts_instance_id"); ok {
		request["DtsInstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	request["product"] = "Dts"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func resourceAlibabacloudStackDtsSynchronizationJobStatusFlow(d *schema.ResourceData, meta interface{}, target string) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSynchronizationJob(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if object["Status"].(string) != target {
		var err error
		var dtsClient *sts.Client
		if client.Config.SecurityToken == "" {
			dtsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
		} else {
			dtsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
		}
		dtsClient.Domain = client.Config.AscmEndpoint
		if err != nil {
			return WrapError(err)
		}
		if target == "Synchronizing" || target == "Suspending" {
			action := "StartDtsJob"
			request := requests.NewCommonRequest()
			request.Version = "2020-01-01"
			request.Method = "POST"
			request.Product = "Dts"
			request.RegionId = client.RegionId
			request.Domain = client.Domain
			request.ApiName = action
			request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
			request.Headers["x-acs-regionid"] = client.RegionId
			request.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
			request.Headers["x-acs-organizationid"] = client.Department
			request.Headers["x-acs-content-type"] = "application/json"
			request.Headers["Content-type"] = "application/json"
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.QueryParams = map[string]string{
				"DtsJobId":       d.Id(),
				"RegionId":       client.RegionId,
				"Product":        "Dts",
				"Version":        "2020-01-01",
				"Action":         action,
				"OrganizationId": client.Department,
				"ResourceId":     client.ResourceGroup,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			response := make(map[string]interface{})
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				raw, err := dtsClient.ProcessCommonRequest(request)
				addDebug(action, raw, request, request.QueryParams)
				if err != nil {
					if NeedRetry(err) {
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
				return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, AlibabacloudStackSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			stateConf := BuildStateConf([]string{}, []string{"Synchronizing"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{"InitializeFailed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		if target == "Suspending" {
			action := "SuspendDtsJob"
			request := requests.NewCommonRequest()
			request.Version = "2020-01-01"
			request.Method = "POST"
			request.Product = "Dts"
			request.RegionId = client.RegionId
			request.Domain = client.Domain
			request.ApiName = action
			request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
			request.Headers["x-acs-regionid"] = client.RegionId
			request.Headers["x-acs-resourcegroupid"] = client.ResourceGroup
			request.Headers["x-acs-organizationid"] = client.Department
			request.Headers["x-acs-content-type"] = "application/json"
			request.Headers["Content-type"] = "application/json"
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.QueryParams = map[string]string{
				"DtsJobId":       d.Id(),
				"RegionId":       client.RegionId,
				"Product":        "Dts",
				"Version":        "2020-01-01",
				"Action":         action,
				"OrganizationId": client.Department,
				"ResourceId":     client.ResourceGroup,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			response := make(map[string]interface{})
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				raw, err := dtsClient.ProcessCommonRequest(request)
				addDebug(action, raw, request, request.QueryParams)
				if err != nil {
					if NeedRetry(err) {
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
				return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dts_subscription_job", action, AlibabacloudStackSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
			stateConf := BuildStateConf([]string{}, []string{"Suspending"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dtsService.DtsSynchronizationJobStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		// d.SetPartial("status")
	}

	return nil
}
