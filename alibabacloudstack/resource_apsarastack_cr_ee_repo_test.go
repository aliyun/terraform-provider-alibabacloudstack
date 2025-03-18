package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCheckCrEeRepoDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_cr_ee_repo" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		crService := CrService{client}
		log.Printf("repo ID %s", rs.Primary.ID)
		_, err := crService.DescribeCrEeRepo(rs.Primary.ID)

		if err == nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
func TestAccAlibabacloudStackCREERepo_Basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cr_ee_repo.default"
	ra := resourceAttrInit(resourceId, crEERepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ee-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCREERepoConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCrEeRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id" : "${alibabacloudstack_cr_ee_namespace.default.instance_id}",
					"namespace": "${alibabacloudstack_cr_ee_namespace.default.name}",
					"name":      "${var.name}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"depends_on" : []string{"alibabacloudstack_cr_ee_repo.fake",},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace": name,
						"name":      name,
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
					"detail": "detail",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detail": "detail",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"summary": "summary update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"summary": "summary update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repo_type": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repo_type": "PRIVATE",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCREERepo_Multi(t *testing.T) {
	var v GetRepoResponse
	resourceId := "alibabacloudstack_cr_ee_repo.default.1"
	ra := resourceAttrInit(resourceId, crEERepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ee-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCREERepoConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCrEeRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${alibabacloudstack_cr_ee_namespace.default.name}",
					"name":      "${var.name}${count.index}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"count":     "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceCREERepoConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
	instance_id = "cri-private"
	name = "${var.name}"
	auto_create = false
	default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_repo" "fake" {
  # 干扰项测试
  instance_id = "${alibabacloudstack_cr_ee_namespace.default.instance_id}"
  name = "${var.name}_fake"
  summary = "summary"
  repo_type = "PUBLIC"
  namespace = "${alibabacloudstack_cr_ee_namespace.default.name}"
}

`, name)
}

var crEERepoMap = map[string]string{
	"namespace": CHECKSET,
	"name":      CHECKSET,
	"summary":   "summary",
	"repo_type": "PUBLIC",
}
