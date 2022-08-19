package apsarastack

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDBConnectionConfigUpdate(t *testing.T) {
	var v *rds.DBInstanceNetInfo
	var rdsEndpoint string
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)

	if rdsEndpoint = os.Getenv("RDS_ENDPOINT"); rdsEndpoint == "" {
		if rdsEndpoint = os.Getenv("APSARASTACK_DOMAIN"); rdsEndpoint == "" {
			t.Fatal("APSARASTACK_DOMAIN must be set for acceptance tests")
		}
		rdsEndpoint = regexp.MustCompile(`.*\.(intra\..*\.com)\/.*`).FindStringSubmatch(rdsEndpoint)[1]
	}

	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("tf-testacc%s.mysql.rds.%s", rand, rdsEndpoint),
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "apsarastack_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${apsarastack_db_instance.instance.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
		},
	})
}

func resourceDBConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
provider "apsarastack" {
	assume_role {}
}
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}

	resource "apsarastack_db_instance" "instance" {
	  engine               = "MySQL"
	  engine_version       = "5.6"
	  instance_type        = "rds.mysql.s2.large"
	  instance_storage     = "5"
	  instance_charge_type = "Postpaid"
	  instance_name        = "${var.name}"
	  vswitch_id           = "${apsarastack_vswitch.default.id}"
	  monitoring_period    = "60"
	  storage_type         = "local_ssd"
	}
	`, RdsCommonTestCase, name)
}
