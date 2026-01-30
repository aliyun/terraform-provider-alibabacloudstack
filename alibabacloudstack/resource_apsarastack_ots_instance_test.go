package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOtsInstance_clusterName(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_ots_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-otsinst-%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceDependence)
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
					"name":         "${var.name}",
					"alias_name":   "${var.name}_alias",
					"cluster_name": "${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}",
					"description":  "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"alias_name":  name + "_alias",
						"description": name + "_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name":  "${var.name}_alias_update",
					"description": "${var.name}_desc_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"alias_name":  name + "_alias_update",
						"description": name + "_desc_update",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackOtsInstance_clusterType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_ots_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-otsinst-%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceDependence)
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
					"name":          "${var.name}",
					"alias_name":    "${var.name}_alias",
					"specification": "${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_type}",
					"description":   "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"alias_name":  name + "_alias",
						"description": name + "_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func resourceOtsInstanceDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	data "alibabacloudstack_ots_clusters" "anyone" {}
	`, name)
}
