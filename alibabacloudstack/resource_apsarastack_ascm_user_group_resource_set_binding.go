package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
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
		},
	}
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	//ascmService := AscmService{client}

	resourceSetId := d.Get("resource_set_id").(string)
	userGroupId := d.Get("user_group_id").(string)

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
	//check, err := ascmService.DescribeAscmUserGroupResourceSetBinding(RgId)
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
	request.ApiName = "AddResourceSetToUserGroup"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "Ascm",
		"Action":          "AddResourceSetToUserGroup",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"ascmRoleId":      "2",
		"userGroupId":     userGroupId,
		"resourceSetId":   resourceSetId,
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddResourceSetToUserGroup is : %s", raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", raw)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_user_group_resource_set_binding", "AddResourceSetToUserGroup", AlibabacloudStackSdkGoERROR)
	}
	addDebug("AddResourceSetToUserGroup", raw, requestInfo, bresponse.GetHttpContentString())
	//}
	//err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	//	check, err = ascmService.DescribeAscmUserGroupResourceSetBinding(RgId)
	//	if err != nil {
	//		return resource.NonRetryableError(err)
	//	}
	//	return resource.RetryableError(err)
	//})
	//log.Printf("CreateOrganization Test %+v", check)
	d.SetId(resourceSetId)
	return resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupResourceSetBindingRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)

	ascmService := &AscmService{client: client}
	obj, err := ascmService.DescribeAscmUserGroupResourceSetBinding(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"resourceGroupId": d.Id()})
	userGroupId := d.Get("user_group_id").(string)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "RemoveResourceSetFromUserGroup",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"userGroupId":     userGroupId,
			"resourceSetId":   d.Id(),
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
		request.ApiName = "RemoveResourceSetFromUserGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
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
