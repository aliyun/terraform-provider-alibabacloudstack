package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDbsBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDbsBackupPlanCreate,
		Read:   resourceApsaraStackDbsBackupPlanRead,
		Update: resourceApsaraStackDbsBackupPlanUpdate,
		Delete: resourceApsaraStackDbsBackupPlanDelete,
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

func resourceApsaraStackDbsBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	conn, err := client.NewDbsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Region":   client.RegionId,
		"Period":   "Year",
		"UsedTime": 1,
	}
	if v, ok := d.GetOk("backup_method"); ok {
		request["BackupMethod"] = v.(string)
	}

	if v, ok := d.GetOk("database_type"); ok {
		request["DatabaseType"] = v.(string)
	}

	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v.(string)
	}

	if v, ok := d.GetOk("database_region"); ok {
		request["DatabaseRegion"] = v.(string)
	}

	if v, ok := d.GetOk("storage_region"); ok {
		request["StorageRegion"] = v.(string)
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v.(string)
	}

	if v, ok := d.GetOk("from_app"); ok {
		request["FromApp"] = v.(string)
	}

	request["ClientToken"] = buildClientToken("CreateBackupPlan")
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_dbs_backup_plan", action, ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BackupPlanId"]))

	return resourceApsaraStackDbsBackupPlanRead(d, meta)
}

func resourceApsaraStackDbsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dbsService := DbsService{client}
	object, err := dbsService.DescribeDbsBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_dbs_backup_plan dbsService.DescribeDbsBackupPlan Failed!!! %s", err)
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

func resourceApsaraStackDbsBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := make(map[string]interface{})
	var response map[string]interface{}

	request["BackupPlanId"] = d.Id()

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return resourceApsaraStackDbsBackupPlanRead(d, meta)
}

func resourceApsaraStackDbsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	// 没有接口
	return nil
}
