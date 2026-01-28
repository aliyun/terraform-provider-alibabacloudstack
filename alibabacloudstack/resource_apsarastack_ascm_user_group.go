package alibabacloudstack

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"organization_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field 'organization_id' has been deprecated. Use the organization to which the current user belongs",
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_in_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Deprecated:    "Field 'role_in_ids' is deprecated and will be removed in a future release. Please use 'role_ids' instead.",
				ConflictsWith: []string{"role_ids"},
			},
			"role_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"role_in_ids"},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmUserGroupCreate,
		resourceAlibabacloudStackAscmUserGroupRead, resourceAlibabacloudStackAscmUserGroupUpdate,
		resourceAlibabacloudStackAscmUserGroupDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	groupName := d.Get("group_name").(string)
	var organizationId string
	if v, ok := d.GetOk("organization_id"); ok {
		organizationId = v.(string)
	} else {
		organizationId = client.Department
	}
	var roleIdList []string
	if v, ok := connectivity.GetResourceDataOk(d, "role_in_ids", "role_ids"); ok {
		roleIds := expandStringList(v.(*schema.Set).List())
		for _, roleId := range roleIds {
			roleIdList = append(roleIdList, roleId)
		}
	}
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateUserGroup", "/ascm/auth/user/createUserGroup")
	requeststring, _ := json.Marshal(map[string]interface{}{"roleIdList": roleIdList})
	request.QueryParams["groupName"] = groupName
	request.QueryParams["organizationId"] = organizationId
	request.QueryParams["OrganizationId"] = organizationId
	request.SetContent(requeststring)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", errmsg)
	}
	addDebug("CreateUserGroup", bresponse, request, request.QueryParams)
	d.SetId(groupName)
	return nil
}

func resourceAlibabacloudStackAscmUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	userGroupId := d.Get("user_group_id").(string)
	var organizationId string
	if v, ok := d.GetOk("organization_id"); ok {
		organizationId = v.(string)
	} else {
		organizationId = client.Department
	}
	if _, ok := d.GetOk("role_ids"); ok && !d.IsNewResource() {
		oldV, newV := d.GetChange("role_ids")
		newSet, okNew := newV.(*schema.Set)
		if !okNew {
			return nil
		}
		oldSet, okOld := oldV.(*schema.Set)
		if !okOld {
			return nil
		}
		remove := oldSet.Difference(newSet).List()
		create := newSet.Difference(oldSet).List()
		for _, roleId := range create {
			req := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRoleToUserGroup", "/ascm/auth/user/addRoleToUserGroup")
			req.QueryParams["userGroupId"] = userGroupId
			req.QueryParams["OrganizationId"] = organizationId
			req.QueryParams["roleId"] = roleId.(string)
			bresp, err := client.ProcessCommonRequest(req)
			if err != nil || bresp.GetHttpStatus() != 200 {
				errmsg := ""
				if bresp != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresp.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "AddRoleToUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug("AddRoleToUser", bresp, req, req.QueryParams)
		}
		for _, roleId := range remove {
			req := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUserGroup", "/ascm/auth/user/removeRoleFromUserGroup")
			req.QueryParams["userGroupId"] = userGroupId
			req.QueryParams["OrganizationId"] = organizationId
			req.QueryParams["roleId"] = roleId.(string)
			bresp, err := client.ProcessCommonRequest(req)
			if err != nil || bresp.GetHttpStatus() != 200 {
				errmsg := ""
				if bresp != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresp.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "RemoveRoleFromUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug("RemoveRoleFromUser", bresp, req, req.QueryParams)
		}
	}
	return nil
}

func resourceAlibabacloudStackAscmUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroup(d.Id())
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
	d.Set("user_group_id", strconv.Itoa(object.Data[0].Id))
	d.Set("group_name", object.Data[0].GroupName)
	d.Set("organization_id", strconv.Itoa(object.Data[0].OrganizationId))
	var roleIds []string
	for _, role := range object.Data[0].Roles {
		roleIds = append(roleIds, strconv.Itoa(role.Id))
	}
	connectivity.SetResourceData(d, roleIds, "role_ids", "role_in_ids")
	return nil
}

func resourceAlibabacloudStackAscmUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	check, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsUserGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		req := client.NewCommonRequest("POST", "ascm", "2019-05-10", "DeleteUserGroup", "/ascm/auth/user/deleteUserGroup")
		req.QueryParams["userGroupId"] = strconv.Itoa(check.Data[0].Id)
		req.QueryParams["OrganizationId"] = d.Get("organization_id").(string)
		req.QueryParams["Department"] = d.Get("organization_id").(string)
		bresp, err := client.ProcessCommonRequest(req)
		if err != nil {
			errmsg := ""
			if bresp != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresp.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "DeleteUserGroup", errmsg))
		}
		return nil
	})
	return err
}
