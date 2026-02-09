package alibabacloudstack

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackImportImage(t *testing.T) {
	var v ecs.Image

	const size = 10 * 1024 * 1024
	tmpFile, err := ioutil.TempFile("", "tf-oss-object-test-acc-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// first write some data to the tempfile just so it's not 0 bytes.
	_, err = tmpFile.Write(make([]byte, size))
	if err != nil {
		panic(err)
	}

	resourceId := "alibabacloudstack_image_import.default"
	ra := resourceAttrInit(resourceId, testAccImageImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(1000, 9999)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-imageimport%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageImageBasicConfigDependence(tmpFile.Name()))
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
					"description":  "${var.name}_desc",
					"image_name":   "${var.name}",
					"architecture": "x86_64",
					"license_type": "Auto",
					"platform":     "Ubuntu",
					"os_type":      "linux",
					"disk_device_mapping": []map[string]interface{}{
						{
							"disk_image_size": "10",
							"format":          "RAW",
							"oss_bucket":      "${alibabacloudstack_oss_bucket.default.bucket}",
							"oss_object":      "${alibabacloudstack_oss_bucket_object.default.key}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name + "_desc",
						"image_name":            name,
						"architecture":          "x86_64",
						"license_type":          "Auto",
						"platform":              "Ubuntu",
						"os_type":               "linux",
						"disk_device_mapping.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_type"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_desc update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_desc update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
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

func resourceImageImageBasicConfigDependence(filename string) func(string) string {
	return func(name string) string {
		return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alibabacloudstack_oss_bucket" "default" {
	bucket = "${var.name}"
	acl = "public-read-write"
	}

	resource "alibabacloudstack_oss_bucket_object" "default" {
	bucket=       "${alibabacloudstack_oss_bucket.default.bucket}"
	key=          "test-object-source-key"
	source=       "%s"
	}
`, name, filename)
	}
}
