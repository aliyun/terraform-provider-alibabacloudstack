package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceAlibabacloudStackImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackImageSharePermissionCreate,
		Read:   resourceAlibabacloudStackImageSharePermissionRead,
		Delete: resourceAlibabacloudStackImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlibabacloudStackImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	imageId := d.Get("image_id").(string)
	accountId := d.Get("account_id").(string)
	request := ecs.CreateModifyImageSharePermissionRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.ImageId = imageId
	accountSli := []string{accountId}
	request.AddAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_image_share_permission", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(imageId + ":" + accountId)
	return resourceAlibabacloudStackImageSharePermissionRead(d, meta)
}

func resourceAlibabacloudStackImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client: client}
	object, err := ecsService.DescribeImageShareByImageId(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	d.Set("image_id", object.ImageId)
	d.Set("account_id", parts[1])
	return WrapError(err)
}

func resourceAlibabacloudStackImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := ecs.CreateModifyImageSharePermissionRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{ "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	parts, err := ParseResourceId(d.Id(), 2)
	request.ImageId = parts[0]
	accountSli := []string{parts[1]}
	request.RemoveAccount = &accountSli
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyImageSharePermission(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_image_share_permission", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
