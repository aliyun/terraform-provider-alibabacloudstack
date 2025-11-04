package alibabacloudstack

import (

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

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
			"private_address": {
				Type:     schema.TypeBool,
				Default:  false,
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
	csService := CsService{client}
	clusterId := d.Get("cluster_id").(string)
	if config, err := csService.GetK8sCluterKubeConfig(clusterId, d.Get("private_address").(bool)); err != nil {
		return err
	} else {
		d.Set("kubeconfig", config)
		d.SetId(clusterId)
		return nil
	}
}
