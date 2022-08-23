package apsarastack

import (
	"fmt"

	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDTSSynchronizationInstance_basic0(t *testing.T) {

	var v map[string]interface{}
	resourceId := "apsarastack_dts_synchronization_instance.default"

	ra := resourceAttrInit(resourceId, ApsaraStackDTSSynchronizationInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDtsSynchronizationInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDTSSynchronizationInstanceBasicDependence0)
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

					"payment_type":                     "PostPaid",
					"source_endpoint_engine_name":      "MySQL",
					"source_endpoint_region":           "cn-qingdao-env17-d01",
					"destination_endpoint_engine_name": "MySQL",
					"destination_endpoint_region":      "cn-qingdao-env17-d01",
					"instance_class":                   "small",
					"sync_architecture":                "bidirectional",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":      "PostPaid",
						"instance_class":    "small",
						"sync_architecture": "bidirectional",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"destination_endpoint_region", "source_endpoint_engine_name", "source_endpoint_region", "destination_endpoint_engine_name",
					"database_count", "status", "quantity", "sync_architecture", "auto_start", "compute_unit", "period", "used_time", "auto_pay", "order_type", "synchronization_direction"},
			},
		},
	})
}

var ApsaraStackDTSSynchronizationInstanceMap0 = map[string]string{
	"sync_architecture":         NOSET,
	"auto_start":                NOSET,
	"compute_unit":              NOSET,
	"period":                    NOSET,
	"used_time":                 NOSET,
	"auto_pay":                  NOSET,
	"order_type":                NOSET,
	"synchronization_direction": NOSET,
	"database_count":            NOSET,
	"status":                    NOSET,
	"quantity":                  NOSET,
}

func ApsaraStackDTSSynchronizationInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
