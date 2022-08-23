package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCheckRepoDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_cr_repo" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		crService := CrService{client}
		log.Printf("repo ID %s", rs.Primary.ID)
		_, err := crService.DescribeCrRepo(rs.Primary.ID)

		if err == nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}
func TestAccApsaraStackCRRepo_Basic(t *testing.T) {
	var v GetRepoResponse
	resourceId := "apsarastack_cr_repo.default"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${apsarastack_cr_namespace.default.name}",
					"name":      "${var.name}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
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

func TestAccApsaraStackCRRepo_Multi(t *testing.T) {
	var v GetRepoResponse
	resourceId := "apsarastack_cr_repo.default.2"
	ra := resourceAttrInit(resourceId, crRepoMap)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCRRepoConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace": "${apsarastack_cr_namespace.default.name}",
					"name":      "${var.name}${count.index}",
					"summary":   "summary",
					"repo_type": "PUBLIC",
					"count":     "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceCRRepoConfigDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
variable "name" {
	default = "%s"
}

resource "apsarastack_cr_namespace" "default" {
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PRIVATE"
}
`, name)
}

var crRepoMap = map[string]string{
	"namespace": CHECKSET,
	"name":      CHECKSET,
	"summary":   "summary",
	"repo_type": "PUBLIC",
}
