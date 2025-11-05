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

func dataSourceAlibabacloudStackAPIGatewayV2Domains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2DomainsRead,

		Schema: map[string]*schema.Schema{
			"domain_regex": {
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
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTPS", "HTTP"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values.
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_auth": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer_dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceAlibabacloudStackAPIGatewayV2DomainDescriptionRead performs the AlibabacloudStack Image lookup.
func dataSourceAlibabacloudStackAPIGatewayV2DomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	instanceId := d.Get("instance_id").(string)
	request := map[string]interface{}{
		"gwInstanceId": instanceId,
		"regionId":     client.RegionId,
	}
	if v, ok := d.GetOk("domain"); ok {
		request["domain"] = v.(string)
	}
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListDomains", "/domain/listDomains", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	records, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	var ids []string
	var domains []map[string]interface{}
	for _, v := range records.([]interface{}) {
		record := v.(map[string]interface{})
		id := fmt.Sprintf("%s:%v", instanceId, record["domainId"])
		if domain_regex, ok := d.GetOk("domain_regex"); ok {
			r := regexp.MustCompile(domain_regex.(string))
			if !r.MatchString(record["domain"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[id]; !exist {
				continue
			}
		}
		if protocol, ok := d.GetOk("protocol"); ok && record["protocol"].(string) != protocol.(string) {
			continue
		}
		domains = append(domains, map[string]interface{}{
			"id":                id,
			"instance_id":       instanceId,
			"domain":            record["domain"],
			"domain_id":         record["domainId"],
			"protocol":          record["protocol"],
			"certificate_id":    record["certificateId"],
			"client_auth":       record["clientAuth"],
			"ca_certificate_id": record["caCertificateId"],
			"subject_dn":        record["subjectDn"],
			"issuer_dn":         record["issuerDn"],
		})
		ids = append(ids, id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("domains", domains); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
