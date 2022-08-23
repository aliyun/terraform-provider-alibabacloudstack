package apsarastack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApsaraStackAlikafkaSaslUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAlikafkaSaslUserCreate,
		Read:   resourceApsaraStackAlikafkaSaslUserRead,
		Update: resourceApsaraStackAlikafkaSaslUserUpdate,
		Delete: resourceApsaraStackAlikafkaSaslUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"type": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ValidateFunc: validation.StringInSlice([]string{
					"plain",
					"scram",
				}, false),
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
		},
	}
}

func resourceApsaraStackAlikafkaSaslUserCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	regionId := client.RegionId
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	usertype := d.Get("type").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	request := alikafka.CreateCreateSaslUserRequest()
	request.InstanceId = instanceId
	request.RegionId = regionId
	request.Username = username
	request.Type = usertype
	request.Domain = client.Config.AlikafkaOpenAPIEndpoint
	if password != "" {
		request.Password = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt2(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.Password = decryptResp
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateSaslUser(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_alikafka_sasl_user", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	// Server may have cache, sleep a while.
	time.Sleep(2 * time.Second)
	d.SetId(instanceId + ":" + username + ":" + usertype)
	return resourceApsaraStackAlikafkaSaslUserUpdate(d, meta)
}

func resourceApsaraStackAlikafkaSaslUserRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := alikafkaService.DescribeAlikafkaSaslUser(d.Id())
	if err != nil {
		// Handle exceptions
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("username", object.Username)
	d.Set("password", object.Password)
	d.Set("type", object.Type)

	return nil
}

func resourceApsaraStackAlikafkaSaslUserUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	usertype := parts[2]

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {

		request := alikafka.CreateCreateSaslUserRequest()
		request.InstanceId = instanceId
		request.RegionId = client.RegionId
		request.Username = username
		request.Domain = client.Config.AlikafkaOpenAPIEndpoint
		request.Type = usertype
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.Password = password
		} else {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt2(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.Password = decryptResp
		}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.CreateSaslUser(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(2 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "ApsaraStack_alikafka_sasl_user", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		// Server may have cache, sleep a while.
		time.Sleep(1000)
	}
	return resourceApsaraStackAlikafkaSaslUserRead(d, meta)
}

func resourceApsaraStackAlikafkaSaslUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	usertype := parts[2]

	request := alikafka.CreateDeleteSaslUserRequest()
	request.RegionId = client.RegionId
	request.InstanceId = instanceId
	request.Username = username
	request.Type = usertype
	request.Domain = client.Config.AlikafkaOpenAPIEndpoint
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteSaslUser(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return WrapError(alikafkaService.WaitForAlikafkaSaslUser(d.Id(), Deleted, DefaultTimeoutMedium))
}
