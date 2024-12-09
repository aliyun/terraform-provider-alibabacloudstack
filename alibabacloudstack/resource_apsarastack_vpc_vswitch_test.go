package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcVswitch0(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.test_default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"cidr_block": "192.168.1.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"cidr_block": "192.168.1.0/24",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vswitch_name": "modify_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vswitch_name": "modify_name",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcVswitchCheckmap = map[string]string{

	"description": "测试111111",

	"cidr_block": "192.168.1.0/24",

	"vpc_id": "${alibabacloudstack_vpc.default.id}",

	"vswitch_name": "${var.name}",

	"zone_id": "cn-wulan-env212-amtest212001-a",

	"ipv6_cidr_block": CHECKSET,

	"tags": "",
}

func AlibabacloudTestAccVpcVswitchBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s



`, name, VSwichCommonTestCase)
}
func TestAccAlibabacloudStackVpcVswitch1(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"cidr_block": "192.168.1.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"cidr_block": "192.168.1.0/24",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vswitch_name": "modify_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vswitch_name": "modify_name",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch2(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch3(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::default.VpcId)}}",

					"cidr_block": "${{cidrsubnet(${{ref(resource, VPC::VPC::2.0.0.5.pre::default.CidrBlock)}} ,4,2)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::default.VpcId)}}",

						"cidr_block": "${{cidrsubnet(${{ref(resource, VPC::VPC::2.0.0.5.pre::default.CidrBlock)}} ,4,2)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch4(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"cidr_block": "172.16.0.0/24",

					"vswitch_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"cidr_block": "172.16.0.0/24",

						"vswitch_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"vswitch_name": "rdktestname",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"vswitch_name": "rdktestname",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch5(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

					"cidr_block": "172.16.0.0/24",

					"vswitch_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

						"cidr_block": "172.16.0.0/24",

						"vswitch_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"vswitch_name": "rdktestname",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"vswitch_name": "rdktestname",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch6(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

					"cidr_block": "172.16.0.0/24",

					"vswitch_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

						"cidr_block": "172.16.0.0/24",

						"vswitch_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"vswitch_name": "rdktestname",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"vswitch_name": "rdktestname",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch7(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"cidr_block": "10.0.10.0/24",

					"vswitch_name": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"cidr_block": "10.0.10.0/24",

						"vswitch_name": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"cidr_block": "10.0.10.0/24",

					"vswitch_name": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"cidr_block": "10.0.10.0/24",

						"vswitch_name": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch8(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"cidr_block": "10.1.1.0/24",

					"vswitch_name": "slb_test_clb_core",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"description": "${{ref(resource, VPC::VPC::2.0.0.5.pre::1SYZIP.Ipv6CidrBlocks[0].Ipv6Isp)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"cidr_block": "10.1.1.0/24",

						"vswitch_name": "slb_test_clb_core",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"description": "${{ref(resource, VPC::VPC::2.0.0.5.pre::1SYZIP.Ipv6CidrBlocks[0].Ipv6Isp)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"vswitch_name": "test_network_acl_vpc_multi_cidr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vswitch_name": "test_network_acl_vpc_multi_cidr",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "vSwitch",

					"vswitch_name": "vSwitch-1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "vSwitch",

						"vswitch_name": "vSwitch-1",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_id": "alibabacloudstack_vpc.default.id",

					"cidr_block": "10.1.1.0/24",

					"vswitch_name": "slb_test_clb_core",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"description": "${{ref(resource, VPC::VPC::2.0.0.5.pre::1SYZIP.Ipv6CidrBlocks[0].Ipv6Isp)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "alibabacloudstack_vpc.default.id",

						"cidr_block": "10.1.1.0/24",

						"vswitch_name": "slb_test_clb_core",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"description": "${{ref(resource, VPC::VPC::2.0.0.5.pre::1SYZIP.Ipv6CidrBlocks[0].Ipv6Isp)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch9(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

					"cidr_block": "172.16.0.0/24",

					"vswitch_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

						"cidr_block": "172.16.0.0/24",

						"vswitch_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"vswitch_name": "rdktestname",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"vswitch_name": "rdktestname",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

					"cidr_block": "172.16.0.0/24",

					"vswitch_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"vpc_id": "${{ref(resource, VPC::VPC::2.0.0.5.pre::OeB4be.VpcId)}}",

						"cidr_block": "172.16.0.0/24",

						"vswitch_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackVpcVswitch10(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"cidr_block": "10.0.10.0/24",

					"vswitch_name": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"cidr_block": "10.0.10.0/24",

						"vswitch_name": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"cidr_block": "10.0.10.0/24",

					"vswitch_name": "alibabacloudstack_vswitch.default.id",

					"vpc_id": "alibabacloudstack_vpc.default.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"cidr_block": "10.0.10.0/24",

						"vswitch_name": "alibabacloudstack_vswitch.default.id",

						"vpc_id": "alibabacloudstack_vpc.default.id",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}
