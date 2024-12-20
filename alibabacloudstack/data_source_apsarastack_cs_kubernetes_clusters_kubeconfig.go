package alibabacloudstack

import (
	"encoding/json"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/denverdino/aliyungo/cs"

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

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	cluster_id := d.Get("cluster_id").(string)
	request.Method = "GET"
	request.Product = "Cs"
	request.Version = "2015-12-15"

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	request.ServiceCode = "cs"
	request.ApiName = "DescribeClusterUserKubeconfig"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"Product":       "Cs",
		"RegionId":      client.RegionId,
		"Action":        "DescribeClusterUserKubeconfig",
		"Version":       cs.CSAPIVersion,
		"Department":    client.Department,
		"ResourceGroup": client.ResourceGroup,
		"ClusterId":     cluster_id,
	}
	request.RegionId = client.RegionId
	raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
		return csClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_cs_kubernetes_clusters_kubeconfig", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request, request.QueryParams)
	response := make(map[string]interface{})
	resp, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(resp.GetHttpContentBytes(), &response)
	d.SetId(cluster_id)
	if err := d.Set("kubeconfig", response["config"]); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}
