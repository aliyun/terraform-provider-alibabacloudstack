package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

func resourceAlibabacloudStackAscmResourceGroupUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmResourceGroupUserAttachmentCreate,
		Read:   resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead,
		Delete: resourceAlibabacloudStackAscmResourceGroupUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rg_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	RgId := d.Get("rg_id").(string)
	userIds := d.Get("user_id").(string)

	request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "BindAscmUserAndResourceGroup", "")
	request.QueryParams["ascm_user_ids"] = fmt.Sprintf("%s", userIds)
	request.QueryParams["resource_group_id"] = RgId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw BindAscmUserAndResourceGroup is : %s", raw)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "BindAscmUserAndResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("BindAscmUserAndResourceGroup", raw, requestInfo, request)
	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "BindAscmUserAndResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("BindAscmUserAndResourceGroup", raw, requestInfo, bresponse.GetHttpContentString())
	d.SetId(RgId)
	return resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead(d, meta)
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	ascmService := &AscmService{client: client}
	obj, err := ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("rg_id", obj.ResourceGroupID)

	return nil
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	check, err := ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"resourceGroupId": d.Id()})

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UnbindAscmUserAndResourceGroup", "")
		request.QueryParams["resourceGroupId"] = d.Id()

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "UnbindAscmUserAndResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		_, err = ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
