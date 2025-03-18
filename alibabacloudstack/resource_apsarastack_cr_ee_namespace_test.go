package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCrEeNamespace_Basic(t *testing.T) {
	var v map[string]interface {}
	resourceId := "alibabacloudstack_cr_ee_namespace.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEeNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEeNamespaceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEe(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
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
					"auto_create": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_create": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_visibility": "PRIVATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_visibility": "PRIVATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               name,
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackCrEeNamespace_Multi(t *testing.T) {
	var v *cr_ee.GetNamespaceResponse
	resourceId := "alibabacloudstack_cr_ee_namespace.default.4"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeCrEeNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-cr-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCrEeNamespaceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithCrEe(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
					"name":               name + "${count.index}",
					"auto_create":        "false",
					"default_visibility": "PUBLIC",
					"count":              "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name + fmt.Sprint(4),
						"auto_create":        "false",
						"default_visibility": "PUBLIC",
					}),
				),
			},
		},
	})
}

func resourceCrEeNamespaceConfigDependence(name string) string {
	return `
	data "alibabacloudstack_cr_ee_instances" "default" {
	}
	`
}

