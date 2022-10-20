package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackExpressConnectAccessPointsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)

	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_express_connect_access_points.default"
	name := fmt.Sprintf("tf-testacc-expressConnectAccessPoints%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectAccessPointsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"ap-cn-qingdao-env17-d01-amtest17"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"ap-cn-qingdao-env17-d01-amtest17"},
			"status": "recommended",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"fake"},
			"status": "full",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"ap-cn-qingdao-env17-d01-amtest17"},
			"status": "recommended",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"fake"},
			"status": "full",
		}),
	}

	var existExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"points.#":                    "1",
			"points.0.id":                 CHECKSET,
			"points.0.access_point_id":    "ap-cn-qingdao-env17-d01-amtest17",
			"points.0.access_point_name":  "",
			"points.0.attached_region_no": "",
			"points.0.description":        "",
			"points.0.host_operator":      CHECKSET,
			"points.0.location":           "",
			"points.0.status":             "recommended",
			"points.0.type":               "VPC",
		}
	}

	var fakeExpressConnectAccessPointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"points.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var ExpressConnectAccessPointsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectAccessPointsMapFunc,
		fakeMapFunc:  fakeExpressConnectAccessPointsMapFunc,
	}

	ExpressConnectAccessPointsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}

func dataSourceExpressConnectAccessPointsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		`, name)
}
