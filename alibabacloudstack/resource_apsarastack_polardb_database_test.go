package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbDatabase_Update(t *testing.T) {
	var database *PolardbDescribedatabasesResponse
	resourceId := "alibabacloudstack_polardb_database.default"
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testaccdatabse%d", rand)

	var dbDatabaseBasicMap = map[string]string{
		"data_base_instance_id": CHECKSET,
		"data_base_name":        name,
		"character_set_name":    "utf8",
		"data_base_description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbDatabaseConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_base_instance_id": "${local.polardb_dbinstance_id}",
					"data_base_name":        name,
					"character_set_name":    "utf8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_base_description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"data_base_description": "from terraform"}),
				),
			},
		},
	})

}

func resourcePolardbDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s
	
	%s
	`, name, VSwitchCommonTestCase, PolarDBCommonTestCase("MySQL", true))
}
