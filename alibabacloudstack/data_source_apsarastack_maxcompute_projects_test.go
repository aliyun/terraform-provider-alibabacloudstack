package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMaxcomputeProjectsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999) * 2
	resourceId := "data.alibabacloudstack_maxcompute_projects.default"
	name := fmt.Sprintf("tf_testAcck%d", rand)
	userid, userpk := InitPreCreateMaxcomputeUser(t, name)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return dataSourceMaxcomputeProjectsConfigDependence(name, userid, userpk)
	})

	// Test with name filter (should return the created project)
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name": "${alibabacloudstack_maxcompute_project.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name": "fake-nonexistent-project",
		}),
	}

	// Test with status filter (default status is "0")
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "0",
			"name":   "${alibabacloudstack_maxcompute_project.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "1",
			"name":   "${alibabacloudstack_maxcompute_project.default.name}",
		}),
	}

	var existMaxcomputeProjectsMapFunc = func(rand int) map[string]string {
		projectName := fmt.Sprintf("tf_testAcck%d", rand)
		return map[string]string{
			"ids.#":                 "1",
			"projects.#":            "1",
			"projects.0.id":         CHECKSET,
			"projects.0.name":       projectName,
			"projects.0.account":    CHECKSET,
			"projects.0.account_pk": CHECKSET,
			"projects.0.quota_id":   CHECKSET,
			"projects.0.disk":       "50",
			"projects.0.core_arch":  CHECKSET,
			"projects.0.vpc_ids.#":  "1",
		}
	}

	var fakeMaxcomputeProjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"projects.#": "0",
		}
	}

	var maxcomputeProjectsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMaxcomputeProjectsMapFunc,
		fakeMapFunc:  fakeMaxcomputeProjectsMapFunc,
	}
	maxcomputeProjectsCheckInfo.dataSourceTestCheck(t, rand, nameConf, statusConf)
}

func dataSourceMaxcomputeProjectsConfigDependence(name, userid, userpk string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "userid" {
  default = "%s"
}

variable "userpk" {
  default = "%s"
}


data "alibabacloudstack_maxcompute_clusters" "default" {
  name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_maxcompute_cu" "default" {
  cu_name      = var.name
  cu_num       = 2
  cluster_name = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}

resource "alibabacloudstack_maxcompute_project" "default" {
  account_pk     = var.userpk
  quota_id       = alibabacloudstack_maxcompute_cu.default.id
  external_table = "true"
  vpc_ids        = [alibabacloudstack_vpc_vpc.default.id]
  name           = var.name
  disk           = 50
  account        = var.userid
}
`, name, userid, userpk)
}
