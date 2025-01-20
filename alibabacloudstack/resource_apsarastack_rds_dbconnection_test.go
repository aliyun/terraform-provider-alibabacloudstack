package alibabacloudstack

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBConnectionConfigUpdate(t *testing.T) {
	var v *rds.DBInstanceNetInfo
	var rdsEndpoint string
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-testAccDBconnection%d", rand)

	if rdsEndpoint = os.Getenv("RDS_ENDPOINT"); rdsEndpoint == "" {
		if rdsEndpoint = os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"); rdsEndpoint == "" {
			t.Fatal("ALIBABACLOUDSTACK_POPGW_DOMAIN must be set for acceptance tests")
		}
		rdsEndpoint = regexp.MustCompile(`.*\.(intra\..*\.com)\/.*`).FindStringSubmatch(rdsEndpoint)[1]
	}

	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("tf-testacc%d.mysql.rds.%s", rand, rdsEndpoint),
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "alibabacloudstack_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
	ResourceTest(t, resource.TestCase{
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
					"instance_id":       "${alibabacloudstack_db_instance.instance.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%d", rand),
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

	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}

	resource "alibabacloudstack_db_instance" "instance" {
	  engine               = "MySQL"
	  engine_version       = "5.6"
	  instance_type        = "rds.mysql.s2.large"
	  instance_storage     = "5"
	  instance_charge_type = "Postpaid"
	  instance_name        = "${var.name}"
	  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
	  monitoring_period    = "60"
	  storage_type         = "local_ssd"
	}
	`, RdsCommonTestCase, name)
}
