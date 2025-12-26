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

func dataSourceAlibabacloudStackCloudfwAddressBooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCloudfwAddressBooksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of group UUIDs to filter the results.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter address books by group name.",
			},
			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the address book (ip or port).",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query keyword for filtering address books by name.",
			},
			"contain_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter address books that contain a specific port.",
			},
			"address_books": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of address books. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the address book.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the address book.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the address book.",
						},
						"group_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the address book.",
						},
						"group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the address book (ip or port).",
						},
						"address_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of IP addresses or ports in the address book.",
						},
						"address_list_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of entries in the address list.",
						},
						"reference_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of times this address book is referenced.",
						},
						"global": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the address book is global (0: no, 1: yes).",
						},
						"auto_add_tag_ecs": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether ECS instances are automatically added based on tags (0: no, 1: yes).",
						},
						"tag_relation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag relation configuration.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCloudfwAddressBooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeAddressBook"

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	pageNumber := 1
	request := make(map[string]interface{})
	request["SourceCode"] = "yundun"
	request["GroupType"] = d.Get("group_type")

	if v, ok := d.GetOk("query"); ok {
		request["Query"] = v.(string)
	}

	if v, ok := d.GetOk("contain_port"); ok {
		request["ContainPort"] = v.(string)
	}

	var objects []interface{}
	var response map[string]interface{}
	var err error

	// Get all pages of results
	for {
		request["CurrentPage"] = pageNumber
		request["PageSize"] = "100"
		response, err = client.DoTeaRequest("POST", "cloudfw", "2017-12-07", action, "", nil, request, nil)

		if err != nil {
			return errmsgs.WrapError(err)
		}

		// Parse response
		acls, ok := response["Acls"]
		if ok && len(acls.([]interface{})) > 0 {
			objects = append(objects, acls.([]interface{})...)
		}
		totalCount, _ := response["TotalCount"].(json.Number).Int64()
		if totalCount <= int64(pageNumber*100) {
			break
		}
	}

	// Set computed values
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)

	for _, object := range objects {
		objectMap := object.(map[string]interface{})
		mapping := map[string]interface{}{}
		id := fmt.Sprintf("%s:%s", objectMap["GroupType"].(string), objectMap["GroupUuid"].(string))
		if nameRegex != nil {
			groupName, ok := objectMap["GroupName"].(string)
			if !ok || !nameRegex.MatchString(groupName) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}
		mapping["id"] = id
		mapping["group_name"] = objectMap["GroupName"]
		mapping["description"] = objectMap["Description"]
		mapping["group_uuid"] = objectMap["GroupUuid"]
		mapping["group_type"] = objectMap["GroupType"]
		mapping["address_list_count"] = objectMap["AddressListCount"]
		mapping["reference_count"] = objectMap["ReferenceCount"]
		mapping["global"] = objectMap["Global"]
		mapping["auto_add_tag_ecs"] = objectMap["AutoAddTagEcs"]
		mapping["tag_relation"] = objectMap["TagRelation"]
		mapping["address_list"] = objectMap["AddressList"]
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("address_books", s); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
