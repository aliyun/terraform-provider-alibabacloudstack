package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackUniversalDnsDomain() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackUniversalDnsDomainCreate, resourceAlibabacloudStackUniversalDnsDomainRead, resourceAlibabacloudStackUniversalDnsDomainUpdate, resourceAlibabacloudStackUniversalDnsDomainDelete)
	return resource
}

func resourceAlibabacloudStackUniversalDnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Name"] = d.Get("name").(string)

	response, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "AddUniversalZone", "", nil, nil, request)
	if err != nil {
		return err
	}

	resourceId, ok := response["Id"].(string)
	if !ok || resourceId == "" {
		return fmt.Errorf("failed to retrieve Id from AddUniversalZone response")
	}
	d.SetId(resourceId)

	return nil
}

func resourceAlibabacloudStackUniversalDnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	universalDnsService := UniversalDnsService{client}
	target, err := universalDnsService.DescribeUniversalZones(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	// Set fields into ResourceData
	d.Set("id", target["Id"])
	d.Set("name", target["Name"])
	d.Set("remark", target["Remark"])
	d.Set("record_count", target["RecordCount"])
	d.Set("create_timestamp", target["CreateTimestamp"])
	d.Set("update_timestamp", target["UpdateTimestamp"])

	return nil
}

func resourceAlibabacloudStackUniversalDnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.HasChange("remark") {
		request := make(map[string]interface{})
		request["Id"] = d.Id()
		request["Remark"] = d.Get("remark").(string)
		request["Name"] = d.Get("name").(string) // Name is required in the API but immutable, sent for completeness

		_, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "UpdateUniversalZoneRemark", "", nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_universal_dns_domain", "UpdateUniversalZoneRemark", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackUniversalDnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"Id": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DeleteUniversalZone", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteUniversalZone", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if success, ok := raw["success"].(bool); !ok || !success {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteUniversalZone", raw["code"])
	}

	return nil
}
