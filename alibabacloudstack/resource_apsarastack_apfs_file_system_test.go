package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApfsFileSystem_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_apfs_file_system.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"status":       "Running",
		"metered_size": CHECKSET,
		"quota_size":   CHECKSET,
		"volume_size":  CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApfsFileSystemDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":      "${data.alibabacloudstack_apfs_zones.default.zones.0.zone_id}",
					"cluster_id":   "${data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.cluster_id}",
					"storage_type": "${data.alibabacloudstack_apfs_zones.default.zones.0.clusters.0.storage_type}",
					"volume_size":  256,
					"description":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"volume_size": "256",
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"volume_size": 512,
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"volume_size": "512",
						"description": name + "_update",
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

func ApfsFileSystemDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alibabacloudstack_apfs_zones" "default" {
}

`, name)
}
