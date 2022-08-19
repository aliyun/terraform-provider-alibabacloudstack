package apsarastack

import (
	"fmt"
	"testing"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackGPDBAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_gpdb_account.default"
	ra := resourceAttrInit(resourceId, ApsaraStackGPDBAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeGpdbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackGPDBAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":      "${apsarastack_gpdb_instance.default.id}",
					"account_name":        name,
					"account_password":    "inputYourCodeHere",
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":        name,
						"account_description": name,
						"db_instance_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "inputYourCodeHere" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var ApsaraStackGPDBAccountMap0 = map[string]string{}

func ApsaraStackGPDBAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "apsarastack_gpdb_zones" "default" {}
data "apsarastack_zones" "default" {}
data "apsarastack_vpcs" "default" {
  name_regex = "default-NODELETING"
}
resource "apsarastack_vpc" "default" {
name       = var.name
cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  
  vpc_id       = "${apsarastack_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
   name              = "${var.name}"
}

resource "apsarastack_gpdb_instance" "default" {
  availability_zone      = "${data.apsarastack_zones.default.zones.0.id}"
  engine                 = "gpdb"
  engine_version         = "4.3"
  instance_class         = "gpdb.group.segsdx2"
  instance_group_count   = 2
  description            = "tf-testAccGpdbInstance_new"
  vswitch_id             = "${apsarastack_vswitch.default.id}"
}
`, name)
}
