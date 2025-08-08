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

func dataSourceAlibabacloudStackPolardbxBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbxBackupsRead,
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
						"backup_model": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_set_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"backup_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbxBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbx_backups := make([]PolarDbXBackupData, 0)
	request := client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeBackupSetList", "")
	PolarDbXBackupResponseObj := &PolarDbXBackupResponse{}
	page := 1
	page_size := 50
	db_instance_id := d.Get("db_instance_id").(string)
	request.QueryParams["DBInstanceName"] = db_instance_id
	if v, ok := d.GetOk("start_time"); ok {
		startTime, err := DateTimeToTimeStamp(v.(string))
		if err != nil {
			return err
		}
		request.QueryParams["StartTime"] = startTime
	}
	if v, ok := d.GetOk("end_time"); ok {
		endTime, err := DateTimeToTimeStamp(v.(string))
		if err != nil {
			return err
		}
		request.QueryParams["EndTime"] = endTime
	}
	for {
		request.QueryParams["PageNumber"] = fmt.Sprint(page)
		request.QueryParams["PageSize"] = fmt.Sprint(page_size)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardbx_backups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolarDbXBackupResponseObj)

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_polardbx_backups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		polardbx_backups = append(polardbx_backups, PolarDbXBackupResponseObj.Data...)
		if page*page_size >= PolarDbXBackupResponseObj.TotalNumber {
			break
		}
		page++
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	datas := make([]interface{}, 0)
	ids := make([]string, 0)
	for _, data := range polardbx_backups {
		if len(idsMap) > 0 {
			if _, exist := idsMap[data.BackupSetId]; !exist {
				continue
			}
		}
		i := map[string]interface{}{
			"id":              data.BackupSetId,
			"backup_model":    data.BackupModel,
			"backup_set_size": data.BackupSetSize,
			"backup_type":     data.BackupType,
			"backup_set_id":   data.BackupSetId,
			"status":          data.Status,
			"end_time":        TimeStampFormat(data.EndTime),
			"begin_time":      TimeStampFormat(data.BeginTime),
		}
		log.Printf("[DEBUG] %s: %#v", "backups", i)
		datas = append(datas, i)
		ids = append(ids, data.BackupSetId)
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

func DateTimeToTimeStamp(timeStr string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04Z", timeStr)
	if err != nil {
		fmt.Println("解析错误:", err)
		return "", err
	}
	timestampMs := t.UnixNano() / 1000000
	return fmt.Sprint(timestampMs), nil
}
