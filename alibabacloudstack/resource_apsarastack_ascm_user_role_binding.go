package alibabacloudstack

import (
	"fmt"
	"log"

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
				ForceNew: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
		DeprecationMessage: "ascm_user already includes corresponding functions. This resource may be removed in future versions.",
	}

	setResourceFunc(resource, resourceAlibabacloudStackAscmUserRoleBindingCreate, resourceAlibabacloudStackAscmUserRoleBindingRead, resourceAlibabacloudStackAscmUserRoleBindingUpdate, resourceAlibabacloudStackAscmUserRoleBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)
	flag := false
	var roleids []int
	if v, ok := d.GetOk("role_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			roleids = append(roleids, id.(int))
		}
	}
	log.Printf("roleids is %v", roleids)
	flag = true
	if flag {
		for i := range roleids {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRoleToUser", "/ascm/auth/role/addRoleToUser")
			request.QueryParams["loginName"] = lname
			request.QueryParams["roleId"] = fmt.Sprint(roleids[i])

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil || bresponse.GetHttpStatus() != 200 {
				errmsg := ""
				if bresponse != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_role_binding", "AddRoleToUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			addDebug("AddRoleToUser", bresponse, request, request.QueryParams)
			log.Printf("response of queryparams AddRoleToUser is : %s", request.QueryParams)
		}
	}

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

	loginName := object.Data[0].LoginName
	d.SetId(loginName)
	d.Set("login_name", loginName)

	role_ids := make([]int, 0)
	for _, role := range object.Data[0].Roles {
		role_ids = append(role_ids, role.ID)
	}
	d.Set("role_ids", role_ids)
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	lname := d.Get("login_name").(string)

	if d.HasChange("role_ids") {
		o, n := d.GetChange("role_ids")
		oldValue := make(map[int]struct{})
		newValue := make(map[int]struct{})
		for _, v := range o.(*schema.Set).List() {
			oldValue[v.(int)] = struct{}{}
		}
		if len(oldValue) == 0 {
			if object, err := ascmService.DescribeAscmUserRoleBinding(d.Id()); err == nil && len(object.Data) > 0 {
				for _, role := range object.Data[0].Roles {
					oldValue[role.ID] = struct{}{}
				}
			}
		}
		for _, v := range n.(*schema.Set).List() {
			newValue[v.(int)] = struct{}{}
		}
		for key := range newValue {
			if _, exist := oldValue[key]; !exist {
				requestBody := map[string]interface{}{
					"loginName": lname,
					"roleId":    key,
				}
				if _, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddRoleToUser", "/ascm/auth/role/addRoleToUser", nil, nil, requestBody); err != nil {
					return err
				}
			}
		}
		for key := range oldValue {
			if _, exist := newValue[key]; !exist {
				requestBody := map[string]interface{}{
					"loginName": lname,
					"roleId":    key,
				}
				if _, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUser", "/ascm/auth/role/removeRoleFromUser", nil, nil, requestBody); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
