package apsarastack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackSlbCACertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSlbCACertificatesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackSlbCACertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := slb.CreateDescribeCACertificatesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_slb_ca_certificates", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeCACertificatesResponse)
	var filteredTemp []slb.CACertificate
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, certificate := range response.CACertificates.CACertificate {
			if r != nil && !r.MatchString(certificate.CACertificateName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[certificate.CACertificateId]; !ok {
					continue
				}
			}

			filteredTemp = append(filteredTemp, certificate)
		}
	} else {
		filteredTemp = response.CACertificates.CACertificate
	}

	return slbCACertificatesDescriptionAttributes(d, filteredTemp, meta)
}

func slbCACertificatesDescriptionAttributes(d *schema.ResourceData, certificates []slb.CACertificate, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, certificate := range certificates {

		mapping := map[string]interface{}{
			"id":                certificate.CACertificateId,
			"name":              certificate.CACertificateName,
			"fingerprint":       certificate.Fingerprint,
			"created_timestamp": certificate.CreateTimeStamp,
			"region_id":         certificate.RegionId,
		}
		ids = append(ids, certificate.CACertificateId)
		names = append(names, certificate.CACertificateName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("certificates", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
