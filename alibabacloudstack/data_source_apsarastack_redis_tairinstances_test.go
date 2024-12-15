package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"fmt"
)

func TestAccAlibabacloudStackKVStoreInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceKVStoreInstancesConfigDependence(),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_kvstore_instances.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.charge_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.region_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.create_time"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.instance_class"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.status"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.availability_zone"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.bandwidth"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.user_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.connections"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_kvstore_instances.default", "instances.connection_domain"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_kvstore_instances.default", "ids.#"),
				),
			},
		},
	})
}

func dataSourceKVStoreInstancesConfigDependence() string {
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstancesDataSource%d"
}

%s

resource "alibabacloudstack_kvstore_instance" "default" {
instance_class = "redis.master.small.default"
instance_name  = var.name
vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
security_ips   = ["10.0.0.1"]
instance_type  = "Redis"
engine_version = "4.0"
}
data "alibabacloudstack_kvstore_instances" "default" {
  name_regex = alibabacloudstack_kvstore_instance.default.instance_name
}
`, getAccTestRandInt(10000, 99999), VSwitchCommonTestCase)
}
