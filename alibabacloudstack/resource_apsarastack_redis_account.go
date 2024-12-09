package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackKVstoreAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackKVStoreAccountCreate,
		Read:   resourceAlibabacloudStackKVStoreAccountRead,
		Update: resourceAlibabacloudStackKVStoreAccountUpdate,
		Delete: resourceAlibabacloudStackKVStoreAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
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
				ValidateFunc: validation.StringInSlice([]string{"Normal"}, false),
				ForceNew:     true,
				Default:      "Normal",
			},
			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RoleReadOnly", "RoleReadWrite", "RoleRepl"}, false),
				Default:      "RoleReadWrite",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "Field 'description' is deprecated and will be removed in a future release. Please use new field 'account_description' instead.",
				ConflictsWith: []string{"account_description"},
			},
			"account_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"description"},
			},
		},
	}
}

func resourceAlibabacloudStackKVStoreAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	request := r_kvstore.CreateCreateAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Get("instance_id").(string)
	request.AccountName = d.Get("account_name").(string)
	request.AccountPrivilege = d.Get("account_privilege").(string)
	password := d.Get("account_password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return errmsgs.WrapError(errmsgs.Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
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
	request.AccountType = d.Get("account_type").(string)

	if v, ok := connectivity.GetResourceDataOk(d, "account_description", "description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := kvstoreService.WaitForKVstoreInstance(request.InstanceId, Normal, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.CreateAccount(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*r_kvstore.CreateAccountResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kvstore_account", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.InstanceId, COLON_SEPARATED, request.AccountName))

	if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return resourceAlibabacloudStackKVStoreAccountRead(d, meta)
}

func resourceAlibabacloudStackKVStoreAccountRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreAccount(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	connectivity.SetResourceData(d, object.AccountDescription, "account_description", "description")
	d.Set("account_privilege", object.DatabasePrivileges.DatabasePrivilege[0].AccountPrivilege)

	return nil
}

func resourceAlibabacloudStackKVStoreAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChanges("account_description", "description") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := r_kvstore.CreateModifyAccountDescriptionRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountDescription = connectivity.GetResourceData(d, "account_description", "description").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyAccountDescription(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*r_kvstore.ModifyAccountDescriptionResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("account_privilege") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := r_kvstore.CreateGrantAccountPrivilegeRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountPrivilege = d.Get("account_privilege").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*r_kvstore.GrantAccountPrivilegeResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChanges("account_password", "kms_encrypted_password") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
		request := r_kvstore.CreateResetAccountPasswordRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return errmsgs.WrapError(errmsgs.Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
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

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ResetAccountPassword(request)
		})
		if err != nil {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*r_kvstore.ResetAccountPasswordResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAlibabacloudStackKVStoreAccountRead(d, meta)
}

func resourceAlibabacloudStackKVStoreAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := r_kvstore.CreateDeleteAccountRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteAccount(request)
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
			return nil
		} else {
			errmsg := ""
			if raw != nil {
				response, ok := raw.(*r_kvstore.DeleteAccountResponse)
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return kvstoreService.WaitForKVstoreAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
