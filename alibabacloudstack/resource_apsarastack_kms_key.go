package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKmsKeyCreate,
		Read:   resourceAlibabacloudStackKmsKeyRead,
		Update: resourceAlibabacloudStackKmsKeyUpdate,
		Delete: resourceAlibabacloudStackKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"automatic_rotation": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
				Default:      "Disabled",
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled", "PendingDeletion"}, false),
				Default:      "Enabled",
			},
			"is_enabled": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'is_enabled' has been deprecated from provider version 1.85.0. New field 'key_state' instead.",
			},
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENCRYPT/DECRYPT", "SIGN/VERIFY"}, false),
				Default:      "ENCRYPT/DECRYPT",
			},
			"last_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"material_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_rotation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Aliyun_KMS", "EXTERNAL"}, false),
				Default:      "Aliyun_KMS",
			},
			"pending_window_in_days": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(7, 30),
				Optional:     true,
				Default:      7,
			},
			"deletion_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(7, 30),
				Deprecated:   "Field 'deletion_window_in_days' has been deprecated from provider version 1.85.0. New field 'pending_window_in_days' instead.",
			},
			"primary_key_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM"}, false),
				Default:      "SOFTWARE",
			},
			"rotation_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := kms.CreateCreateKeyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v, ok := d.GetOk("automatic_rotation"); ok {
		request.EnableAutomaticRotation = requests.NewBoolean(convertAutomaticRotationRequest(v.(string)))
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("key_usage"); ok {
		request.KeyUsage = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok {
		request.Origin = v.(string)
	}
	if v, ok := d.GetOk("protection_level"); ok {
		request.ProtectionLevel = v.(string)
	}
	if v, ok := d.GetOk("rotation_interval"); ok {
		request.RotationInterval = v.(string)
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateKey(request)
	})
	bresponse, ok := raw.(*kms.CreateKeyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kms_key", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw)
	d.SetId(fmt.Sprintf("%v", bresponse.KeyMetadata.KeyId))

	return resourceAlibabacloudStackKmsKeyRead(d, meta)
}

func resourceAlibabacloudStackKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("arn", object.Arn)
	d.Set("automatic_rotation", object.AutomaticRotation)
	d.Set("creation_date", object.CreationDate)
	d.Set("creator", object.Creator)
	d.Set("delete_date", object.DeleteDate)
	d.Set("description", object.Description)
	d.Set("key_state", object.KeyState)
	d.Set("key_usage", object.KeyUsage)
	d.Set("last_rotation_date", object.LastRotationDate)
	d.Set("material_expire_time", object.MaterialExpireTime)
	d.Set("next_rotation_date", object.NextRotationDate)
	d.Set("origin", object.Origin)
	d.Set("primary_key_version", object.PrimaryKeyVersion)
	d.Set("protection_level", object.ProtectionLevel)
	d.Set("rotation_interval", object.RotationInterval)
	return nil
}

func resourceAlibabacloudStackKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kmsService := KmsService{client}
	d.Partial(true)

	if d.HasChange("description") {
		request := kms.CreateUpdateKeyDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.KeyId = d.Id()
		request.Description = d.Get("description").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateKeyDescription(request)
		})
		bresponse, ok := raw.(*kms.UpdateKeyDescriptionResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	update := false
	request := kms.CreateUpdateRotationPolicyRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.KeyId = d.Id()
	if d.HasChange("automatic_rotation") {
		update = true
	}
	request.EnableAutomaticRotation = requests.NewBoolean(convertAutomaticRotationRequest(d.Get("automatic_rotation").(string)))
	if d.HasChange("rotation_interval") {
		update = true
		request.RotationInterval = d.Get("rotation_interval").(string)
	}
	if update {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateRotationPolicy(request)
		})
		bresponse, ok := raw.(*kms.UpdateRotationPolicyResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	if d.HasChange("key_state") || d.HasChange("is_enabled") {
		object, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			return errmsgs.WrapError(err)
		}
		var target = ""
		if k, ok := d.GetOk("key_state"); ok {
			target = k.(string)
		} else {
			if k, ok := d.GetOk("is_enabled"); ok {
				if k.(bool) {
					target = "Enable"
				} else {
					target = "Disabled"
				}
			}
		}

		if object.KeyState != target {
			if target == "Disabled" {
				request := kms.CreateDisableKeyRequest()
				client.InitRpcRequest(*request.RpcRequest)
				request.KeyId = d.Id()
				raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
					return kmsClient.DisableKey(request)
				})
				bresponse, ok := raw.(*kms.DisableKeyResponse)
				addDebug(request.GetActionName(), raw)
				if err != nil {
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
			}
			if target == "Enabled" {
				request := kms.CreateEnableKeyRequest()
				client.InitRpcRequest(*request.RpcRequest)
				request.KeyId = d.Id()
				raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
					return kmsClient.EnableKey(request)
				})
				bresponse, ok := raw.(*kms.EnableKeyResponse)
				addDebug(request.GetActionName(), raw)
				if err != nil {
					errmsg := ""
					if ok {
						errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
					}
					return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
				}
			}
		}
	}
	d.Partial(false)
	return resourceAlibabacloudStackKmsKeyRead(d, meta)
}

func resourceAlibabacloudStackKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := kms.CreateScheduleKeyDeletionRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.KeyId = d.Id()
	if v, ok := d.GetOk("pending_window_in_days"); ok {
		request.PendingWindowInDays = requests.NewInteger(v.(int))
	} else {
		if v, ok := d.GetOk("deletion_window_in_days"); ok {
			request.PendingWindowInDays = requests.NewInteger(v.(int))
		}
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.ScheduleKeyDeletion(request)
	})
	bresponse, ok := raw.(*kms.ScheduleKeyDeletionResponse)
	addDebug(request.GetActionName(), raw)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return nil
}

func convertAutomaticRotationRequest(source string) bool {
	switch source {
	case "Disabled":
		return false
	case "Enabled":
		return true
	}
	return false
}
