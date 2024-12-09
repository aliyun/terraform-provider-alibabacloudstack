package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAlikafkaSaslUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAlikafkaSaslUserCreate,
		Read:   resourceAlibabacloudStackAlikafkaSaslUserRead,
		Update: resourceAlibabacloudStackAlikafkaSaslUserUpdate,
		Delete: resourceAlibabacloudStackAlikafkaSaslUserDelete,
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

func resourceAlibabacloudStackAlikafkaSaslUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	instanceId := d.Get("instance_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	usertype := d.Get("type").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	request := alikafka.CreateCreateSaslUserRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Username = username
	request.Type = usertype
	if password != "" {
		request.Password = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt2(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.Password = decryptResp
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.CreateSaslUser(request)
		})
		bresponse, ok := raw.(*alikafka.CreateSaslUserResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "AlibabacloudStack_alikafka_sasl_user", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	// Server may have cache, sleep a while.
	time.Sleep(2 * time.Second)
	d.SetId(instanceId + ":" + username + ":" + usertype)
	return resourceAlibabacloudStackAlikafkaSaslUserUpdate(d, meta)
}

func resourceAlibabacloudStackAlikafkaSaslUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := alikafkaService.DescribeAlikafkaSaslUser(d.Id())
	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", parts[0])
	d.Set("username", object.Username)
	d.Set("password", object.Password)
	d.Set("type", object.Type)

	return nil
}

func resourceAlibabacloudStackAlikafkaSaslUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	usertype := parts[2]

	if d.HasChanges("password", "kms_encrypted_password") {
		request := alikafka.CreateCreateSaslUserRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = instanceId
		request.Username = username
		request.Type = usertype
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.Password = password
		} else {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt2(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request.Password = decryptResp
		}

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.CreateSaslUser(request)
			})
			bresponse, ok := raw.(*alikafka.CreateSaslUserResponse)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
					time.Sleep(2 * time.Second)
					return resource.RetryableError(err)
				}
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "AlibabacloudStack_alikafka_sasl_user", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return err
		}

		// Server may have cache, sleep a while.
		time.Sleep(1000)
	}
	return resourceAlibabacloudStackAlikafkaSaslUserRead(d, meta)
}

func resourceAlibabacloudStackAlikafkaSaslUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	instanceId := parts[0]
	username := parts[1]
	usertype := parts[2]

	request := alikafka.CreateDeleteSaslUserRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = instanceId
	request.Username = username
	request.Type = usertype
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteSaslUser(request)
		})
		bresponse, ok := raw.(*alikafka.DeleteSaslUserResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser, "ONS_SYSTEM_FLOW_CONTROL"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	return errmsgs.WrapError(alikafkaService.WaitForAlikafkaSaslUser(d.Id(), Deleted, DefaultTimeoutMedium))
}
