package alibabacloudstack

import (
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDataWorksFileTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDataWorksFileTypesRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the DataWorks project.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The keyword used to filter file types by name.",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"file_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the node type.",
						},
						"node_type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The numeric identifier of the node type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDataWorksFileTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"ProjectId": d.Get("project_id"),
		"PageSize":  100,
		"locale":    "en_US",
	}

	var nameFilter string
	idsMaps := getIdsStringFilter(d)
	if v, ok := d.GetOk("name"); ok {
		reqQuery["Keyword"] = v.(string)
		nameFilter = v.(string)
	}

	var names []string
	var ids []string
	var fileTypes []interface{}

	pageNum := 1

	for {
		reqQuery["PageNumber"] = pageNum
		response, err := client.DoTeaRequest("GET", "dataworks-public", "2020-05-18", "ListFileType", "", nil, reqQuery, nil)
		if err != nil {
			return err
		}

		nodeTypeInfoList, err := jsonpath.Get("$.NodeTypeInfoList.NodeTypeInfo", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg,
				"alibabacloudstack_dataworks_file_types", "$.NodeTypeInfoList.NodeTypeInfo", response)
		}

		for _, item := range nodeTypeInfoList.([]interface{}) {
			if itemMap, ok := item.(map[string]interface{}); ok {
				var name string
				var id int
				if v, exists := itemMap["NodeTypeName"]; !exists || v.(string) == "" {
					continue
				} else {
					name = v.(string)
				}
				if nameFilter != "" && name != nameFilter {
					continue
				}

				if v, err := toInt(itemMap["NodeType"]); err != nil {
					continue
				} else {
					id = v
				}
				if _, existed := idsMaps[strconv.Itoa(id)]; len(idsMaps) > 0 && !existed {
					continue
				}

				fileTypes = append(fileTypes, map[string]interface{}{
					"node_type_name": name,
					"node_type_id":   id,
				})
				ids = append(ids, strconv.Itoa(id))
				names = append(names, name)
			}
		}
		if len(nodeTypeInfoList.([]interface{})) < 100 {
			break
		}
		pageNum += 1
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("file_types", fileTypes); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
