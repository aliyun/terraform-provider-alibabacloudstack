package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasLifecyclePolicy_basic(t *testing.T) {
	var v *NasDescribelifecyclepoliciesResponse
	resourceId := "alibabacloudstack_nas_lifecycle_policy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribelifecyclepoliciesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccnaslifecycle%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasLifecyclePolicyDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOss(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_policy_name": name,
					"file_system_id":        "${alibabacloudstack_nas_file_system.default.id}",
					"path":                  "/",
					"recursive":             "false",
					"lifecycle_rule_name":   "DEFAULT_ATIME_14",
					"oss_bucket":            "${alibabacloudstack_oss_bucket.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_policy_name": name,
						"file_system_id":        CHECKSET,
						"path":                  "/",
						"recursive":             "false",
						"lifecycle_rule_name":   "DEFAULT_ATIME_14",
						"oss_bucket":            CHECKSET,
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
					"recursive":           "true",
					"lifecycle_rule_name": "DEFAULT_ATIME_30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recursive":           "true",
						"lifecycle_rule_name": "DEFAULT_ATIME_30",
					}),
				),
			},
		},
	})
}

func AlibabacloudStackNasLifecyclePolicyDependence(name string) string {
	clusterFilter := GetOssClusterFilter()
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s

%s

resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "${var.name}"
  acl    = "public-read"
  oss_cluster = local.cluster_filter
}

`, name, NasCommonTestCase, clusterFilter)
}
