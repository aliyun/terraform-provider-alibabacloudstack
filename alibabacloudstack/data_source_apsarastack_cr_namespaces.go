package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCRNamespacesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_create": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_visibility": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCRNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := CrService{client}

	request := client.NewCommonRequest("POST", "cr", "2016-06-07", "GetNamespaceList", "")

	raw, err := client.WithEcsClient(func(crClient *ecs.Client) (interface{}, error) {
		return crClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var crResp crListResponse
	log.Printf("response for datasource %v", bresponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &crResp)
	log.Printf("unmarshalled response for datasource %v", crResp)

	addDebug(request.GetActionName(), bresponse)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var names []string
	var s []map[string]interface{}

	for _, ns := range crResp.Data.Namespaces {
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r := regexp.MustCompile(nameRegex.(string))
			if !r.MatchString(ns.Namespace) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"name": ns.Namespace,
		}

		raw, err := crService.DescribeCrNamespace(ns.Namespace)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cr_namespace", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		mapping["auto_create"] = raw.Data.Namespace.AutoCreate
		mapping["default_visibility"] = raw.Data.Namespace.DefaultVisibility

		names = append(names, ns.Namespace)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("namespaces", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
