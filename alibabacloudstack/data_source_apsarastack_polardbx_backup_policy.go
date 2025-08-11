package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackPolardbxBackupPolicys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbxBackupPolicysRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_set_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_plan_begin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remove_log_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cold_data_backup_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local_log_retention_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cold_data_backup_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"force_clean_on_high_space_usage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_way": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_log_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_local_retention_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbxBackupPolicysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbxService := PolardbXService{client}
	db_instance_id := d.Get("db_instance_id").(string)
	object, err := polardbxService.DescribePolarDbXBackupConfig(db_instance_id)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("backup_period", GetBackupPeriod(object.BackupPeriod))
	d.Set("backup_set_retention", object.BackupSetRetention)
	d.Set("backup_plan_begin", object.BackupPlanBegin)
	d.Set("cold_data_backup_interval", object.ColdDataBackupInterval)
	d.Set("remove_log_retention", object.RemoveLogRetention)
	d.Set("local_log_retention_number", object.LocalLogRetentionNumber)
	d.Set("cold_data_backup_retention", object.ColdDataBackupRetention)
	d.Set("force_clean_on_high_space_usage", object.ForceCleanOnHighSpaceUsage)
	d.Set("backup_way", object.BackupWay)
	d.Set("local_log_retention", object.LocalLogRetention)
	d.Set("backup_type", object.BackupType)
	d.Set("log_local_retention_space", object.LogLocalRetentionSpace)
	d.SetId(db_instance_id)
	return nil
}
