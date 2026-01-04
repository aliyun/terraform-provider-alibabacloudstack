package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCrEEArtifactLifecycleRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cr_ee_attestor_lifecycle_rule.default"
	ra := resourceAttrInit(resourceId, crRepoMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCrEEArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cree-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEArtifactLifecycleRuleDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scope":               "NAMESPACE",
					"instance_id":         "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
					"retention_tag_count": "30",
					"tag_regexp":          "release-v.*",
					"enable_delete_tag":   "true",
					"namespace_name":      "${alibabacloudstack_cr_ee_namespace.default.name}",
					"recent_pull_keep":    30,
					"recent_push_keep":    20,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":               "NAMESPACE",
						"instance_id":         CHECKSET,
						"retention_tag_count": "30",
						"tag_regexp":          "release-v.*",
						"enable_delete_tag":   "true",
						"namespace_name":      name,
						"recent_pull_keep":    "30",
						"recent_push_keep":    "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope":               "REPO",
					"namespace_name":      "${alibabacloudstack_cr_ee_namespace.default2.name}",
					"repo_name":           "${var.name}2",
					"retention_tag_count": "25",
					"tag_regexp":          "release-2.*",
					"enable_delete_tag":   "false",
					"recent_pull_keep":    "10",
					"recent_push_keep":    "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":               "REPO",
						"repo_name":           name + "2",
						"namespace_name":      name + "2",
						"retention_tag_count": "25",
						"tag_regexp":          "release-2.*",
						"enable_delete_tag":   "false",
						"recent_pull_keep":    "10",
						"recent_push_keep":    "10",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCrEEArtifactLifecycleRule_Repo(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cr_ee_attestor_lifecycle_rule.default"
	ra := resourceAttrInit(resourceId, crRepoMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCrEEArtifactLifecycleRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cree-rule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEEArtifactLifecycleRuleDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scope":               "REPO",
					"instance_id":         "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
					"retention_tag_count": "30",
					"tag_regexp":          "release-v.*",
					"enable_delete_tag":   "true",
					"namespace_name":      "${alibabacloudstack_cr_ee_namespace.default.name}",
					"repo_name":           "${var.name}",
					"recent_pull_keep":    30,
					"recent_push_keep":    20,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":               "REPO",
						"instance_id":         CHECKSET,
						"retention_tag_count": "30",
						"tag_regexp":          "release-v.*",
						"enable_delete_tag":   "true",
						"namespace_name":      name,
						"repo_name":           name,
						"recent_pull_keep":    "30",
						"recent_push_keep":    "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope":               "NAMESPACE",
					"namespace_name":      "${alibabacloudstack_cr_ee_namespace.default2.name}",
					"retention_tag_count": "25",
					"tag_regexp":          "release-2.*",
					"enable_delete_tag":   "false",
					"recent_pull_keep":    "10",
					"recent_push_keep":    "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope":               "NAMESPACE",
						"namespace_name":      name + "2",
						"retention_tag_count": "25",
						"tag_regexp":          "release-2.*",
						"enable_delete_tag":   "false",
						"recent_pull_keep":    "10",
						"recent_push_keep":    "10",
					}),
				),
			},
		},
	})
}

func resourceCrEEArtifactLifecycleRuleDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
	instance_id = "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_namespace" "default2" {
	instance_id = "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"
	name = "${var.name}2"
	auto_create	= true
	default_visibility = "PRIVATE"
}

`, name)
}

var CrEEArtifactLifecycleRuleBasicMap = map[string]string{
	"schedule":                        CHECKSET,
	"auto":                            CHECKSET,
	"rule_id":                         CHECKSET,
	"enable_delete_untagged_manifest": CHECKSET,
	"modified_time":                   CHECKSET,
	"create_time":                     CHECKSET,
}
