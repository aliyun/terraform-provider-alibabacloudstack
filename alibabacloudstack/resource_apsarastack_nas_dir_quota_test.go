package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasDirQuota_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_nas_dir_quota.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasDirQuotaCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasDirQuota")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasdirquota%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasDirQuotaBasicdependence)
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
					"file_system_id": "${alibabacloudstack_nas_file_system.default.id}",
					"path":           "/",
					"quotas": []map[string]interface{}{
						{
							"quota_type":       "Enforcement",
							"user_type":        "Uid",
							"user_id":          "500",
							"size_limit":       100,
							"file_count_limit": 10000,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id":            CHECKSET,
						"path":                      "/",
						"quotas.#":                  "1",
						"quotas.0.quota_type":       "Enforcement",
						"quotas.0.user_type":        "Uid",
						"quotas.0.user_id":          "500",
						"quotas.0.size_limit":       "100",
						"quotas.0.file_count_limit": "10000",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"quotas": []map[string]interface{}{
						{
							"quota_type": "Accounting",
							"user_type":  "AllUsers",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quotas.0.quota_type":       "Accounting",
						"quotas.0.user_type":        "AllUsers",
						"quotas.0.user_id":          REMOVEKEY,
						"quotas.0.size_limit":       REMOVEKEY,
						"quotas.0.file_count_limit": REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"quotas": []map[string]interface{}{
						{
							"quota_type": "Accounting",
							"user_type":  "AllUsers",
						},
						{
							"quota_type":       "Enforcement",
							"user_type":        "Uid",
							"user_id":          "500",
							"size_limit":       100,
							"file_count_limit": 10000,
						},
						{
							"quota_type": "Accounting",
							"user_type":  "Uid",
							"user_id":    "550",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quotas.#":            "3",
						"quotas.0.quota_type": REMOVEKEY,
						"quotas.0.user_type":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quotas": []map[string]interface{}{
						{
							"quota_type": "Accounting",
							"user_type":  "AllUsers",
						},
						{
							"quota_type":       "Enforcement",
							"user_type":        "Uid",
							"user_id":          "500",
							"size_limit":       100,
							"file_count_limit": 10000,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quotas.#": "2",
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

var AlibabacloudTestAccNasDirQuotaCheckmap = map[string]string{}

func AlibabacloudTestAccNasDirQuotaBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

`, name, NasCommonTestCase)
}
