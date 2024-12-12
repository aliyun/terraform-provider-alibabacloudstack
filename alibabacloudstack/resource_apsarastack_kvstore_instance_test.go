package alibabacloudstack

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRkvstoreInstanceBaisc0(t *testing.T) {

	var v r_kvstore.DBInstanceAttribute

	resourceId := "alibabacloudstack_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccKvstoreCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	name := fmt.Sprintf("tf-testacc%skvstore12345", defaultRegionToTest)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccKvstoreBasicdependence)
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
					instance_name        = name
					instance_class       = "redis.amber.logic.sharding.1g.2db.0rodb.6proxy.multithread"
					engine_version       = "5.0"
					availability_zone    = var.azone
					cpu_type             = "intel"
					architecture_type    = "cluster"
					node_type            = "MASTER_SLAVE"
					series               = "enterprise"
					vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
					instance_charge_type = "PostPaid"
					instance_type        = "Redis"
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_name": "modify_description",

						"vswitch_name": name,
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

					"vpc_auth_mode": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_auth_mode": "Open",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccKvstoreCheckmap = map[string]string{}

func AlibabacloudTestAccKvstoreBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s
`, name, VSwichCommonTestCase)
}
