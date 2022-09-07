package alibabacloudstack

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeClusterQutaos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeClusterQutaosRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cu_total": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"disk_available": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"cu_available": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"disk_total": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackMaxcomputeClusterQutaosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	conn, err := client.NewAscmClient()
	if err != nil {
		return WrapError(err)
	}

	action := "GetOdpsQuota"

	cluster := d.Get("cluster").(string)

	roleId, err := client.RoleIds()
	if err != nil {
		err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
		return err
	}

	request := map[string]interface{}{
		"RegionName":      client.RegionId,
		"Region":          client.RegionId,
		"RegionId":        client.RegionId,
		"Product":         "ascm",
		"OrganizationId":  client.Department,
		"ResourceGroupId": client.ResourceGroup,
		"CurrentRoleId":   roleId,
		"Cluster":         cluster,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, cluster, action, AlibabacloudStackSdkGoERROR)
		return err
	}
	addDebug(action, response, request)

	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("Maxcompute Cluster", cluster)), NotFoundMsg, ProviderERROR)
		return err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, cluster, "$", response)
	}
	object := v.(map[string]interface{})["data"].(map[string]interface{})

	d.Set("cu_total", object["cuTotal"].(string))
	d.Set("disk_available", object["diskAvailable"].(string))
	d.Set("cu_available", object["cuAvailable"].(string))
	d.Set("disk_total", object["diskTotal"].(string))

	d.SetId(d.Get("cluster").(string))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), d)
	}
	return nil
}
