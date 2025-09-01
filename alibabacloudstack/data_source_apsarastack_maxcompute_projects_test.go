package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeProjectsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceAlibabacloudstackMaxcomputeProjects(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_maxcompute_projects.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_projects.default", "projects.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_projects.default", "projects.name"),
				),
			},
		},
	})
}

func datasourceAlibabacloudstackMaxcomputeProjects(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}


resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = var.name
  description           = "maxcomput project test"
}

resource "alibabacloudstack_maxcompute_cu" "default" {
	cu_name =      "${var.name}"
	cu_num =       2
	cluster_name = "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}"
}


resource "alibabacloudstack_maxcompute_project" "default" {
  account_pk = "${alibabacloudstack_maxcompute_user.default.user_pk}"
  quota_id = "${alibabacloudstack_maxcompute_cu.default.id}"
  external_table = "true"
  vpc_ids = [
              "${alibabacloudstack_vpc_vpc.default.id}"
            ]
  name = "${var.name}"
  disk = "50"
  account = "${alibabacloudstack_maxcompute_user.default.user_id}"
}

data "alibabacloudstack_maxcompute_projects" "default"{
	name = var.name
}
`, name)
}
