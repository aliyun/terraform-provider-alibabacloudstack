package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_network_acl_attachment", &resource.Sweeper{
		Name: "alibabacloudstack_network_acl_attachment",
		F:    testSweepNetworkAclAttachment,
	})
}

func testSweepNetworkAclAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var networkAcls []vpc.NetworkAcl
	request := vpc.CreateDescribeNetworkAclsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNetworkAcls(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		response, _ := raw.(*vpc.DescribeNetworkAclsResponse)
		if len(response.NetworkAcls.NetworkAcl) < 1 {
			break
		}
		networkAcls = append(networkAcls, response.NetworkAcls.NetworkAcl...)

		if len(response.NetworkAcls.NetworkAcl) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	for _, nacl := range networkAcls {
		name := nacl.NetworkAclName
		id := nacl.NetworkAclId
		resources := nacl.Resources.Resource
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Unassociating Network Acl: %s (%s)", name, id)
		request := vpc.CreateUnassociateNetworkAclRequest()
		request.NetworkAclId = id

		unassociateNetworkAclResource := []vpc.UnassociateNetworkAclResource{}
		for i := 0; i < len(resources); i++ {
			vpcSource := vpc.UnassociateNetworkAclResource{
				ResourceId:   resources[i].ResourceId,
				ResourceType: resources[i].ResourceType,
			}
			unassociateNetworkAclResource = append(unassociateNetworkAclResource, vpcSource)
		}
		request.Resource = &unassociateNetworkAclResource

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateNetworkAcl(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to unassociate Network Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

// Skip the test because 'resources' is conflict with 'alibabacloudstack_network_acl'.
func SkipTestAccAlibabacloudStackVpcNetworkAclAttachment_basic(t *testing.T) {
	resourceId := "alibabacloudstack_network_acl_attachment.default"
	ra := resourceAttrInit(resourceId, testAccNaclAttachmentCheckMap)
	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAttachment_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_associate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "2",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_unassociate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
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

func testAccCheckNetworkAclAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errmsgs.WrapError(errmsgs.Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return errmsgs.WrapError(errmsgs.Error("No Network Acl Attachment ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		vpcService := VpcService{client}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		networkAclId := parts[0]

		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		res, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		resources := make([]vpc.Resource, 0)
		for _, v := range res {
			item := v.(map[string]interface{})
			resources = append(resources, vpc.Resource{
				Status:       fmt.Sprint(item["Status"]),
				ResourceId:   fmt.Sprint(item["ResourceId"]),
				ResourceType: fmt.Sprint(item["ResourceType"]),
			})
		}
		err = vpcService.DescribeNetworkAclAttachment(networkAclId, resources)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		return nil
	}
}

func testAccCheckNetworkAclAttachmentDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_network_acl_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		networkAclId := parts[0]

		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		vpcResource := []vpc.Resource{}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		for _, e := range resources {
			item := e.(map[string]interface{})

			vpcResource = append(vpcResource, vpc.Resource{
				ResourceId:   item["ResourceId"].(string),
				ResourceType: item["ResourceType"].(string),
			})
		}
		err = vpcService.WaitForNetworkAclAttachment(networkAclId, vpcResource, Deleted, DefaultTimeout)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func testAccNetworkAclAttachment_create(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default2" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_network_acl_attachment" "default" {
	network_acl_id = "${alibabacloudstack_network_acl.default.id}"
    resources {
          resource_id = "${alibabacloudstack_vswitch.default.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

func testAccNetworkAclAttachment_associate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default2" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_network_acl_attachment" "default" {
	network_acl_id = "${alibabacloudstack_network_acl.default.id}"
    resources {
          resource_id = "${alibabacloudstack_vswitch.default.id}"
          resource_type = "VSwitch"
        }
	resources {
          resource_id = "${alibabacloudstack_vswitch.default2.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

func testAccNetworkAclAttachment_unassociate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default2" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_network_acl_attachment" "default" {
	network_acl_id = "${alibabacloudstack_network_acl.default.id}"
    resources {
          resource_id = "${alibabacloudstack_vswitch.default.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

var testAccNaclAttachmentCheckMap = map[string]string{
	"network_acl_id": CHECKSET,
}
