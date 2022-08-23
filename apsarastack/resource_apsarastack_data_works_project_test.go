package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDataWorksProject_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_data_works_project.default"
	ra := resourceAttrInit(resourceId, ApsaraStackDataWorksProjectMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDataWorksProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDataWorksProjectBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name":   name,
					"task_auth_type": "PROJECT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":   name,
						"task_auth_type": "PROJECT",
					}),
				),
			},
		},
	})
}

var ApsaraStackDataWorksProjectMap0 = map[string]string{}

func ApsaraStackDataWorksProjectBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
