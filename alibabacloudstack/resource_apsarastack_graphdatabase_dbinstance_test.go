package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackGraphdatabaseDbinstance0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_graphdatabase_dbinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccGraphdatabaseDbinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoGdbDescribedbinstanceaccesswhitelistRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraph_databasedb_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccGraphdatabaseDbinstanceBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"db_version": "1.0",

					"db_instance_network_type": "VPC",

					"payment_type": "PayAsYouGo",

					"region_id": "cn-hangzhou",

					"db_instance_storage_type": "cloud_ssd",

					"db_instance_description": "ssd测试",

					"db_node_class": "gdb.r.2xlarge",

					"db_instance_category": "ha",

					"zone_id": "cn-hangzhou-h",

					"vpc_id": "vpc-bp1bvsykm9f9hkfeikfi5",

					"vswitch_id": "vsw-bp152wgftimgq80eiii6k",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"db_version": "1.0",

						"db_instance_network_type": "VPC",

						"payment_type": "PayAsYouGo",

						"region_id": "cn-hangzhou",

						"db_instance_storage_type": "cloud_ssd",

						"db_instance_description": "ssd测试",

						"db_node_class": "gdb.r.2xlarge",

						"db_instance_category": "ha",

						"zone_id": "cn-hangzhou-h",

						"vpc_id": "vpc-bp1bvsykm9f9hkfeikfi5",

						"vswitch_id": "vsw-bp152wgftimgq80eiii6k",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"db_instance_storage_type": "cloud_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"db_instance_storage_type": "cloud_ssd",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccGraphdatabaseDbinstanceCheckmap = map[string]string{

	"db_instance_network_type": CHECKSET,

	"port": CHECKSET,

	"db_instance_type": CHECKSET,

	"db_instance_storage_type": CHECKSET,

	"db_node_storage": CHECKSET,

	"master_db_instance_id": CHECKSET,

	"db_instance_category": CHECKSET,

	"db_version": CHECKSET,

	"current_minor_version": CHECKSET,

	"payment_type": CHECKSET,

	"public_connection_string": CHECKSET,

	"db_instance_id": CHECKSET,

	"db_node_class": CHECKSET,

	"lock_reason": CHECKSET,

	"maintain_time": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"zone_id": CHECKSET,

	"create_time": CHECKSET,

	"vswitch_id": CHECKSET,

	"db_node_count": CHECKSET,

	"expired": CHECKSET,

	"latest_minor_version": CHECKSET,

	"lock_mode": CHECKSET,

	"db_instance_memory": CHECKSET,

	"read_only_db_instance_ids": CHECKSET,

	"db_instance_ip_array": CHECKSET,

	"vpc_id": CHECKSET,

	"db_instance_cpu": CHECKSET,

	"db_instance_description": CHECKSET,

	"region_id": CHECKSET,

	"connection_string": CHECKSET,

	"expire_time": CHECKSET,

	"public_port": CHECKSET,
}

func AlibabacloudTestAccGraphdatabaseDbinstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
