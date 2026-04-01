package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserRoleBinding() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		DeprecationMessage: "ascm_user already includes corresponding functions",
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmUserRoleBindingCreate, resourceAlibabacloudStackAscmUserRoleBindingRead, resourceAlibabacloudStackAscmUserRoleBindingUpdate, resourceAlibabacloudStackAscmUserRoleBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)
	d.SetId(lname)

	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}
	d.Set("login_name", object.Data[0].LoginName)
	role_ids := make([]string, 0)
	if len(object.Data[0].Roles) > 0 {
		for _, role := range object.Data[0].Roles {
			role_ids = append(role_ids, fmt.Sprint(role.ID))
		}
	}
	d.Set("role_ids", role_ids)
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	var roleIdList []string

	if v, ok := d.GetOk("role_ids"); ok {
		roleids := expandStringList(v.(*schema.Set).List())

		for _, roleid := range roleids {
			roleIdList = append(roleIdList, roleid)
		}
	}
	if len(roleIdList) < 1 {
		return errmsgs.Error("User role_ids cannot be empty!")
	}
	lname := d.Get("login_name").(string)
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetRolesForUserByLoginName", "/ascm/auth/user/ResetRolesForUserByLoginName")
	request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])

	QueryParams := map[string]interface{}{
		"loginName":  lname,
		"roleIdList": roleIdList,
	}

	requeststring, err := json.Marshal(QueryParams)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request.SetContent(requeststring)
	request.Headers["Content-Type"] = requests.Json

	bresponse, err := client.ProcessCommonRequest(request)

	log.Printf("response of raw ResetRolesForUserByLoginName is : %s", bresponse)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug("ResetRolesForUserByLoginName", bresponse, request, request.QueryParams)
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
