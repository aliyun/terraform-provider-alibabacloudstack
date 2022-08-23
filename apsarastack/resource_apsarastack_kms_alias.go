package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackKmsAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackKmsAliasCreate,
		Read:   resourceApsaraStackKmsAliasRead,
		Update: resourceApsaraStackKmsAliasUpdate,
		Delete: resourceApsaraStackKmsAliasDelete,
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

func resourceApsaraStackKmsAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := kms.CreateCreateAliasRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.AliasName = d.Get("alias_name").(string)
	request.KeyId = d.Get("key_id").(string)
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateAlias(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kms_alias", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(request.AliasName)

	return resourceApsaraStackKmsAliasRead(d, meta)
}
func resourceApsaraStackKmsAliasRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsAlias(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alias_name", d.Id())
	d.Set("key_id", object.KeyId)
	return nil
}
func resourceApsaraStackKmsAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	if d.HasChange("key_id") {
		request := kms.CreateUpdateAliasRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.AliasName = d.Id()
		request.KeyId = d.Get("key_id").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateAlias(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	return resourceApsaraStackKmsAliasRead(d, meta)
}
func resourceApsaraStackKmsAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := kms.CreateDeleteAliasRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.AliasName = d.Id()
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DeleteAlias(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return nil
}
