package alibabacloudstack

import (
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDmsEnterpriseUser() *schema.Resource {
	resource := &schema.Resource{
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
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The field `mobile` has been deprecated and is scheduled for removal in version 3.21.0.",
			},
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "NORMAL"}, false),
				Default:      "NORMAL",
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
				Deprecated:    "Field 'nick_name' has been deprecated.  Please use new field 'user_name' instead.",
				ConflictsWith: []string{"user_name"},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDmsEnterpriseUserCreate, resourceAlibabacloudStackDmsEnterpriseUserRead, resourceAlibabacloudStackDmsEnterpriseUserUpdate, resourceAlibabacloudStackDmsEnterpriseUserDelete)
	return resource
}

func resourceAlibabacloudStackDmsEnterpriseUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dmsService := DmsService{client}
	d.SetId(d.Get("uid").(string))
	object, err := dmsService.DescribeDmsEnterpriseUser(d.Id())
	if object != nil && object["State"].(string) == "DELETE" {
		request := map[string]interface{}{
			"Uid": d.Id(),
		}
		action := "EnableUser"
		_, err = client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	} else {
		action := "RegisterUser"
		request := make(map[string]interface{})

		if v, ok := d.GetOk("role_names"); ok && v != nil {
			request["RoleNames"] = convertListToCommaSeparate(v.(*schema.Set).List())
		}

		request["Uid"] = d.Get("uid")
		request["UserNick"] = connectivity.GetResourceData(d, "user_name", "nick_name").(string)

		_, err = client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceAlibabacloudStackDmsEnterpriseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dmsService := DmsService{client}
	object, err := dmsService.DescribeDmsEnterpriseUser(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_dms_enterprise_user dms_enterpriseService.DescribeDmsEnterpriseUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("uid", d.Id())
	d.Set("role_names", object["RoleNameList"].(map[string]interface{})["RoleNames"])
	d.Set("status", object["State"])
	d.Set("max_execute_count", object["MaxExecuteCount"])
	d.Set("max_result_count", object["MaxResultCount"])
	connectivity.SetResourceData(d, object["NickName"], "user_name", "nick_name")
	return nil
}

func resourceAlibabacloudStackDmsEnterpriseUserUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dmsService := DmsService{client}

	if d.HasChanges("role_names", "user_name", "nick_name", "max_execute_count", "max_result_count") {
		request := map[string]interface{}{
			"Uid":       d.Id(),
			"RoleNames": convertListToCommaSeparate(d.Get("role_names").(*schema.Set).List()),
			"UserNick":  connectivity.GetResourceData(d, "user_name", "nick_name").(string),
		}
		if _, ok := d.GetOk("max_execute_count"); ok {
			request["MaxExecuteCount"] = d.Get("max_execute_count")
		}
		if _, ok := d.GetOk("max_result_count"); ok {
			request["MaxResultCount"] = d.Get("max_result_count")
		}
		action := "UpdateUser"
		_, err = client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}

	if d.HasChange("status") {
		object, err := dmsService.DescribeDmsEnterpriseUser(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "DISABLE" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				action := "DisableUser"
				_, err = client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
				if err != nil {
					return err
				}
			}
			if target == "NORMAL" {
				request := map[string]interface{}{
					"Uid": d.Id(),
				}
				action := "EnableUser"
				_, err = client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackDmsEnterpriseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteUser"
	request := map[string]interface{}{
		"Uid": d.Get("uid"),
	}

	_, err := client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
