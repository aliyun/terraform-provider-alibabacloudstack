package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("apsarastack_vpc", &resource.Sweeper{
		Name: "apsarastack_vpc",
		F:    testSweepVpcs,
		Dependencies: []string{
			"apsarastack_vswitch",
			"apsarastack_nat_gateway",
			"apsarastack_security_group",
			"apsarastack_ots_instance",
			"apsarastack_router_interface",
			"apsarastack_route_table",
			"apsarastack_cen_instance",
		},
	})
}

func testSweepVpcs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting apsarastack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var vpcs []vpc.Vpc
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVpcs(request)
			})
			return err
		}); err != nil {
			log.Printf("[ERROR] Error retrieving VPCs: %s", WrapError(err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 {
			break
		}
		vpcs = append(vpcs, response.Vpcs.Vpc...)

		if len(response.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		}
		request.PageNumber = page
	}

	for _, v := range vpcs {
		name := v.VpcName
		id := v.VpcId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPC: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VPC: %s (%s)", name, id)
		service := VpcService{client}
		err := service.sweepVpc(id)
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_vpc" {
			continue
		}

		// Try to find the VPC
		instance, err := vpcService.DescribeVpc(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("VPC %s still exist", instance.VpcId))
	}

	return nil
}

func TestAccApsaraStackVpcBasic(t *testing.T) {
	var v vpc.Vpc
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "apsarastack_vpc.default"
	ra := resourceAttrInit(resourceId, testAccCheckVpcCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsarastackVpcBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs": []string{"106.11.62.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "172.16.0.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "172.16.0.0/16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
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
					"cidr_block":            "172.16.0.0/12",
					"vpc_name":              name + "update",
					"description":           name + "update",
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":              "172.16.0.0/12",
						"vpc_name":                name + "update",
						"description":             name + "update",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
		},
	})

}

func TestAccApsarastackVpc_enableIpv6(t *testing.T) {
	var v vpc.Vpc
	resourceId := "apsarastack_vpc.default"
	ra := resourceAttrInit(resourceId, ApsarastackVpcMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeVpc")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsarastackVpcBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VpcIpv6SupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_cidrs":  []string{"106.11.62.0/24"},
					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_cidrs.#": "1",
						"enable_ipv6":  "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "enable_ipv6"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secondary_cidr_blocks.#": "1",
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
					"vpc_name":              name + "update",
					"description":           name + "update",
					"secondary_cidr_blocks": []string{"10.0.0.0/8"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_name":                name + "update",
						"description":             name + "update",
						"secondary_cidr_blocks.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackVpcMulti(t *testing.T) {
	var v vpc.Vpc
	rand := acctest.RandInt()
	resourceId := "apsarastack_vpc.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckVpcCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVpcConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf_testAccVpcConfigName%d", rand),
					}),
				),
			},
		},
	})

}

func testAccCheckVpcConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
provider "apsarastack" {
	assume_role {}
}
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	count = 10
	cidr_block = "172.16.0.0/12"
}
`, rand)
}

var testAccCheckVpcCheckMap = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": "",
}

func ApsarastackVpcBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
`)
}

var ApsarastackVpcMap1 = map[string]string{
	"status":          CHECKSET,
	"router_id":       CHECKSET,
	"router_table_id": CHECKSET,
	"route_table_id":  CHECKSET,
	"ipv6_cidr_block": CHECKSET,
}

func ApsarastackVpcBasicDependence1(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}

`)
}
