package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackUniversalDnsLine() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"v4_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"v6_addresses": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackUniversalDnsLineCreate, resourceAlibabacloudStackUniversalDnsLineRead, resourceAlibabacloudStackUniversalDnsLineUpdate, resourceAlibabacloudStackUniversalDnsLineDelete)
	return resource
}

func resourceAlibabacloudStackUniversalDnsLineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name")
	request["Priority"] = "0"
	v4Addresses := d.Get("v4_addresses").(*schema.Set).List()
	v6Addresses := d.Get("v6_addresses").(*schema.Set).List()
	if len(v4Addresses) == 0 && len(v6Addresses) == 0 {
		return fmt.Errorf("only one of v4_addresses and v6_addresses can be set")
	}
	request["V4Addresses"] = v4Addresses
	request["V6Addresses"] = v6Addresses
	response, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "AddUniversalLine", "", nil, request, nil)
	if err != nil {
		return err
	}

	resourceId, ok := response["Id"].(string)
	if !ok {
		return fmt.Errorf("failed to get Id from response")
	}

	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackUniversalDnsLineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	universalDnsService := UniversalDnsService{client}
	object, err := universalDnsService.DescribeUniversalDnsLine(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("name", object["Name"])
	d.Set("v4_addresses", object["V4Addresses"])
	d.Set("v6_addresses", object["V6Addresses"])
	d.Set("priority", object["Priority"])

	return nil
}

func resourceAlibabacloudStackUniversalDnsLineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// If it is a new resource, there is no need to update.
	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("name", "v4_addresses", "v6_addresses") {
		updateReq := make(map[string]interface{})
		updateReq["Id"] = d.Id()
		updateReq["Name"] = d.Get("name")
		v4Addresses := d.Get("v4_addresses").(*schema.Set).List()
		v6Addresses := d.Get("v6_addresses").(*schema.Set).List()
		if len(v4Addresses) == 0 && len(v6Addresses) == 0 {
			return fmt.Errorf("only one of v4_addresses and v6_addresses can be set")
		}
		updateReq["V4Addresses"] = v4Addresses
		updateReq["V6Addresses"] = v6Addresses
		_, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "UpdateUniversalLine", "", nil, nil, updateReq)
		if err != nil {
			return fmt.Errorf("failed to update Universal DNS Line: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackUniversalDnsLineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DeleteUniversalLine", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteUniversalLine", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if !raw["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteUniversalLine", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
