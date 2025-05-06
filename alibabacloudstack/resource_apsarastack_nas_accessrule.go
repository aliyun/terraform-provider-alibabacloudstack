package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNasAccessRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_cidr_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rw_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDWR", "RDONLY"}, false),
				Default:      "RDWR",
			},
			"user_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"no_squash", "root_squash", "all_squash"}, false),
				Default:      "no_squash",
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasAccessRuleCreate, resourceAlibabacloudStackNasAccessRuleRead, resourceAlibabacloudStackNasAccessRuleUpdate, resourceAlibabacloudStackNasAccessRuleDelete)
	return resource
}

func resourceAlibabacloudStackNasAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateAccessRule"
	request := make(map[string]interface{})
	request["AccessGroupName"] = d.Get("access_group_name")
	request["SourceCidrIp"] = d.Get("source_cidr_ip")
	if v, ok := d.GetOk("rw_access_type"); ok && v.(string) != "" {
		request["RWAccessType"] = v
	}
	if v, ok := d.GetOk("user_access_type"); ok && v.(string) != "" {
		request["UserAccessType"] = v
	}
	request["Priority"] = d.Get("priority")

	response, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
	
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(request["AccessGroupName"], ":", response["AccessRuleId"]))
	return nil
}

func resourceAlibabacloudStackNasAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return err
	}
	request := map[string]interface{}{
		"AccessGroupName": parts[0],
		"AccessRuleId":    parts[1],
	}
	request["SourceCidrIp"] = d.Get("source_cidr_ip")
	request["RWAccessType"] = d.Get("rw_access_type")
	request["UserAccessType"] = d.Get("user_access_type")
	request["Priority"] = d.Get("priority")

	if d.HasChanges("source_cidr_ip", "rw_access_type", "user_access_type", "priority") {
		action := "ModifyAccessRule"
		_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackNasAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasAccessRule(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_access_rule nasService.DescribeNasAccessRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("access_rule_id", object["AccessRuleId"])
	d.Set("source_cidr_ip", object["SourceCidrIp"])
	d.Set("access_group_name", parts[0])
	d.Set("priority", formatInt(object["Priority"]))
	d.Set("rw_access_type", object["RWAccess"])
	d.Set("user_access_type", object["UserAccess"])

	return nil
}

func resourceAlibabacloudStackNasAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteAccessRule"
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return err
	}
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"AccessGroupName": parts[0],
		"AccessRuleId":    parts[1],
	}

	_, err = client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return err
	}
	return nil
}
