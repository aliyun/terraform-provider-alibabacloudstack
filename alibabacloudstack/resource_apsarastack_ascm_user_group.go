package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserGroupCreate,
		Read:   resourceAlibabacloudStackAscmUserGroupRead,
		Update: resourceAlibabacloudStackAscmUserGroupUpdate,
		Delete: resourceAlibabacloudStackAscmUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_in_ids": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: "Attribute role_in_ids has been deprecated and replaced with role_ids.",
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackAscmUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	groupName := d.Get("group_name").(string)
	organizationId := d.Get("organization_id").(string)

	var loginNamesList []string

	if v, ok := d.GetOk("role_in_ids"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}

	if v, ok := d.GetOk("role_ids"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}
	requeststring, err := json.Marshal(map[string]interface{}{"roleIdList": loginNamesList})
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateUserGroup", "/ascm/auth/user/createUserGroup")

	request.QueryParams["groupName"] = groupName
	request.QueryParams["organizationId"] = organizationId
	// request.QueryParams["roleIdList"] = string(requeststring)
	request.SetContent(requeststring)

	bresponse, err := client.ProcessCommonRequest(request)
	log.Printf("response of raw CreateUserGroup is : %s", bresponse)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", errmsg)
	}
	addDebug("CreateUserGroup", bresponse, request, request.QueryParams)

	d.SetId(groupName)

	return resourceAlibabacloudStackAscmUserGroupUpdate(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackAscmUserGroupRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

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
	if len(object.Data[0].Roles) >= 0 {
		for _, role := range object.Data[0].Roles {
			roleIds = append(roleIds, strconv.Itoa(role.Id))
		}
	}
	d.Set("role_ids", roleIds)

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
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "DeleteUserGroup", "/ascm/auth/user/deleteUserGroup")
		request.QueryParams["userGroupId"] = strconv.Itoa(check.Data[0].Id)

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "DeleteUserGroup", errmsg))
		}
		return nil
	})
	return nil
}
