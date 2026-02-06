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

func TestAccAlibabacloudStackEdasSlbAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alibabacloudstack_edas_slb_attachment.default"

	ra := resourceAttrInit(resourceId, edasSLBAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSLBAttachmentDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testEdasCheckSLBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":        "${alibabacloudstack_edas_application.default.id}",
					"slb_id":        "${alibabacloudstack_slb.default.id}",
					"slb_ip":        "${alibabacloudstack_slb.default.address}",
					"type":          "${alibabacloudstack_slb.default.address_type}",
					"listener_port": "9999",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vserver_group_id", "listener_port"},
			},
		},
	})
}

var edasSLBAttachmentMap = map[string]string{
	"app_id":        CHECKSET,
	"slb_id":        CHECKSET,
	"slb_ip":        CHECKSET,
	"type":          CHECKSET,
	"listener_port": CHECKSET,
	"slb_status":    CHECKSET,
	"vswitch_id":    CHECKSET,
}

func testEdasCheckSLBAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasSLBAttachmentDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

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

resource "alibabacloudstack_slb" "default" {
	name          = "${var.name}"
	vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
	address_type  = "intranet"
	specification = "slb.s1.small"
}

resource "alibabacloudstack_slb_vservergroup" "default" {
	load_balancer_id = "${alibabacloudstack_slb.default.id}"
	vserver_group_name = "vserver_group_name"
}
`, name, EdasEcsClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
