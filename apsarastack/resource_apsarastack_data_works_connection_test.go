package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDataWorksConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_data_works_connection.default"
	ra := resourceAttrInit(resourceId, ApsaraStackDataWorksConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDataWorksConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDataWorksConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
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
					"content":         ApsaraStackDataWorksRdsContentMap,
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

var ApsaraStackDataWorksConnectionMap0 = map[string]string{}
var ApsaraStackDataWorksRdsContentMap = map[string]interface{}{
	"password":     "inputYourCodeHere@ascm",
	"instanceName": "rm-qd8ba0rn156zu20iu",
	"rdsOwnerId":   "1640757090422435",
	"username":     "cxt",
	"database":     "cxt_test",
	"tag":          "rds",
}

func ApsaraStackDataWorksConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
