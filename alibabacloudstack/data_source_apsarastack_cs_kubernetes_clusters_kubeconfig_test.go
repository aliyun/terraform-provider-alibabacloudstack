package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSKubernetesClustersKubeConfigDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacckubernetesconfig-%d", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: testdataSourceCSKubernetesClustersKubeconfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig", "kubeconfig"),
				),
			},
		},
	})
}

func testdataSourceCSKubernetesClustersKubeconfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
		default = "%s"
	}
	
%s

data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
	cluster_id = local.k8s_cluster_id
	private_address = true
  }
`, name, AckK8sCommonTestCase())
}
