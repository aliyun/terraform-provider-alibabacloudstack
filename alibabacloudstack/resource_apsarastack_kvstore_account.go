package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackKVStoreAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	request := r_kvstore.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = d.Get("instance_id").(string)
	request.AccountName = d.Get("account_name").(string)
	request.AccountPrivilege = d.Get("account_privilege").(string)
	password := d.Get("account_password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)
	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		request.AccountPassword = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.AccountPassword = decryptResp.Plaintext
	}
	request.AccountType = d.Get("account_type").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := kvstoreService.WaitForKVstoreInstance(request.InstanceId, Normal, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.CreateAccount(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_kvstore_account", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.InstanceId, COLON_SEPARATED, request.AccountName))

	if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlibabacloudStackKVStoreAccountRead(d, meta)
}

func resourceAlibabacloudStackKVStoreAccountRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("description", object.AccountDescription)
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

	if d.HasChange("description") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateModifyAccountDescriptionRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("description").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyAccountDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("description")
	}

	if d.HasChange("account_privilege") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateGrantAccountPrivilegeRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountPrivilege = d.Get("account_privilege").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("account_privilege")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateResetAccountPasswordRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.InstanceId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.AccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AlibabacloudStackClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext

		}

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ResetAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		//d.SetPartial("kms_encrypted_password")
		//d.SetPartial("kms_encryption_context")
		//d.SetPartial("account_password")
	}

	d.Partial(false)
	return resourceAlibabacloudStackKVStoreAccountRead(d, meta)
}

func resourceAlibabacloudStackKVStoreAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := r_kvstore.CreateDeleteAccountRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.InstanceId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteAccount(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
			return nil
		} else {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return kvstoreService.WaitForKVstoreAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
