package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackPolardbBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbBackupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"db_instance_id": {
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
						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"slave_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_instance_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_db_names": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"store_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_scale": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_status": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbbackup_policyservice := PolardbService{client}
	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeBackups", "")
	PolardbDescribebackupsResponseObj := &PolardbDescribebackupsResponse{}
	request.QueryParams["DBInstanceId"] = d.Get("db_instance_id").(string)
	request.QueryParams["PageSize"] = "100"
	now := time.Now().UTC()
	if v, ok := d.GetOk("start_time"); ok {
		request.QueryParams["StartTime"] = v.(string)
	} else {
		sevenDaysAgo := now.AddDate(0, 0, -7)
		request.QueryParams["StartTime"] = sevenDaysAgo.Format("2006-01-02T15:04Z")
	}
	if v, ok := d.GetOk("end_time"); ok {
		request.QueryParams["EndTime"] = v.(string)
	} else {
		request.QueryParams["EndTime"] = now.Format("2006-01-02T15:04Z")
	}
	bresponse, err := polardbbackup_policyservice.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribebackupsResponseObj)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	idsMap := getIdsStringFilter(d)
	datas := make([]interface{}, 0)
	ids := make([]string, 0)
	for _, data := range PolardbDescribebackupsResponseObj.Items.Backup {
		if len(idsMap) > 0 {
			if _, exist := idsMap[fmt.Sprint(data.BackupId)]; !exist {
				continue
			}
		}
		i := map[string]interface{}{
			"id":                fmt.Sprint(data.BackupId),
			"backup_method":     data.BackupMethod,
			"backup_id":         data.BackupId,
			"backup_mode":       data.BackupMode,
			"backup_status":     data.BackupStatus,
			"backup_size":       data.BackupSize,
			"slave_status":      data.SlaveStatus,
			"host_instance_id":  data.HostInstanceID,
			"backup_db_names":   data.BackupDBNames,
			"store_status":      data.StoreStatus,
			"backup_end_time":   data.BackupEndTime,
			"backup_start_time": data.BackupStartTime,
			"meta_status":       data.MetaStatus,
			"backup_scale":      data.BackupScale,
			"backup_location":   data.BackupLocation,
		}
		log.Printf("[DEBUG] %s: %#v", "backups", i)
		datas = append(datas, i)
		backup_id := fmt.Sprint(data.BackupId)
		ids = append(ids, backup_id)
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
