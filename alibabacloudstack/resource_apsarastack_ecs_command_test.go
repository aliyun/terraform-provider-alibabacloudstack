package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/helper/sdk_patch/datahub_patch"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsCommand0(t *testing.T) {
	var v *datahub_patch.EcsDescribeEcsCommandResult

	resourceId := "alibabacloudstack_ecs_command.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsCommandCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribecommandsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ecscommand%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsCommandBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"command_content": "systemctl stop kubelet.service; systemctl disable kubelet.service; systemctl daemon-reload; yum -y remove kubeadm kubelet kubectl;",
					"type": "RunShellScript",
					"description": "testDescription",
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content": "systemctl stop kubelet.service; systemctl disable kubelet.service; systemctl daemon-reload; yum -y remove kubeadm kubelet kubectl;",
						"type": "RunShellScript",
						"description": "testDescription",
						"name": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"command_content": "echo",

					"type": "RunShellScript",

					"description": "testDescription-update",

					"name": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"command_content": "echo",

						"type": "RunShellScript",

						"description": "testDescription-update",

						"name": name+"-update",
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

var AlibabacloudTestAccEcsCommandCheckmap = map[string]string{
	"description": CHECKSET,
	"timeout": CHECKSET,
	"command_content": CHECKSET,
	"type": CHECKSET,
	"enable_parameter": CHECKSET,
	"name": CHECKSET,
}

func AlibabacloudTestAccEcsCommandBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
