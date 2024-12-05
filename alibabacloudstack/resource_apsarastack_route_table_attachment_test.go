package alibabacloudstack

import (
	"fmt"
	"log"
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
	resource.AddTestSweepers("alibabacloudstack_route_table_attachment", &resource.Sweeper{
		Name: "alibabacloudstack_route_table_attachment",
		F:    testSweepRouteTableAttachment,
	})
}

func testSweepRouteTableAttachment(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.RouteTableNoSupportedRegions) {
		log.Printf("[INFO] Skipping Route Table unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var routeTables []vpc.RouterTableListType
	req := vpc.CreateDescribeRouteTableListRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.QueryParams["Department"] = client.Department
	req.QueryParams["ResourceGroup"] = client.ResourceGroup
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{ "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving RouteTables: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeRouteTableListResponse)
		if resp == nil || len(resp.RouterTableList.RouterTableListType) < 1 {
			break
		}
		routeTables = append(routeTables, resp.RouterTableList.RouterTableListType...)

		if len(resp.RouterTableList.RouterTableListType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, vtb := range routeTables {
		name := vtb.RouteTableName
		id := vtb.RouteTableId
		for _, vswitch := range vtb.VSwitchIds.VSwitchId {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Route Table: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Unassociating Route Table: %s (%s)", name, id)
			req := vpc.CreateUnassociateRouteTableRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				req.Scheme = "https"
			} else {
				req.Scheme = "http"
			}
			req.Headers = map[string]string{"RegionId": client.RegionId}
			req.QueryParams = map[string]string{ "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			req.RouteTableId = id
			req.VSwitchId = vswitch
			_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.UnassociateRouteTable(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to unassociate Route Table (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}

func testAccCheckRouteTableAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_route_table_attachment" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if len(parts) != 2 {
			return fmt.Errorf("invalid resource id")
		}
		_, err := vpcService.DescribeRouteTableAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Route Table attachment error %#v", err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackRouteTableAttachmentBasic(t *testing.T) {
	var v vpc.RouterTableListType
	resourceId := "alibabacloudstack_route_table_attachment.default"
	rand := acctest.RandIntRange(1000, 9999)
	ra := resourceAttrInit(resourceId, testAccRouteTableAttachmentBasicCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAttachmentConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccAlibabacloudStackRouteTableAttachmentMulti(t *testing.T) {
	var v vpc.RouterTableListType
	resourceId := "alibabacloudstack_route_table_attachment.default.1"
	rand := acctest.RandIntRange(1000, 9999)
	ra := resourceAttrInit(resourceId, testAccRouteTableAttachmentBasicCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAttachmentConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccRouteTableAttachmentConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTableAttachment%d"
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

resource "alibabacloudstack_route_table" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
    name = "${var.name}"
    description = "${var.name}_description"
}

resource "alibabacloudstack_route_table_attachment" "default" {
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	route_table_id = "${alibabacloudstack_route_table.default.id}"
}
`, rand)
}

func testAccRouteTableAttachmentConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteTableAttachment%d"
}

variable "number" {
	default = "2"
}

resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}
 data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
 count = "${var.number}"
  vpc_id = "${ alibabacloudstack_vpc.default.id }"
  cidr_block = "172.16.${count.index}.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_route_table" "default" {
	count = "${var.number}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
    name = "${var.name}"
    description = "${var.name}_description"
}

resource "alibabacloudstack_route_table_attachment" "default" {
    count = "${var.number}"
	vswitch_id = "${element(alibabacloudstack_vswitch.default.*.id,count.index)}"
	route_table_id = "${element(alibabacloudstack_route_table.default.*.id,count.index)}"
}
`, rand)
}

var testAccRouteTableAttachmentBasicCheckMap = map[string]string{
	"vswitch_id":     CHECKSET,
	"route_table_id": CHECKSET,
}
