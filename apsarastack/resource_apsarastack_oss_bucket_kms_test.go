package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackOssBucketKms_basic(t *testing.T) {

	//var v http.Header
	resourceId := "apsarastack_oss_bucket_kms.default"
	ra := resourceAttrInit(resourceId, ossBucketKmsBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-kms-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketKmsConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudOssBucketKmsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":              "${apsarastack_oss_bucket.default.bucket}",
					"sse_algorithm":       "KMS",
					"kms_data_encryption": "SM4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":              "${apsarastack_oss_bucket.default.bucket}",
					"sse_algorithm":       "KMS",
					"kms_data_encryption": "SM4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
					}),
				),
			},
		},
	})
}

func resourceOssBucketKmsConfigDependence(name string) string {

	return fmt.Sprintf(`
resource "apsarastack_oss_bucket" "default" {
	bucket = "%s"
}
data "apsarastack_kms_keys" "enabled" {
	status = "%s"
}
`, name, string(EnabledStatus))
}

var ossBucketKmsBasicMap = map[string]string{
	"bucket":              CHECKSET,
	"sse_algorithm":       "KMS",
	"kms_data_encryption": "SM4",
}

func testAccCheckAlicloudOssBucketKmsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ossService := OssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_oss_bucket" || rs.Type != "apsarastack_oss_bucket" {
			continue
		}
		bucket, err := ossService.DescribeOssBucket(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if bucket.BucketInfo.Name != "" {
			return WrapError(Error("bucket still exist"))
		}
	}

	return nil
}
