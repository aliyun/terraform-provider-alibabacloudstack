package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCloudFirewallControlPolicyOrder() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"acl_uuid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"in", "out"}, false),
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackCloudFirewallControlPolicyOrderCreate,
		resourceAlibabacloudStackCloudFirewallControlPolicyOrderRead,
		resourceAlibabacloudStackCloudFirewallControlPolicyOrderUpdate,
		resourceAlibabacloudStackCloudFirewallControlPolicyOrderDelete)
	return resource
}

func resourceAlibabacloudStackCloudFirewallControlPolicyOrderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "ModifyControlPolicyPosition"
	request := make(map[string]interface{})
	request["SourceCode"] = "yundun"
	request["Direction"] = d.Get("direction")
	request["NewOrder"] = d.Get("order")
	request["AclUuid"] = d.Get("acl_uuid")

	cloudfwService := CloudfwService{client}
	controlpolicy_id := fmt.Sprint(d.Get("acl_uuid"), ":", d.Get("direction"))
	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(controlpolicy_id)
	oldorder, err := toInt(object["Order"])
	if err != nil {
		return err
	}
	request["OldOrder"] = oldorder

	_, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(request["AclUuid"], ":", request["Direction"]))
	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyOrderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "ModifyControlPolicyPosition"
	request := map[string]interface{}{
		"AclUuid":    parts[0],
		"Direction":  parts[1],
		"SourceCode": "yundun",
	}

	if d.HasChange("order") {
		object, _ := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
		oldorder, err := toInt(object["Order"])
		if err != nil {
			return err
		}
		request["OldOrder"] = oldorder
		request["NewOrder"] = d.Get("order")
		response, err = client.DoTeaRequest("POST", "Cloudfw", "2017-12-07", action, "", nil, nil, request)
		addDebug(action, response, request)
		if err != nil {
			return err
		}
	}

	d.SetId(fmt.Sprint(request["AclUuid"], ":", request["Direction"]))
	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}
	object, err := cloudfwService.DescribeCloudFirewallControlPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_cloud_firewall_control_policy_order cloudfwService.DescribeCloudFirewallControlPolicy Failed!!! %s", err)
			// d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("acl_uuid", parts[0])
	d.Set("direction", parts[1])
	d.Set("order", formatInt(object["Order"]))

	return nil
}

func resourceAlibabacloudStackCloudFirewallControlPolicyOrderDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource alibabacloudstack_cloud_firewall_control_policy_order [%s]  will not be deleted", d.Id())
	return nil
}
