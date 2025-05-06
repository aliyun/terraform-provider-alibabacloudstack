package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackKmsCiphertext() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackKmsCiphertextCreate, schema.Noop, nil, resourceAlibabacloudStackKmsCiphertextDelete)
	return resource
}

func resourceAlibabacloudStackKmsCiphertextCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Since a ciphertext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	request := kms.CreateEncryptRequest()
	client.InitRpcRequest(*request.RpcRequest)

	request.Plaintext = d.Get("plaintext").(string)
	request.KeyId = d.Get("key_id").(string)

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convertMaptoJsonString(cm)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_kms_ciphertext", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		request.EncryptionContext = string(contextJson)
	}

	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Encrypt(request)
	})
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*kms.EncryptResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kms_ciphertext", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.EncryptResponse)
	d.Set("ciphertext_blob", response.CiphertextBlob)

	return nil
}

func resourceAlibabacloudStackKmsCiphertextDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
