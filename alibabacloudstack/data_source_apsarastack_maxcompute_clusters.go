package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeClustersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"core_arch": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeClustersRead(d *schema.ResourceData, meta interface{}) error {
	objects, err := DescribeMaxcomputeProject(meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	var r *regexp.Regexp
	if rt, ok := d.GetOk("name_regex"); ok && rt.(string) != "" {
		r = regexp.MustCompile(rt.(string))
	}
	var t []map[string]interface{}
	var ids []string
	for _, object := range objects {
		cluster_raw := object.(map[string]interface{})
		if r != nil && !r.MatchString(cluster_raw["cluster"].(string)) {
			continue
		}
		cluster := map[string]interface{}{
			"cluster":   cluster_raw["cluster"].(string),
			"core_arch": cluster_raw["core_arch"].(string),
			"project":   cluster_raw["project"].(string),
			"region":    cluster_raw["region"].(string),
		}
		t = append(t, cluster)
		ids = append(ids, cluster["cluster"].(string))
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("clusters", t); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), t)
	}
	return nil
}

func DescribeMaxcomputeProject(meta interface{}) ([]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}

	roleId, err := client.RoleIds()
	if err != nil {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return nil, err
	}

	request := map[string]interface{}{
		"CurrentRoleId":   roleId,
	}

	response, err = client.DoTeaRequest("POST", "ASCM", "2019-05-10", "ListOdpsClusters", "", nil, request)

	if err != nil {
		if errmsgs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute Cluster", "ListOdpsClusters")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		return nil, err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = errmsgs.Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return nil, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(response)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListOdpsClusters", "$", response, errmsg)
	}
	return v.(map[string]interface{})["data"].([]interface{}), nil
}
