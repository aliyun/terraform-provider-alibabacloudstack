package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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
	var requestInfo *ecs.Client

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

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateUserGroup", "/roa/ascm/auth/user/createUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"groupName":      groupName,
		"organizationId": organizationId,
		"roleIdList":     loginNamesList,
	}

	requeststring, err := json.Marshal(QueryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", "json format failed")
	}
	request.SetContent(requeststring)
	request.Headers["Content-Type"] = requests.Json

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	log.Printf("response of raw CreateUserGroup is : %s", raw)

	if err != nil {
		errmsg := ""
		bresponse, ok := raw.(*responses.CommonResponse)
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", errmsg)
	}

	addDebug("CreateUserGroup", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)

	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		bresponse, ok := raw.(*responses.CommonResponse)
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", errmsg)
	}
	addDebug("CreateUserGroup", raw, requestInfo, bresponse.GetHttpContentString())

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
	var requestInfo *ecs.Client

	check, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsUserGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsUserGroupExist", check, requestInfo, map[string]string{"groupName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "DeleteUserGroup", "/ascm/auth/user/deleteUserGroup")
		request.QueryParams["userGroupId"] = strconv.Itoa(check.Data[0].Id)

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			errmsg := ""
			bresponse, ok := raw.(*responses.CommonResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group", "DeleteUserGroup", errmsg))
		}
		check, err = ascmService.DescribeAscmUserGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteUserGroup", raw, request)
		return nil
	})
	return nil
}
