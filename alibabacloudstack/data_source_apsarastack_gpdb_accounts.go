package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackGpdbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackGpdbAccountsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Creating", "Deleting"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackGpdbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeAccounts"
	request := make(map[string]interface{})
	request["DBInstanceId"] = d.Get("db_instance_id")
	var objects []map[string]interface{}
	var accountNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		accountNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")

	response, err := client.DoTeaRequest("POST", "gpdb", "2016-05-03", action, "", nil, request)
	if err != nil {
		return err
	}
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_gpdb_accounts", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Accounts.DBInstanceAccount", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if accountNameRegex != nil && !accountNameRegex.MatchString(fmt.Sprint(item["AccountName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"], ":", item["AccountName"])]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != item["AccountStatus"].(string) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_description": object["AccountDescription"],
			"id":                  fmt.Sprint(object["DBInstanceId"], ":", object["AccountName"]),
			"account_name":        fmt.Sprint(object["AccountName"]),
			"db_instance_id":      object["DBInstanceId"],
			"status":              convertGpdbAccountStatusResponse(object["AccountStatus"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AccountName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("accounts", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
