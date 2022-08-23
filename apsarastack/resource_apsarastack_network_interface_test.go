package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("apsarastack_network_interface", &resource.Sweeper{
		Name: "apsarastack_network_interface",
		F:    testApsaraStackNetworkInterface,
	})
}

func testApsaraStackNetworkInterface(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %#v", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	req := ecs.CreateDescribeNetworkInterfacesRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	var enis []ecs.NetworkInterfaceSet
	for {
		raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
			return client.DescribeNetworkInterfaces(req)
		})
		if err != nil {
			return fmt.Errorf("Describe NetworkInterfaces failed, %#v", err)
		}

		resp := raw.(*ecs.DescribeNetworkInterfacesResponse)
		if resp == nil || len(resp.NetworkInterfaceSets.NetworkInterfaceSet) == 0 {
			break
		}

		enis = append(enis, resp.NetworkInterfaceSets.NetworkInterfaceSet...)

		if len(resp.NetworkInterfaceSets.NetworkInterfaceSet) < PageSizeLarge {
			break
		}

		pageNumber, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = pageNumber
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	sweeped := false
	service := VpcService{client}
	for _, eni := range enis {
		name := eni.NetworkInterfaceName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			if need, err := service.needSweepVpc(eni.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping NetworkInterface %s (id: %s; instanceId: %s).", name, eni.NetworkInterfaceId, eni.InstanceId)
			continue
		}
		sweeped = true
		if eni.InstanceId != "" {
			req := ecs.CreateDetachNetworkInterfaceRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				req.Scheme = "https"
			} else {
				req.Scheme = "http"
			}
			req.Headers = map[string]string{"RegionId": client.RegionId}
			req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			req.InstanceId = eni.InstanceId
			req.NetworkInterfaceId = eni.NetworkInterfaceId
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DetachNetworkInterface(req)
			})

			if err != nil {
				log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
				continue
			}

			if err := ecsService.WaitForNetworkInterface(eni.NetworkInterfaceId, Available, DefaultTimeout); err != nil {
				log.Printf("[ERROR] Detach NetworkInterface failed, %#v", err)
				continue
			}
		}

		log.Printf("[INFO] Deleting NetworkInterface %s", name)
		req := ecs.CreateDeleteNetworkInterfaceRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.NetworkInterfaceId = eni.NetworkInterfaceId
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteNetworkInterface(req)
		})

		if err != nil {
			log.Printf("[ERROR] Delete NetworkInterface failed, %#v", err)
			continue
		}
	}

	if sweeped {
		time.Sleep(30 * time.Second)
	}

	return nil
}

func testAccCheckNetworkInterfaceDestroy(t *terraform.State) error {
	for _, rs := range t.RootModule().Resources {
		if rs.Type != "apsarastack_network_interface" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterface(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func TestAccApsaraStackNetworkInterfaceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	var v ecs.NetworkInterfaceSet
	resourceId := "apsarastack_network_interface.default"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNetworkInterface%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkInterfaceConfig_privateIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccNetworkInterfaceConfig_private_ips(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ips.#":     "3",
						"private_ips_count": "3",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackNetworkInterfaceMulti(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "apsarastack_network_interface.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckNetworkInterfaceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkInterfaceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccNetworkInterfaceConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccNetworkInterface%d", rand),
					}),
				),
			},
		},
	})
}

func testAccNetworkInterfaceConfigBasic(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "apsarastack_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
	
}
`, rand)
}

func testAccNetworkInterfaceConfig_privateIp(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccNetworkInterface"
}
resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "apsarastack_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
    private_ip = "192.168.0.2"
	
}
`, rand)
}

func testAccNetworkInterfaceConfig_private_ips(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "apsarastack_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
	name = "${var.name}%d"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
	private_ip = "192.168.0.2"
	private_ips = ["192.168.0.3", "192.168.0.5", "192.168.0.6"]	
}
`, rand)
}

func testAccNetworkInterfaceConfig_multi(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNetworkInterface"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
}

data "apsarastack_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "192.168.0.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_network_interface" "default" {
	name = "${var.name}%d"
    count = 3
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
	
}
`, rand)
}

var testAccCheckNetworkInterfaceCheckMap = map[string]string{
	"vswitch_id":        CHECKSET,
	"security_groups.#": "1",
	"private_ip":        CHECKSET,
	"private_ips.#":     "0",
	"private_ips_count": "0",
	"description":       "",
	"tags.%":            NOSET,
}
