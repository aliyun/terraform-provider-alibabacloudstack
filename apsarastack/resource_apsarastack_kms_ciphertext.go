package apsarastack

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackKmsCiphertext() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackKmsCiphertextCreate,
		Read:   schema.Noop,
		Delete: resourceApsaraStackKmsCiphertextDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
}

func resourceApsaraStackKmsCiphertextCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	// Since a ciphertext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	request := kms.CreateEncryptRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.Plaintext = d.Get("plaintext").(string)
	request.KeyId = d.Get("key_id").(string)
	request.RegionId = client.RegionId

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convertMaptoJsonString(cm)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kms_ciphertext", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		request.EncryptionContext = string(contextJson)
	}

	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Encrypt(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kms_ciphertext", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.EncryptResponse)
	d.Set("ciphertext_blob", response.CiphertextBlob)

	return nil
}

func resourceApsaraStackKmsCiphertextDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
