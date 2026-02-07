package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEdasInstanceApplicationScale_basic(t *testing.T) {
	var v *edas.AppInfo
	resourceId := "alibabacloudstack_edas_application_scale.default"

	ra := resourceAttrInit(resourceId, edasIAAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeApplicationStatus")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasIAAttachmentDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testEdasCheckIAAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alibabacloudstack_edas_application.default.id}",
					"deploy_group": "${data.alibabacloudstack_edas_deploy_groups.default.groups.0.group_id}",
					"ecu_info":     []string{"${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_ecs_instance.default1.id]}"},
					"force_status": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_status"},
			},
		},
	})
}

var edasIAAttachmentMap = map[string]string{
	"app_id":       CHECKSET,
	"deploy_group": CHECKSET,
}

func testEdasCheckIAAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasIAAttachmentDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

resource "alibabacloudstack_ecs_instance" "default1" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs1"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type,
	  system_disk_category
    ]
  }
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  network_mode = "2"
  cluster_type = "2"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
	cluster_id 	 = "${alibabacloudstack_edas_cluster.default.id}"
	instance_ids = ["${alibabacloudstack_ecs_instance.default.id}", "${alibabacloudstack_ecs_instance.default1.id}"]
}

resource "alibabacloudstack_edas_application" "default" {
	application_name = "${var.name}"
	package_type = "JAR"
	cluster_id = "${alibabacloudstack_edas_instance_cluster_attachment.default.cluster_id}"
	group_id = "all"
	component_id = "8"
	descriotion = "Test Description"
	ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_ecs_instance.default.id]}"]
	package_version = "v1.0.0"
	war_url = "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
}

data "alibabacloudstack_edas_deploy_groups" "default" {
	app_id = "${alibabacloudstack_edas_application.default.id}"
}

`, name, ECSInstanceCommonTestCase, os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
