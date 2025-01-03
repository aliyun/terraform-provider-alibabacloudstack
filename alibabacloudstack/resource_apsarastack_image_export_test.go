package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImageExport(t *testing.T) {
	var v ecs.Image
	resourceId := "alibabacloudstack_image_export.default"
	ra := resourceAttrInit(resourceId, testAccExportImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccecsimageexportconfigbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageExportBasicConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${alibabacloudstack_image.default.id}",
					"oss_bucket": "${alibabacloudstack_oss_bucket.default.bucket}",
					"oss_prefix": "ecsExport",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_prefix": "ecsExport",
					}),
				),
			},
		},
	})
}

var testAccExportImageCheckMap = map[string]string{
	"image_id":   CHECKSET,
	"oss_bucket": CHECKSET,
}

func resourceImageExportBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s

	resource "alibabacloudstack_image" "default" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
	image_name        = "${var.name}"
	}

	resource "alibabacloudstack_oss_bucket" "default" {
	bucket = "${var.name}"
	acl = "public-read-write"
	}
`, name, ECSInstanceCommonTestCase)
}
