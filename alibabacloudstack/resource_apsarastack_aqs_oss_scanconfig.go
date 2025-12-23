package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAqsOssScanconfig() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scan_day_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(1, 7),
				},
				Optional: true,
			},
			"scan_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "2"}, false),
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"decryption": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OSS", "No"}, false),
			},
			"key_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modified_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					oldTime, _ := ParseToTimestampWithLocalTimezone(oldValue)
					newTime, _ := ParseToTimestampWithLocalTimezone(newValue)
					if oldTime == newTime || RoundToNearestThousand(newTime) == RoundToNearestThousand(oldTime) {
						return true
					}
					return false
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAqsOssScanconfigCreate, resourceAlibabacloudStackAqsOssScanconfigRead, resourceAlibabacloudStackAqsOssScanconfigUpdate, resourceAlibabacloudStackAqsOssScanconfigDelete)

	return resource
}

func resourceAlibabacloudStackAqsOssScanconfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	enable := 0
	if v, ok := d.GetOk("enable"); ok && v.(bool) {
		enable = 1
	}
	lastModifiedStartTime := d.Get("last_modified_start_time").(string)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", lastModifiedStartTime, time.Local)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "CreateOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	timestamp := t.Unix()

	reqBody := map[string]interface{}{
		"From":                  "sas",
		"Enable":                enable,
		"ScanMode":              d.Get("scan_mode"),
		"StartTime":             d.Get("start_time"),
		"EndTime":               d.Get("end_time"),
		"BucketNameList":        []string{d.Get("bucket_name").(string)},
		"DecryptionList":        []string{d.Get("decryption").(string)},
		"LastModifiedStartTime": timestamp,
		"ScanDayList":           d.Get("scan_day_list").(*schema.Set).List(),
	}

	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		reqBody["KeyPrefixList"] = []string{v.(string)}
		reqBody["AllKeyPrefix"] = false
	} else {
		reqBody["AllKeyPrefix"] = true
	}

	if v, ok := d.GetOk("key_suffix"); ok && v.(string) != "" {
		reqBody["KeySuffixList"] = []string{v.(string)}
	} else {
		reqBody["KeySuffixList"] = []string{"all"}
	}

	_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "CreateOssScanConfig", "", nil, reqBody, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	getReq := map[string]interface{}{
		"From":        "sas",
		"CurrentPage": 1,
		"PageSize":    100,
		"bucketName":  d.Get("bucket_name"),
	}
	response, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ListOssScanConfig", "", nil, getReq, nil)
	configs, ok := response["Data"].([]interface{})
	if !ok {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "ListOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	resourceId := ""
	for _, v := range configs {
		config := v.(map[string]interface{})
		bucketNameList := config["BucketNameList"].([]interface{})
		if bucketNameList[0].(string) == d.Get("bucket_name").(string) {
			resourceId = fmt.Sprint(config["Id"])
		}
	}
	if resourceId == "" {
		return fmt.Errorf("failed to find oss scanconfig with bucket name %s", d.Get("bucket_name"))
	} else {
		d.SetId(resourceId)
	}
	return nil
}

func resourceAlibabacloudStackAqsOssScanconfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	aqsService := AqsService{client}
	data, err := aqsService.DescribeAqsOssScanConfig(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "GetOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("enable", fmt.Sprint(data["Enable"]) == "1")
	d.Set("start_time", data["StartTime"])
	d.Set("end_time", data["EndTime"])
	d.Set("scan_mode", data["ScanMode"])
	if v, ok := data["ScanDayList"].([]interface{}); ok {
		scanDayList := make([]int, 0)
		for _, item := range v {
			if val, err := item.(json.Number).Int64(); err == nil {
				scanDayList = append(scanDayList, int(val))
			} else {
				return errmsgs.WrapError(err)
			}
		}
		d.Set("scan_day_list", scanDayList)
	}

	if v, ok := data["BucketNameList"].([]interface{}); ok {
		name := v[0].(string)
		d.Set("bucket_name", name)
	}
	if v, ok := data["DecryptionList"].([]interface{}); ok {
		name := v[0].(string)
		d.Set("decryption", name)
	}
	if v, ok := data["KeySuffixList"].([]interface{}); ok {
		name := v[0].(string)
		d.Set("key_suffix", name)
	}
	if v, ok := data["KeyPrefixList"].([]interface{}); ok {
		name := v[0].(string)
		d.Set("key_prefix", name)
	}
	lastModifiedStartTime, err := data["LastModifiedStartTime"].(json.Number).Int64()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	formattedTime := time.Unix(lastModifiedStartTime, 0).Format("2006-01-02 15:04:05")
	d.Set("last_modified_start_time", formattedTime)
	return nil
}

func resourceAlibabacloudStackAqsOssScanconfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	enable := 0
	if v, ok := d.GetOk("enable"); ok && v.(bool) {
		enable = 1
	}
	lastModifiedStartTime := d.Get("last_modified_start_time").(string)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", lastModifiedStartTime, time.Local)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "UpdateOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	timestamp := t.Unix()

	reqBody := map[string]interface{}{
		"Id":                    d.Id(),
		"From":                  "sas",
		"Enable":                enable,
		"ScanMode":              d.Get("scan_mode"),
		"StartTime":             d.Get("start_time"),
		"EndTime":               d.Get("end_time"),
		"BucketNameList":        []string{d.Get("bucket_name").(string)},
		"DecryptionList":        []string{d.Get("decryption").(string)},
		"LastModifiedStartTime": timestamp,
		"ScanDayList":           d.Get("scan_day_list").(*schema.Set).List(),
	}

	if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
		reqBody["KeyPrefixList"] = []string{v.(string)}
		reqBody["AllKeyPrefix"] = false
	} else {
		reqBody["AllKeyPrefix"] = true
	}

	if v, ok := d.GetOk("key_suffix"); ok && v.(string) != "" {
		reqBody["KeySuffixList"] = []string{v.(string)}
	} else {
		reqBody["KeySuffixList"] = []string{"all"}
	}

	_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "UpdateOssScanConfig", "", nil, reqBody, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "UpdateOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackAqsOssScanconfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	aqsService := AqsService{client}
	data, err := aqsService.DescribeAqsOssScanConfig(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "GetOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(data["Enable"]) == "1" {
		// if the config is enabled, we need to disable it first
		lastModifiedStartTime := d.Get("last_modified_start_time").(string)
		t, err := time.ParseInLocation("2006-01-02 15:04:05", lastModifiedStartTime, time.Local)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "UpdateOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		timestamp := t.Unix()

		reqBody := map[string]interface{}{
			"Id":                    d.Id(),
			"From":                  "sas",
			"Enable":                0,
			"ScanMode":              d.Get("scan_mode"),
			"StartTime":             d.Get("start_time"),
			"EndTime":               d.Get("end_time"),
			"BucketNameList":        []string{d.Get("bucket_name").(string)},
			"DecryptionList":        []string{d.Get("decryption").(string)},
			"LastModifiedStartTime": timestamp,
			"ScanDayList":           d.Get("scan_day_list").(*schema.Set).List(),
		}

		if v, ok := d.GetOk("key_prefix"); ok && v.(string) != "" {
			reqBody["KeyPrefixList"] = []string{v.(string)}
			reqBody["AllKeyPrefix"] = false
		} else {
			reqBody["AllKeyPrefix"] = true
		}

		if v, ok := d.GetOk("key_suffix"); ok && v.(string) != "" {
			reqBody["KeySuffixList"] = []string{v.(string)}
		} else {
			reqBody["KeySuffixList"] = []string{"all"}
		}
		_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "UpdateOssScanConfig", "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfig", "UpdateOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}
	reqQuery := map[string]interface{}{
		"From": "sas",
		"Id":   d.Id(),
	}

	_, err = client.DoTeaRequest("POST", "aegis", "2016-11-11", "DeleteOssScanConfig", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

func ParseToTimestampWithLocalTimezone(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}
	return t.In(time.Local).Unix(), nil
}

func RoundToNearestThousand(value int64) int64 {
	remainder := value % 1000
	if remainder >= 500 {
		return ((value / 1000) + 1) * 1000
	}
	return (value / 1000) * 1000
}
