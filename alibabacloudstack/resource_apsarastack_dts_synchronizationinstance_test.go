package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDtsSynchronizationinstance0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_dts_synchronizationinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDtsSynchronizationinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDtsDescribedtsjobsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronization_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDtsSynchronizationinstanceBasicdependence)
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

					"source_endpoint_engine_name": "MySQL",

					"destination_endpoint_engine_name": "MySQL",

					"destination_endpoint_region": "cn-hangzhou",

					"source_endpoint_region": "cn-hangzhou",

					"payment_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"source_endpoint_engine_name": "MySQL",

						"destination_endpoint_engine_name": "MySQL",

						"destination_endpoint_region": "cn-hangzhou",

						"source_endpoint_region": "cn-hangzhou",

						"payment_type": "PostPaid",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDtsSynchronizationinstanceCheckmap = map[string]string{

	"create_time": CHECKSET,

	"instance_class": CHECKSET,

	"source_endpoint_engine_name": CHECKSET,

	"destination_endpoint_engine_name": CHECKSET,

	"dts_instance_id": CHECKSET,

	"destination_endpoint_region": CHECKSET,

	"source_endpoint_region": CHECKSET,

	"type": CHECKSET,

	"payment_type": CHECKSET,
}

func AlibabacloudTestAccDtsSynchronizationinstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
