package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNetwokerAclAttachment0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_network_acl_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNetworkAclAttachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetwork_acl_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNetworkAclAttachmentdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"network_acl_id": "${alibabacloudstack_network_acl.default.id}",
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alibabacloudstack_vswitch.default.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_acl_id": CHECKSET,

						"resources.#": "1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"network_acl_id": "${alibabacloudstack_network_acl.default.id}",
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alibabacloudstack_vswitch.default.id}",
							"resource_type": "VSwitch",
						},
						{
							"resource_id":   "${alibabacloudstack_vswitch.default2.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_acl_id": CHECKSET,

						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"network_acl_id": "${alibabacloudstack_network_acl.default.id}",
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${alibabacloudstack_vswitch.default2.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_acl_id": CHECKSET,

						"resources.#": "1",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccNetworkAclAttachmentCheckmap = map[string]string{

	// "status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccNetworkAclAttachmentdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_zones" "default" {
		available_resource_creation= "VSwitch"
	}
	
	resource "alibabacloudstack_vpc" "default" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
	
	resource "alibabacloudstack_network_acl" "default" {
		vpc_id = "${alibabacloudstack_vpc.default.id}"
		network_acl_name = "${var.name}"
	}
	
	
	resource "alibabacloudstack_vswitch" "default" {
		vpc_id = "${alibabacloudstack_vpc.default.id}"
		cidr_block = "172.16.0.0/24"
		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		name = "${var.name}"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		vpc_id = "${alibabacloudstack_vpc.default.id}"
		cidr_block = "172.16.1.0/24"
		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		name = "${var.name}"
	}
	
	`, name)
}
