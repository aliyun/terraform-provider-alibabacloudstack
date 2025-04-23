package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackKmsSecret() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_delete_without_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"planned_delete_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recovery_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("force_delete_without_recovery").(bool)
				},
			},
			"secret_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"secret_data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"text", "binary"}, false),
				Default:      "text",
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_stages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackKmsSecretCreate, resourceAlibabacloudStackKmsSecretRead, resourceAlibabacloudStackKmsSecretUpdate, resourceAlibabacloudStackKmsSecretDelete)
	return resource
}

func resourceAlibabacloudStackKmsSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := kms.CreateCreateSecretRequest()
	client.InitRpcRequest(*request.RpcRequest)

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("encryption_key_id"); ok {
		request.EncryptionKeyId = v.(string)
	}
	request.SecretData = d.Get("secret_data").(string)
	if v, ok := d.GetOk("secret_data_type"); ok {
		request.SecretDataType = v.(string)
	}
	request.SecretName = d.Get("secret_name").(string)
	if v, ok := d.GetOk("tags"); ok {
		addTags := make([]JsonTag, 0)
		for key, value := range v.(map[string]interface{}) {
			addTags = append(addTags, JsonTag{
				TagKey:   key,
				TagValue: value.(string),
			})
		}
		tags, err := json.Marshal(addTags)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.Tags = string(tags)
	}
	request.VersionId = d.Get("version_id").(string)
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateSecret(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*kms.CreateSecretResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kms_secret", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*kms.CreateSecretResponse)
	d.SetId(response.SecretName)

	return nil
}

func resourceAlibabacloudStackKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsSecret(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("secret_name", d.Id())
	d.Set("arn", object.Arn)
	d.Set("description", object.Description)
	d.Set("encryption_key_id", object.EncryptionKeyId)
	d.Set("planned_delete_time", object.PlannedDeleteTime)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)

	getSecretValueObject, err := kmsService.GetSecretValue(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("secret_data", getSecretValueObject.SecretData)
	d.Set("secret_data_type", getSecretValueObject.SecretDataType)
	d.Set("version_id", getSecretValueObject.VersionId)
	d.Set("version_stages", getSecretValueObject.VersionStages.VersionStage)
	return nil
}

func resourceAlibabacloudStackKmsSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kmsService := KmsService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := kmsService.setResourceTags(d, "secret"); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	if d.HasChange("description") {
		request := kms.CreateUpdateSecretRequest()
		client.InitRpcRequest(*request.RpcRequest)

		request.SecretName = d.Id()
		request.Description = d.Get("description").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateSecret(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*kms.UpdateSecretResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	update := false
	request := kms.CreatePutSecretValueRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.SecretName = d.Id()
	if d.HasChange("secret_data") {
		update = true
	}
	request.SecretData = d.Get("secret_data").(string)
	if d.HasChange("version_id") {
		update = true
	}
	request.VersionId = d.Get("version_id").(string)
	if d.HasChange("secret_data_type") {
		update = true
		request.SecretDataType = d.Get("secret_data_type").(string)
	}
	if d.HasChange("version_stages") {
		update = true
		request.VersionStages = convertListToJsonString(d.Get("version_stages").(*schema.Set).List())
	}
	if update {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.PutSecretValue(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*kms.PutSecretValueResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	d.Partial(false)
	return nil
}

func resourceAlibabacloudStackKmsSecretDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := kms.CreateDeleteSecretRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.SecretName = d.Id()
	if v, ok := d.GetOkExists("force_delete_without_recovery"); ok {
		request.ForceDeleteWithoutRecovery = fmt.Sprintf("%v", v.(bool))
	}
	if v, ok := d.GetOk("recovery_window_in_days"); ok {
		request.RecoveryWindowInDays = fmt.Sprintf("%v", v.(int))
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DeleteSecret(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.errmsgs.ResourceNotfound"}) {
			return nil
		}
		errmsg := ""
		if response, ok := raw.(*kms.DeleteSecretResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}