package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRdsAccount0(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_rds_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRdsAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoRdsDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRdsAccountBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"data_base_instance_id": "alibabacloudstack_db_instance.default.id",

					"account_name": "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"data_base_instance_id": "alibabacloudstack_db_instance.default.id",

						"account_name": "test1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"account_description": "test-AccountDescription",

					"account_name": "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"account_description": "test-AccountDescription",

						"account_name": "test1",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRdsAccountCheckmap = map[string]string{

	"account_description": CHECKSET,

	"data_base_instance_id": CHECKSET,

	"account_type": CHECKSET,

	"account_name": CHECKSET,
}

func AlibabacloudTestAccRdsAccountBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s



`, name, DBInstanceCommonTestCase)
}
func TestAccAlibabacloudStackRdsAccount1(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_rds_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRdsAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoRdsDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srdsaccount%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRdsAccountBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"data_base_instance_id": "alibabacloudstack_db_instance.default.id",

					"account_name": "test2",

					"account_description": "ccapi-test-create",

					"account_type": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"data_base_instance_id": "alibabacloudstack_db_instance.default.id",

						"account_name": "test2",

						"account_description": "ccapi-test-create",

						"account_type": "Normal",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"account_description": "ccapi-update",

					"account_name": "test1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"account_description": "ccapi-update",

						"account_name": "test1",
					}),
				),
			},
		},
	})
}

