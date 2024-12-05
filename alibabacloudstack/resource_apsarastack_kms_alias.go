package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKmsAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKmsAliasCreate,
		Read:   resourceAlibabacloudStackKmsAliasRead,
		Update: resourceAlibabacloudStackKmsAliasUpdate,
		Delete: resourceAlibabacloudStackKmsAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alias_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlibabacloudStackKmsAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := kms.CreateCreateAliasRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.AliasName = d.Get("alias_name").(string)
	request.KeyId = d.Get("key_id").(string)
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateAlias(request)
	})
	response, ok := raw.(*kms.CreateAliasResponse)
	addDebug(request.GetActionName(), raw)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kms_alias", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	d.SetId(request.AliasName)

	return resourceAlibabacloudStackKmsAliasRead(d, meta)
}

func resourceAlibabacloudStackKmsAliasRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsAlias(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("alias_name", d.Id())
	d.Set("key_id", object.KeyId)
	return nil
}

func resourceAlibabacloudStackKmsAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChange("key_id") {
		request := kms.CreateUpdateAliasRequest()
		client.InitRpcRequest(*request.RpcRequest)

		request.AliasName = d.Id()
		request.KeyId = d.Get("key_id").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateAlias(request)
		})
		response, ok := raw.(*kms.UpdateAliasResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return resourceAlibabacloudStackKmsAliasRead(d, meta)
}

func resourceAlibabacloudStackKmsAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := kms.CreateDeleteAliasRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.AliasName = d.Id()
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DeleteAlias(request)
	})
	response, ok := raw.(*kms.DeleteAliasResponse)
	addDebug(request.GetActionName(), raw)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}
