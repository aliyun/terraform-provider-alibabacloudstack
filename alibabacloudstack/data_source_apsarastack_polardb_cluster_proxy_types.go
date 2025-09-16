package alibabacloudstack

import (
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbClusterProxyTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterProxyTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PostgreSQL", "Oracle"}, false),
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"core_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed values.
			"proxy_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbClusterProxyTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := map[string]string{}
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			filterIds[vv.(string)] = ""
		}
	}

	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"label":        "true",
		"resourceType": "POLARDB",
		"status":       "Available",
		"groupFiled":   "proxyClass",
	}
	if v, ok := d.GetOk("db_version"); ok {
		reqQuery["dbVersion"] = v
	}

	if v, ok := d.GetOk("db_type"); ok {
		reqQuery["dbType"] = v
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "GroupCommonSpec", "/ascm/manage/saleconf/commonSpec/group", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	for _, v := range response["data"].([]interface{}) {
		data := v.(map[string]interface{})
		id := data["proxyClass"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := existedId[id]; exists {
			continue
		}
		classLabel := data["proxyClassLabel"].(string)
		core_count, err := strconv.Atoi(strings.Split(classLabel, " ")[0])
		if err != nil {
			return errmsgs.WrapErrorf(err, "error data is :%#v", data)
		}
		if v, ok := d.GetOk("core_count"); ok && core_count != v.(int) {
			continue
		}

		types = append(types, map[string]interface{}{
			"id":         id,
			"core_count": core_count,
		})
		existedId[id] = ""
		ids = append(ids, id)
	}
	sort.SliceStable(types, func(i, j int) bool {
		return types[i]["core_count"].(int) < types[j]["core_count"].(int)
	})

	d.Set("ids", ids)
	d.Set("proxy_classes", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
