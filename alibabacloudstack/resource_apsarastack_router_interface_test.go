package alibabacloudstack

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_router_interface", &resource.Sweeper{
		Name: "alibabacloudstack_router_interface",
		F:    testSweepRouterInterfaces,
	})
}

func testSweepRouterInterfaces(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ris []vpc.RouterInterfaceType
	req := vpc.CreateDescribeRouterInterfacesRequest()
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
			return vpcClient.DescribeRouterInterfaces(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Router Interfaces: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) < 1 {
			break
		}
		ris = append(ris, resp.RouterInterfaceSet.RouterInterfaceType...)

		if len(resp.RouterInterfaceSet.RouterInterfaceType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	service := VpcService{client}
	for _, v := range ris {
		name := v.Name
		id := v.RouterInterfaceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a RI name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcInstanceId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Router Interface: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Router Interface: %s (%s)", name, id)
		req := vpc.CreateDeleteRouterInterfaceRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.RouterInterfaceId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Router Interface (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckRouterInterfaceExists(n string, ri *vpc.RouterInterfaceType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		vpcService := VpcService{client}

		response, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		ri = &response
		return nil
	}
}

func testAccCheckRouterInterfaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_router_interface" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		vpcService := VpcService{client}

		ri, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.RouterInterfaceId == rs.Primary.ID {
			return WrapError(Error("Interface %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccAlibabacloudStackRouterInterfaceBasic(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alibabacloudstack_router_interface.default"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccRouterInterfaceConfig_role(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"role":          "InitiatingSide",
						"specification": "Large.1",
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccRouterInterfaceConfig_specification(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"role":          "InitiatingSide",
						"specification": "Large.2",
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccRouterInterfaceConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d_change", rand),
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccRouterInterfaceConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d_description", rand),
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccRouterInterfaceConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name":          fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
						"description":   fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
						"role":          "InitiatingSide",
						"specification": "Large.2",
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})

}

func TestAccAlibabacloudStackRouterInterfaceMulti(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alibabacloudstack_router_interface.default.2"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand),
					}),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccRouterInterfaceConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "AcceptingSide"
	name = "${var.name}"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	count = 3
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "AcceptingSide"
	name = "${var.name}"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_role(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "InitiatingSide"
	specification = "Large.1"
	name = "${var.name}"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_specification(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "InitiatingSide"
	name = "${var.name}"
	specification = "Large.2"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_name(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "InitiatingSide"
	name = "${var.name}_change"
	specification = "Large.2"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_description(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "InitiatingSide"
	name = "${var.name}_change"
	specification = "Large.2"
	description = "${var.name}_description"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func testAccRouterInterfaceConfig_all(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterfaceConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "%s"
}

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "${var.region}"
	router_type = "VRouter"
	router_id = "${alibabacloudstack_vpc.default.router_id}"
	role = "InitiatingSide"
	name = "${var.name}"
	specification = "Large.2"
	description = "${var.name}"
}`, rand, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

var testAccRouterInterfaceCheckMap = map[string]string{
	"opposite_region":        CHECKSET,
	"router_type":            "VRouter",
	"router_id":              CHECKSET,
	"role":                   "AcceptingSide",
	"description":            "",
	"health_check_source_ip": "",
	"health_check_target_ip": "",
}
