package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceAlibabacloudStackAscmUserGroupResourceSetBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate,
		Read:   resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead,
		Delete: resourceAlibabacloudStackAscmUserGroupResourceSetBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	resourceSetId := d.Get("resource_set_id").(string)
	userGroupId := d.Get("user_group_id").(string)

	ascmRoleId := d.Get("ascm_role_id").(string)

	request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "AddResourceSetToUserGroup", "")
	mergeMaps(request.QueryParams, map[string]string{
		"ProductName":  "ascm",
		"userGroupId":  userGroupId,
		"resourceSetId": resourceSetId,
		"ascmRoleId":   ascmRoleId,
	})

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddResourceSetToUserGroup is : %s", raw)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, request)
	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, bresponse.GetHttpContentString())
	d.SetId(resourceSetId)
	return resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	ascmService := &AscmService{client: client}
	obj, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("resource_set_id", strconv.Itoa(obj.Data[0].Id))

	return nil
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"resourceGroupId": d.Id()})
	userGroupId := d.Get("user_group_id").(string)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "RemoveResourceSetFromUserGroup", "")
		request.QueryParams["ProductName"] = "ascm"
		request.QueryParams["userGroupId"] = userGroupId
		request.QueryParams["resourceSetId"] = d.Id()

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "RemoveResourceSetFromUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		addDebug("RemoveResourceSetFromUserGroup", raw, request)
		_, err = ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
