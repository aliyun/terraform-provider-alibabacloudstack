package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_vpn_gateway", &resource.Sweeper{
		Name: "alibabacloudstack_vpn_gateway",
		F:    testSweepVPNGateways,
		Dependencies: []string{
			"alibabacloudstack_ssl_vpn_server",
			"alibabacloudstack_ssl_vpn_client_cert",
		},
	})
}

func testSweepVPNGateways(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var gws []vpc.VpnGateway
	req := vpc.CreateDescribeVpnGatewaysRequest()
	req.Headers["x-ascm-product-name"] = "Vpc"
	req.Headers["x-acs-organizationId"] = client.Department
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVpnGateways(req)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving VPN Gateways: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeVpnGatewaysResponse)
		if resp == nil || len(resp.VpnGateways.VpnGateway) < 1 {
			break
		}
		gws = append(gws, resp.VpnGateways.VpnGateway...)

		if len(resp.VpnGateways.VpnGateway) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range gws {
		name := v.Name
		id := v.VpnGatewayId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPN Gateway: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting VPN Gateway: %s (%s)", name, id)
		req := vpc.CreateDeleteVpnGatewayRequest()
		req.VpnGatewayId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnGateway(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPN Gateway (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(10 * time.Second)
	}
	return nil
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_vpn" {
			continue
		}

		instance, err := vpnGatewayService.DescribeVpnGateway(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if instance.VpnGatewayId != "" {
			return WrapError(Error("VPN %s still exist", instance.VpnGatewayId))
		}
	}

	return nil
}

// At present, some properties of this resource do not support modification, including: period, bandwidth, enable_ipsec,
// enable_ssl, ssl_connections etc.
func TestAccAlibabacloudStackVpnGatewayBasic(t *testing.T) {
	var v vpc.DescribeVpnGatewayResponse

	resourceId := "alibabacloudstack_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, testAccVpnGatewayCheckMap)
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVpnConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVpnConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVpnConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccVpnConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccVpnConfig%d_description", rand),
					}),
				),
			},
			{
				Config: testAccVpnConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testAccVpnConfig%d", rand),
						"description":  fmt.Sprintf("tf-testAccVpnConfig%d", rand),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackVpnGatewayMulti(t *testing.T) {
	var v vpc.DescribeVpnGatewayResponse

	resourceId := "alibabacloudstack_vpn_gateway.default.4"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(testAccVpnGatewayCheckMap),
				),
			},
		},
	})

}

func testAccVpnConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default =  "tf-testAccVpnConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, rand)
}

func testAccVpnConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default =  "tf-testAccVpnConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	count = 5
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, rand)
}

func testAccVpnConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default =  "tf-testAccVpnConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	name = "${var.name}_change"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, rand)
}
func testAccVpnConfig_description(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default =  "tf-testAccVpnConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	name = "${var.name}_change"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	description = "${var.name}_description"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, rand)
}

func testAccVpnConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default =  "tf-testAccVpnConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = false
	instance_charge_type = "PostPaid"
	description = "${var.name}"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	tags = {
		Created= "TF",
		For=     "Test",
	}
}
`, rand)
}

var testAccVpnGatewayCheckMap = map[string]string{
	"vpc_id":       CHECKSET,
	"bandwidth":    "10",
	"enable_ssl":   "false",
	"enable_ipsec": "true",
	"description":  "",
	"vswitch_id":   CHECKSET,
}
