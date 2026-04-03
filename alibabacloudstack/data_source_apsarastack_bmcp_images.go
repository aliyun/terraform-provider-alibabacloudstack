package alibabacloudstack

import (
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackBmcpImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpImagesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"customize", "public"}, false),
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"part": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reboot_script": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"signature_file_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_script": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"kernel_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_sum": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request query parameters
	reqQuery := map[string]interface{}{
		"QueryKey":    "Name",
		"QueryValue":  "",
		"PageNumber":  1,
		"PageSize":    10,
		"PackageType": "image",
	}

	reqQuery["Type"] = d.Get("type").(string)

	// Call ListImage API
	resp, err := client.DoTeaRequest("POST", "bms", "2022-05-30", "ListImage", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListImage", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response data
	data, err := jsonpath.Get("$.data", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListImage", "$.data", resp)
	}

	// Extract the actual items from the data field
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListImage", "data is not a map", resp)
	}

	var items []interface{}
	if items, ok = dataMap["data"].([]interface{}); !ok {
		items = []interface{}{}
	}

	// Prepare filtering maps
	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredItems []interface{}
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := itemMap["name"].(string)
		idStr := itemMap["uniqueKey"].(string)

		// Apply id filter (fuzzy match)
		if len(idsMap) > 0 {
			matchFound := false
			for id := range idsMap {
				if strings.Contains(idStr, id) || strings.Contains(id, idStr) {
					matchFound = true
					break
				}
			}
			if !matchFound {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	// Construct final result
	result := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range filteredItems {
		itemMap := item.(map[string]interface{})
		mapping := map[string]interface{}{}
		mapping["id"] = itemMap["name"].(string)
		mapping["name"] = itemMap["name"]
		mapping["description"] = itemMap["description"]
		mapping["type"] = itemMap["type"]
		mapping["platform"] = itemMap["platform"]
		mapping["zone"] = itemMap["zone"]
		mapping["region"] = itemMap["region"]
		mapping["status"] = itemMap["status"]
		mapping["size"] = itemMap["size"]
		mapping["create_time"] = itemMap["createTime"]
		mapping["update_time"] = itemMap["updateTime"]
		mapping["os_arch"] = itemMap["osArch"]
		mapping["disk_format"] = itemMap["diskFormat"]
		mapping["package_type"] = itemMap["packageType"]
		mapping["url"] = itemMap["URL"]
		mapping["unique_key"] = itemMap["uniqueKey"]
		mapping["template"] = itemMap["template"]
		mapping["resource_group"] = itemMap["resourceGroup"]
		mapping["file_name"] = itemMap["fileName"]
		mapping["snapshot_id"] = itemMap["snapshotID"]
		mapping["error_description"] = itemMap["errorDescription"]
		mapping["part"] = itemMap["part"]
		mapping["reboot_script"] = itemMap["rebootScript"]
		mapping["signature_file_url"] = itemMap["signatureFileURL"]
		mapping["owner"] = itemMap["owner"]
		mapping["component_type"] = itemMap["componentType"]
		mapping["version"] = itemMap["version"]
		mapping["custom_script"] = itemMap["customScript"]
		mapping["deleted"] = itemMap["deleted"]
		mapping["disable"] = itemMap["disable"]
		mapping["kernel_version"] = itemMap["kernelVersion"]
		mapping["organization"] = itemMap["organization"]
		mapping["check_sum"] = itemMap["checkSum"]

		result = append(result, mapping)
		ids = append(ids, mapping["id"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("images", result); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
