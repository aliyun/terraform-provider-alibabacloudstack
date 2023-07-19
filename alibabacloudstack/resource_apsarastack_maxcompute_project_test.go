package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcompute_basic(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccack%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMaxcomputeConfigBasic, name, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
						"disk": "50",
					}),
				),
			},
			{
				Config: fmt.Sprintf(testAccMaxcomputeConfigBasic, name, 55),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
						"disk": "55",
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeConfigBasic = `

variable "name" {
	default = "%s"
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_maxcompute_cu" "default"{
  cu_name      = "${var.name}"
  cu_num       = "1"
  cluster_name = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}

resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = "${var.name}"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,       
    ]
  }
}

resource "alibabacloudstack_maxcompute_project" "default"{
	cluster        = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
	external_table = "false"
	quota_id       = alibabacloudstack_maxcompute_cu.default.id
	disk           = %d
	name           = "${var.name}"
	aliyun_account = "${alibabacloudstack_maxcompute_user.default.user_name}"
    pk = "1075451910171540"
}
`

func TestAccAlibabacloudStackMaxcompute_advance(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_project.default.4"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := "tf_testAccMCProject"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMaxcomputeConfigAdvance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":       name + "4",
						"specification_type": "OdpsStandard",
						"order_type":         "PayAsYouGo",
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeConfigAdvance = `
variable "name" {
	default = "%s"
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_maxcompute_cu" "default"{
  cu_name      = "${var.name}"
  cu_num       = "1"
  cluster_name = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}

resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = "${var.name}"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
}

resource "alibabacloudstack_maxcompute_project" "default"{
	cluster        = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
	external_table = "false"
	quota_id       = alibabacloudstack_maxcompute_cu.default.id
	disk           = 50
	name           = "${var.name}"
	aliyun_account = "${alibabacloudstack_maxcompute_user.default.user_name}"
    pk = "1075451910171540"
}
`
