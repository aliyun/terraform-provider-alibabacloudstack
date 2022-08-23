package apsarastack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApsaraStackMaxcomputeClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackMaxcomputeClustersRead,
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

func dataSourceApsaraStackMaxcomputeClustersRead(d *schema.ResourceData, meta interface{}) error {

	objects, err := DescribeMaxcomputeProject(meta)
	if err != nil {
		return WrapError(err)
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
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), t)
	}
	return nil
}

func DescribeMaxcomputeProject(meta interface{}) ([]interface{}, error) {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	conn, err := client.NewAscmClient()
	if err != nil {
		return nil, WrapError(err)
	}

	action := "ListOdpsClusters"

	roleId, err := client.RoleIds()
	if err != nil {
		err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
		return nil, err
	}

	request := map[string]interface{}{
		"RegionName":      client.RegionId,
		"Region":          client.RegionId,
		"RegionId":        client.RegionId,
		"Product":         "ascm",
		"OrganizationId":  client.Department,
		"ResourceGroupId": client.ResourceGroup,
		"CurrentRoleId":   roleId,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, "ListOdpsClusters", action, ApsaraStackSdkGoERROR)
		return nil, err
	}
	addDebug(action, response, request)

	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("Maxcompute Cluster", "ListOdpsClusters")), NotFoundMsg, ProviderERROR)
		return nil, err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return nil, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, "ListOdpsClusters", "$", response)
	}
	return v.(map[string]interface{})["data"].([]interface{}), nil
}
