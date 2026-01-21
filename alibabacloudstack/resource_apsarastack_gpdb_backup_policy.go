package alibabacloudstack

import (
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackGpdbBackupPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable_recovery_point": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"recovery_point_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackGpdbBackupPolicyCreate, resourceAlibabacloudStackGpdbBackupPolicyRead, resourceAlibabacloudStackGpdbBackupPolicyUpdate, resourceAlibabacloudStackGpdbBackupPolicyDelete)
	return resource
}

func resourceAlibabacloudStackGpdbBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	// No create API exists for this resource, only set the resource ID
	dbInstanceId := d.Get("db_instance_id").(string)
	d.SetId(dbInstanceId)
	return nil
}

func resourceAlibabacloudStackGpdbBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gpdbService := GpdbService{client}

	// Query the backup policy using only the resource ID as required
	object, err := gpdbService.DescribeBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("db_instance_id", d.Id())

	if v, ok := object["PreferredBackupTime"]; ok {
		d.Set("preferred_backup_time", v)
	}

	if v, ok := object["PreferredBackupPeriod"]; ok {
		d.Set("preferred_backup_period", v)
	}

	if v, ok := object["BackupRetentionPeriod"]; ok {
		d.Set("backup_retention_period", v)
	}

	if v, ok := object["EnableRecoveryPoint"]; ok {
		d.Set("enable_recovery_point", v)
	}

	if v, ok := object["RecoveryPointPeriod"]; ok {
		period, err := strconv.Atoi(v.(string))
		if err == nil {
			d.Set("recovery_point_period", period)
		}
	}
	return nil
}

func resourceAlibabacloudStackGpdbBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("preferred_backup_time", "preferred_backup_period", "backup_retention_period", "enable_recovery_point", "recovery_point_period") {
		request := map[string]interface{}{
			"DBInstanceId":          d.Id(),
			"PreferredBackupTime":   d.Get("preferred_backup_time").(string),
			"PreferredBackupPeriod": d.Get("preferred_backup_period").(string),
			"BackupRetentionPeriod": d.Get("backup_retention_period").(int),
			"EnableRecoveryPoint":   d.Get("enable_recovery_point").(bool),
		}
		if d.Get("enable_recovery_point").(bool) {
			request["RecoveryPointPeriod"] = d.Get("recovery_point_period").(int)
		}
		_, err := client.DoTeaRequest("POST", "gpdb", "2016-05-03", "ModifyBackupPolicy", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_gpdb_backup_policy", "ModifyBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackGpdbBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
