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

func dataSourceAlibabacloudStackAPIGatewayV2Certificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2CertificatesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sni": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cert_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
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
						"snis": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// dataSourceAlibabacloudStackAPIGatewayV2CertificateDescriptionRead performs the AlibabacloudStack Image lookup.
func dataSourceAlibabacloudStackAPIGatewayV2CertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"gwInstanceId": d.Get("instance_id").(string),
		"regionId":     client.RegionId,
	}
	if v, ok := d.GetOk("sni"); ok {
		request["sni"] = v.(string)
	}
	if v, ok := d.GetOk("cert_type"); ok {
		request["certType"] = v.(string)
	}
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListCertificates", "/certificate/listCertificates", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	records, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap :=getIdsStringFilter(d)
	var ids []string
	var certificates []map[string]interface{}
	for _, v := range records.([]interface{}) {
		record := v.(map[string]interface{})
		id := fmt.Sprintf("%s:%v", d.Get("instance_id").(string), record["certificateId"])
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(record["certificateName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[id]; !exist {
				continue
			}
		}
		certificates = append(certificates, map[string]interface{}{
			"id":               id,
			"certificate_id":   record["certificateId"],
			"certificate_name": record["certificateName"],
			"expire_time":      record["expireTime"],
			"update_time":      record["updateTime"],
			"create_time":      record["createTime"],
			"cert_type":        record["certType"],
			"snis":             record["snis"],
		})
		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("certificates", certificates); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
