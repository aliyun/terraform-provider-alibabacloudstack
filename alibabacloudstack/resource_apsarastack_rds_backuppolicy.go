package alibabacloudstack

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDBBackupPolicyCreate,
		Read:   resourceAlibabacloudStackDBBackupPolicyRead,
		Update: resourceAlibabacloudStackDBBackupPolicyUpdate,
		Delete: resourceAlibabacloudStackDBBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Optional:     true,
				Default:      "02:00Z-03:00Z",
			},

			"backup_retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(7, 730),
			},

			"enable_backup_log": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"log_backup_retention_period": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(7, 730),
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: logRetentionPeriodDiffSuppressFunc,
			},

			"local_log_retention_hours": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(0, 7*24),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"local_log_retention_space": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(5, 50),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"high_space_usage_protection": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Default:          "Enable",
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"log_backup_frequency": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"compress_type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"1", "4", "8"}, false),
				Computed:     true,
				Optional:     true,
			},

			"archive_backup_retention_period": {
				Type:             schema.TypeInt,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: archiveBackupPeriodDiffSuppressFunc,
			},

			"archive_backup_keep_count": {
				Type:             schema.TypeInt,
				ValidateFunc:     validation.IntBetween(1, 31),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},

			"archive_backup_keep_policy": {
				Type:             schema.TypeString,
				ValidateFunc:     validation.StringInSlice([]string{"ByMonth", "ByWeek", "KeepAll"}, false),
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: enableBackupLogDiffSuppressFunc,
			},
		},
	}
}

func resourceAlibabacloudStackDBBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("instance_id").(string))
	return resourceAlibabacloudStackDBBackupPolicyUpdate(d, meta)
}

func resourceAlibabacloudStackDBBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("instance_id", d.Id())
	d.Set("preferred_backup_time", object.PreferredBackupTime)
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("backup_retention_period", object.BackupRetentionPeriod)
	d.Set("enable_backup_log", object.EnableBackupLog == "1")
	d.Set("log_backup_retention_period", object.LogBackupRetentionPeriod)
	d.Set("local_log_retention_hours", object.LocalLogRetentionHours)
	d.Set("local_log_retention_space", object.LocalLogRetentionSpace)
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if instance.Engine == "SQLServer" {
		d.Set("high_space_usage_protection", "Enable")
	} else {
		d.Set("high_space_usage_protection", object.HighSpaceUsageProtection)
	}
	d.Set("log_backup_frequency", object.LogBackupFrequency)
	d.Set("compress_type", object.CompressType)
	d.Set("archive_backup_retention_period", object.ArchiveBackupRetentionPeriod)
	d.Set("archive_backup_keep_count", object.ArchiveBackupKeepCount)
	d.Set("archive_backup_keep_policy", object.ArchiveBackupKeepPolicy)
	return nil
}

func resourceAlibabacloudStackDBBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}

	updateForData := false
	updateForLog := false
	if d.HasChanges("preferred_backup_period","preferred_backup_time", "backup_retention_period",
			"compress_type","log_backup_frequency", "archive_backup_retention_period", 
			"archive_backup_keep_count", "archive_backup_keep_policy") {
		updateForData = true
	}

	if d.HasChanges("enable_backup_log", "log_backup_retention_period", "local_log_retention_hours",
			"local_log_retention_space", "high_space_usage_protection") {
		updateForLog = true
	}

	if updateForData || updateForLog {
		if err := rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := rdsService.ModifyDBBackupPolicy(d, updateForData, updateForLog); err != nil {
				if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return resourceAlibabacloudStackDBBackupPolicyRead(d, meta)
}

func resourceAlibabacloudStackDBBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	request := rds.CreateModifyBackupPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = "Tuesday,Thursday,Saturday"
	request.BackupRetentionPeriod = "7"
	request.PreferredBackupTime = "02:00Z-03:00Z"
	request.EnableBackupLog = "1"
	instance, err := rdsService.DescribeDBInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if instance.Engine != "SQLServer" {
		request.LogBackupRetentionPeriod = "7"
	}
	if instance.Engine == "MySQL" && instance.DBInstanceStorageType == "local_ssd" {
		request.ArchiveBackupRetentionPeriod = "0"
		request.ArchiveBackupKeepCount = "1"
		request.ArchiveBackupKeepPolicy = "ByMonth"
	}

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifyBackupPolicy(request)
	})
	if err != nil {
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*rds.ModifyBackupPolicyResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return rdsService.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium)
}
