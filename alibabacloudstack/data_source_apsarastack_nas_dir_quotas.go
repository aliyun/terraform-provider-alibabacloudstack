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

func dataSourceAlibabacloudStackNasDirQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackNasDirQuotasRead,

		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
			"dir_quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dir_inode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_quotas": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"quota_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size_limit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_count_limit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"size_real": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_count_real": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasDirQuotasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeDirQuotas"
	request := make(map[string]interface{})
	request["FileSystemId"] = d.Get("file_system_id").(string)
	if v, ok := d.GetOk("path"); ok {
		request["Path"] = v.(string)
	}
	var pageNumber int
	pageNumber = 1
	request["PageSize"] = 100
	request["PageNumber"] = pageNumber

	idsMap := getIdsStringFilter(d)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var dirQuotas []interface{}

	for {
		response, err := client.DoTeaRequest("POST", "nas", "2017-06-26", action, "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_dir_quotas", action, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		pageDirQuotas, ok := response["DirQuotaInfos"].([]interface{})
		if !ok {
			break
		}

		dirQuotas = append(dirQuotas, pageDirQuotas...)

		totalCount, _ := response["TotalCount"].(json.Number).Int64()
		if pageNumber*100 >= int(totalCount) {
			break
		}
		request["PageNumber"] = pageNumber + 1
	}

	var filteredDirQuotas []interface{}

	for _, quota := range dirQuotas {
		quotaInfo := quota.(map[string]interface{})
		path := quotaInfo["Path"].(string)
		fileSystemId := d.Get("file_system_id").(string)
		id := fmt.Sprintf("%s:%s", fileSystemId, path)

		// Check if id is in ids list
		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		// For name_regex, we use the path as the name to match against
		if nameRegex != nil {
			if !nameRegex.MatchString(path) {
				continue
			}
		}

		filteredDirQuotas = append(filteredDirQuotas, quotaInfo)
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)

	for _, object := range filteredDirQuotas {
		objectMap := object.(map[string]interface{})
		mapping := map[string]interface{}{}

		mapping["file_system_id"] = d.Get("file_system_id").(string)
		mapping["path"] = objectMap["Path"]
		mapping["dir_inode"] = objectMap["DirInode"]

		// Process user quotas
		userQuotaInfos, ok := objectMap["UserQuotaInfos"].([]interface{})
		if !ok {
			mapping["user_quotas"] = []map[string]interface{}{}
		} else {
			userQuotas := make([]map[string]interface{}, 0)
			for _, userQuota := range userQuotaInfos {
				userQuotaMap := userQuota.(map[string]interface{})
				userQuotaItem := map[string]interface{}{
					"user_type":        userQuotaMap["UserType"],
					"user_id":          userQuotaMap["UserId"],
					"quota_type":       userQuotaMap["QuotaType"],
					"size_limit":       userQuotaMap["SizeLimit"],
					"file_count_limit": userQuotaMap["FileCountLimit"],
					"size_real":        userQuotaMap["SizeReal"],
					"file_count_real":  userQuotaMap["FileCountReal"],
				}
				userQuotas = append(userQuotas, userQuotaItem)
			}
			mapping["user_quotas"] = userQuotas
		}

		id := fmt.Sprintf("%s:%s", d.Get("file_system_id").(string), mapping["path"].(string))
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("dir_quotas", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
