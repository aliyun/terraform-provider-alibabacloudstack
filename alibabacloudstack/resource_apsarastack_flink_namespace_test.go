package alibabacloudstack

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_flink_namespace", &resource.Sweeper{
		Name: "alibabacloudstack_flink_namespace",
		F:    testSweepFlinkNamespace,
	})
}

func testSweepFlinkNamespace(region string) error {
	// skip not supported region

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("error getting AlibabacloudStack client: %s", err))
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	response, err := client.DoTeaRequest("GET", "ververica", "2020-05-01", "DescribeNamespaces", "/flink/namespace/list", nil, nil, nil)
	if err != nil {
		return err
	}
	vs, err := jsonpath.Get("$.data", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Process Common Request Failed")
	}

	data := vs.([]interface{})
	ns := make([]string, 0)
	for _, v := range data {
		namespace := v.(map[string]interface{})
		name := namespace["Name"].(string)
		for _, p := range prefixes {
			if strings.HasPrefix(name, strings.ToLower(p)) {
				ns = append(ns, name)
			}
		}
	}
	log.Printf("namespace ray %v", ns)
	for _, n := range ns {
		requestBody := map[string]interface{}{
			"name": n,
		}
		_, err := client.DoTeaRequest("DELETE", "ververica", "2020-05-01", "CreateNamespace", "/flink/namespace/create", nil, nil, requestBody)
		if err != nil {
			return err
		}
		log.Printf("response for delete %v", response)
	}
	return nil
}

func TestAccAlibabacloudStackFlinkNamespace_Basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_flink_namespace.default"
	ra := resourceAttrInit(resourceId, flinkNamespaceMap)
	serviceFunc := func() interface{} {
		return &FlinkService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-flink-ns-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceFlinkNamespaceConfigDependence)

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
					"name":      name,
					"cu":        "1",
					"cpu_type":  "Intel",
					"owner_uid": "${alibabacloudstack_ascm_user.default.user_uid}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"cu":       "1",
						"cpu_type": "Intel",
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
					"cu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu": "2",
					}),
				),
			},
		},
	})
}

func resourceFlinkNamespaceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}

resource "alibabacloudstack_ascm_user" "default" {
  display_name = var.name
  mobile_nation_code = "86"
  login_name = var.name
  login_policy_id = "1"
  role_ids = [
               "8",
               "9"
             ]
  cellphone_number = "13612345678"
  email = "${var.name}@gmail.com"
}
`, name)
}

var flinkNamespaceMap = map[string]string{
	"name":     CHECKSET,
	"cu":       CHECKSET,
	"cpu_type": CHECKSET,
}
