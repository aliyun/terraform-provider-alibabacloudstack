package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNasAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackNasAccessRuleCreate,
		Read:   resourceAlibabacloudStackNasAccessRuleRead,
		Update: resourceAlibabacloudStackNasAccessRuleUpdate,
		Delete: resourceAlibabacloudStackNasAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceAlibabacloudStackNasAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateAccessRule"
	request := make(map[string]interface{})
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.Region
	request["AccessGroupName"] = d.Get("access_group_name")
	request["SourceCidrIp"] = d.Get("source_cidr_ip")
	request["Product"] = "Nas"
	request["OrganizationId"] = client.Department
	request["Department"] = client.Department
	request["ResourceGroup"] = client.ResourceGroup
	if v, ok := d.GetOk("rw_access_type"); ok && v.(string) != "" {
		request["RWAccessType"] = v
	}
	if v, ok := d.GetOk("user_access_type"); ok && v.(string) != "" {
		request["UserAccessType"] = v
	}
	request["Priority"] = d.Get("priority")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_nas_access_rule", action, AlibabacloudStackSdkGoERROR)
	}
	d.SetId(fmt.Sprint(request["AccessGroupName"], ":", response["AccessRuleId"]))
	return resourceAlibabacloudStackNasAccessRuleRead(d, meta)
}

func resourceAlibabacloudStackNasAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = WrapError(err)
		return err
	}
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"AccessGroupName": parts[0],
		"AccessRuleId":    parts[1],
	}
	request["Product"] = "Nas"
	request["OrganizationId"] = client.Department
	request["Department"] = client.Department
	request["ResourceGroup"] = client.ResourceGroup
	update := false
	if d.HasChange("source_cidr_ip") {
		update = true
	}
	request["SourceCidrIp"] = d.Get("source_cidr_ip")

	if d.HasChange("rw_access_type") {
		update = true
	}
	request["RWAccessType"] = d.Get("rw_access_type")

	if d.HasChange("user_access_type") {
		update = true
	}
	request["UserAccessType"] = d.Get("user_access_type")

	if d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")

	if update {
		action := "ModifyAccessRule"
		conn, err := client.NewNasClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackNasAccessRuleRead(d, meta)
}

func resourceAlibabacloudStackNasAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasAccessRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_access_rule nasService.DescribeNasAccessRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		err = WrapError(err)
		return err
	}
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"AccessGroupName": parts[0],
		"AccessRuleId":    parts[1],
	}
	request["Product"] = "Nas"
	request["OrganizationId"] = client.Department
	request["Department"] = client.Department
	request["ResourceGroup"] = client.ResourceGroup
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.NasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
