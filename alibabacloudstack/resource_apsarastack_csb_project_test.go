package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackCSBProject_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_csb_project.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackCSBProjectMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCsbProjectDetail")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackCSBProjectBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data":         "{\\\"projectName\\\":\\\"test3\\\",\\\"projectOwnerName\\\":\\\"test3\\\",\\\"projectOwnerEmail\\\":\\\"\\\",\\\"projectOwnerPhoneNum\\\":\\\"\\\",\\\"description\\\":\\\"\\\"}",
					"data2":        "{\\\"projectName\\\":\\\"test15\\\",\\\"projectOwnerName\\\":\\\"test15\\\",\\\"projectOwnerEmail\\\":\\\"\\\",\\\"projectOwnerPhoneNum\\\":\\\"\\\",\\\"description\\\":\\\"\\\",\\\"gmtModified\\\":1672912101000,\\\"csbId\\\":134,\\\"gmtCreate\\\":1672912101000,\\\"ownerId\\\":\\\"1827872887260637\\\",\\\"apiNum\\\":0,\\\"userId\\\":\\\"1827872887260637\\\",\\\"srcType\\\":0,\\\"deleteFlag\\\":0,\\\"id\\\":259,\\\"status\\\":1}",
					"csb_id":       "134",
					"project_name": "test3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"csb_id":       "134",
						"project_name": "test3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name": "test15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name": "test15",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackCSBProjectMap0 = map[string]string{
	"csb_id":       CHECKSET,
	"project_name": CHECKSET,
}

func AlibabacloudStackCSBProjectBasicDependence0(name string) string {
	return fmt.Sprintf(`
`)
}
