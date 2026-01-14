package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackBmsKeypairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmsKeypairsRead,

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
			"keypairs": {
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
						"public_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_pair_fingerprint": {
							Type:     schema.TypeString,
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
						"shared": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmsKeypairsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request query parameters
	reqQuery := map[string]interface{}{
		"DeployType": "bms",
		"PageNumber": 1,
		"PageSize":   50,
	}

	// Call ListKeyPair API
	resp, err := client.DoTeaRequest("POST", "bms", "2024-03-01", "ListKeyPair", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "ListKeyPair", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Parse response data
	data, err := jsonpath.Get("$.data.data", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListKeyPair", "$.data.data", resp)
	}

	items, ok := data.([]interface{})
	if !ok {
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
		idStr := name

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[idStr]; !exists {
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
		mapping["public_key"] = itemMap["publicKey"]
		mapping["key_pair_fingerprint"] = itemMap["keyPairFingerPrint"]
		mapping["create_time"] = itemMap["createTime"]
		mapping["update_time"] = itemMap["updateTime"]

		if v, ok := itemMap["shared"]; ok {
			mapping["shared"] = v
		} else {
			mapping["shared"] = 0
		}

		result = append(result, mapping)
		ids = append(ids, mapping["name"].(string))
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("keypairs", result); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
