package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbBackendServers_classic(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(), //testAccCheckSlbBackendServersDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_instance.new.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from the remote
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				//Config: testAccSlbBackendServersClassicUpdateServer,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "70",
						},
						{
							"server_id": "${alibabacloudstack_instance.new.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
		},
	})
}

func buildBackendServersMap(count int) []map[string]interface{} {
	var result []map[string]interface{}

	str := `${alibabacloudstack_instance.instance.%d.id}`
	for i := 0; i < count; i++ {
		tmp := make(map[string]interface{}, 2)
		tmp["server_id"] = fmt.Sprintf(str, i)
		tmp["weight"] = "10"
		result = append(result, tmp)
	}
	return result
}

func resourceBackendServerConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_instance" "new" {
	image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type        = "${local.default_instance_type_id}"
	system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	system_disk_size     = 40
	system_disk_name     = "test_sys_diskv2"
	security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
	instance_name        = "${var.name}_ecsv2"
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	zone_id    = data.alibabacloudstack_zones.default.zones.0.id
	is_outdated          = false
	lifecycle {
	  ignore_changes = [
		instance_type
	  ]
	}
  }

resource "alibabacloudstack_slb" "default" {
  	name = "${var.name}"
  	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

`, name, ECSInstanceCommonTestCase)
}

var slb_vpc = map[string]string{
	"backend_servers.#": "2",
}
