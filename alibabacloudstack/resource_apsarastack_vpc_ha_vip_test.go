package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcHavip_basic(t *testing.T) {
	var v *VpcDescribehavipsResponse
	resourceId := "alibabacloudstack_vpc_ha_vip.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribehavipsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testAccVpcHavipBasic_%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpcHavipBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ha_vip_name":              "${var.name}",
					"description":              "${var.name}",
					"ip_address":               "172.16.1.88",
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.default.id}",
					"vpc_id":                   "${alibabacloudstack_vpc_vpc.default.id}",
					"associated_instance_type": "EcsInstance",
					"associated_instances": []string{
						"${alibabacloudstack_ecs_instance.default.0.id}",
						"${alibabacloudstack_ecs_instance.default.1.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name":              name,
						"description":              name,
						"ip_address":               "172.16.1.88",
						"associated_instance_type": "EcsInstance",
						"associated_instances.#":   "2",
						"vswitch_id":               CHECKSET,
						"vpc_id":                   CHECKSET,
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
					"ha_vip_name": "${var.name}_change",
					"description": "${var.name}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ha_vip_name": fmt.Sprintf("%s_change", name),
						"description": fmt.Sprintf("%s_change", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"associated_instance_type": "NetworkInterface",
					"associated_instances": []string{
						"${alibabacloudstack_ecs_networkinterface.default.0.id}",
						"${alibabacloudstack_ecs_networkinterface.default.1.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"associated_instance_type": "NetworkInterface",
						"associated_instances.#":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"associated_instance_type": "NetworkInterface",
					"associated_instances": []string{
						"${alibabacloudstack_ecs_networkinterface.default.0.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"associated_instance_type": "NetworkInterface",
						"associated_instances.#":   "1",
						"associated_instances.1":   REMOVEKEY,
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "Test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "Test",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF-update",
			// 			"For":     "Test-update",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF-update",
			// 			"tags.For":     "Test-update",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": REMOVEKEY,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "0",
			// 			"tags.Created": REMOVEKEY,
			// 			"tags.For":     REMOVEKEY,
			// 		}),
			// 	),
			// },
		},
	})
}

func resourceVpcHavipBasicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s

%s

resource "alibabacloudstack_ecs_instance" "default" {
  count                = 2
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs_${count.index}"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      instance_type,
	  system_disk_category
    ]
  }
}


resource "alibabacloudstack_ecs_networkinterface" "default" {
  	count                	= 2
	network_interface_name 	= "${var.name}_eni_${count.index}"
    vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
	security_groups      	= [alibabacloudstack_ecs_securitygroup.default.id]
}
`, name, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes, SecurityGroupCommonTestCase)
}
