package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSKubernetesClustersKubeConfigDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testdataSourceCSKubernetesClustersKubeconfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig", "kubeconfig"),
				),
			},
		},
	})
}

const testdataSourceCSKubernetesClustersKubeconfig = `
data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
	cluster_id = "c47f06f706de24b21bc43c33ee07f2163"
  }
`
