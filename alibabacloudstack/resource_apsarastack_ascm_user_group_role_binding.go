package alibabacloudstack

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroupRoleBinding() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
		DeprecationMessage: "ascm_user_group already includes corresponding functions",
		Importer: &schema.ResourceImporter{
			State: resourceAlibabacloudStackAscmUserGroupRoleBindingImportState,
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackAscmUserGroupRoleBindingCreate,
		resourceAlibabacloudStackAscmUserGroupRoleBindingRead,
		resourceAlibabacloudStackAscmUserGroupRoleBindingUpdate,
		resourceAlibabacloudStackAscmUserGroupRoleBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	id = strings.TrimPrefix(id, "group:")

	userGroupId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user_group_id in import ID: must be integer")
	}

	d.Set("user_group_id", userGroupId)
	d.SetId(strconv.Itoa(userGroupId))

	return []*schema.ResourceData{d}, nil
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	userGroupId := d.Get("user_group_id").(int)
	var roleids []int
	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandIntList(v.(*schema.Set).List())
	}
	log.Printf("roleids is %v", roleids)

	for i := range roleids {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRoleToUserGroup", "/ascm/auth/user/addRoleToUserGroup")
		mergeMaps(request.QueryParams, map[string]string{
			"ProductName":      "ascm",
			"userGroupId":      strconv.Itoa(userGroupId),
			"RoleId":           fmt.Sprint(roleids[i]),
			"SecurityToken":    client.Config.SecurityToken,
			"SignatureVersion": "1.0",
			"SignatureMethod":  "HMAC-SHA1",
		})
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf("response of raw AddRoleToUserGroup Role(%d) is : %s", roleids[i], bresponse)

		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("AddRoleToUserGroup", bresponse, request, request.QueryParams)
		if bresponse.GetHttpStatus() != 200 {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		log.Printf("response of queryparams AddRoleToUserGroup is : %s", request.QueryParams)
	}

	d.SetId(strconv.Itoa(userGroupId))
	return nil
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroupRoleBinding(d.Id())
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

	userGroupId, _ := strconv.Atoi(d.Id())
	d.Set("user_group_id", userGroupId)

	role_ids := make([]int, 0)
	for _, role := range object.Data[0].Roles {
		role_ids = append(role_ids, role.Id)
	}
	d.Set("role_ids", role_ids)

	return nil
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	user_group_id := d.Get("user_group_id").(int)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	if d.HasChange("role_ids") {
		o, n := d.GetChange("role_ids")
		oldValue := make(map[int]struct{})
		newValue := make(map[int]struct{})
		for _, v := range o.(*schema.Set).List() {
			oldValue[v.(int)] = struct{}{}
		}
		if len(oldValue) == 0 {
			if object, err := ascmService.DescribeAscmUserGroupRoleBinding(d.Id()); err == nil && len(object.Data) > 0 {
				for _, role := range object.Data[0].Roles {
					oldValue[role.Id] = struct{}{}
				}
			}
		}
		for _, v := range n.(*schema.Set).List() {
			newValue[v.(int)] = struct{}{}
		}

		for key := range newValue {
			if _, exist := oldValue[key]; !exist {
				requestBody := map[string]interface{}{
					"userGroupId": user_group_id,
					"roleId":      key,
				}
				if _, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddRoleToUserGroup", "/ascm/auth/user/addRoleToUserGroup", nil, nil, requestBody); err != nil {
					return err
				}
			}
		}
		for key := range oldValue {
			if _, exist := newValue[key]; !exist {
				requestBody := map[string]interface{}{
					"userGroupId": user_group_id,
					"roleId":      key,
				}
				if _, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUserGroup", "/ascm/auth/user/removeRoleFromUserGroup", nil, nil, requestBody); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
