package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackQuickBiUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackQuickBiUserCreate,
		Read:   resourceAlibabacloudStackQuickBiUserRead,
		Update: resourceAlibabacloudStackQuickBiUserUpdate,
		Delete: resourceAlibabacloudStackQuickBiUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					parts := strings.Split(new, ":")
					if len(parts) < 2 {
						return false
					}
					return parts[1] == old
				},
			},
			"admin_user": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"auth_admin_user": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"nick_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Analyst", "Developer", "Visitor"}, false),
			},
		},
	}
}

func resourceAlibabacloudStackQuickBiUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "AddUser"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("account_id"); ok {
		request["AccountId"] = v
	}
	request["AccountName"] = d.Get("account_name")
	request["AdminUser"] = d.Get("admin_user")
	request["AuthAdminUser"] = d.Get("auth_admin_user")
	request["NickName"] = d.Get("nick_name")
	request["UserType"] = convertQuickBiUserUserTypeRequest(d.Get("user_type").(string))

	response, err = client.DoTeaRequest("POST", "quickbi-user", "2022-03-01", action, "", nil, request)
	if err != nil {
		return err
	}
	responseResult := response["Result"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseResult["UserId"]))

	return resourceAlibabacloudStackQuickBiUserRead(d, meta)
}

func resourceAlibabacloudStackQuickBiUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	quickbiPublicService := QuickbiPublicService{client}
	object, err := quickbiPublicService.DescribeQuickBiUser(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_quick_bi_user quickbiPublicService.DescribeQuickBiUser Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("account_id", object["AccountId"])
	d.Set("account_name", object["AccountName"])
	d.Set("admin_user", object["AdminUser"])
	d.Set("auth_admin_user", object["AuthAdminUser"])
	d.Set("nick_name", object["NickName"])
	d.Set("user_type", convertQuickBiUserUserTypeResponse(formatInt(object["UserType"])))
	return nil
}

func resourceAlibabacloudStackQuickBiUserUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"UserId": d.Id(),
	}
	if d.HasChange("admin_user") || d.IsNewResource() {
		update = true
	}
	request["AdminUser"] = d.Get("admin_user")
	if d.HasChange("auth_admin_user") || d.IsNewResource() {
		update = true
	}
	request["AuthAdminUser"] = d.Get("auth_admin_user")

	request["NickName"] = d.Get("nick_name")
	if d.HasChange("user_type") {
		update = true
	}
	request["UserType"] = convertQuickBiUserUserTypeRequest(d.Get("user_type").(string))
	if update {
		action := "UpdateUser"
		_, err = client.DoTeaRequest("POST", "quickbi-user", "2022-03-01", action, "", nil, request)
		if err != nil {
			return err
		}
	}
	return resourceAlibabacloudStackQuickBiUserRead(d, meta)
}

func resourceAlibabacloudStackQuickBiUserDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteUser"
	request := map[string]interface{}{
		"UserId": d.Id(),
	}

	_, err = client.DoTeaRequest("POST", "quickbi-user", "2022-03-01", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return nil
		}
		return err
	}
	return nil
}

func convertQuickBiUserUserTypeRequest(source interface{}) interface{} {
	switch source {
	case "Analyst":
		return 3
	case "Developer":
		return 1
	case "Visitor":
		return 2
	}
	return 0
}

func convertQuickBiUserUserTypeResponse(source interface{}) interface{} {
	switch source {
	case 3:
		return "Analyst"
	case 1:
		return "Developer"
	case 2:
		return "Visitor"
	}
	return ""
}
