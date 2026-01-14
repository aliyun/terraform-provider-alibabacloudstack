package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAqsOssScanconfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAqsOssScanconfigsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of OssScanConfig IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by bucket name.",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of bucket names corresponding to the scan configs.",
			},
			"scan_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_day_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"scan_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"decryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bucket_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"real_time_incr": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"all_key_prefix": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAqsOssScanconfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request query parameters
	request := map[string]interface{}{
		"From":        "sas",
		"CurrentPage": 1,
		"PageSize":    100,
	}

	// Call ListOssScanConfig API to retrieve all scan configs
	response, err := client.DoTeaRequest("POST", "aegis", "2016-11-11", "ListOssScanConfig", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_aqs_oss_scanconfigs", "ListOssScanConfig", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Extract Data from response
	data, ok := response["Data"].([]interface{})
	if !ok {
		// If no data returned, set empty values and return
		d.SetId("")
		d.Set("scan_configs", []interface{}{})
		d.Set("names", []interface{}{})
		return nil
	}

	// Prepare filtering maps
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Process each scan config
	var filteredConfigs []map[string]interface{}
	var names []interface{}
	var ids []string

	for _, item := range data {
		config, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		idStr := fmt.Sprint(config["Id"])
		bucketNameList, ok := config["BucketNameList"].([]interface{})
		if !ok || len(bucketNameList) == 0 {
			continue
		}
		bucketName := bucketNameList[0].(string)

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[idStr]; !exists {
				continue
			}
		}

		// Apply name_regex filter
		if nameRegex != nil && !nameRegex.MatchString(bucketName) {
			continue
		}

		// Convert fields as needed
		enable := false
		if fmt.Sprint(config["Enable"]) == "1" {
			enable = true
		}

		scanMode := fmt.Sprint(config["ScanMode"])
		startTime := config["StartTime"].(string)
		endTime := config["EndTime"].(string)

		// Handle ScanDayList
		var scanDayList []int
		if sdl, ok := config["ScanDayList"].([]interface{}); ok {
			for _, day := range sdl {
				if val, err := day.(json.Number).Int64(); err == nil {
					scanDayList = append(scanDayList, int(val))
				}
			}
		}

		// Handle DecryptionList
		decryption := ""
		if dl, ok := config["DecryptionList"].([]interface{}); ok && len(dl) > 0 {
			decryption = dl[0].(string)
		}

		// Handle KeySuffixList
		keySuffix := ""
		if ksl, ok := config["KeySuffixList"].([]interface{}); ok && len(ksl) > 0 {
			keySuffix = ksl[0].(string)
		}

		// Handle KeyPrefixList
		keyPrefix := ""
		if kpl, ok := config["KeyPrefixList"].([]interface{}); ok && len(kpl) > 0 {
			keyPrefix = kpl[0].(string)
		}

		// Handle RealTimeIncr
		realTimeIncr := false
		if rti, ok := config["RealTimeIncr"].(bool); ok {
			realTimeIncr = rti
		}

		// Handle AllKeyPrefix
		allKeyPrefix := false
		if akp, ok := config["AllKeyPrefix"].(bool); ok {
			allKeyPrefix = akp
		}

		// Construct result map
		result := map[string]interface{}{
			"id":                       idStr,
			"enable":                   enable,
			"start_time":               startTime,
			"end_time":                 endTime,
			"scan_day_list":            scanDayList,
			"scan_mode":                scanMode,
			"bucket_name":              bucketName,
			"decryption":               decryption,
			"key_suffix":               keySuffix,
			"key_prefix":               keyPrefix,
			"last_modified_start_time": config["LastModifiedStartTime"],
			"bucket_count":             config["BucketCount"],
			"real_time_incr":           realTimeIncr,
			"all_key_prefix":           allKeyPrefix,
			"last_update_time":         config["LastUpdateTime"],
		}

		filteredConfigs = append(filteredConfigs, result)
		names = append(names, bucketName)
		ids = append(ids, idStr)
	}

	// Set computed attributes
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("scan_configs", filteredConfigs); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
