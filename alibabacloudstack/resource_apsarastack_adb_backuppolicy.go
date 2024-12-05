package alibabacloudstack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAdbBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAdbBackupPolicyCreate,
		Read:   resourceAlibabacloudStackAdbBackupPolicyRead,
		Update: resourceAlibabacloudStackAdbBackupPolicyUpdate,
		Delete: resourceAlibabacloudStackAdbBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"preferred_backup_period": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"preferred_backup_time": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(BACKUP_TIME, false),
				Required:     true,
			},
			"backup_retention_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackAdbBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("db_cluster_id").(string))
	return resourceAlibabacloudStackAdbBackupPolicyUpdate(d, meta)
}

func resourceAlibabacloudStackAdbBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbBackupPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("db_cluster_id", d.Id())
	d.Set("backup_retention_period", strconv.Itoa(object.BackupRetentionPeriod))
	d.Set("preferred_backup_period", strings.Split(object.PreferredBackupPeriod, ","))
	d.Set("preferred_backup_time", object.PreferredBackupTime)

	return nil
}

func resourceAlibabacloudStackAdbBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}

	if d.HasChange("preferred_backup_period") || d.HasChange("preferred_backup_time") {
		periodList := expandStringList(d.Get("preferred_backup_period").(*schema.Set).List())
		preferredBackupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
		preferredBackupTime := d.Get("preferred_backup_time").(string)

		// wait instance running before modifying
		if err := adbService.WaitForCluster(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := adbService.ModifyAdbBackupPolicy(d.Id(), preferredBackupTime, preferredBackupPeriod); err != nil {
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

	return resourceAlibabacloudStackAdbBackupPolicyRead(d, meta)
}

func resourceAlibabacloudStackAdbBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// In case of a delete we are resetting to default values which is Tuesday,Friday each 1am-2am
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := adb.CreateModifyBackupPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = d.Id()
	request.PreferredBackupTime = "01:00Z-02:00Z"
	request.PreferredBackupPeriod = "Tuesday,Friday"

	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.ModifyBackupPolicy(request)
	})
	response, ok := raw.(*adb.ModifyBackupPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
