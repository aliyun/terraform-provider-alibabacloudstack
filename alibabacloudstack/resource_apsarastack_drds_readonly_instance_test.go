package alibabacloudstack

import (
	"fmt"
	"testing"


	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"


	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_drds_readonly_instance", &resource.Sweeper{
		Name: "alibabacloudstack_drds_readonly_instance",
		F:    testSweepDrdsInstances,
	})
}

func TestAccAlibabacloudStackDrdsReadonlyInstance_Vpc(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alibabacloudstack_drds_readonly_instance.default"
	ra := resourceAttrInit(resourceId, drdsReadonlyInstancebasicMap)

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDrdsInstance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testacc-readonlyinstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsReadonlyInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "${var.name}",
					"master_instance_id": "${alibabacloudstack_drds_instance.default.id}",
					"zone_id":              "${alibabacloudstack_vpc_vswitch.default.availability_zone}",
					"instance_charge_type": "PostPaid",
					"vswitch_id":           "${alibabacloudstack_vpc_vswitch.default.id}",
					"specification":        "${data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
					"description": "${var.name}_u",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_u",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackDrdsReadonlyInstance_Classic(t *testing.T) {
	var v *drds.DescribeDrdsInstanceResponse

	resourceId := "alibabacloudstack_drds_readonly_instance.default"
	ra := resourceAttrInit(resourceId, drdsInstancebasicMap)

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDrdsInstance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testacc-readonly-instance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsReadonlyInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "${var.name}_readonly",
					"master_instance_id": "${alibabacloudstack_drds_instance.default.id}",
					"zone_id":       "${alibabacloudstack_vpc_vswitch.default.availability_zone}",
					"specification": "${data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name+"_readonly",
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
					"description":   "${var.name}_update",
					"specification": "${data.alibabacloudstack_drds_instance_specifications.default.specifications.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
		},
	})
}

func resourceDrdsReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}
	
	data "alibabacloudstack_drds_instance_specifications" "default" {
		sorted_by = "CPU"
	}

	resource "alibabacloudstack_drds_instance" "default" {
		description = var.name
		zone_id = alibabacloudstack_vpc_vswitch.default.availability_zone
		vswitch_id = alibabacloudstack_vpc_vswitch.default.id
		specification = data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id
	}	

	%s	
`, name, VSwitchCommonTestCase)
}

var drdsReadonlyInstancebasicMap = map[string]string{
	"master_instance_id":   CHECKSET,
	"description":          CHECKSET,
	"zone_id":              CHECKSET,
	"instance_charge_type": "PostPaid",
	"specification":        CHECKSET,
}
