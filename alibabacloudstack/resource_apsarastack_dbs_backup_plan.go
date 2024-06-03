package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDbsBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDbsBackupPlanCreate,
		Read:   resourceAlibabacloudStackDbsBackupPlanRead,
		Update: resourceAlibabacloudStackDbsBackupPlanUpdate,
		Delete: resourceAlibabacloudStackDbsBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backup_plan_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"backup_method": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"logical", "physical"}, false),
			},
			"database_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "MSSQL", "Oracle", "MongoDB", "Redis"}, false),
			},
			"instance_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"small", "large"}, false),
			},
			"database_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDS", "PolarDB", "DDS", "Kvstore", "Other"}, false),
			},
			"from_app": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OpenApi",
			},
			"backup_plan_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackDbsBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := requests.NewCommonRequest()
	request.ApiName = action
	request.Version = "2019-03-06"
	request.Method = "POST"
	request.Product = "dbs"
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
		"Product":        "dbs",
		"Version":        "2019-03-06",
		"Action":         action,
		"OrganizationId": client.Department,
		"Period":         "Year",
		"UsedTime":       "1",
		"ClientToken":    buildClientToken("CreateBackupPlan"),
	}
	if v, ok := d.GetOk("backup_method"); ok {
		request.QueryParams["BackupMethod"] = v.(string)
	}

	if v, ok := d.GetOk("database_type"); ok {
		request.QueryParams["DatabaseType"] = v.(string)
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request.QueryParams["InstanceClass"] = v.(string)
	}

	if v, ok := d.GetOk("database_region"); ok {
		request.QueryParams["DatabaseRegion"] = v.(string)
	}

	if v, ok := d.GetOk("storage_region"); ok {
		request.QueryParams["StorageRegion"] = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.QueryParams["InstanceType"] = v.(string)
	}

	if v, ok := d.GetOk("from_app"); ok {
		request.QueryParams["FromApp"] = v.(string)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	request.Domain = client.Config.AscmEndpoint
	var err error
	var dbsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dbsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dbsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dbsClient.Domain = client.Config.AscmEndpoint
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dbsClient.ProcessCommonRequest(request)
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dbs_backup_plan", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BackupPlanId"]))

	return resourceAlibabacloudStackDbsBackupPlanRead(d, meta)
}

func resourceAlibabacloudStackDbsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dbsService := DbsService{client}
	object, err := dbsService.DescribeDbsBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dbs_backup_plan dbsService.DescribeDbsBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if err != nil {
		return WrapError(err)
	}

	for key, value := range object {
		fmt.Println(key, value)
	}

	d.Set("backup_plan_id", d.Id())
	d.Set("backup_plan_name", object["BackupPlanName"].(string))

	return nil
}

func resourceAlibabacloudStackDbsBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := make(map[string]interface{})
	var response map[string]interface{}

	request["BackupPlanId"] = d.Id()
	request["product"] = "dbs"
	request["Product"] = "dbs"
	if d.HasChange("backup_plan_name") {
		request["BackupPlanName"] = d.Get("backup_plan_name").(string)
	}

	action := "ModifyBackupPlanName"
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("ModifyBackupPlanName")
	conn, err := client.NewDbsClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return resourceAlibabacloudStackDbsBackupPlanRead(d, meta)
}

func resourceAlibabacloudStackDbsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	// 没有接口
	return nil
}
