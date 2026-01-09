package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackOssBucketKms_basic(t *testing.T) {

	//var v http.Header
	resourceId := "alibabacloudstack_oss_bucket_kms.default"
	ra := resourceAttrInit(resourceId, ossBucketKmsBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-kms-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketKmsConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudOssBucketKmsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":              "${alibabacloudstack_oss_bucket.default.bucket}",
					"sse_algorithm":       "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
						"sse_algorithm":       "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":              "${alibabacloudstack_oss_bucket.default.bucket}",
					"sse_algorithm":       "KMS",
					"kms_master_key_id": "${alibabacloudstack_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": name,
						"sse_algorithm":       "KMS",
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

func resourceOssBucketKmsConfigDependence(name string) string {

	return fmt.Sprintf(`
resource "alibabacloudstack_oss_bucket" "default" {
	bucket = "%s"
	lifecycle {
	    ignore_changes = [
	      sse_algorithm,
		  kms_key_id
	    ]
	}
}

%s
`, name, KeyCommonTestCase)
}

var ossBucketKmsBasicMap = map[string]string{
	"bucket":              CHECKSET,
}

func testAccCheckAlicloudOssBucketKmsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ossService := OssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_oss_bucket" || rs.Type != "alibabacloudstack_oss_bucket" {
			continue
		}
		bucket, err := ossService.DescribeOssBucket(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if bucket.BucketInfo.Name != "" {
			return errmsgs.WrapError(errmsgs.Error("bucket still exist"))
		}
	}

	return nil
}
