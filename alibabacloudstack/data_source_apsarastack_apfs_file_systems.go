package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackApfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApfsFileSystemsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of file system IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by file system description.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the file system, such as Pending, Running, etc.",
			},
			"file_systems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the file system.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the file system.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone ID where the file system is located.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the file system.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID of the file system.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster location of the file system.",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol type used by the file system.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage type of the file system.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume size in GB.",
						},
						"capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total capacity of the file system.",
						},
						"quota_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quota size assigned to the file system.",
						},
						"metered_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Metered usage size of the file system.",
						},
						"encrypt_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Encryption type of the file system.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing method of the file system.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of the file system.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth limit of the file system.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackApfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	request := make(map[string]interface{})
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	// Handle ids filter
	idsMap :=getIdsStringFilter(d)

	// Handle name_regex filter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return err
		}
		nameRegex = r
	}

	// Call API
	resp, err := client.DoTeaRequest("POST", "EFS", "2017-06-26", "DescribeFileSystems", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_apfs_file_systems", "DescribeFileSystems", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	fileSystemList, err := jsonpath.Get("$.FileSystems.FileSystem", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_apfs_file_systems", "$.FileSystems.FileSystem", resp)
	}

	var ids []string
	var s []map[string]interface{}
	for _, fs := range fileSystemList.([]interface{}) {
		fsMap, ok := fs.(map[string]interface{})
		if !ok {
			continue
		}

		id, _ := fsMap["FileSystemId"].(string)
		description, _ := fsMap["Description"].(string)

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name_regex filter
		if nameRegex != nil && !nameRegex.MatchString(description) {
			continue
		}

		mapping := make(map[string]interface{})
		mapping["id"] = id
		mapping["file_system_id"] = id
		mapping["description"] = fsMap["Description"]
		mapping["zone_id"] = fsMap["ZoneId"]
		mapping["create_time"] = fsMap["CreateTime"]
		mapping["region_id"] = fsMap["RegionId"]
		mapping["location"] = fsMap["Location"]
		mapping["protocol_type"] = fsMap["ProtocolType"]
		mapping["storage_type"] = fsMap["StorageType"]
		mapping["volume_size"] = fsMap["VolumeSize"]
		mapping["capacity"] = fsMap["Capacity"]
		mapping["quota_size"] = fsMap["QuotaSize"]
		mapping["metered_size"] = fsMap["MeteredSize"]
		mapping["encrypt_type"] = fsMap["EncryptType"]
		mapping["charge_type"] = fsMap["ChargeType"]
		mapping["status"] = fsMap["Status"]
		mapping["bandwidth"] = fsMap["Bandwidth"]
		ids = append(ids, id)

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("file_systems", s); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
