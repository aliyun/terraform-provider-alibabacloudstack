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
	resource.AddTestSweepers("alibabacloudstack_common_bandwidth_package_attachment", &resource.Sweeper{
		Name: "alibabacloudstack_common_bandwidth_package_attachment",
		F:    testSweepCommonBandwidthPackageAttachment,
	})
}

func testSweepCommonBandwidthPackageAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting alibabacloudstack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var commonBandwidthPackages []vpc.CommonBandwidthPackage
	req := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.RegionId = client.RegionId
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams["Department"] = client.Department
	req.QueryParams["ResourceGroup"] = client.ResourceGroup
	req.QueryParams = map[string]string{ "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CommonBandwidthPackage: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		if resp == nil || len(resp.CommonBandwidthPackages.CommonBandwidthPackage) < 1 {
			break
		}
		commonBandwidthPackages = append(commonBandwidthPackages, resp.CommonBandwidthPackages.CommonBandwidthPackage...)

		if len(resp.CommonBandwidthPackages.CommonBandwidthPackage) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, cbwp := range commonBandwidthPackages {
		name := cbwp.Name
		id := cbwp.BandwidthPackageId
		for _, eip := range cbwp.PublicIpAddresses.PublicIpAddresse {
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Common Bandwidth Package: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Unassociating Common Bandwidth Package: %s (%s)", name, id)
			req := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				req.Scheme = "https"
			} else {
				req.Scheme = "http"
			}
			req.RegionId = client.RegionId
			req.Headers = map[string]string{"RegionId": client.RegionId}
			req.QueryParams = map[string]string{ "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			req.BandwidthPackageId = id
			req.IpInstanceId = eip.AllocationId
			_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.RemoveCommonBandwidthPackageIp(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to unassociate Common Bandwidth Package (%s (%s)): %s", name, id, err)
			}
		}
	}
	return nil
}

func TestAccAlibabacloudStackCommonBandwidthPackageAttachmentBasic(t *testing.T) {
	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alibabacloudstack_common_bandwidth_package_attachment.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth_package_id": CHECKSET,
		"instance_id":          CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfigBasic(rand),
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

func TestAccAlibabacloudStackCommonBandwidthPackageAttachmentMulti(t *testing.T) {
	var v vpc.CommonBandwidthPackage
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alibabacloudstack_common_bandwidth_package_attachment.default.1"
	ra := resourceAttrInit(resourceId, map[string]string{
		"bandwidth_package_id": CHECKSET,
		"instance_id":          CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCommonBandwidthPackageAttachmentConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckCommonBandwidthPackageAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	VpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_common_bandwidth_package_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if len(parts) != 2 {
			return WrapError(Error("invalid resource id"))
		}
		_, err = VpcService.DescribeCommonBandwidthPackageAttachment(rs.Primary.ID)
		if err != nil {
			return WrapErrorf(err, "DescribeCommonBandwidthPackageAttachment")
		}

	}
	return nil
}
func testAccCommonBandwidthPackageAttachmentConfigBasic(rand int) string {
	return fmt.Sprintf(`
    variable "name"{
    	default = "tf-testAccBandwidtchPackage%d"
    }

	resource "alibabacloudstack_common_bandwidth_package" "default" {
		bandwidth = "2"
		name = "${var.name}"
		description = "${var.name}_description"
	}

	resource "alibabacloudstack_eip" "default" {
		name = "${var.name}"
		bandwidth            = "2"
	}

	resource "alibabacloudstack_common_bandwidth_package_attachment" "default" {
		bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.default.id}"
		instance_id = "${alibabacloudstack_eip.default.id}"
	}
	`, rand)
}

func testAccCommonBandwidthPackageAttachmentConfigMulti(rand int) string {
	return fmt.Sprintf(`
    variable "name"{
    	default = "tf-testAccBandwidtchPackage%d"
    }

	variable "number" {
    	default = "2"
    }

	resource "alibabacloudstack_common_bandwidth_package" "default" {
		count = "${var.number}"
		bandwidth = "2"
		name = "${var.name}"
		description = "${var.name}_description"
	}

	resource "alibabacloudstack_eip" "default" {
		count = "${var.number}"
		name = "${var.name}"
		bandwidth            = "2"
	}

	resource "alibabacloudstack_common_bandwidth_package_attachment" "default" {
		count = "${var.number}"
		bandwidth_package_id = "${element(alibabacloudstack_common_bandwidth_package.default.*.id,count.index)}"
		instance_id = "${element(alibabacloudstack_eip.default.*.id,count.index)}"
	}
	`, rand)
}
