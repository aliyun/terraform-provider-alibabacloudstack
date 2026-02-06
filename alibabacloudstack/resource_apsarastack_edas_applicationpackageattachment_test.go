package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackEdasApplicationPackageAttachment_basic(t *testing.T) {
	var v string
	resourceId := "alibabacloudstack_application_deployment.default"
	ra := resourceAttrInit(resourceId, edasAPAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasAPAttachmentDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testEdasCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":          "${alibabacloudstack_edas_application.default.id}",
					"group_id":        "${data.alibabacloudstack_edas_deploy_groups.default.groups[0].group_id}",
					"war_url":         fmt.Sprintf("http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar", os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN")),
					"package_version": "v2.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"war_url"},
			},
		},
	})
}

func testEdasCheckDeploymentDestroy(s *terraform.State) error {
	return nil
}

var edasAPAttachmentBasicMap = map[string]string{
	"app_id":   CHECKSET,
	"group_id": CHECKSET,
	"war_url":  CHECKSET,
}

func resourceEdasAPAttachmentDependence(name string) string {
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

data "alibabacloudstack_edas_deploy_groups" "default" {
	app_id = "${alibabacloudstack_edas_application.default.id}"
}
`, name, EdasEcsClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
