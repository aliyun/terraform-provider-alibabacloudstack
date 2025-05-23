package alibabacloudstack

// Generated By apsara-orchestration-generator
// Product POLARDB Resouce BackupPolicy
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbBackuppolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{

			"backup_log": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"backup_policy_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"backup_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"compress_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enable_backup_log": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"high_space_usage_protection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_log_retention_hours": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_log_retention_space": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"log_backup_frequency": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"log_backup_local_retention_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"log_backup_retention_period": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"preferred_backup_period": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"released_keep_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbBackuppolicyCreate,
		resourceAlibabacloudStackPolardbBackuppolicyRead, resourceAlibabacloudStackPolardbBackuppolicyUpdate, resourceAlibabacloudStackPolardbBackuppolicyDelete)
	return resource
}

func resourceAlibabacloudStackPolardbBackuppolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyBackupPolicy", "")
	PolardbModifybackuppolicyResponse := PolardbModifybackuppolicyResponse{}

	if v, ok := d.GetOk("backup_log"); ok {
		request.QueryParams["BackupLog"] = v.(string)
	}

	if v, ok := d.GetOk("backup_policy_mode"); ok {
		request.QueryParams["BackupPolicyMode"] = v.(string)
	}

	if v, ok := d.GetOk("backup_retention_period"); ok {
		request.QueryParams["BackupRetentionPeriod"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("compress_type"); ok {
		request.QueryParams["CompressType"] = v.(string)
	}

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return fmt.Errorf("DBInstanceId is required")
	}

	if v, ok := d.GetOk("enable_backup_log"); ok {
		request.QueryParams["EnableBackupLog"] = v.(string)
	}

	if v, ok := d.GetOk("high_space_usage_protection"); ok {
		request.QueryParams["HighSpaceUsageProtection"] = v.(string)
	}

	if v, ok := d.GetOk("local_log_retention_hours"); ok {
		request.QueryParams["LocalLogRetentionHours"] = v.(string)
	}

	if v, ok := d.GetOk("local_log_retention_space"); ok {
		request.QueryParams["LocalLogRetentionSpace"] = v.(string)
	}

	if v, ok := d.GetOk("log_backup_frequency"); ok {
		request.QueryParams["LogBackupFrequency"] = v.(string)
	}

	if v, ok := d.GetOk("log_backup_local_retention_number"); ok {
		request.QueryParams["LogBackupLocalRetentionNumber"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("log_backup_retention_period"); ok {
		request.QueryParams["LogBackupRetentionPeriod"] = v.(string)
	}

	if v, ok := d.GetOk("preferred_backup_period"); ok {
		request.QueryParams["PreferredBackupPeriod"] = v.(string)
	}

	if v, ok := d.GetOk("preferred_backup_time"); ok {
		request.QueryParams["PreferredBackupTime"] = v.(string)
	}

	if v, ok := d.GetOk("released_keep_policy"); ok {
		request.QueryParams["ReleasedKeepPolicy"] = v.(string)
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_backup_policy", "ModifyBackupPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifybackuppolicyResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_backup_policy", "ModifyBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	db_instance_id := d.Get("db_instance_id").(string)

	d.SetId(fmt.Sprintf("%s", db_instance_id))
	return nil
}

func resourceAlibabacloudStackPolardbBackuppolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChanges("backup_log", "backup_policy_mode", "backup_retention_period", "compress_type", "enable_backup_log", "high_space_usage_protection", "local_log_retention_hours", "local_log_retention_space", "log_backup_frequency", "log_backup_local_retention_number", "log_backup_retention_period", "preferred_backup_period", "preferred_backup_time", "released_keep_policy") {
		request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyBackupPolicy", "")
		PolardbModifybackuppolicyResponse := PolardbModifybackuppolicyResponse{}

		if v, ok := d.GetOk("backup_log"); ok {
			request.QueryParams["BackupLog"] = v.(string)
		}

		if v, ok := d.GetOk("backup_policy_mode"); ok {
			request.QueryParams["BackupPolicyMode"] = v.(string)
		}

		if v, ok := d.GetOk("backup_retention_period"); ok {
			request.QueryParams["BackupRetentionPeriod"] = strconv.Itoa(v.(int))
		}

		if v, ok := d.GetOk("compress_type"); ok {
			request.QueryParams["CompressType"] = v.(string)
		}

		if v, ok := d.GetOk("db_instance_id"); ok {
			request.QueryParams["DBInstanceId"] = v.(string)
		} else {
			return fmt.Errorf("DBInstanceId is required")
		}

		if v, ok := d.GetOk("enable_backup_log"); ok {
			request.QueryParams["EnableBackupLog"] = v.(string)
		}

		if v, ok := d.GetOk("high_space_usage_protection"); ok {
			request.QueryParams["HighSpaceUsageProtection"] = v.(string)
		}

		if v, ok := d.GetOk("local_log_retention_hours"); ok {
			request.QueryParams["LocalLogRetentionHours"] = v.(string)
		}

		if v, ok := d.GetOk("local_log_retention_space"); ok {
			request.QueryParams["LocalLogRetentionSpace"] = v.(string)
		}

		if v, ok := d.GetOk("log_backup_frequency"); ok {
			request.QueryParams["LogBackupFrequency"] = v.(string)
		}

		if v, ok := d.GetOk("log_backup_local_retention_number"); ok {
			request.QueryParams["LogBackupLocalRetentionNumber"] = strconv.Itoa(v.(int))
		}

		if v, ok := d.GetOk("log_backup_retention_period"); ok {
			request.QueryParams["LogBackupRetentionPeriod"] = v.(string)
		}

		if v, ok := d.GetOk("preferred_backup_period"); ok {
			request.QueryParams["PreferredBackupPeriod"] = v.(string)
		}

		if v, ok := d.GetOk("preferred_backup_time"); ok {
			request.QueryParams["PreferredBackupTime"] = v.(string)
		}

		if v, ok := d.GetOk("released_keep_policy"); ok {
			request.QueryParams["ReleasedKeepPolicy"] = v.(string)
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_polardb_backup_policy", "ModifyBackupPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifybackuppolicyResponse)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_backup_policy", "ModifyBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
		}

	}

	return nil
}

func resourceAlibabacloudStackPolardbBackuppolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbbackup_policyservice := PolardbService{client}
	response, err := polardbbackup_policyservice.DoPolardbDescribebackuppolicyRequest(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardb_backuppolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	data := response
	
	d.Set("db_instance_id", d.Id())
	
	d.Set("backup_log", data.BackupLog)

	d.Set("backup_retention_period", data.BackupRetentionPeriod)

	d.Set("compress_type", data.CompressType)

	d.Set("enable_backup_log", data.EnableBackupLog)

	d.Set("high_space_usage_protection", data.HighSpaceUsageProtection)

	d.Set("local_log_retention_hours", data.LocalLogRetentionHours)

	d.Set("local_log_retention_space", data.LocalLogRetentionSpace)

	d.Set("log_backup_frequency", data.LogBackupFrequency)

	d.Set("log_backup_local_retention_number", data.LogBackupLocalRetentionNumber)

	d.Set("log_backup_retention_period", data.LogBackupRetentionPeriod)

	d.Set("preferred_backup_period", data.PreferredBackupPeriod)

	d.Set("preferred_backup_time", data.PreferredBackupTime)

	d.Set("released_keep_policy", data.ReleasedKeepPolicy)

	return nil
}

func resourceAlibabacloudStackPolardbBackuppolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

type PolardbModifybackuppolicyResponse struct {
	RequestId                     string `json:"RequestId"`
	DBInstanceID                  string `json:"DBInstanceID"`
	EnableBackupLog               string `json:"EnableBackupLog"`
	LocalLogRetentionHours        int    `json:"LocalLogRetentionHours"`
	LocalLogRetentionSpace        string `json:"LocalLogRetentionSpace"`
	HighSpaceUsageProtection      string `json:"HighSpaceUsageProtection"`
	CompressType                  string `json:"CompressType"`
	LogBackupLocalRetentionNumber int    `json:"LogBackupLocalRetentionNumber"`
}