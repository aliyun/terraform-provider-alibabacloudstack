package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeProject_basic(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)
	userId, userPk := InitPreCreateMaxcomputeUser(t, name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMaxcomputeProjectDependenceWithUser)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           "${var.name}",
					"disk":           "50",
					"account":        userId,
					"account_pk":     userPk,
					"quota_id":       "${alibabacloudstack_maxcompute_cu.default.id}",
					"external_table": "true",
					"vpc_ids":        []string{"${alibabacloudstack_vpc_vpc.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"quota_id":  CHECKSET,
						"disk":      "50",
						"vpc_ids.#": CHECKSET,
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"encryption":        "true",
			// 		"encrypt_algorithm": "AES256",
			// 		"encryption_key":    "6cc5b591-0e73-4168-ae54-a1a6e38bae36",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"encryption":        "true",
			// 			"encrypt_algorithm": "AES256",
			// 			"encryption_key":    "6cc5b591-0e73-4168-ae54-a1a6e38bae36",
			// 		}),
			// 	),
			// },
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cpu_type", "external_table", "account"},
			},
		},
	})
}

func resourceMaxcomputeProjectDependenceWithUser(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}
%s

resource "alibabacloudstack_maxcompute_cu" "default" {
	cu_name =      "${var.name}"
	cu_num =       2
	cluster_name = "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}"
}
`, name, VpcCommonTestCase)
}
