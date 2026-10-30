package alibabacloudstack

// Product PolarDB Resource Cluster Backup (cluster edition, polardb / 2017-08-01)
import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// polardbClusterBackupSuccess is the BackupStatus a finished backup set carries.
const polardbClusterBackupSuccess = "Success"

// polardbClusterBackupPending stands for every status the cluster edition reports while
// the backup job has not finished, so that the waiter keeps polling instead of failing
// on an intermediate status it has never seen.
const polardbClusterBackupPending = "Pending"

func resourceAlibabacloudStackPolardbClusterBackup() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backup_job_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backup_id": {
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
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbClusterBackupCreate,
		resourceAlibabacloudStackPolardbClusterBackupRead,
		nil,
		resourceAlibabacloudStackPolardbClusterBackupDelete)
	return resource
}

func resourceAlibabacloudStackPolardbClusterBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardb_service := PolardbService{client}
	db_cluster_id := d.Get("db_cluster_id").(string)

	// A backup set gets its BackupId only when the backup job finishes, and that id is not
	// the BackupJobId answered by CreateBackup, so remember the backup sets which already
	// exist in order to tell the new one apart afterwards.
	before, err := polardb_service.DescribePolardbClusterBackups(db_cluster_id)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_cluster_backup", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	before_ids := make([]string, 0)
	for _, v := range before {
		before_ids = append(before_ids, v.BackupId)
	}

	// Automatic backup jobs of the cluster policy run on their own schedule; a manual
	// CreateBackup issued while such a job is running is rejected with BackupJobExists,
	// so keep retrying until the running job finishes, within the create timeout.
	var bresponse *responses.CommonResponse
	retryErr := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "polardb", "2017-08-01", "CreateBackup", "")
		request.QueryParams["DBClusterId"] = db_cluster_id
		response, rerr := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), response, request, request.QueryParams)
		if rerr != nil {
			if response == nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(rerr, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			if strings.Contains(rerr.Error(), "BackupJobExists") || strings.Contains(errmsg, "BackupJobExists") {
				return resource.RetryableError(errmsgs.WrapErrorf(rerr, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(rerr, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_cluster_backup", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		bresponse = response
		return nil
	})
	if retryErr != nil {
		return errmsgs.WrapError(retryErr)
	}

	PolardbClusterCreatebackupResponseObj := PolardbClusterCreatebackupResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbClusterCreatebackupResponseObj)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_cluster_backup", "CreateBackup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.Set("backup_job_id", PolardbClusterCreatebackupResponseObj.BackupJobId)

	stateConf := BuildStateConf([]string{polardbClusterBackupPending}, []string{polardbClusterBackupSuccess}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polardb_service.PolardbClusterBackupStateRefreshFunc(db_cluster_id, before_ids, []string{"Failed"}))
	object, err := stateConf.WaitForState()
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, db_cluster_id)
	}
	backup := object.(*PolardbClusterBackupData)
	d.SetId(fmt.Sprintf("%s:%s", db_cluster_id, backup.BackupId))
	return resourceAlibabacloudStackPolardbClusterBackupRead(d, meta)
}

func resourceAlibabacloudStackPolardbClusterBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardb_service := PolardbService{client}
	data, err := polardb_service.DescribePolardbClusterBackup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_polardb_cluster_backup polardb_service.DescribePolardbClusterBackup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_cluster_backup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	d.Set("db_cluster_id", data.DBClusterId)
	d.Set("backup_id", data.BackupId)
	d.Set("backup_method", data.BackupMethod)
	d.Set("backup_mode", data.BackupMode)
	d.Set("backup_type", data.BackupType)
	d.Set("backup_status", data.BackupStatus)
	d.Set("backups_level", data.BackupsLevel)
	d.Set("store_status", data.StoreStatus)
	d.Set("backup_set_size", data.BackupSetSize)
	d.Set("is_avail", data.IsAvail)
	d.Set("consistent_time", data.ConsistentTime)
	d.Set("backup_start_time", data.BackupStartTime)
	d.Set("backup_end_time", data.BackupEndTime)
	return nil
}

// DeleteBackup of the cluster edition removes a backup set which Terraform was asked to
// create, and the primary/standby edition backup resource leaves the backup sets behind
// for the same reason; keeping the backup set is what the retention policy is for.
func resourceAlibabacloudStackPolardbClusterBackupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

type PolardbClusterCreatebackupResponse struct {
	RequestId   string `json:"RequestId"`
	BackupJobId int    `json:"BackupJobId"`
}
