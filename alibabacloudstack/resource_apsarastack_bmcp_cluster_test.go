package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackBmcpCluster_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_bmcp_cluster.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccBmcpClusterCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeBmcpCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccBmcpCluster%d", rand)
	password := getAccTestPassword(12)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccBmcpClusterBasicDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":        name,
					"vpc_id":              "${alibabacloudstack_vpc.default.id}",
					"evpc_id":             "${alibabacloudstack_evpc_evpc.default.id}",
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"standard_vswitch_id": "${alibabacloudstack_vswitch.standard.id}",
					"password":            password,
					"machine_type":        "${data.alibabacloudstack_bmcp_machinetypes.all.machinetypes.0.name}",
					"node_count":          1,
					"vswitch_id":          "${alibabacloudstack_vswitch.default.id}",
					"lifecycle": []map[string]interface{}{{
						"ignore_changes": TfRawString("[machine_type]"),
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name,
						"node_count":   "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "is_create_cpfs_cluster", "is_install_yundun_aegis", "standard_vswitch_id", "cluster_arch_type", "machine_type", "vswitch_id", "machine_type"},
			},
		},
	})
}

var AlibabacloudTestAccBmcpClusterCheckmap = map[string]string{
	"cluster_id":   CHECKSET,
	"cluster_name": CHECKSET,
	"status":       CHECKSET,
	"region_id":    CHECKSET,
	"node_count":   CHECKSET,
	"cpu_count":    CHECKSET,
	"mem_count":    CHECKSET,
	"flops_count":  CHECKSET,
	"video_memory": CHECKSET,
	"create_time":  CHECKSET,
	"update_time":  CHECKSET,
}

func AlibabacloudTestAccBmcpClusterBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	name              = "${var.name}"
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "192.168.40.0/24"
	is_cgw            = true
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id # Availability zone
}

resource "alibabacloudstack_vswitch" "standard" {
	name              = "${var.name}-std"
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "192.168.50.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id # Availability zone
}

resource "alibabacloudstack_evpc_evpc" "default" {
	evpc_name   = "${var.name}"
	description = "${var.name}"
}

data "alibabacloudstack_bmcp_machinetypes" "all" {
    min_standard_instance_count = 1
}
`, name)
}
