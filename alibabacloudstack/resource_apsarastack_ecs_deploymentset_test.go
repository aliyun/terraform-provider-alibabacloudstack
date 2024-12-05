package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsDeploymentset0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_deploymentset.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDeploymentsetCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribedeploymentsetsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdeployment_set%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDeploymentsetBasicdependence)
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

					"description": "RDKTest",

					"strategy": "Availability",

					"deployment_set_name": "RDKTest",

					"granularity": "Host",

					"domain": "Default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDKTest",

						"strategy": "Availability",

						"deployment_set_name": "RDKTest",

						"granularity": "Host",

						"domain": "Default",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "RDK-Test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDK-Test",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEcsDeploymentsetCheckmap = map[string]string{

	"description": CHECKSET,

	"group_count": CHECKSET,

	"create_time": CHECKSET,

	"granularity": CHECKSET,

	"deployment_set_id": CHECKSET,

	"instance_amount": CHECKSET,

	"page_total": CHECKSET,

	"strategy": CHECKSET,

	"deployment_set_name": CHECKSET,

	"region_id": CHECKSET,

	"domain": CHECKSET,

	"instance_ids": CHECKSET,
}

func AlibabacloudTestAccEcsDeploymentsetBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
