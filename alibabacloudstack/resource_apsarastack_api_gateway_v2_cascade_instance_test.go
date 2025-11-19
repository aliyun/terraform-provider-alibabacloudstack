package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2CascadeInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_cascade_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCascadeInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccApiGwV2%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGatewayV2CascadeInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// "instance_type":       "0",
					"instance_name":       "${var.name}",
					"cascade_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// "instance_type":       "0",
						"instance_name":       name,
						"cascade_instance_id": CHECKSET,
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

func APIGatewayV2CascadeInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

`, name)
}
