package alibabacloudstack

import (
	"fmt"
	"testing"


	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCpfsFileSystem(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cpfs_file_system.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasFileSystem0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCpfsFileSystemDependence1)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":     "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}",
					"zone_id":          "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}",
					"cluster_id":       "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}",
					"description":      "${var.name}",
					"capacity":         20480,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":    CHECKSET,
						"storage_type":     CHECKSET,
						"zone_id":          CHECKSET,
						"cluster_id":       CHECKSET,
						"description":      name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				//When the instance state changes to running, the cluster_id may change during queries.
				ImportStateVerifyIgnore: []string{"cluster_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
					"capacity":    30720,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
						"capacity":    "30720",
					}),
				),
			},
		},
	})
}

func AlibabacloudStackCpfsFileSystemDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_nas_zones" "default" {
	file_system_type = "bmcpfs"
}
`, name)
}
