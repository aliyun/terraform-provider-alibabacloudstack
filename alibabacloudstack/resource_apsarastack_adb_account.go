package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAdbAccountCreate,
		Read:   resourceAlibabacloudStackAdbAccountRead,
		Update: resourceAlibabacloudStackAdbAccountUpdate,
		Delete: resourceAlibabacloudStackAdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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

			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string("Super")}, false),
				Default:      "Super",
				ForceNew:     true,
				//Removed:      "Field 'account_type' has been removed from provider version 1.81.0.",
			},

			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackAdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}
	request := adb.CreateCreateAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = d.Get("db_cluster_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)
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

	if v, ok := d.GetOk("account_description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.CreateAccount(request)
		})
		response, ok := raw.(*adb.CreateAccountResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_adb_account", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.NonRetryableError(err)

		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.AccountName))

	if err := adbService.WaitForAdbAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackAdbAccountRead(d, meta)
}

func resourceAlibabacloudStackAdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbAccount(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_description", object.AccountDescription)
	d.Set("account_type", object.AccountType)

	return nil
}

func resourceAlibabacloudStackAdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChanges("account_password", "kms_encrypted_password") {
		if err := adbService.WaitForAdbAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := adb.CreateResetAccountPasswordRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.DBClusterId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)
		if password == "" && kmsPassword == "" {
			return errmsgs.WrapError(errmsgs.Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.AccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errmsgs.WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext
		}

		raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
			return adbClient.ResetAccountPassword(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*adb.ResetAccountPasswordResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAlibabacloudStackAdbAccountRead(d, meta)
}

func resourceAlibabacloudStackAdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	adbService := AdbService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := adb.CreateDeleteAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.DBClusterId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithAdbClient(func(adbClient *adb.Client) (interface{}, error) {
		return adbClient.DeleteAccount(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
			return nil
		}
		errmsg := ""
		if raw != nil {
			response, ok := raw.(*adb.DeleteAccountResponse)
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return adbService.WaitForAdbAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
