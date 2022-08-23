package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_db_instance" {
			continue
		}
		request := rds.CreateDescribeBackupPolicyRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = rs.Primary.ID
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func TestAccApsaraStackDBBackupPolicy_mysql(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "apsarastack_db_instance.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyMysqlConfigDependence)
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
					"engine":           "MySQL",
					"engine_version":   "5.6",
					"storage_type":     "local_ssd",
					"instance_type":    "rds.mysql.s2.large",
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

func resourceDBBackupPolicyMysqlConfigDependence(name string) string {
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
