package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmResourceGroupUserAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rg_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmResourceGroupUserAttachmentCreate, resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead, resourceAlibabacloudStackAscmResourceGroupUserAttachmentUpdate, resourceAlibabacloudStackAscmResourceGroupUserAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	RgId := d.Get("rg_id").(string)
	userIds := d.Get("user_id").(string)

	userIdsArray := fmt.Sprintf("[%s]", userIds)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "BindAscmUserAndResourceGroup", "/ascm/auth/resource_group/add_ascm_users")
	request.QueryParams["ascm_user_ids"] = userIdsArray
	request.QueryParams["resource_group_id"] = RgId

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "BindAscmUserAndResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	id := fmt.Sprintf("%s:%s", RgId, userIds)
	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) < 2 {
		return errmsgs.WrapError(fmt.Errorf("Invalid ID format for resource group user attachment: %s", d.Id()))
	}
	rgId := parts[0]
	userId := parts[1]

	ascmService := &AscmService{client: client}
	response, err := ascmService.DescribeAscmResourceGroupUserAttachment(rgId)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	userFound := false
	for _, user := range response.Data {
		if fmt.Sprintf("%d", user.ID) == userId {
			userFound = true
			break
		}
	}

	if !userFound {
		d.SetId("")
		return nil
	}
	d.Set("rg_id", rgId)
	d.Set("user_id", userId)

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	parts := strings.Split(d.Id(), ":")
	if len(parts) < 2 {
		return errmsgs.WrapError(fmt.Errorf("Invalid ID format for resource group user attachment: %s", d.Id()))
	}
	rgId := parts[0]
	userId := parts[1]
	check, err := ascmService.DescribeAscmResourceGroupUserAttachment(rgId)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	userFound := false
	for _, user := range check.Data {
		if fmt.Sprintf("%d", user.ID) == userId {
			userFound = true
			break
		}
	}

	if !userFound {
		return nil
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UnbindAscmUserAndResourceGroup", "/ascm/auth/resource_group/remove_ascm_users")
		request.QueryParams["ascm_user_ids"] = fmt.Sprintf("[%s]", userId)
		request.QueryParams["resourceGroupId"] = rgId

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "UnbindAscmUserAndResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		if bresponse.GetHttpStatus() != 200 {
			return resource.RetryableError(fmt.Errorf("UnbindAscmUserAndResourceGroup failed with status: %d", bresponse.GetHttpStatus()))
		}

		response, err := ascmService.DescribeAscmResourceGroupUserAttachment(rgId)
		if err != nil {
			if !errmsgs.NotFoundError(err) {
				return resource.RetryableError(err)
			}
		} else {
			userStillExists := false
			for _, user := range response.Data {
				if fmt.Sprintf("%d", user.ID) == userId {
					userStillExists = true
					break
				}
			}

			if !userStillExists {
				return resource.NonRetryableError(nil)
			}
		}

		return resource.RetryableError(fmt.Errorf("User %s still exists in resource group %s", userId, rgId))
	})

	return err
}
