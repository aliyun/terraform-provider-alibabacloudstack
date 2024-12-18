package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAdbDbcluster0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_adb_dbcluster.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAdbDbclusterCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAdbDescribebackuppolicyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbdb_cluster%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAdbDbclusterBasicdependence)
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

					"region_id": "cn-hangzhou",

					"payment_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-hangzhou",

						"payment_type": "Postpaid",
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

var AlibabacloudTestAccAdbDbclusterCheckmap = map[string]string{

	"resource_group_id": CHECKSET,

	"db_cluster_type": CHECKSET,

	"log_backup_retention_period": CHECKSET,

	"db_node_count": CHECKSET,

	"executor_count": CHECKSET,

	"lock_reason": CHECKSET,

	"engine": CHECKSET,

	"tags": CHECKSET,

	"db_node_storage": CHECKSET,

	"status": CHECKSET,

	"vpc_id": CHECKSET,

	"vswitch_id": CHECKSET,

	"compute_resource": CHECKSET,

	"expired": CHECKSET,

	"lock_mode": CHECKSET,

	"pay_type": CHECKSET,

	"db_cluster_name": CHECKSET,

	"vpc_cloud_instance_id": CHECKSET,

	"expire_time": CHECKSET,

	"storage_resource": CHECKSET,

	"db_cluster_id": CHECKSET,

	"db_cluster_network_type": CHECKSET,

	"db_cluster_version": CHECKSET,

	"commodity_code": CHECKSET,

	"payment_type": CHECKSET,

	"maintain_time": CHECKSET,

	"db_cluster_category": CHECKSET,

	"security_ips": CHECKSET,

	"zone_id": CHECKSET,

	"create_time": CHECKSET,

	"mode": CHECKSET,

	"db_node_class": CHECKSET,

	"enable_backup_log": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccAdbDbclusterBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
