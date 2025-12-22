package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAqsOssScanconfig_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_aqs_oss_scanconfig.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AqsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAqsOssScanConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAqsOssScanconfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: nil,
		Providers:    testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable":                   true,
					"start_time":               "00:00:00",
					"end_time":                 "03:00:00",
					"scan_day_list":            []int{4, 5},
					"scan_mode":                "1",
					"bucket_name":              "${alibabacloudstack_oss_bucket.default.bucket}",
					"decryption":               "OSS",
					"key_suffix":               "all",
					"key_prefix":               "test",
					"last_modified_start_time": "2025-12-17 14:04:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable":                   "true",
						"start_time":               "00:00:00",
						"end_time":                 "03:00:00",
						"scan_day_list.#":          "2",
						"scan_day_list.0":          "4",
						"scan_day_list.1":          "5",
						"scan_mode":                "1",
						"bucket_name":              name,
						"decryption":               "OSS",
						"key_suffix":               "all",
						"key_prefix":               "test",
						"last_modified_start_time": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable":                   false,
					"start_time":               "02:02:00",
					"end_time":                 "05:00:00",
					"scan_day_list":            []int{1, 3, 5},
					"scan_mode":                "2",
					"decryption":               "No",
					"key_suffix":               ".bz2",
					"key_prefix":               "test222",
					"last_modified_start_time": "2025-12-17 14:30:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable":                   "false",
						"start_time":               "02:02:00",
						"end_time":                 "05:00:00",
						"scan_day_list.#":          "3",
						"scan_day_list.0":          "1",
						"scan_day_list.1":          "3",
						"scan_day_list.2":          "5",
						"scan_mode":                "2",
						"decryption":               "No",
						"key_suffix":               ".bz2",
						"key_prefix":               "test222",
						"last_modified_start_time": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceAqsOssScanconfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_oss_clusters" "default" {
  provider = alibabacloudstack-common
}

resource "alibabacloudstack_oss_bucket" "default" {
  provider = alibabacloudstack-common
  bucket = "${var.name}"
  oss_cluster = "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}"
//   oss_cluster = "OssHybridCluster-A-20251125-0096"
}
	`, name)
}
