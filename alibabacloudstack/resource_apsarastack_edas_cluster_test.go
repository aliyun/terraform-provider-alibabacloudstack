package alibabacloudstack

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccPreCheckEdasNamespace(t *testing.T) {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping test: failed to initialize client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"RegionTag":  fmt.Sprintf("%s:test_probe", region),
		"RegionName": "test_probe",
	}
	response, err := client.DoTeaRequest("POST", "Edas", "2017-08-01", "InsertOrUpdateRegion", "/pop/v5/user_region_def", nil, request, nil)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "user.invalid") || strings.Contains(errMsg, "invalidaction.notfound") || strings.Contains(errMsg, "notsupportapi") {
			t.Skipf("Skipping test: EDAS namespace API is not available in the environment: %s", err)
		}
		return
	}
	// EDAS POP API returns error in response body, not as Go error
	if code, ok := response["Code"]; ok {
		codeStr := fmt.Sprintf("%v", code)
		if codeStr != "200" {
			msg := ""
			if m, ok := response["Message"]; ok {
				msg = fmt.Sprintf("%v", m)
			}
			if strings.Contains(strings.ToLower(msg), "user.invalid") || strings.Contains(strings.ToLower(msg), "invalidaction") {
				t.Skipf("Skipping test: EDAS namespace API is not available in the environment (Code: %s, Message: %s)", codeStr, msg)
			}
		}
	}
}

func TestAccAlibabacloudStackEdasCluster_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "alibabacloudstack_edas_cluster.default"
	ra := resourceAttrInit(resourceId, EdasClusterBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasGetCluster")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasClusterConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckEdasNamespace(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name":      "${var.name}",
					"cluster_type":      2,
					"network_mode":      2,
					"vpc_id":            "${alibabacloudstack_vpc_vpc.default.id}",
					"logical_region_id": "${alibabacloudstack_edas_namespace.default.namespace_logical_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type": "2",
						"network_mode": "2",
						"vpc_id":       CHECKSET,
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

var EdasClusterBasicMap = map[string]string{}

func resourceEdasClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "logical_id" {
  default = "%s:%s"
}

%s

resource "alibabacloudstack_edas_namespace" "default" {
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

`, name, defaultRegionToTest, name, VpcCommonTestCase)
}
