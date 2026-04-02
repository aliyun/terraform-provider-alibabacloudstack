package alibabacloudstack

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackBcmpKeyPair() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackBcmpKeyPairCreate,
		resourceAlibabacloudStackBcmpKeyPairRead, resourceAlibabacloudStackBcmpKeyPairUpdate, resourceAlibabacloudStackBcmpKeyPairDelete)

	return resource
}

func resourceAlibabacloudStackBcmpKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	keyName := d.Get("key_pair_name").(string)

	request := map[string]interface{}{
		"cloudType": "private",
		"Name":      keyName,
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		request["PublicKey"] = publicKey.(string)
	}

	response, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "CreateKeyPair", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_keypair", "CreateKeyPair", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	if !response["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_bcmp_keypair", "CreateKeyPair", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data := response["data"].(map[string]interface{})
	keyPairData := data["data"].(map[string]interface{})

	d.SetId(keyPairData["name"].(string))

	if file, ok := d.GetOk("key_file"); ok {
		privateKey := keyPairData["privateKey"].(string)
		if privateKey != "" {
			ioutil.WriteFile(file.(string), []byte(privateKey), 0600)
			os.Chmod(file.(string), 0400)
		}
	}

	return nil
}

func resourceAlibabacloudStackBcmpKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlibabacloudStackBcmpKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}

	response, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListKeyPair", "", nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_keypair", "ListKeyPair", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	if !response["success"].(bool) {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_bcmp_keypair", "ListKeyPair", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data := response["data"].([]interface{})
	keyPairId := d.Id()

	for _, item := range data {
		keyPair := item.(map[string]interface{})
		if keyPair["name"].(string) == keyPairId {
			d.Set("key_pair_name", keyPair["name"])
			d.Set("finger_print", keyPair["keyPairFingerPrint"])
			d.Set("region", keyPair["region"])
			d.Set("resource_group_name", keyPair["ResourceGroupName"])
			d.Set("create_time", keyPair["createTime"])
			d.Set("update_time", keyPair["updateTime"])
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceAlibabacloudStackBcmpKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"KeyPairName": d.Id(),
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "DeleteKeyPair", "", nil, request, nil)
		if err != nil {
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteKeyPair", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}

		if !response["success"].(bool) {
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteKeyPair", errmsgs.AlibabacloudStackSdkGoERROR))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
