package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hsm_cluster.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HsmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHsmCluster")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_cluster%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildHsmClusterDependencies)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName:     resourceId,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":       "${var.name}",
					"master_instance_id": "${alibabacloudstack_hsm_instance.default0.id}",
					"vpc_id":             "${alibabacloudstack_vpc.default.id}",
					"vswitch_ids":        "${alibabacloudstack_vswitch.default.id}",
					"zone_nos":           "${data.alibabacloudstack_zones.default.zones.0.id}",
					"ip_white_list":      "123.12.13.1/16,124.13.14.2/18",
					"password":           "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":  name,
						"vpc_id":        CHECKSET,
						"vswitch_ids":   CHECKSET,
						"zone_nos":      CHECKSET,
						"ip_white_list": "123.12.13.1/16,124.13.14.2/18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":  "${var.name}_updated",
					"ip_white_list": "123.12.13.1/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name":  name + "_updated",
						"ip_white_list": "123.12.13.1/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_instance_ids": []string{"${alibabacloudstack_hsm_instance.default1.id}", "${alibabacloudstack_hsm_instance.default2.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sub_instance_ids": []string{"${alibabacloudstack_hsm_instance.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sub_instance_ids.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func buildHsmClusterDependencies(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" "default" {
  provider = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

%s

data "alibabacloudstack_hsm_vendors" "default" {
}

resource "alibabacloudstack_hsm_instance" "default0" {
	product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_no = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark = "${var.name}"
}

resource "alibabacloudstack_hsm_instance" "default1" {
	product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_no = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark = "${var.name}"
}

resource "alibabacloudstack_hsm_instance" "default2" {
	product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_no = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark = "${var.name}"
}
`, name, RandomPasswordTestCase(12, 1))
}
