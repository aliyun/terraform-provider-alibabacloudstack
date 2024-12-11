package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsCommand0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_command.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsCommandCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribecommandsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secscommand%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsCommandBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"command_content": "systemctl stop kubelet.service; systemctl disable kubelet.service; systemctl daemon-reload; yum -y remove kubeadm kubelet kubectl;",

					"type": "RunShellScript",

					"description": "testDescription",

					"command_name": "testName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"command_content": "systemctl stop kubelet.service; systemctl disable kubelet.service; systemctl daemon-reload; yum -y remove kubeadm kubelet kubectl;",

						"type": "RunShellScript",

						"description": "testDescription",

						"command_name": "testName",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"command_content": "echo",

					"type": "RunShellScript",

					"description": "testDescription-update",

					"command_name": "testName-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"command_content": "echo",

						"type": "RunShellScript",

						"description": "testDescription-update",

						"command_name": "testName-update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEcsCommandCheckmap = map[string]string{

	"category": CHECKSET,

	"description": CHECKSET,

	"parameter_names": CHECKSET,

	"timeout": CHECKSET,

	"create_time": CHECKSET,

	"provider": CHECKSET,

	"command_content": CHECKSET,

	"working_dir": CHECKSET,

	"type": CHECKSET,

	"invoke_times": CHECKSET,

	"enable_parameter": CHECKSET,

	"latest": CHECKSET,

	"command_id": CHECKSET,

	"command_name": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccEcsCommandBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
