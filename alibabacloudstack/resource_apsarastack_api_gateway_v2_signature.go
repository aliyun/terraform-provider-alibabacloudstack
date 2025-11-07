package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAPIGatewayV2Signature() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sig_scheme_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sig_alg": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HmacSHA256", "HmacSHA1", "HmacSM3"}, false),
			},
			"sig_scheme_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2SignatureCreate, resourceAlibabacloudStackAPIGatewayV2SignatureRead, resourceAlibabacloudStackAPIGatewayV2SignatureUpdate, resourceAlibabacloudStackAPIGatewayV2SignatureDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2SignatureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gwInstanceId := d.Get("gw_instance_id").(string)

	reqBody := map[string]interface{}{
		"sigSchemeName": d.Get("sig_scheme_name").(string),
		"sigAlg":        d.Get("sig_alg").(string),
		"gwInstanceId":  gwInstanceId,
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "CreateSignatureScheme", "/signatureScheme/createSignatureScheme", nil, nil, reqBody)
	if err != nil {
		return err
	}

	sigSchemeId, ok := resp["data"].(string)
	if !ok {
		return fmt.Errorf("failed to get sigSchemeId from response")
	}

	// Generate resource ID
	resourceId := fmt.Sprintf("%s:%s", gwInstanceId, sigSchemeId)
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2SignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGateWayV2Service{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}
	gwInstanceId := parts[0]
	object, err := apiGatewayV2Service.DescribeApiGatewayV2Signature(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_api_gateway_v2_signature apiGatewayV2Service.DescribeApiGatewayV2Signature Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	var APIGateWayV2SignatureAlgMap = map[string]string{
		"HMAC-SHA256": "HmacSHA256",
		"HMAC-SHA1":   "HmacSHA1",
		"HMAC-SM3":    "HmacSM3",
	}
	sigAlg, ok := APIGateWayV2SignatureAlgMap[object["sigAlg"].(string)]
	if ok {
		d.Set("sig_alg", sigAlg)
	}
	d.Set("gw_instance_id", gwInstanceId)
	d.Set("sig_scheme_name", object["sigSchemeName"])
	d.Set("sig_scheme_id", object["sigSchemeId"])
	d.Set("secret_key", object["secretKey"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("status", object["sigSchemeStatus"])

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2SignatureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return err
	}
	gwInstanceId := parts[0]
	sigSchemeId := parts[1]

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("sig_scheme_name") {
		reqBody := map[string]interface{}{
			"sigSchemeId":   sigSchemeId,
			"sigSchemeName": d.Get("sig_scheme_name"),
			"gwInstanceId":  gwInstanceId,
		}

		if _, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifySignatureScheme", "/signatureScheme/modifySignatureScheme", nil, nil, reqBody); err != nil {
			return fmt.Errorf("failed to modify signature scheme: %v", err)
		}
	}
	if d.HasChanges("status") {
		reqBody := map[string]interface{}{
			"sigSchemeId":     sigSchemeId,
			"sigSchemeStatus": d.Get("status"),
			"gwInstanceId":    gwInstanceId,
		}

		if _, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ModifySignatureSchemeStatus", "/signatureScheme/modifySignatureSchemeStatus", nil, nil, reqBody); err != nil {
			return fmt.Errorf("failed to modify signature scheme: %v", err)
		}
	}

	return resourceAlibabacloudStackAPIGatewayV2SignatureRead(d, meta)
}

func resourceAlibabacloudStackAPIGatewayV2SignatureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[0]
	sigSchemeId := parts[1]

	reqQuery := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"sigSchemeId":  sigSchemeId,
	}

	raw, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "DeleteSignatureScheme", "/signatureScheme/deleteSignatureScheme", nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteSignatureScheme", raw)
	}

	return nil
}
