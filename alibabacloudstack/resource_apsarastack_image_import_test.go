package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImportImage(t *testing.T) {
	var v ecs.Image
	resourceId := "alibabacloudstack_image_import.default"
	ra := resourceAttrInit(resourceId, testAccImageImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(1000, 9999)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageImageBasicConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOSSForImageImport(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
					"image_name":   name,
					"architecture": "x86_64",
					"license_type": "Auto",
					"platform":     "Ubuntu",
					"os_type":      "linux",
					"disk_device_mapping": []map[string]interface{}{
						{
							"disk_image_size": "10",
							"format":          "RAW",
							"oss_bucket":      os.Getenv("ALIBABACLOUDSTACK_OSS_BUCKET_FOR_IMAGE"),
							"oss_object":      os.Getenv("ALIBABACLOUDSTACK_OSS_OBJECT_FOR_IMAGE"),
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
						"image_name":                       name,
						"architecture":                     "x86_64",
						"license_type":                     "Auto",
						"platform":                         "Ubuntu",
						"os_type":                          "linux",
						"disk_device_mapping.#":            "1",
						"disk_device_mapping.0.oss_bucket": CHECKSET,
						"disk_device_mapping.0.oss_object": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescriptionchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescriptionchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%dchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
					"image_name":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%ddescription", rand),
						"image_name":  fmt.Sprintf("tf-testAccEcsImageImportConfigBasic%d", rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_type"},
			},
		},
	})

}

var testAccImageImageCheckMap = map[string]string{}

func resourceImageImageBasicConfigDependence(name string) string {
	return ""
}
