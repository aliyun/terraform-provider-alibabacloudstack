package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackBmsKeypair() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"key_pair_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource, resourceAlibabacloudStackBmsKeypairCreate,
		resourceAlibabacloudStackBmsKeypairRead,
		nil,
		resourceAlibabacloudStackBmsKeypairDelete)
	return resource
}

func resourceAlibabacloudStackBmsKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqBody := map[string]interface{}{
		"cloudType":  "private",
		"DeployType": "bms",
		"Name":       d.Get("name").(string),
	}

	if v, ok := d.GetOk("name"); ok {
		reqBody["Name"] = v.(string)
	}
	if v, ok := d.GetOk("public_key"); ok {
		reqBody["PublicKey"] = v.(string)
	}

	resp, err := client.DoTeaRequest("POST", "bms", "2024-03-01", "CreateKeyPair", "", nil, nil, reqBody)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "CreateKeyPair", "POST", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	data, err := jsonpath.Get("$.data.data", resp)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "CreateKeyPair", "$.data.data", resp)
	}
	keypairData := data.(map[string]interface{})
	name, ok := keypairData["name"].(string)
	if !ok || name == "" {
		return fmt.Errorf("failed to get keypair name from response")
	}

	if v, ok := keypairData["privateKey"].(string); ok && v != "" {
		d.Set("private_key", v)
	}
	if v, ok := keypairData["publicKey"].(string); ok && v != "" {
		d.Set("public_key", v)
	}
	d.SetId(name)
	return nil
}

func resourceAlibabacloudStackBmsKeypairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	bmsService := BmsService{client}
	target, err := bmsService.DescribeKeyPair(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "resource_apsarastack_bms_keypair", "ReadKeyPair", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("name", target["name"])
	if v, ok := target["privateKey"].(string); ok && v != "" {
		d.Set("private_key", v)
	}
	if v, ok := target["publicKey"].(string); ok && v != "" {
		d.Set("public_key", v)
	}
	d.Set("key_pair_fingerprint", target["keyPairFingerPrint"])

	return nil
}

func resourceAlibabacloudStackBmsKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"KeyPairName": d.Get("name"),
	}

	_, err := client.DoTeaRequest("POST", "bms", "2024-03-01", "DeleteKeyPair", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteKeyPair", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
