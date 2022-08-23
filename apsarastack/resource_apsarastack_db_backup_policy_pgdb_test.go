package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDBBackupPolicy_pgdb(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_instance.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicyPgdb"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyPostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
					"instance_name":    "${var.name}",
					"engine":           "PostgreSQL",
					"engine_version":   "10.0",
					"storage_type":     "local_ssd",
					"instance_type":    "pg.x4.large.2",
					"instance_storage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encryption": "false",
					}),
				),
			},
		},
	})
}

func resourceDBBackupPolicyPostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}
data "apsarastack_zones" "default" {
  available_resource_creation = "Rds"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}`, name)
}
