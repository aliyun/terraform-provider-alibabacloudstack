package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDBAccountCreate,
		Read:   resourceAlibabacloudStackDBAccountRead,
		Update: resourceAlibabacloudStackDBAccountUpdate,
		Delete: resourceAlibabacloudStackDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"data_base_instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: validation.StringLenBetween(6, 32),
			},
			"kms_encrypted_password": {
				Type:            schema.TypeString,
				Optional:        true,
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
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Super"}, false),
				ForceNew:     true,
				Computed:     true,
			},

			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				Deprecated:  "Field 'instance_id' is deprecated and will be removed in a future release. Please use 'data_base_instance_id' instead.",
			},
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				Deprecated:  "Field 'name' is deprecated and will be removed in a future release. Please use 'account_name' instead.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Super"}, false),
				ForceNew:     true,
				Computed:     true,
				Deprecated:  "Field 'type' is deprecated and will be removed in a future release. Please use 'account_type' instead.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:  "Field 'description' is deprecated and will be removed in a future release. Please use 'account_description' instead.",
			},
		},
	}
}

func resourceAlibabacloudStackDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	request := rds.CreateCreateAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "data_base_instance_id", "instance_id"); err == nil {
		request.DBInstanceId = v.(string)
	} else {
		return err
	}
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "account_name", "name"); err == nil {
		request.AccountName = v.(string)
	} else {
		return err
	}

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		request.AccountPassword = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.AccountPassword = decryptResp.Plaintext
	}
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "account_type", "type"); err == nil {
		request.AccountType = v.(string)
	} else {
		return err
	}

	// Description will not be set when account type is normal and it is a API bug
	if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "account_description", "description"); err == nil && v.(string) != "" {
		request.AccountDescription = v.(string)
	} else if err != nil {
		return err
	}
	// wait instance running before modifying
	if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.CreateAccount(request)
		})
		bresponse, ok := raw.(*rds.CreateAccountResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_db_account", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBInstanceId, COLON_SEPARATED, request.AccountName))

	if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackDBAccountRead(d, meta)
}

func resourceAlibabacloudStackDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBAccount(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	connectivity.SetResourceData(d, object.DBInstanceId, "data_base_instance_id", "instance_id")
	connectivity.SetResourceData(d, object.AccountName, "account_name", "name")
	connectivity.SetResourceData(d, object.AccountType, "account_type", "type")
	connectivity.SetResourceData(d, object.AccountDescription, "account_description", "description")

	return nil
}

func resourceAlibabacloudStackDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("account_description") || d.HasChange("description") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := rds.CreateModifyAccountDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = instanceId
		request.AccountName = accountName
		if v, err := connectivity.GetResourceData(d, reflect.TypeOf(""), "account_description", "description"); err == nil {
			request.AccountDescription = v.(string)
		} else {
			return err
		}

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyAccountDescription(request)
		})
		bresponse, ok := raw.(*rds.ModifyAccountDescriptionResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("account_description")
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := rds.CreateResetAccountPasswordRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = instanceId
		request.AccountName = accountName

		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			//d.SetPartial("password")
			request.AccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext
			//d.SetPartial("kms_encrypted_password")
			//d.SetPartial("kms_encryption_context")
		}

		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ResetAccountPassword(request)
		})
		bresponse, ok := raw.(*rds.ResetAccountPasswordResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlibabacloudStackDBAccountRead(d, meta)
}

func resourceAlibabacloudStackDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := rds.CreateDeleteAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DeleteAccount(request)
	})
	bresponse, ok := raw.(*rds.DeleteAccountResponse)
	if err != nil && !errmsgs.IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return rdsService.WaitForAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
