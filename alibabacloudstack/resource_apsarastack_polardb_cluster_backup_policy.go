package alibabacloudstack

import (
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbClusterBackupPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_level1_backup_frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_level1_backup_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_level1_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_level1_backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"data_level2_backup_retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_backup_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_next_backup_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_level2_backup_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_level2_backup_another_region_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_level2_backup_another_region_retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_backup_another_region_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_backup_another_region_retention_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_backup_log": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterBackupPolicyCreate, resourceAlibabacloudStackPolardbClusterBackupPolicyRead, resourceAlibabacloudStackPolardbClusterBackupPolicyUpdate, nil)

	return resource
}

func resourceAlibabacloudStackPolardbClusterBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	dbClusterId := d.Get("db_cluster_id").(string)
	d.SetId(dbClusterId)
	return resourceAlibabacloudStackPolardbClusterBackupPolicyRead(d, meta)
}

func resourceAlibabacloudStackPolardbClusterBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}
	object, err := polardbService.DescribePolardbClusterBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_cluster_backup_policy polardbService.DescribePolardbClusterBackupPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())

	d.Set("data_level1_backup_frequency", object["DataLevel1BackupFrequency"])
	d.Set("data_level1_backup_period", object["DataLevel1BackupPeriod"])
	d.Set("data_level1_backup_time", object["DataLevel1BackupTime"])

	d.Set("data_level1_backup_retention_period", object["DataLevel1BackupRetentionPeriod"])
	d.Set("data_level2_backup_retention_period", object["DataLevel2BackupRetentionPeriod"])

	d.Set("backup_retention_policy_on_cluster_deletion", object["BackupRetentionPolicyOnClusterDeletion"])
	d.Set("preferred_backup_time", object["PreferredBackupTime"])
	d.Set("preferred_backup_period", object["PreferredBackupPeriod"])
	d.Set("backup_retention_period", object["BackupRetentionPeriod"])

	d.Set("backup_frequency", object["BackupFrequency"])
	d.Set("preferred_next_backup_time", object["PreferredNextBackupTime"])
	d.Set("data_level2_backup_period", object["DataLevel2BackupPeriod"])
	d.Set("data_level2_backup_another_region_region", object["DataLevel2BackupAnotherRegionRegion"])
	d.Set("data_level2_backup_another_region_retention_period", object["DataLevel2BackupAnotherRegionRetentionPeriod"])
	logObject, err := polardbService.DescribePolardbClusterLogBackupPolicy(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("log_backup_retention_period", logObject["LogBackupRetentionPeriod"])
	d.Set("log_backup_another_region_region", logObject["LogBackupAnotherRegionRegion"])
	d.Set("log_backup_another_region_retention_period", logObject["LogBackupAnotherRegionRetentionPeriod"])
	d.Set("enable_backup_log", logObject["LogBackupAnotherRegionRetentionPeriod"])

	return nil
}

func resourceAlibabacloudStackPolardbClusterBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// if d.IsNewResource() {
	// 	return nil
	// }

	if d.HasChanges(
		"data_level1_backup_period",
		"data_level1_backup_time",
		"data_level1_backup_retention_period",
	) {
		reqQuery := map[string]interface{}{
			"DBClusterId":                            d.Id(),
			"DataLevel1BackupFrequency":              "Normal",
			"BackupRetentionPolicyOnClusterDeletion": "NONE",
			"DataLevel1BackupPeriod":                 d.Get("data_level1_backup_period"),
			"DataLevel1BackupTime":                   d.Get("data_level1_backup_time"),
			"DataLevel1BackupRetentionPeriod":        d.Get("data_level1_backup_retention_period"),
			"DataLevel2BackupRetentionPeriod":        0,
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyBackupPolicy", "", nil, nil, reqQuery); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ModifyBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.HasChange("log_backup_retention_period") {
		reqQuery := map[string]interface{}{
			"DBClusterId":              d.Get("db_cluster_id"),
			"LogBackupRetentionPeriod": d.Get("log_backup_retention_period"),
		}

		if _, err := client.DoTeaRequest("POST", "polardb", "2017-08-01", "ModifyLogBackupPolicy", "", nil, nil, reqQuery); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ModifyLogBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}
