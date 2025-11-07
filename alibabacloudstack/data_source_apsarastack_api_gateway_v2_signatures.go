package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2Signatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2SignaturesRead,

		Schema: map[string]*schema.Schema{
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"signatures": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gw_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sig_scheme_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sig_scheme_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sig_alg": {
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
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2SignaturesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	gwInstanceId := d.Get("gw_instance_id").(string)

	// Prepare filters
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Call API to list signature schemes
	requestBody := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"current":      1,
		"size":         100,
	}

	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListSignatureSchemes", "/signatureScheme/listSignatureSchemes", nil, nil, requestBody)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	records, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	signatures := make([]map[string]interface{}, 0)
	names := make([]string, 0)
	ids := make([]string, 0)

	var APIGateWayV2SignatureAlgMap = map[string]string{
		"HMAC-SHA256": "HmacSHA256",
		"HMAC-SHA1":   "HmacSHA1",
		"HMAC-SM3":    "HmacSM3",
	}

	for _, record := range records.([]interface{}) {
		r := record.(map[string]interface{})

		sigSchemeId := r["sigSchemeId"].(string)
		sigSchemeName := r["sigSchemeName"].(string)
		id := fmt.Sprintf("%s:%s", gwInstanceId, sigSchemeId)
		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		if nameRegex != nil && !nameRegex.MatchString(sigSchemeName) {
			continue
		}

		sigAlg, ok := APIGateWayV2SignatureAlgMap[r["sigAlg"].(string)]
		if !ok {
			sigAlg = r["sigAlg"].(string)
		}

		signature := map[string]interface{}{
			"id":              id,
			"gw_instance_id":  gwInstanceId,
			"sig_scheme_id":   r["sigSchemeId"],
			"sig_scheme_name": r["sigSchemeName"],
			"sig_alg":         sigAlg,
			"secret_key":      r["secretKey"],
			"create_time":     r["createTime"],
			"update_time":     r["updateTime"],
			"status":          r["sigSchemeStatus"],
		}

		signatures = append(signatures, signature)
		names = append(names, r["sigSchemeName"].(string))
		ids = append(ids, id)
	}

	// Set computed fields
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("signatures", signatures); err != nil {
		return err
	}
	if err := d.Set("names", names); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}
