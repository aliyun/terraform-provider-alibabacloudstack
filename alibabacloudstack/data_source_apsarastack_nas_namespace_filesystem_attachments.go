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

func dataSourceAlibabacloudStackNasNamespaceFilesystemAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasNamespaceFilesystemAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of attachment IDs in the form of NasNamespaceId:FileSystemId.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by mapped path.",
			},
			"nas_namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the namespace.",
			},
			"file_system_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the file system.",
			},
			"mapped_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The mapped path of the file system in the namespace.",
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Field 'id' has been deprecated from provider version 1.200.0. Use 'attachment_id' instead.",
						},
						"attachment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the attachment, formatted as NasNamespaceId:FileSystemId.",
						},
						"nas_namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the namespace.",
						},
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the file system.",
						},
						"mapped_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mapped path of the file system in the namespace.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the attachment.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage type of the file system.",
						},
						"file_system_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the file system.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasNamespaceFilesystemAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "ListFileSystemsInNamespace"

	idsMap := getIdsStringFilter(d)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	request := make(map[string]interface{})
	request["NasNamespaceId"] = d.Get("nas_namespace_id").(string)

	if v, ok := d.GetOk("file_system_id"); ok && v.(string) != "" {
		request["FileSystemId"] = v.(string)
	}

	if v, ok := d.GetOk("mapped_path"); ok && v.(string) != "" {
		request["MappedPath"] = v.(string)
	}

	request["PageSize"] = 100

	var objects []interface{}
	pageNumber := 1

	for {
		request["PageNumber"] = pageNumber

		response, err := client.DoTeaRequest("POST", "nas", "2017-06-26", action, "", nil, request, nil)

		if err != nil {
			return errmsgs.WrapError(err)
		}
		nsMembers, ok := response["NSMembers"].([]interface{})
		if !ok {
			break
		}
		objects = append(objects, nsMembers...)
		totalCount, _ := response["TotalCount"].(json.Number).Int64()
		if pageNumber*100 > int(totalCount) {
			break
		}
		pageNumber++
	}

	attachments := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, object := range objects {
		obj := object.(map[string]interface{})
		nasNamespaceId := d.Get("nas_namespace_id").(string)
		fileSystemId := obj["FileSystemId"].(string)
		attachmentId := fmt.Sprintf("%s:%s", nasNamespaceId, fileSystemId)
		if nameRegex != nil {
			mappedPath, exists := obj["MappedPath"]
			if !exists || !nameRegex.MatchString(fmt.Sprintf("%v", mappedPath)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[attachmentId]; !ok {
				continue
			}
		}
		mapping := map[string]interface{}{}
		mapping["id"] = attachmentId
		mapping["attachment_id"] = attachmentId
		mapping["nas_namespace_id"] = d.Get("nas_namespace_id").(string)
		mapping["file_system_id"] = fileSystemId
		mapping["mapped_path"] = obj["MappedPath"]
		mapping["create_time"] = obj["CreateTime"]
		mapping["storage_type"] = obj["StorageType"]
		mapping["file_system_type"] = obj["FileSystemType"]

		attachments = append(attachments, mapping)
		ids = append(ids, attachmentId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("attachments", attachments); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
