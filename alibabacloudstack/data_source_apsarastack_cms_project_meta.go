package alibabacloudstack

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCmsProjectMeta() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCmsProjectMetaRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"namespace": {
				Type:         schema.TypeString,
				Optional:     true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCmsProjectMetaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "Cms", "2019-01-01", "DescribeProjectMeta", "")

	response := Data{}

	for {
		log.Printf(" request of DescribeProjectMeta : %v", request)
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		log.Printf(" response of raw DescribeProjectMeta : %s", bresponse)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cms_project_meta", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		if response.Success == true || len(response.Resources.Resource) < 1 {
			break
		}
	}

	var r *regexp.Regexp
	if rt, ok := d.GetOk("name_regex"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var namespace string
	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
	}
	var ids []string
	var s []map[string]interface{}
	for _, rg := range response.Resources.Resource {
		if r != nil && !r.MatchString(rg.Description) {
			continue
		}
		if namespace != "" && namespace != rg.Namespace {
			continue
		}
		mapping := map[string]interface{}{
			"description": rg.Description,
			"namespace":   rg.Namespace,
		}
		if len(rg.Labels) > 0 {
			labels := make([]map[string]string, 0)
			for _, label := range rg.Labels {
				label_map := map[string]string{
					"name":  label.Name,
					"value": label.Value,
				}
				labels = append(labels, label_map)
			}
			mapping["labels"] = labels
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("resources", s); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}

type Data struct {
	PageSize   int    `json:"PageSize"`
	RequestID  string `json:"RequestId"`
	PageNumber int    `json:"PageNumber"`
	Total      int    `json:"Total"`
	Resources  struct {
		Resource []struct {
			Description string `json:"Description"`
			Labels      []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"Labels"`
			Namespace string `json:"Namespace"`
		} `json:"Resource"`
	} `json:"Resources"`
	Code    int  `json:"Code"`
	Success bool `json:"Success"`
}
