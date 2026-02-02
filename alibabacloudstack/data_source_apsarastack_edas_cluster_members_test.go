package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEdasClusterMembersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_edas_cluster_members.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf%d", rand),
		dataSourceEdasClusterMembersDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_cluster_member.default.id}"},
			"cluster_id": "${alibabacloudstack_edas_cluster_member.default.cluster_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_cluster_member.default.id}-fake"},
			"cluster_id": "${alibabacloudstack_edas_cluster_member.default.cluster_id}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alibabacloudstack_edas_cluster_member.default.cluster_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alibabacloudstack_edas_cluster_member.default.cluster_id}-fake",
		}),
	}

	var existEdasClusterMembersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"ids.0":                 CHECKSET,
			"members.#":             "1",
			"members.0.cluster_id":  CHECKSET,
			"members.0.instance_id": CHECKSET,
			"members.0.ecu_id":      CHECKSET,
			"members.0.ecs_id":      CHECKSET,
			"members.0.status":      CHECKSET,
			"members.0.create_time": CHECKSET,
			"members.0.update_time": CHECKSET,
		}
	}

	var fakeEdasClusterMembersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"members.#": "0",
		}
	}

	var EdasClusterMembersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasClusterMembersMapFunc,
		fakeMapFunc:  fakeEdasClusterMembersMapFunc,
	}

	EdasClusterMembersCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

func dataSourceEdasClusterMembersDependence(name string) string {
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

resource "alibabacloudstack_edas_cluster_member" "default" { 
    cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

`, name, defaultRegionToTest, name, ECSInstanceCommonTestCase)
}
