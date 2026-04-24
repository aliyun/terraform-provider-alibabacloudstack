package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOssBucketKms_basic(t *testing.T) {

	var v *oss.ApplyServerSideEncryptionByDefault
	resourceId := "alibabacloudstack_oss_bucket_kms.default"
	ra := resourceAttrInit(resourceId, ossBucketKmsBasicMap)
	serviceFunc := func() interface{} {
		return &OssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketKmsConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOss(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":        "${alibabacloudstack_oss_bucket.default.bucket}",
					"sse_algorithm": "AES256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        name,
						"sse_algorithm": "AES256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":            "${alibabacloudstack_oss_bucket.default.bucket}",
					"sse_algorithm":     "KMS",
					"kms_master_key_id": "${alibabacloudstack_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":        name,
						"sse_algorithm": "KMS",
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
	clusterFilter := GetOssClusterFilter()
	return fmt.Sprintf(`
variable name {
	default = "%s"
}

%s

resource "alibabacloudstack_oss_bucket" "default" {
	bucket = var.name
	oss_cluster = local.cluster_filter
	lifecycle {
	    ignore_changes = [
	      sse_algorithm,
		  kms_key_id
	    ]
	}
}

%s
`, name, clusterFilter, KeyCommonTestCase)
}

var ossBucketKmsBasicMap = map[string]string{
	"bucket": CHECKSET,
}
