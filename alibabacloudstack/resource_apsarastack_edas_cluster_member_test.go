package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEdasClusterMember_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alibabacloudstack_edas_cluster_member.default"
	ra := resourceAttrInit(resourceId, EdasClusterMemberBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeClusterMember")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasClusterMemberConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":  "${alibabacloudstack_edas_cluster.default.id}",
					"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":        CHECKSET,
						"instance_id":       CHECKSET,
						"cluster_member_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var EdasClusterMemberBasicMap = map[string]string{}

func resourceEdasClusterMemberConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "logical_id" {
  default = "%s:%s"
}

%s

resource "alibabacloudstack_edas_namespace" "default" {
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  logical_region_id = "${alibabacloudstack_edas_namespace.default.namespace_logical_id}"
  network_mode = "2"
  cluster_type = "2"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

`, name, defaultRegionToTest, name, ECSInstanceCommonTestCase)
}
