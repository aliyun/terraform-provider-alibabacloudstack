package alibabacloudstack

import (
	"time"
	"fmt"
	"strings"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
	setResourceFunc(resource, 
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead,
		nil,
		resourceAlibabacloudStackAscmUserGroupResourceSetBindingDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	resourceSetId := d.Get("resource_set_id").(string)
	userGroupId := d.Get("user_group_id").(string)

	ascmRoleId := d.Get("ascm_role_id").(string)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddResourceSetToUserGroup", "/ascm/auth/user/addResourceSetToUserGroup")
	mergeMaps(request.QueryParams, map[string]string{
		"ProductName":   "ascm",
		"userGroupId":   userGroupId,
		"resourceSetId": resourceSetId,
		"ascmRoleId":    ascmRoleId,
	})

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddResourceSetToUserGroup", bresponse, requestInfo, request)
	if bresponse.GetHttpStatus() != 200 {
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddResourceSetToUserGroup", bresponse, requestInfo, bresponse.GetHttpContentString())
	d.SetId(fmt.Sprintf("%s:%s:%s",resourceSetId, userGroupId, ascmRoleId))
	return nil
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	
	id_infos := strings.Split(d.Id(), ":")
	if len(id_infos) == 1 {
		// 此时为老数据d.SetId(resourceSetId)， 强制修改一次数据的格式
		d.SetId(fmt.Sprintf("%s:%s:%s",d.Get("resource_set_id").(string), d.Get("user_group_id").(string), d.Get("ascm_role_id").(string)))
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
	// 确保import动作时可以正常加载
	id_infos = strings.Split(d.Id(), ":")
	d.Set("resource_set_id", id_infos[0])
	d.Set("user_group_id", id_infos[1])
	d.Set("ascm_role_id", strconv.Itoa(resp.Data[0].AuthorizedRoleId))

	return nil
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
