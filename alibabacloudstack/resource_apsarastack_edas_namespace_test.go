package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"


	"github.com/PaesslerAG/jsonpath"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func init() {
	resource.AddTestSweepers(
		"alibabacloudstack_edas_namespace",
		&resource.Sweeper{
			Name: "alibabacloudstack_edas_namespace",
			F:    testSweepEdasNamespace,
		})
}

func testSweepEdasNamespace(region string) error {
// 	if testSweepPreCheckWithRegions(region, true, connectivity.EdasSupportedRegions) {
// 		log.Printf("[INFO] Skipping Edas Namespace unsupported region: %s", region)
// 		return nil
// 	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tftestacc",
	}
	request := map[string]interface{}{}
	action:= "ListUserDefineRegion"

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
	
		response, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", action, "/pop/v5/user_region_defs", nil, request, nil)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.UserDefineRegionList.UserDefineRegionEntity", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.UserDefineRegionList.UserDefineRegionEntity", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["RegionName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Edas Namespace: %s", item["RegionName"].(string))
			continue
		}
		action := "DeleteUserDefineRegion"
		request := map[string]interface{}{
			"Id": StringPointer(fmt.Sprint(item["Id"])),
		}
		_, err = client.DoTeaRequest("DELETE", "Edas", "2017-08-01", action, "/pop/v5/user_region_def", nil, request, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Edas Namespace (%s): %s", item["RegionName"].(string), err)
		}
		log.Printf("[INFO] Delete Edas Namespace success: %s ", item["RegionName"].(string))
	}

	return nil
}

func TestAccAlibabacloudStackEDASNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_edas_namespace.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEDASNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(100, 999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEDASNamespaceBasicDependence0)
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
// 					"debug_enable":         "false",
					"description":          "${var.name}",
					"namespace_name":       "${var.name}",
					"namespace_logical_id": TfRawString("substr(local.logical_id, 0, min(length(local.logical_id), 32))"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
// 						"debug_enable":         "false",
						"description":          name,
						"namespace_name":       name,
						"namespace_logical_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
// 			{
// 				Config: testAccConfig(map[string]interface{}{
// 					"debug_enable": "true",
// 				}),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheck(map[string]string{
// 						"debug_enable": "true",
// 					}),
// 				),
// 			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudStackEDASNamespaceMap0 = map[string]string{
	"namespace_logical_id": CHECKSET,
	"namespace_name":       CHECKSET,
// 	"debug_enable":         CHECKSET,
}

func AlibabacloudStackEDASNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alibabacloudstack_account" "current" {
}

locals {
  logical_id = "${data.alibabacloudstack_account.current.region}:${var.name}"
}



`, name)
}
