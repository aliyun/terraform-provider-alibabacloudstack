package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSlbCACertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbCACertificatesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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

func dataSourceAlibabacloudStackSlbCACertificatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateDescribeCACertificatesRequest()
	client.InitRpcRequest(*request.RpcRequest)

	idsMap := getIdsStringFilter(d)
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	response, ok := raw.(*slb.DescribeCACertificatesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_ca_certificates", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var filteredTemp []slb.CACertificate
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}
	for _, certificate := range response.CACertificates.CACertificate {
		if nameRegex != nil && !nameRegex.MatchString(certificate.CACertificateName) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[certificate.CACertificateId]; !ok {
				continue
			}
		}

		filteredTemp = append(filteredTemp, certificate)
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
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
