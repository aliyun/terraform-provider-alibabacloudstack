package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmUserGroupResourceSetBinding() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ascm_role_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_auth_expire": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expire_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"LENGTH", "VALID_UNTIL"}, false), // validUntilUtc  // lengthByDay
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"length_by_day": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingUpdate,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	resourceSetId := d.Get("resource_set_id").(string)
	userGroupId := d.Get("user_group_id").(string)
	ascmRoleId := d.Get("ascm_role_id").(string)
	enableAuthExpire := d.Get("enable_auth_expire").(bool)
	expire_type := d.Get("expire_type").(string)
	body := map[string]interface{}{
		"userGroupIdList":  []string{userGroupId},
		"resourceSetId":    resourceSetId,
		"ascmRoleId":       ascmRoleId,
		"enableAuthExpire": enableAuthExpire,
	}
	var config string
	if enableAuthExpire && expire_type == "" {
		return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", "when `enable_auth_expire` = true, `expire_type` is required")
	}
	if expire_type == "VALID_UNTIL" {
		if d.Get("expiration_time").(string) == "" {
			return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", "when `expire_type` = VALID_UNTIL, `expiration_time` is required")
		} else {
			config = fmt.Sprintf("{\"type\":\"%s\",\"validUntilUtc\":\"%s\"}", expire_type, d.Get("expiration_time").(string))
			log.Printf("[DEBUG] ===============================================expiration_time %#v", d.Get("expiration_time").(string))
		}
	}
	if expire_type == "LENGTH" {
		if d.Get("length_by_day").(int) == 0 {
			return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", "when `expire_type` = LENGTH, `length_by_day` is required")
		} else {
			config = fmt.Sprintf("{\"type\":\"%s\",\"lengthByDay\":%d}", expire_type, d.Get("length_by_day").(int))

		}
	}
	if config != "" {
		body["timeConfig"] = config
	}

	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddResourceSetToUserGroup", "/ascm/auth/user/addResourceSetToUserGroup", nil, nil, body)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(response)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(fmt.Sprintf("%s:%s:%s", resourceSetId, userGroupId, ascmRoleId))
	return nil
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	id_infos := strings.Split(d.Id(), ":")
	if len(id_infos) == 1 {
		// At this point, it's old data d.SetId(resourceSetId), forcibly modify the data format once
		d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("resource_set_id").(string), d.Get("user_group_id").(string), d.Get("ascm_role_id").(string)))
	}

	ascmService := &AscmService{client: client}
	resp, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	// Ensure that import actions can load normally
	id_infos = strings.Split(d.Id(), ":")
	d.Set("resource_set_id", id_infos[0])
	d.Set("user_group_id", id_infos[1])
	d.Set("ascm_role_id", strconv.Itoa(resp.Data[0].AuthorizedRoleId))
	// d.Set("enable_auth_expire", resp.Data[0].EnableAuthExpire)
	if resp.Data[0].TimeConfig != "" {
		d.Set("enable_auth_expire", resp.Data[0].EnableAuthExpire)
		config := make(map[string]interface{})
		err := json.Unmarshal([]byte(resp.Data[0].TimeConfig), &config)
		if err != nil {
			return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "DescribeAscmUserGroupResourceSetBinding", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		d.Set("expire_type", config["type"])
		if v, ok := config["validUntilUtc"]; ok {
			d.Set("expiration_time", v)
		}
		if v, ok := config["lengthByDay"]; ok {
			d.Set("length_by_day", v)
		}
	}

	return nil
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	return noUpdatesAllowedCheck(d, []string{"enable_auth_expire", "expire_type", "expiration_time", "length_by_day"})
	// if d.IsNewResource() {
	// 	return nil
	// }
	// client := meta.(*connectivity.AlibabacloudStackClient)

	// if d.HasChanges("enable_auth_expire", "expire_type", "expiration_time", "length_by_day") {
	// 	id_infos := strings.Split(d.Id(), ":")
	// 	resourceSetId := id_infos[0]
	// 	userGroupId := id_infos[1]
	// 	ascmRoleId := id_infos[2]
	// 	enableAuthExpire := d.Get("enable_auth_expire").(bool)
	// 	expire_type := d.Get("expire_type").(string)
	// 	body := map[string]interface{}{
	// 		"userGroupId":      userGroupId,
	// 		"resourceGroupId":  resourceSetId,
	// 		"ascmRoleId":       ascmRoleId,
	// 		"enableAuthExpire": enableAuthExpire,
	// 	}
	// 	var config string
	// 	if enableAuthExpire && expire_type == "" {
	// 		return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "UpdateUserGroupResourceExpireConfig", "when `enable_auth_expire` = true, `expire_type` is required")
	// 	}
	// 	if expire_type == "VALID_UNTIL" {
	// 		if d.Get("expiration_time").(string) == "" {
	// 			return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "UpdateUserGroupResourceExpireConfig", "when `expire_type` = VALID_UNTIL, `expiration_time` is required")
	// 		} else {
	// 			config = fmt.Sprintf("{\"type\":\"%s\",\"validUntilUtc\":\"%s\"}", expire_type, d.Get("expiration_time").(string))
	// 			log.Printf("[DEBUG] ===============================================expiration_time %#v", d.Get("expiration_time").(string))
	// 		}
	// 	}
	// 	if expire_type == "LENGTH" {
	// 		if d.Get("length_by_day").(int) == 0 {
	// 			return errmsgs.Error(errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "UpdateUserGroupResourceExpireConfig", "when `expire_type` = LENGTH, `length_by_day` is required")
	// 		} else {
	// 			config = fmt.Sprintf("{\"type\":\"%s\",\"lengthByDay\":%d}", expire_type, d.Get("length_by_day").(int))

	// 		}
	// 	}
	// 	if config != "" {
	// 		body["timeConfig"] = config
	// 	}

	// 	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "UpdateUserGroupResourceExpireConfig", "/ascm/auth/user/updateUserGroupResourceExpireConfig", nil, nil, body)
	// 	if err != nil {
	// 		errmsg := ""
	// 		if response != nil {
	// 			errmsg = errmsgs.GetAsapiErrorMessage(response)
	// 		}
	// 		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "UpdateUserGroupResourceExpireConfig", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	// 	}
	// }
	// return nil
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	_, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", requestInfo, map[string]string{"resourceGroupId": d.Id()})
	var resourceSetId, userGroupId string
	id_infos := strings.Split(d.Id(), ":")
	resourceSetId = id_infos[0]
	userGroupId = id_infos[1]

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveResourceSetFromUserGroup", "/ascm/auth/user/removeResourceSetFromUserGroup")
		request.QueryParams["ProductName"] = "ascm"
		request.QueryParams["userGroupId"] = userGroupId
		request.QueryParams["resourceSetId"] = resourceSetId
		delete(request.QueryParams, "ResourceGroup")
		delete(request.QueryParams, "OrganizationId")
		delete(request.QueryParams, "Department")
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "RemoveResourceSetFromUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		addDebug("RemoveResourceSetFromUserGroup", bresponse, request)
		_, err = ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
