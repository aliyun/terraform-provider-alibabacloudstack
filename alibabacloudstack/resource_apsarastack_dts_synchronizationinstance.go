package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDtsSynchronizationInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compute_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sync_architecture": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bidirectional", "oneway"}, false),
			},
			"destination_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"destination_endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_endpoint_engine_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"source_endpoint_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"xxlarge", "xlarge", "large", "medium", "small"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "Subscription"}, false),
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dts_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDtsSynchronizationInstanceCreate, 
		resourceAlibabacloudStackDtsSynchronizationInstanceRead, 
		nil, 
		resourceAlibabacloudStackDtsSynchronizationInstanceDelete)
	return resource
}

func resourceAlibabacloudStackDtsSynchronizationInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateDtsInstance"
	request := client.NewCommonRequest("POST", "Dts", "2020-01-01", action, "")
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	mergeMaps(request.QueryParams, map[string]string{
		"AutoPay":       "false",
		"AutoStart":     "true",
		"InstanceClass": "small",
		"Type":          "SYNC",
	})
	if v, ok := d.GetOk("compute_unit"); ok {
		request.QueryParams["ComputeUnit"] = v.(string)
	}
	if v, ok := d.GetOk("database_count"); ok {
		request.QueryParams["DatabaseCount"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request.QueryParams["DestinationEndpointEngineName"] = v.(string)
	}

	if v, ok := d.GetOk("destination_endpoint_region"); ok {
		request.QueryParams["DestinationRegion"] = v.(string)
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request.QueryParams["PayType"] = v.(string)
	}
	if v, ok := d.GetOk("payment_duration_unit"); ok {
		request.QueryParams["Period"] = v.(string)
	}
	if v, ok := d.GetOk("quantity"); ok {
		request.QueryParams["Quantity"] = strconv.Itoa(v.(int))
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
	if v, ok := d.GetOk("payment_duration"); ok {
		request.QueryParams["UsedTime"] = strconv.Itoa(v.(int))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
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
			errmsg := ""
			if raw != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dts_synchronization_instance", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
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
	d.SetId(fmt.Sprint(response["InstanceId"]))

	return nil
}

func resourceAlibabacloudStackDtsSynchronizationInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsSynchronizationInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dts_synchronization_instance dtsService.DescribeDtsSynchronizationInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("instance_class", object["SynchronizationJobClass"])
	d.Set("payment_type", convertDtsSyncPaymentTypeResponse(object["PayType"]))
	d.Set("dts_job_id", object["SynchronizationJobId"])
	return nil
}

func resourceAlibabacloudStackDtsSynchronizationInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("payment_type"); ok {
		if v.(string) == "Subscription" {
			log.Printf("[WARN] Cannot destroy resource: alibabacloudstack_dts_synchronization_job because it's s. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}

	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteSynchronizationJob"
	request := make(map[string]interface{})
	request["SynchronizationJobId"] = d.Id()
	_, err := client.DoTeaRequest("POST", "Dts", "2020-01-01", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidJobId"}) {
			return nil
		}
		return err
	}
	return nil
}

func convertDtsSyncPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}

func convertDtsSyncPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
