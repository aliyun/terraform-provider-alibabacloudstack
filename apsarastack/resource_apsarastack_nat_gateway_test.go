package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("apsarastack_nat_gateway", &resource.Sweeper{
		Name: "apsarastack_nat_gateway",
		F:    testSweepNatGateways,
		Dependencies: []string{
			"apsarastack_cs_cluster",
		},
	})
}

func testSweepNatGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Apsarastack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
	}

	var gws []vpc.NatGateway
	req := vpc.CreateDescribeNatGatewaysRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNatGateways(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Nat Gateways: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeNatGatewaysResponse)
		if resp == nil || len(resp.NatGateways.NatGateway) < 1 {
			break
		}
		gws = append(gws, resp.NatGateways.NatGateway...)

		if len(resp.NatGateways.NatGateway) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	service := VpcService{client}
	for _, v := range gws {
		name := v.Name
		id := v.NatGatewayId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Nat Gateway: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Nat Gateway: %s (%s)", name, id)
		if err := service.sweepNatGateway(id); err != nil {
			log.Printf("[ERROR] Failed to delete Nat Gateway (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_nat_gateway" {
			continue
		}

		if _, err := vpcService.DescribeNatGateway(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Nat gateway %s still exist", rs.Primary.ID)
	}

	return nil
}

func TestAccApsaraStackNatGatewayBasic(t *testing.T) {
	var v vpc.NatGateway
	resourceId := "apsarastack_nat_gateway.default"
	ra := resourceAttrInit(resourceId, testAccCheckNatGatewayBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNatGatewayConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccNatGatewayConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNatGatewayConfig%d", rand),
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_type(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "Small",
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_specification(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "Middle",
					}),
				),
			},
			{
				Config: testAccNatGatewayConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"specification": "Small",
						"name":          fmt.Sprintf("tf-testAccNatGatewayConfig%d_all", rand),
					}),
				),
			},
		},
	})
}

func testAccNatGatewayConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	name = "${var.name}"
}
`, rand)
}

func testAccNatGatewayConfig_type(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	name = "${var.name}"
}
`, rand)
}

func testAccNatGatewayConfig_name(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	name = "${var.name}"
}
`, rand)
}

func testAccNatGatewayConfig_specification(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	name = "${var.name}"
	specification = "Middle"
}
`, rand)
}

func testAccNatGatewayConfig_all(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNatGatewayConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	name = "${var.name}_all"
	specification = "Small"
}
`, rand)
}

var testAccCheckNatGatewayBasicMap = map[string]string{
	"name":                  "tf-testAccNatGatewayConfigSpec",
	"specification":         "Small",
	"bandwidth_package_ids": "",
	"forward_table_ids":     CHECKSET,
	"snat_table_ids":        CHECKSET,
}
