package alibabacloudstack

//import (
//	"fmt"
//	"testing"
//
//	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
//	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
//	
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//)
//
//func TestAccAlibabacloudStackCrEERepo_Basic(t *testing.T) {
//	setTestaccCrEEInstanceId(t)
//	var v *cr_ee.GetRepositoryResponse
//	resourceId := "alibabacloudstack_cr_ee_repo.default"
//	ra := resourceAttrInit(resourceId, nil)
//	serviceFunc := func() interface{} {
//		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//	}
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEERepo")
//	rac := resourceAttrCheckInit(rc, ra)
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := getAccTestRandInt(1000000, 9999999)
//	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEERepoConfigDependence)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheckWithCrEE(t)
//		},
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id": testaccCrEEInstanceId,
//					"namespace":   "${alibabacloudstack_cr_ee_namespace.default.name}",
//					"name":        "${var.name}",
//					"summary":     "summary",
//					"repo_type":   "PUBLIC",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_id": testaccCrEEInstanceId,
//						"namespace":   name,
//						"name":        name,
//						"summary":     "summary",
//						"repo_type":   "PUBLIC",
//						"repo_id":     CHECKSET,
//					}),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"detail": "detail",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"detail": "detail",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"summary": "summary update",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"summary": "summary update",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"repo_type": "PRIVATE",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"repo_type": "PRIVATE",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"detail": "detail update",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"detail": "detail update",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"summary":   "summary",
//					"repo_type": "PUBLIC",
//					"detail":    "detail",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"summary":   "summary",
//						"repo_type": "PUBLIC",
//						"detail":    "detail",
//					}),
//				),
//			},
//		},
//	})
//}
//
//func TestAccAlibabacloudStackCrEERepo_Multi(t *testing.T) {
//	setTestaccCrEEInstanceId(t)
//	var v *cr_ee.GetRepositoryResponse
//	resourceId := "alibabacloudstack_cr_ee_repo.default.4"
//	ra := resourceAttrInit(resourceId, nil)
//	serviceFunc := func() interface{} {
//		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//	}
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEERepo")
//	rac := resourceAttrCheckInit(rc, ra)
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	rand := getAccTestRandInt(1000000, 9999999)
//	name := fmt.Sprintf("tf-testacc-cr-repo-%d", rand)
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEERepoConfigDependence)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheckWithCrEE(t)
//		},
//		IDRefreshName: resourceId,
//		Providers:     testAccProviders,
//		CheckDestroy:  rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id": testaccCrEEInstanceId,
//					"namespace":   "${alibabacloudstack_cr_ee_namespace.default.name}",
//					"name":        "${var.name}${count.index}",
//					"summary":     "summary",
//					"repo_type":   "PUBLIC",
//					"detail":      "detail",
//					"count":       "5",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_id": testaccCrEEInstanceId,
//						"namespace":   name,
//						"name":        name + fmt.Sprint(4),
//						"summary":     "summary",
//						"repo_type":   "PUBLIC",
//						"detail":      "detail",
//					}),
//				),
//			},
//		},
//	})
//}
//
//func resourceCrEERepoConfigDependence(name string) string {
//	fn := func() string {
//		return fmt.Sprintf(`
//variable "name" {
//	default = "%s"
//}
//
//resource "alibabacloudstack_cr_ee_namespace" "default" {
//	instance_id = "%s"
//	name = "${var.name}"
//	auto_create	= false
//	default_visibility = "PRIVATE"
//}
//`, name, testaccCrEEInstanceId)
//	}
//
//	return fn()
//}
