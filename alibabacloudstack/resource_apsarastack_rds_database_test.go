package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRdsDatabase0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_rds_database.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRdsDatabaseCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoRdsDescribedatabasesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsdatabase%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRdsDatabaseBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"character_set_name": "gbk",

					"data_base_instance_id": "rm-bp107i59mi7wvqsf2",

					"data_base_name": "rds_mysql",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"character_set_name": "gbk",

						"data_base_instance_id": "rm-bp107i59mi7wvqsf2",

						"data_base_name": "rds_mysql",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"data_base_name": "rds_mysql",

					"data_base_description": "test-DataBaseDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"data_base_name": "rds_mysql",

						"data_base_description": "test-DataBaseDescription",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRdsDatabaseCheckmap = map[string]string{

	"status": CHECKSET,

	"character_set_name": CHECKSET,

	"data_base_instance_id": CHECKSET,

	"data_base_description": CHECKSET,

	"accounts": CHECKSET,

	"data_base_name": CHECKSET,

	"engine": CHECKSET,
}

func AlibabacloudTestAccRdsDatabaseBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
