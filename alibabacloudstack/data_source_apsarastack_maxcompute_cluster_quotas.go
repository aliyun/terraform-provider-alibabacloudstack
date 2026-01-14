package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackMaxcomputeClusterQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMaxcomputeClusterQutaosRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
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

	cluster := d.Get("cluster").(string)

	roleId, err := client.RoleIds()
	if err != nil {
		return err
	}

	request := map[string]interface{}{
		"CurrentRoleId": roleId,
		"Cluster":       cluster,
	}

	response, err = client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetOdpsQuotaForAscm", "", nil, request, nil)
	if err != nil {
		if errmsgs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute Cluster", cluster)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		return err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		errmsg := ""
		if response != nil {
			errmsg = errmsgs.GetAsapiErrorMessage(response)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, cluster, "$", response, errmsg)
		return err
	}
	object := v.(map[string]interface{})["Data"].(map[string]interface{})

	d.Set("cu_total", object["cuTotal"].(string))
	d.Set("disk_available", object["diskAvailable"].(string))
	d.Set("cu_available", object["cuAvailable"].(string))
	d.Set("disk_total", object["diskTotal"].(string))

	d.SetId(d.Get("cluster").(string))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), d); err != nil {
			return err
		}
	}
	return nil
}
