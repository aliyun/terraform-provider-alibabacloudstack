package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
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

			//"user_ids": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem:     &schema.Schema{Type: schema.TypeString},
			//	MinItems: 1,
			//	ForceNew: true,
			//},
		},
	}
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	//ascmService := AscmService{client}

	RgId := d.Get("rg_id").(string)
	userIds := d.Get("user_id").(string)

	//var userId string
	//var userIds []string
	//if v, ok := d.GetOk("user_ids"); ok {
	//	userIds = expandStringList(v.(*schema.Set).List())
	//	for i, k := range userIds {
	//		if i != 0 {
	//			userId = fmt.Sprintf("%s\",\"%s", userId, k)
	//		} else {
	//			userId = k
	//		}
	//	}
	//}
	//check, err := ascmService.DescribeAscmResourceGroupUserAttachment(RgId)
	//if err != nil {
	//	return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_organization", "ORG alreadyExist", AlibabacloudStackSdkGoERROR)
	//}
	//if len(check.Data) == 0 {

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "BindAscmUserAndResourceGroup"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{
		"RegionId":          client.RegionId,
		"AccessKeySecret":   client.SecretKey,
		"Product":           "Ascm",
		"Action":            "BindAscmUserAndResourceGroup",
		"Version":           "2019-05-10",
		"ProductName":       "ascm",
		"ascm_user_ids":     fmt.Sprintf("%s", userIds),
		"resource_group_id": RgId,
		//"X-acs-body": fmt.Sprintf("\"ascm_user_ids\":\"[\"5249\"]\""),
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw BindAscmUserAndResourceGroup is : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "BindAscmUserAndResourceGroup", raw)
	}
	addDebug("BindAscmUserAndResourceGroup", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "BindAscmUserAndResourceGroup", AlibabacloudStackSdkGoERROR)
	}
	addDebug("BindAscmUserAndResourceGroup", raw, requestInfo, bresponse.GetHttpContentString())
	//}
	//err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	//	check, err = ascmService.DescribeAscmResourceGroupUserAttachment(RgId)
	//	if err != nil {
	//		return resource.NonRetryableError(err)
	//	}
	//	return resource.RetryableError(err)
	//})
	//log.Printf("CreateOrganization Test %+v", check)
	d.SetId(RgId)
	return resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead(d, meta)
}

func resourceAlibabacloudStackAscmResourceGroupUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	ascmService := &AscmService{client: client}
	obj, err := ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("rg_id", obj.ResourceGroupID)
	//d.Set("user_ids", obj.AscmUserIds)

	return nil
}
func resourceAlibabacloudStackAscmResourceGroupUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"resourceGroupId": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "UnbindAscmUserAndResourceGroup",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"resourceGroupId": d.Id(),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "UnbindAscmUserAndResourceGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		_, err = ascmService.DescribeAscmResourceGroupUserAttachment(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
