package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackVpcVpc_basic(t *testing.T) {

	var v map[string]interface{}

	// TODO Describe method，v 的类型

	resourceId := "alibabacloudstack_vpc_vpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpcattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVpcBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "RDK更新",
					"cidr_block":  "192.168.0.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDK更新",
						"cidr_block":  "192.168.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"description": "RDK更新1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDK更新1",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccAlibabacloudStackVpcVpc1(t *testing.T) {

	var v map[string]interface{}

	// TODO Describe method，v 的类型

	resourceId := "alibabacloudstack_vpc_vpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpcattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVpcBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"cidr_block": "172.16.0.0/12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cidr_block": "172.16.0.0/12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"cidr_block": "192.168.0.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cidr_block": "192.168.0.0/16",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackVpcVpc2(t *testing.T) {

	var v map[string]interface{}

	// TODO Describe method，v 的类型

	resourceId := "alibabacloudstack_vpc_vpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpcattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVpcBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"enable_ipv6": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"enable_ipv6": "true",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackVpcVpc3(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_vpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpcattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVpcIpv6Basicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block":            "172.16.0.0/24",
					"secondary_cidr_blocks": []string{"192.168.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":              "172.16.0.0/24",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block":            "172.16.0.0/24",
					"secondary_cidr_blocks": []string{"172.16.1.0/24", "172.16.2.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":              "172.16.0.0/24",
						"secondary_cidr_blocks.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_cidr_blocks": []map[string]interface{}{
						{
							"ipv6_isp":        "${data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[0].service_provider}",
							"ipv6_cidr_block": "${local.subnet0_ipv6_cidr}",
						},
						{
							"ipv6_isp":        "${data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[1].service_provider}",
							"ipv6_cidr_block": "${local.subnet1_ipv6_cidr}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_block":                    CHECKSET,
						"ipv6_cidr_blocks.#":                 "2",
						"ipv6_cidr_blocks.0.ipv6_cidr_block": CHECKSET,
						"ipv6_cidr_blocks.1.ipv6_cidr_block": CHECKSET,
						"ipv6_cidr_blocks.1.ipv6_isp":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_cidr_blocks": []map[string]interface{}{
						{
							"ipv6_isp":        "${data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[0].service_provider}",
							"ipv6_cidr_block": "${local.subnet0_ipv6_cidr}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_block":                    CHECKSET,
						"ipv6_cidr_blocks.#":                 "1",
						"ipv6_cidr_blocks.0.ipv6_cidr_block": CHECKSET,
						"ipv6_cidr_blocks.1.ipv6_cidr_block": REMOVEKEY,
						"ipv6_cidr_blocks.1.ipv6_isp":        REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_cidr_blocks": []map[string]interface{}{
						{
								"ipv6_isp":        "${data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[0].service_provider}",
								"ipv6_cidr_block": "${local.subnet0_ipv6_cidr_update}",
							},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_cidr_block":                    CHECKSET,
						"ipv6_cidr_blocks.#":                 "1",
						"ipv6_cidr_blocks.0.ipv6_cidr_block": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcVpcCheckmap = map[string]string{
	"status": CHECKSET,
}

func AlibabacloudTestAccVpcVpcBasicdependence(name string) string {

	//  TODO  检查依赖变量

	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}


func AlibabacloudTestAccVpcVpcIpv6Basicdependence(name string) string {

	//  TODO  检查依赖变量

	return fmt.Sprintf(`
%s
data "alibabacloudstack_vpc_ipv6_isps" "default" {
}

locals {
  subnet0_ipv6_cidr = cidrsubnet(data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[0].cidr_block, 13, 1644)
  subnet0_ipv6_cidr_update = cidrsubnet(data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[0].cidr_block, 13, 1645)
  subnet1_ipv6_cidr = cidrsubnet(data.alibabacloudstack_vpc_ipv6_isps.default.ipv6_isps[1].cidr_block, 13, 1644)
}
`, AlibabacloudTestAccVpcVpcBasicdependence(name))
}

