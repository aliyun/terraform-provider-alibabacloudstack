package alibabacloudstack

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackPolardbxBackupPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_set_retention": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(7, 730),
			},
			"backup_plan_begin": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remove_log_retention": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(7, 730),
			},
			"cold_data_backup_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local_log_retention_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(6, 100),
			},
			"cold_data_backup_retention": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"force_clean_on_high_space_usage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"backup_way": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "P",
			},
			"local_log_retention": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"log_local_retention_space": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
		},
	}
	// XXX: 逻辑特殊，不建议合并
	setResourceFunc(resource, resourceAlibabacloudStackPolardbxBackupPolicyCreate,
		resourceAlibabacloudStackPolardbxBackupPolicyRead, resourceAlibabacloudStackPolardbxBackupPolicyUpdate, resourceAlibabacloudStackPolardbxBackupPolicyDelete)
	return resource
}

func resourceAlibabacloudStackPolardbxBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("db_instance_id").(string))
	return resourceAlibabacloudStackPolardbxBackupPolicyUpdate(d, meta)
}

func resourceAlibabacloudStackPolardbxBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("backup_period", "backup_set_retention",
		"backup_plan_begin", "remove_log_retention", "local_log_retention_number",
		"force_clean_on_high_space_usage", "local_log_retention", "log_local_retention_space") {
		request := client.NewCommonRequest("POST", "polardbx", "2020-02-02", "UpdateBackupPolicy", "")
		request.QueryParams["IsEnabled"] = "1"
		request.QueryParams["BackupType"] = d.Get("backup_type").(string)
		request.QueryParams["BackupWay"] = d.Get("backup_way").(string)
		request.QueryParams["DBInstanceName"] = d.Get("db_instance_id").(string)
		request.QueryParams["BackupSetRetention"] = fmt.Sprint(d.Get("backup_set_retention").(int))
		request.QueryParams["BackupPlanBegin"] = d.Get("backup_plan_begin").(string)
		request.QueryParams["BackupPeriod"] = SetBackupPeriod(d.Get("backup_period").(string))
		request.QueryParams["RemoveLogRetention"] = fmt.Sprint(d.Get("remove_log_retention").(int))
		request.QueryParams["LocalLogRetentionNumber"] = fmt.Sprint(d.Get("local_log_retention_number").(int))
		request.QueryParams["ForceCleanOnHighSpaceUsage"] = fmt.Sprint(d.Get("force_clean_on_high_space_usage").(int))
		request.QueryParams["LocalLogRetention"] = fmt.Sprint(d.Get("local_log_retention").(int))
		request.QueryParams["LogLocalRetentionSpace"] = fmt.Sprint(d.Get("log_local_retention_space").(int))
		// {"BackupType":"0","BackupSetRetention":30,"RemoveLogRetention":30,"BackupPeriod":"1001000","BackupWay":"P","BackupPlanBegin":"03:00Z","IsEnabled":1,"LogLocalRetentionSpace":30,"LocalLogRetention":7,"ForceCleanOnHighSpaceUsage":1,"LocalLogRetentionNumber":60,"DBInstanceName":"pxc-unr6vlauoszezq","RegionId":"cn-wulan-env17e-d01"}
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardbx_backup_policy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	return nil
}

func resourceAlibabacloudStackPolardbxBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbxService := PolardbXService{client}
	object, err := polardbxService.DescribePolarDbXBackupConfig(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("db_instance_id", d.Id())
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
	return nil
}

func resourceAlibabacloudStackPolardbxBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func SetBackupPeriod(backupPeriod string) string {
	standard := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	backupPeriods := strings.Split(backupPeriod, ",")
	result := make([]string, 0)
	for _, v := range standard {
		if slices.Contains(backupPeriods, v) {
			result = append(result, "1")
		} else {
			result = append(result, "0")

		}
	}
	return strings.Join(result, "")
}

func GetBackupPeriod(backupPeriodCode string) string {
	standard := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	result := make([]string, 0)
	for i, v := range backupPeriodCode {
		if v == '1' {
			result = append(result, standard[i])
		}
	}
	return strings.Join(result, ",")
}
