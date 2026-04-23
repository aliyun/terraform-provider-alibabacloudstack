package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOssSingleTunnel_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_oss_single_tunnel.default"
	ra := resourceAttrInit(resourceId, ossSingleTunnelMap)

	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-oss-single-tunnel-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOssSingleTunnelDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOss(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster":    "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}",
					"shared":     "0",
					"label":      "${var.name}",
					"vpc_id":     "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster":    CHECKSET,
						"shared":     "0",
						"label":      name,
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared", "vswitch_id"},
			},
		},
	})
}

var ossSingleTunnelMap = map[string]string{
	"shared": CHECKSET,
	"vip":    CHECKSET,
}

func testAccOssSingleTunnelDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alibabacloudstack_oss_clusters" "default" {
}

%s

`, name, VSwitchCommonTestCase)
}
