package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsEbsStorageSets_basic(t *testing.T) {
	var v *datahub.EcsDescribeEcsEbsStorageSetsResult
	resourceId := "alibabacloudstack_ecs_ebs_storage_set.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEcsEbsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEcsEbsStorageSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf-testAcc_storage_set%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEcsEbsBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"storage_set_name":    name,
					"maxpartition_number": "2",
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_set_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudStackEcsEbsMap = map[string]string{
	//"enable_parameter": "false",
}

//func AlibabacloudStackEcsEbsBasicDependence(name string) string {
//	return ""
//}

func AlibabacloudStackEcsEbsBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alibabacloudstack_zones" "default" {}

`, name)
}
