package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDmsEnterpriseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDmsEnterpriseUserCreate,
		Read:   resourceAlibabacloudStackDmsEnterpriseUserRead,
		Update: resourceAlibabacloudStackDmsEnterpriseUserUpdate,
		Delete: resourceAlibabacloudStackDmsEnterpriseUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"max_execute_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_result_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "NORMAL"}, false),
				Default:      "NORMAL",
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"nick_name"},
			},
			"nick_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'nick_name' has been deprecated from version 1.100.0. Use 'user_name' instead.",
				ConflictsWith: []string{"user_name"},
			},
		},
	}
}

func resourceAlibabacloudStackDmsEnterpriseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "RegisterUser"
	request := make(map[string]interface{})
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("mobile"); ok {
		request["Mobile"] = v
	}

	if v, ok := d.GetOk("role_names"); ok && v != nil {
		request["RoleNames"] = convertListToCommaSeparate(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}

	request["Uid"] = d.Get("uid")
	if v, ok := d.GetOk("user_name"); ok {
		request["UserNick"] = v
	} else if v, ok := d.GetOk("nick_name"); ok {
		request["UserNick"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["Product"] = "dms-enterprise"
		request["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dms_enterprise_user", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Uid"]))

	return resourceAlibabacloudStackDmsEnterpriseUserUpdate(d, meta)
}
func resourceAlibabacloudStackDmsEnterpriseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	object, err := dms_enterpriseService.DescribeDmsEnterpriseUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dms_enterprise_user dms_enterpriseService.DescribeDmsEnterpriseUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("uid", d.Id())
	d.Set("mobile", object["Mobile"])
	d.Set("role_names", object["RoleNameList"].(map[string]interface{})["RoleNames"])
	d.Set("status", object["State"])
	d.Set("user_name", object["NickName"])
	d.Set("nick_name", object["NickName"])
	return nil
}
func resourceAlibabacloudStackDmsEnterpriseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dms_enterpriseService := Dms_enterpriseService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"Uid": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("mobile") {
		update = true
		request["Mobile"] = d.Get("mobile")
	}
	if !d.IsNewResource() && d.HasChange("role_names") {
		update = true
		request["RoleNames"] = convertListToCommaSeparate(d.Get("role_names").(*schema.Set).List())
	}
	if !d.IsNewResource() && d.HasChange("user_name") {
		update = true
		request["UserNick"] = d.Get("user_name")
	}
	if !d.IsNewResource() && d.HasChange("nick_name") {
		update = true
		request["UserNick"] = d.Get("nick_name")
	}
	if update {
		if _, ok := d.GetOk("max_execute_count"); ok {
			request["MaxExecuteCount"] = d.Get("max_execute_count")
		}
		if _, ok := d.GetOk("max_result_count"); ok {
			request["MaxResultCount"] = d.Get("max_result_count")
		}
		if _, ok := d.GetOk("tid"); ok {
			request["Tid"] = d.Get("tid")
		}
		action := "UpdateUser"
		conn, err := client.NewDmsenterpriseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["RegionId"] = client.RegionId
			request["Product"] = "dms-enterprise"
			request["OrganizationId"] = client.Department
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		//d.SetPartial("mobile")
		//d.SetPartial("role_names")
		//d.SetPartial("nick_name")
		//d.SetPartial("user_name")
	}
	if d.HasChange("status") {
		object, err := dms_enterpriseService.DescribeDmsEnterpriseUser(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "DISABLE" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				if v, ok := d.GetOk("tid"); ok {
					request["Tid"] = v
				}
				action := "DisableUser"
				conn, err := client.NewDmsenterpriseClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					request["RegionId"] = client.RegionId
					request["Product"] = "dms-enterprise"
					request["OrganizationId"] = client.Department
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			if target == "NORMAL" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				if v, ok := d.GetOk("tid"); ok {
					request["Tid"] = v
				}
				action := "EnableUser"
				conn, err := client.NewDmsenterpriseClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					request["RegionId"] = client.RegionId
					request["Product"] = "dms-enterprise"
					request["OrganizationId"] = client.Department
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			//d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackDmsEnterpriseUserRead(d, meta)
}
func resourceAlibabacloudStackDmsEnterpriseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteUser"
	var response map[string]interface{}
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Uid": d.Id(),
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["RegionId"] = client.RegionId
		request["Product"] = "dms-enterprise"
		request["OrganizationId"] = client.Department
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
