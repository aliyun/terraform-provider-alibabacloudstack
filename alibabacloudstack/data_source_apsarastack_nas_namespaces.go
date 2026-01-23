package alibabacloudstack

import (
	"encoding/json"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackNasNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasNamespacesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of namespace IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by namespace description.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone ID that the namespace belongs to.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The storage type of the namespace. Possible values: Performance, Capacity.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol type of the namespace. Possible values: NFS, SMB.",
			},
			"file_system_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "standard",
				Description: "The file system type. Default value: standard.",
			},
			// Computed values
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nas_namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypt_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	pageNumber := 1
	request := map[string]interface{}{
		"PageSize": 10,
	}

	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v.(string)
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request["StorageType"] = v.(string)
	}

	var nasNamespaces []interface{}

	for {
		request["PageNumber"] = pageNumber
		raw, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "DescribeNamespaces", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		namespaces, ok := raw["NasNamespaces"].([]interface{})
		if ok {
			nasNamespaces = append(nasNamespaces, namespaces...)
		}
		totalCount, _ := raw["TotalCount"].(json.Number).Int64()
		if pageNumber*10 >= int(totalCount) {
			break
		}
		pageNumber++
	}

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var filteredNamespaces []interface{}
	for _, item := range nasNamespaces {
		namespace := item.(map[string]interface{})

		// Check if id is in the ids filter
		if len(idsMap) > 0 {
			id := namespace["NasNamespaceId"].(string)
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Check if description matches name_regex filter
		if nameRegex != nil {
			description := ""
			if desc, ok := namespace["Description"]; ok && desc != nil {
				description = desc.(string)
			}
			if !nameRegex.MatchString(description) {
				continue
			}
		}

		if zoneId, ok := d.GetOk("zone_id"); ok && zoneId != namespace["ZoneId"] {
			continue
		}

		if protocolType, ok := d.GetOk("protocol_type"); ok && protocolType != namespace["ProtocolType"] {
			continue
		}

		filteredNamespaces = append(filteredNamespaces, namespace)
	}

	// Set computed values
	var s []map[string]interface{}
	var ids []string

	for _, item := range filteredNamespaces {
		namespace := item.(map[string]interface{})
		mapping := map[string]interface{}{}

		mapping["id"] = namespace["NasNamespaceId"]
		mapping["nas_namespace_id"] = namespace["NasNamespaceId"]
		mapping["description"] = namespace["Description"]
		mapping["create_time"] = namespace["CreateTime"]
		mapping["zone_id"] = namespace["ZoneId"]
		mapping["protocol_type"] = namespace["ProtocolType"]
		mapping["storage_type"] = namespace["StorageType"]
		mapping["file_system_type"] = namespace["FileSystemType"]

		if encryptType, ok := namespace["EncryptType"]; ok {
			if encryptTypeFloat, ok := encryptType.(float64); ok {
				mapping["encrypt_type"] = int(encryptTypeFloat)
			} else if encryptTypeInt, ok := encryptType.(int); ok {
				mapping["encrypt_type"] = encryptTypeInt
			}
		}

		mapping["status"] = namespace["Status"]

		s = append(s, mapping)
		ids = append(ids, namespace["NasNamespaceId"].(string))
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("namespaces", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
