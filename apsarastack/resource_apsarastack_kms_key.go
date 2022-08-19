package apsarastack

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackKmsKeyCreate,
		Read:   resourceApsaraStackKmsKeyRead,
		Update: resourceApsaraStackKmsKeyUpdate,
		Delete: resourceApsaraStackKmsKeyDelete,
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

func resourceApsaraStackKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := kms.CreateCreateKeyRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kms_key", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*kms.CreateKeyResponse)
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	d.SetId(fmt.Sprintf("%v", response.KeyMetadata.KeyId))

	return resourceApsaraStackKmsKeyRead(d, meta)
}
func resourceApsaraStackKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
func resourceApsaraStackKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	kmsService := KmsService{client}
	d.Partial(true)

	if d.HasChange("description") {
		request := kms.CreateUpdateKeyDescriptionRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

		request.KeyId = d.Id()
		request.Description = d.Get("description").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.UpdateKeyDescription(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("description")
	}
	update := false
	request := kms.CreateUpdateRotationPolicyRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		//d.SetPartial("automatic_rotation")
		//d.SetPartial("rotation_interval")
	}
	if d.HasChange("key_state") || d.HasChange("is_enabled") {
		object, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
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
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

				request.KeyId = d.Id()
				raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
					return kmsClient.DisableKey(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
				}
				//d.SetPartial("key_state")
				//d.SetPartial("is_enabled")
			}
			if target == "Enabled" {
				request := kms.CreateEnableKeyRequest()
				request.Headers = map[string]string{"RegionId": client.RegionId}
				request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

				request.KeyId = d.Id()
				raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
					return kmsClient.EnableKey(request)
				})
				addDebug(request.GetActionName(), raw)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
				}
				//d.SetPartial("key_state")
				//d.SetPartial("is_enabled")
			}
		}
	}
	d.Partial(false)
	return resourceApsaraStackKmsKeyRead(d, meta)
}
func resourceApsaraStackKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := kms.CreateScheduleKeyDeletionRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

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
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
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
