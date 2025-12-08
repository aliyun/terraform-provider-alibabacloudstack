package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPrometheusV2Contacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPrometheusV2ContactsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of contact IDs to filter results.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter contacts by username.",
			},
			"contacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mobile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPrometheusV2ContactsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request body for PageContact API
	requestBody := map[string]interface{}{
		"keyword":         "",
		"pageNumber":      1,
		"pageSize":        1000,
		"direction":       "desc",
		"page":            1,
		"selectedRowKeys": []string{},
	}

	// Call PageContact API
	resp, err := client.DoTeaRequest("POST", "prometheus2", "2023-04-13", "PageContact", "/log/api/v2/alert/contact/page", nil, nil, requestBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "PageContact", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	success, ok := resp["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to retrieve contacts: %v", resp)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format when retrieving contacts: %v", resp)
	}

	contactList, ok := data["data"].([]interface{})
	if !ok {
		contactList = []interface{}{}
	}

	// Process filtering
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
		nameRegex, err = regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("invalid name_regex pattern: %v", err)
		}
	}

	var filteredContacts []interface{}
	for _, item := range contactList {
		contact, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract fields
		idVal, _ := contact["id"]
		contactId := fmt.Sprint(idVal)
		username, _ := contact["username"].(string)

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[contactId]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(username) {
			continue
		}

		filteredContacts = append(filteredContacts, contact)
	}

	// Prepare result
	contacts := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, item := range filteredContacts {
		contact := item.(map[string]interface{})
		idVal, _ := contact["id"]
		contactId := fmt.Sprint(idVal)
		groups := contact["groups"].([]interface{})
		var groupIdStrs []string
		for _, v := range groups {
			group := v.(map[string]interface{})
			groupIdStrs = append(groupIdStrs, fmt.Sprint(group["id"]))
		}
		mapping := map[string]interface{}{
			"id":       contactId,
			"username": contact["username"],
			"mobile":   contact["mobile"],
			"mail":     contact["mail"],
			"groups":   groupIdStrs,
		}

		contacts = append(contacts, mapping)
		ids = append(ids, contactId)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("contacts", contacts); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
