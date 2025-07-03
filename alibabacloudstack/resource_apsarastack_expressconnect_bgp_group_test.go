package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectBgpgroup_basic0(t *testing.T) {
	var v *ExpressconnectBgpGroup
	router_id := os.Getenv("ALIBABACLOUDSTACK_EXCONNECT_ROUTER_ID")
	if router_id == "" {
		t.Skip("Skipping TestAccAlibabacloudStackExpressconnectBgpgroup_basic0: The Env:ALIBABACLOUDSTACK_EXCONNECT_ROUTER_ID unset!")
		t.Skipped()
	}
	resourceId := "alibabacloudstack_expressconnect_bgp_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectBgpgroupCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribebgpgroupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccexpressconnect-bgp-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectBgpgroupDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_name": "${var.name}",
					"description":    "${var.name}",
					"local_asn":      "65534",
					"peer_asn":       "10",
					"router_id":      router_id,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_name": name,
						"description":    name,
						"local_asn":      "65534",
						"peer_asn":       "10",
						"router_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_name": "${var.name}_test",
					"description":    "${var.name}_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_name": fmt.Sprintf("%s_test", name),
						"description":    fmt.Sprintf("%s_test", name),
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

var AlibabacloudStackExpressconnectBgpgroupCheckMap = map[string]string{}

func AlibabacloudStackExpressconnectBgpgroupDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

`, name)
}
