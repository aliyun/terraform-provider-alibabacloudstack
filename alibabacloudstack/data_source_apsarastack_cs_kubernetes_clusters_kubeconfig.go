package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackCSKubernetesClustersKubeConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCSKubernetesClustersKubeConfigRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackCSKubernetesClustersKubeConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Get("cluster_id").(string)
	request := client.NewCommonRequest("GET", "CS", "2015-12-15", "DescribeClusterUserKubeconfig", fmt.Sprintf("/k8s/%s/user_config", clusterId))
	resp, err := client.ProcessCommonRequest(request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_cs_kubernetes_clusters_kubeconfig", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), resp, request, request.QueryParams)
	response := make(map[string]interface{})
	err = json.Unmarshal(resp.GetHttpContentBytes(), &response)
	d.SetId(clusterId)
	if err := d.Set("kubeconfig", response["config"]); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
