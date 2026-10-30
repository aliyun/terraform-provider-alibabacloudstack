package alibabacloudstack

// Product PolarDB DataSource Cluster Backups (cluster edition, polardb / 2017-08-01)
import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackPolardbClusterBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterBackupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"db_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backups_level": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"store_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_set_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"is_avail": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"consistent_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbClusterBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardb_service := PolardbService{client}
	now := time.Now().UTC()
	// DescribeBackups of the cluster edition makes both ends of the time window mandatory,
	// so default to the last 7 days when the caller does not narrow the search down.
	// The bounds are minute-granular, and a backup set completing seconds after its
	// own truncated minute would fall outside an end bound of "now", so the default
	// end is padded by a day; the service accepts window ends in the future.
	start_time := now.AddDate(0, 0, -7).Format(polardbClusterBackupTimeLayout)
	end_time := now.AddDate(0, 0, 1).Format(polardbClusterBackupTimeLayout)
	if v, ok := d.GetOk("start_time"); ok {
		start_time = v.(string)
	}
	if v, ok := d.GetOk("end_time"); ok {
		end_time = v.(string)
	}
	backups, err := polardb_service.DoPolardbClusterDescribebackupsRequest(d.Get("db_cluster_id").(string), start_time, end_time)
	if err != nil {
		return err
	}
	idsMap := getIdsStringFilter(d)
	datas := make([]interface{}, 0)
	ids := make([]string, 0)
	for _, data := range backups {
		if len(idsMap) > 0 {
			if _, exist := idsMap[data.BackupId]; !exist {
				continue
			}
		}
		i := map[string]interface{}{
			"id":                data.BackupId,
			"backup_id":         data.BackupId,
			"db_cluster_id":     data.DBClusterId,
			"backup_method":     data.BackupMethod,
			"backup_mode":       data.BackupMode,
			"backup_type":       data.BackupType,
			"backup_status":     data.BackupStatus,
			"backups_level":     data.BackupsLevel,
			"store_status":      data.StoreStatus,
			"backup_set_size":   data.BackupSetSize,
			"is_avail":          data.IsAvail,
			"consistent_time":   data.ConsistentTime,
			"backup_start_time": data.BackupStartTime,
			"backup_end_time":   data.BackupEndTime,
		}
		log.Printf("[DEBUG] %s: %#v", "backups", i)
		datas = append(datas, i)
		ids = append(ids, data.BackupId)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("backups", datas); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
