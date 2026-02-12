package alibabacloudstack

import (
	"fmt"
	"os"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSBService_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_csb_service.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackCSBServicetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCsbServiceDetail")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-csbservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCSBServiceBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_TEST_EXISTED_CSB_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"csb_id":          "${var.csb_id}",
					"project_id":      "${alibabacloudstack_csb_project.default.0.project_id}",
					"service_version": "1.0.0",
					"alias":           "${var.name}",
					"service_name":    "${var.name}",
					"model_version":   "2.0",
					"consume_types":   []string{"Restful"},
					"provide_type":    "Restful",
					"route_conf_json": TfRawString(` jsonencode({
					    serviceRouteStrategy = "DIRECT"
					    importConf = {
					      provideType         = "Restful"
					      accessEndpointJSON  = jsonencode({
							maxCountScan   = null
					        endpoint         = "http://127.0.0.1:8086/monitor/status.1688"
					        method           = "POST"
					        requestFormat    = "JSON"
					        requestTimeout   = 30000
					        responseFormat   = "passThrough"
					        traceEnabled     = false
					        maxCountScan     = null
					      })
					      inputParameterMap  = []
					      outputParameterMap = []
					    }
					    fallbackConf = {
					      fallbackValue = ""
					    }
					    fusingConf = {
					      requestWindowInMilliseconds = 0
					      requestVolumeThreshold      = 0
					      errorThresholdPercent       = 0
					    }
					  })`),
					"scope": 0,
					"qps":   0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_version": "1.0.0",
						"alias":           name,
						"service_name":    name,
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
					"project_id":    "${alibabacloudstack_csb_project.default.1.project_id}",
					"alias":         "${var.name}_alias",
					"service_name":  "${var.name}_update",
					"skip_auth":     true,
					"all_visiable":  false,
					"scope":         "0",
					"consume_types": []string{"Restful"},
					"provide_type":  "Restful",
					"route_conf_json": TfRawString(` jsonencode({
					    serviceRouteStrategy = "DIRECT"
					    importConf = {
					      provideType = "Restful"
					      accessEndpointJSON = jsonencode({
							maxCountScan   = null
					        endpoint       = "http://127.0.0.1:8087/monitor/status.1688"
					        method         = "POST"
					        requestFormat  = "JSON"
					        responseFormat = "passThrough"
					        traceEnabled   = false
					        requestTimeout = 20000
					      })
					      inputParameterMap  = []
					      outputParameterMap = []
					    }
					    fallbackConf = { fallbackValue = "" }
					    fusingConf   = { 
					      requestWindowInMilliseconds = 0
					      errorThresholdPercent       = 0
					      requestVolumeThreshold      = 0
					    }
					  })`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias":        name + "_alias",
						"service_name": name + "_update",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackCSBServicetMap0 = map[string]string{
	"csb_id": CHECKSET,
}

func AlibabacloudStackCSBServiceBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable name {
		default = "%s"
	}
	
	variable csb_id {
		default = "%s"
	}
	
	resource alibabacloudstack_csb_project default {
		count = 2
	csb_id=           "${var.csb_id}"
	project_name=     "${var.name}_${count.index}"
	owner_name=       "${var.name}_user"
	owner_email=      "${var.name}@aliyun.test"
	owner_phone_num=  "13900000000"
	description=      "${var.name} desc"
	}
	
	`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_CSB_ID"))
}
