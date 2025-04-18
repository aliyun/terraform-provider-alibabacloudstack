package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksConnectionBasicDependence0)
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
					"project_id":      "10023",
					"connection_type": "rds",
					"content":         AlibabacloudStackDataWorksRdsContentMap,
					"env_type":        "1",
					"sub_type":        "mysql",
					"name":            name,
					"description":     "description" + name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":      "10023",
						"connection_type": "rds",
						"env_type":        "1",
						"sub_type":        "mysql",
						"name":            name,
						"description":     "description" + name,
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
					"description": "description update" + name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description update" + name,
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDataWorksConnectionMap0 = map[string]string{}
var AlibabacloudStackDataWorksRdsContentMap = map[string]interface{}{
	"password":     "inputYourCodeHere@ascm",
	"instanceName": "rm-qd8ba0rn156zu20iu",
	"rdsOwnerId":   "1640757090422435",
	"username":     "cxt",
	"database":     "cxt_test",
	"tag":          "rds",
}

func AlibabacloudStackDataWorksConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
