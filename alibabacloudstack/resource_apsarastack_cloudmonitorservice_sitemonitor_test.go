package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudmonitorserviceSitemonitor0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cloudmonitorservice_sitemonitor.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCloudmonitorserviceSitemonitorCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCmsDescribesitemonitorattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloud_monitor_servicesite_monitor%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudmonitorserviceSitemonitorBasicdependence)
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

					"option_json": "{'Dnstype': 'A', 'Failurerate': 0.5, 'Pingnum': 10}",

					"interval": "1",

					"address": "www.aliyun.com",

					"task_name": "siteMonitorTest",

					"task_type": "PING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"option_json": "{'Dnstype': 'A', 'Failurerate': 0.5, 'Pingnum': 10}",

						"interval": "1",

						"address": "www.aliyun.com",

						"task_name": "siteMonitorTest",

						"task_type": "PING",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"status": "1",

					"option_json": "{'Dnstype': 'A', 'Failurerate': 1, 'Pingnum': 15}",

					"address": "http://www.aliyun.com",

					"task_name": "RekSiteMonitor",

					"interval": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": "1",

						"option_json": "{'Dnstype': 'A', 'Failurerate': 1, 'Pingnum': 15}",

						"address": "http://www.aliyun.com",

						"task_name": "RekSiteMonitor",

						"interval": "5",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudmonitorserviceSitemonitorCheckmap = map[string]string{

	"status": CHECKSET,

	"option_json": CHECKSET,

	"task_id": CHECKSET,

	"address": CHECKSET,

	"task_name": CHECKSET,

	"create_time": CHECKSET,

	"task_type": CHECKSET,

	"isp_cities": CHECKSET,

	"interval": CHECKSET,
}

func AlibabacloudTestAccCloudmonitorserviceSitemonitorBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
