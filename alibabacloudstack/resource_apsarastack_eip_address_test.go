package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_eip", &resource.Sweeper{
		Name: "alibabacloudstack_eip",
		F:    testSweepEips,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alibabacloudstack_instance",
			"alibabacloudstack_slb",
			"alibabacloudstack_nat_gateway",
		},
	})
}

func testSweepEips(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var eips []vpc.EipAddress
	req := vpc.CreateDescribeEipAddressesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.RegionId = client.RegionId
	req.QueryParams["Department"] = client.Department
	req.QueryParams["ResourceGroup"] = client.ResourceGroup
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeEipAddresses(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving EIPs: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeEipAddressesResponse)
		if resp == nil || len(resp.EipAddresses.EipAddress) < 1 {
			break
		}
		eips = append(eips, resp.EipAddresses.EipAddress...)

		if len(resp.EipAddresses.EipAddress) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	for _, v := range eips {
		name := v.Name
		id := v.AllocationId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping EIP: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting EIP: %s (%s)", name, id)
		req := vpc.CreateReleaseEipAddressRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.QueryParams["Department"] = client.Department
		req.QueryParams["ResourceGroup"] = client.ResourceGroup
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.AllocationId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ReleaseEipAddress(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete EIP (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckEIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_eip" {
			continue
		}

		_, err := vpcService.DescribeEip(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackEip_basic(t *testing.T) {
	var v vpc.EipAddress
	resourceId := "alibabacloudstack_eip.default"
	ra := resourceAttrInit(resourceId, testAccCheckEipCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 999999)
	name := fmt.Sprintf("tf-testAcceEipName%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckEipConfigBasic)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":   "5",
					"name":        "${var.name}",
					"description": "${var.name}",
					"ip_address":  "43.149.7.150",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":   "5",
						"name":        name,
						"description": name,
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
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_desc",
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

func testAccCheckEipConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name"{
	default = "%s"
}

`, name)
}

var testAccCheckEipCheckMap = map[string]string{
	"name":        "",
	"description": "",
	"bandwidth":   "5",
	// read method does't return a value for the period attribute, so it is not tested
	"ip_address": CHECKSET,
	"status":     CHECKSET,
}
