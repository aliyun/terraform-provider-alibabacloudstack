package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasNamespaceFilesystemAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_namespace_filesystem_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasNamespaceFilesystemAttachment)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasNamespaceFilesystemAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sNasnpsFsAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasNamespaceFilesystemAttachmentBasicDependence)
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
					"mapped_path":      "${var.name}",
					"file_system_id":   "${alibabacloudstack_nas_file_system.default.id}",
					"nas_namespace_id": "${alibabacloudstack_nas_namespace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mapped_path":      name,
						"file_system_id":   CHECKSET,
						"nas_namespace_id": CHECKSET,
						"create_time":      CHECKSET,
						"storage_type":     CHECKSET,
						"file_system_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mapped_path": "${var.name}_updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mapped_path": name + "_updated",
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

var AlibabacloudStackNasNamespaceFilesystemAttachment = map[string]string{}

func AlibabacloudStackNasNamespaceFilesystemAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s

resource "alibabacloudstack_nas_namespace" "default" {
	zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
	cluster_id = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id}"
	description = var.name
	storage_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type}"
	protocol_type = "${data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type}"
	encrypt_type = "0"
}

`, name, NasCommonTestCase)
}
