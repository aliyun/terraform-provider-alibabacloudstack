package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackEdasApplicationPackageAttachment_basic(t *testing.T) {
	var v *edas.Applcation
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
					"group_id":        "${alibabacloudstack_edas_deploy_group.default.group_id}",
					"war_url":         fmt.Sprintf("http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar", os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN")),
					"package_version": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

// resource "alibabacloudstack_vpc" "default" {
// 	cidr_block = "172.16.0.0/12"
// 	name       = "${var.name}"
// }

// resource "alibabacloudstack_edas_cluster" "default" {
// 	cluster_name = "${var.name}"
// 	cluster_type = 2
// 	network_mode = 2
// 	vpc_id       = "${alibabacloudstack_vpc.default.id}"
// }

resource "alibabacloudstack_edas_application" "default" {
	application_name = "${var.name}"
	cluster_id = "7ccfd5b3-a164-424f-8de3-906d22262f8c"
	package_type = "JAR"
	component_id = "8"
}

resource "alibabacloudstack_edas_deploy_group" "default" {
	app_id = "${alibabacloudstack_edas_application.default.id}"
	group_name = "${var.name}"
}
`, name)
}
