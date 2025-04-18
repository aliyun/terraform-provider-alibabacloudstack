package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOssBucketQuota_basic(t *testing.T) {
	resourceId := "alibabacloudstack_oss_bucket_quota.default"
	ra := resourceAttrInit(resourceId, ossBucketQuotaBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-quota-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOssBucketQuotaConfig)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alibabacloudstack_oss_bucket.default.bucket}",
					"quota":  "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"quota":  "2048",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"quota": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota": "1024",
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

func testAccOssBucketQuotaConfig(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "%s"
}
`, name)
}

var ossBucketQuotaBasicMap = map[string]string{
	"bucket": CHECKSET,
	"quota":  CHECKSET,
}
