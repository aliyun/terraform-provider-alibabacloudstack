package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayServiceDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_service.current"
	name := fmt.Sprintf("tf-testacc-apigateway-service%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceApiGatewayServiceConfigDependence)

	enableOnConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable": "On",
		}),
	}

	var existApiGatewayServiceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"id":     CHECKSET,
			"status": "Opened",
		}
	}

	var fakeApiGatewayServiceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"id":     "",
			"status": "",
		}
	}

	var apiGatewayServiceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existApiGatewayServiceMapFunc,
		fakeMapFunc:  fakeApiGatewayServiceMapFunc,
		PreCheck: func(){
			testAccPreCheckWithAPIIsNotSupport(t)
		},
	}
	apiGatewayServiceCheckInfo.dataSourceTestCheck(t, rand, enableOnConf)
}

func dataSourceApiGatewayServiceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, DataZoneCommonTestCase)
}
