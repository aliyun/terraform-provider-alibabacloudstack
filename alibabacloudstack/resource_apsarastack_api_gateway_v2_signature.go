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
			"cascade_link_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2SignatureCreate, resourceAlibabacloudStackAPIGatewayV2SignatureRead, resourceAlibabacloudStackAPIGatewayV2SignatureUpdate, resourceAlibabacloudStackAPIGatewayV2SignatureDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2SignatureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	gwInstanceId := d.Get("gw_instance_id").(string)
	cascade_link_ids := d.Get("cascade_link_ids").(*schema.Set).List()
	reqBody := map[string]interface{}{
		"sigSchemeName": d.Get("sig_scheme_name").(string),
		"sigAlg":        d.Get("sig_alg").(string),
		"gwInstanceId":  gwInstanceId,
	}

	action := "CreateSignatureScheme"
	pattern := "/signatureScheme/createSignatureScheme"
	idpre := "sig"
	if len(cascade_link_ids) > 0 {
		cascadeLinkIds := make([]string, 0)
		for _, item := range cascade_link_ids {
			cascadeLinkIds = append(cascadeLinkIds, item.(string))
		}
		reqBody["cascadeLinkIds"] = cascadeLinkIds
		action = "CreateSourceSigScheme"
		pattern = "/sourceSigScheme/createSourceSigScheme"
		idpre = "sourceSig"
	}

	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody)
	if err != nil {
		return err
	}

	sigSchemeId, ok := resp["data"].(string)
	if !ok {
		return fmt.Errorf("failed to get sigSchemeId from response")
	}

	// Generate resource ID
	resourceId := fmt.Sprintf("%s:%s:%s", idpre, gwInstanceId, sigSchemeId)
	d.SetId(resourceId)
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2SignatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	apiGatewayV2Service := ApiGateWayV2Service{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return err
	}
	gwInstanceId := parts[1]
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
	} else {
		d.Set("sig_alg", object["sigAlg"])
	}
	d.Set("gw_instance_id", gwInstanceId)
	d.Set("sig_scheme_name", object["sigSchemeName"])
	d.Set("sig_scheme_id", object["sigSchemeId"])
	d.Set("secret_key", object["secretKey"])
	d.Set("create_time", object["createTime"])
	d.Set("update_time", object["updateTime"])
	d.Set("status", object["sigSchemeStatus"])
	linkEntityRelations, ok := object["linkEntityRelations"]
	if ok {
		linkIds := make([]string, 0)
		for _, v := range linkEntityRelations.([]interface{}) {
			linkEntityRelation := v.(map[string]interface{})
			linkIds = append(linkIds, linkEntityRelation["linkId"].(string))
		}
		d.Set("cascade_link_ids", linkIds)
	}
	return nil
}

func resourceAlibabacloudStackAPIGatewayV2SignatureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return err
	}
	idpre := parts[0]
	gwInstanceId := parts[1]
	sigSchemeId := parts[2]

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("sig_scheme_name") {
		reqBody := map[string]interface{}{
			"sigSchemeId":   sigSchemeId,
			"sigSchemeName": d.Get("sig_scheme_name"),
			"gwInstanceId":  gwInstanceId,
		}
		action := "ModifySignatureScheme"
		pattern := "/signatureScheme/modifySignatureScheme"
		if idpre == "sourceSig" {
			action = "ModifySourceSigScheme"
			pattern = "/sourceSigScheme/modifySourceSigScheme"
		}

		if _, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody); err != nil {
			return fmt.Errorf("failed to modify signature scheme: %v", err)
		}
	}
	if d.HasChanges("status") {
		reqBody := map[string]interface{}{
			"sigSchemeId":     sigSchemeId,
			"sigSchemeStatus": d.Get("status"),
			"gwInstanceId":    gwInstanceId,
		}
		action := "ModifySignatureSchemeStatus"
		pattern := "/signatureScheme/modifySignatureSchemeStatus"
		if idpre == "sourceSig" {
			action = "ModifySourceSigSchemeStatus"
			pattern = "/sourceSigScheme/modifySourceSigSchemeStatus"
		}
		if _, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqBody); err != nil {
			return fmt.Errorf("failed to modify signature scheme: %v", err)
		}
	}

	return resourceAlibabacloudStackAPIGatewayV2SignatureRead(d, meta)
}

func resourceAlibabacloudStackAPIGatewayV2SignatureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	gwInstanceId := parts[1]
	sigSchemeId := parts[2]
	idpre := parts[0]
	reqQuery := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"sigSchemeId":  sigSchemeId,
	}
	
	action := "DeleteSignatureScheme"
	pattern := "/signatureScheme/deleteSignatureScheme"
	if idpre == "sourceSig" {
		action = "DeleteSourceSigScheme"
		pattern = "/sourceSigScheme/deleteSourceSigScheme"
	}

	raw, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, raw)
	}

	return nil
}
