package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"apsarastack_vpc_ipv6_gateway",
		&resource.Sweeper{
			Name: "apsarastack_vpc_ipv6_gateway",
			F:    testSweepVpcIpv6Gateway,
		})
}

func testSweepVpcIpv6Gateway(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.VpcIpv6GatewaySupportRegions) {
		log.Printf("[INFO] Skipping Vpc Ipv6 Gateway unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeIpv6Gateways"
	request := map[string]interface{}{}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["RegionId"] = client.RegionId

	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			request["Product"] = "Vpc"
			request["OrganizationId"] = client.Department
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Ipv6Gateways.Ipv6Gateway", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Ipv6Gateways.Ipv6Gateway", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpc Ipv6 Gateway: %s", item["Name"].(string))
				continue
			}
			action := "DeleteIpv6Gateway"
			deleteRequest := map[string]interface{}{
				"Ipv6GatewayId": item["Ipv6GatewayId"],
			}
			deleteRequest["Product"] = "Vpc"
			deleteRequest["OrganizationId"] = client.Department
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, deleteRequest, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpc Ipv6 Gateway (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Vpc Ipv6 Gateway success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccApsaraStackVPCIpv6Gateway_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VpcIpv6GatewaySupportRegions)
	resourceId := "apsarastack_vpc_ipv6_gateway.default"
	ra := resourceAttrInit(resourceId, ApsaraStackVPCIpv6GatewayMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeVpcIpv6Gateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6gateway%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackVPCIpv6GatewayBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":            "${apsarastack_vpc.default.id}",
					"ipv6_gateway_name": "${var.name}",
					"description":       "${var.name}",
					"spec":              "Small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":            CHECKSET,
						"ipv6_gateway_name": name,
						"description":       name,
						"spec":              "Small",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "Medium",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "Medium",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "Large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "Large",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_gateway_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_gateway_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":              "Small",
					"ipv6_gateway_name": "${var.name}",
					"description":       "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":              "Small",
						"ipv6_gateway_name": name,
						"description":       name,
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

var ApsaraStackVPCIpv6GatewayMap0 = map[string]string{
	"spec":   CHECKSET,
	"status": CHECKSET,
}

func ApsaraStackVPCIpv6GatewayBasicDependence0(name string) string {
	return fmt.Sprintf(` 
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}

resource "apsarastack_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}
`, name)
}
