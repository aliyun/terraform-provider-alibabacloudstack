package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
	request := client.NewCommonRequest("POST", "dbs", "2019-03-06", action, "")
	request.Headers["x-acs-caller-sdk-source"] = "Terraform" // 必填，调用来源说明
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-type"] = "application/json"
	request.QueryParams["Period"] = "Year"
	request.QueryParams["UsedTime"] = "1"
	request.QueryParams["ClientToken"] = buildClientToken("CreateBackupPlan")

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
	var err error
	var dbsClient *sts.Client
	if client.Config.SecurityToken == "" {
		dbsClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		dbsClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	dbsClient.Domain = client.Config.Endpoints[connectivity.DBSCode]
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := dbsClient.ProcessCommonRequest(request)
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
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dbs_backup_plan", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)
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

	d.SetId(fmt.Sprint(response["BackupPlanId"]))

	return resourceAlibabacloudStackDbsBackupPlanRead(d, meta)
}

func resourceAlibabacloudStackDbsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dbsService := DbsService{client}
	object, err := dbsService.DescribeDbsBackupPlan(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dbs_backup_plan dbsService.DescribeDbsBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		errmsg := ""
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DescribeDbsBackupPlan", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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

	request["BackupPlanId"] = d.Id()
	if d.HasChange("backup_plan_name") {
		request["BackupPlanName"] = d.Get("backup_plan_name").(string)
	}

	action := "ModifyBackupPlanName"
	request["ClientToken"] = buildClientToken("ModifyBackupPlanName")

	_, err := client.DoTeaRequest("POST", "DBS", "2019-03-06", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return resourceAlibabacloudStackDbsBackupPlanRead(d, meta)
}

func resourceAlibabacloudStackDbsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	// 没有接口
	return nil
}
