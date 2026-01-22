package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasNamespace_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_namespace.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasNamespace)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sNasNamespace%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasNamespaceBasicDependence)
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
					"protocol_type": "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}",
					"storage_type":  "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}",
					"zone_id":       "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}",
					"cluster_id":    "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}",
					"description":   name,
					"encrypt_type":  "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":      CHECKSET,
						"storage_type":       CHECKSET,
						"encrypt_type":       "0",
						"zone_id":            CHECKSET,
						"cluster_id":         CHECKSET,
						"description":        name,
						"create_time":        CHECKSET,
						"file_system_type":   CHECKSET,
						"mount_target_count": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id"},
			},
		},
	})
}

var AlibabacloudStackNasNamespace = map[string]string{}

func AlibabacloudStackNasNamespaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s

data "alibabacloudstack_nas_zones" "default" {
}

`, name, DataZoneCommonTestCase)
}
