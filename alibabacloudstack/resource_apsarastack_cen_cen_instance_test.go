package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_cen_instance", &resource.Sweeper{
		Name: "alibabacloudstack_cen_instance",
		F: func(region string) error { return SweeperFunc(
			sweeperEntry{fn: testSweepCenInstances, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource"}, broad: broadPrefixes, fatal: true, waitSec: 10},
			sweeperEntry{fn: testSweepVSwitches, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource_vsw"}},
			sweeperEntry{fn: testSweepVPCs, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource_vpc"}},
		) },
	})
}

func testSweepCenInstances(prefixes []string) error {
		rawClient, err := sharedClientForRegion("")
		if err != nil {
			return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
		}
		client := rawClient.(*connectivity.AlibabacloudStackClient)

		// Step 1: List all CEN instances (with pagination)
		type cenEntry struct {
			CenId       string `json:"CenId"`
			Name        string `json:"Name"`
			Description string `json:"Description"`
		}
		var allCens []cenEntry
		pageNumber := 1
		for {
			request := client.NewCommonRequest("GET", "Cbn", "2017-09-12", "DescribeCens", "")
			request.QueryParams["PageSize"] = "50"
			request.QueryParams["PageNumber"] = fmt.Sprintf("%d", pageNumber)
			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				return fmt.Errorf("error listing CEN instances: %s", err)
			}
			var censResponse struct {
				Cens struct {
					Cen []cenEntry `json:"Cen"`
				} `json:"Cens"`
			}
			if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &censResponse); err != nil {
				return fmt.Errorf("error parsing CEN list: %s", err)
			}
			if len(censResponse.Cens.Cen) == 0 {
				break
			}
			allCens = append(allCens, censResponse.Cens.Cen...)
			if len(censResponse.Cens.Cen) < 50 {
				break
			}
			pageNumber++
		}

		for _, cen := range allCens {
			matched := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(cen.Name, prefix) || strings.HasPrefix(cen.Description, prefix) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			log.Printf("[INFO] Sweeping CEN instance: %s (%s)", cen.Name, cen.CenId)

			// Step 2: List and delete transit router VPC attachments
			// Try listing VPC attachments directly with CenId (no TransitRouterId needed)
			vpcAttRequest := client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterVpcAttachments", "")
			vpcAttRequest.QueryParams["CenId"] = cen.CenId
			vpcAttResp, err := client.ProcessCommonRequest(vpcAttRequest)
			if err != nil {
				log.Printf("[WARN] Error listing VPC attachments for CEN %s: %s", cen.CenId, err)
			} else {
				var attResponse struct {
					TransitRouterAttachments []struct {
						TransitRouterAttachmentId string `json:"TransitRouterAttachmentId"`
					} `json:"TransitRouterAttachments"`
				}
				if err := json.Unmarshal(vpcAttResp.GetHttpContentBytes(), &attResponse); err != nil {
					log.Printf("[WARN] Error parsing VPC attachments for CEN %s: %s", cen.CenId, err)
				} else {
					log.Printf("[INFO] Found %d VPC attachments for CEN %s", len(attResponse.TransitRouterAttachments), cen.CenId)
					for _, att := range attResponse.TransitRouterAttachments {
						log.Printf("[INFO] Deleting VPC attachment: %s", att.TransitRouterAttachmentId)
						delRequest := client.NewCommonRequest("POST", "Cbn", "2017-09-12", "DeleteTransitRouterVpcAttachment", "")
						delRequest.QueryParams["CenId"] = cen.CenId
						delRequest.QueryParams["TransitRouterAttachmentId"] = att.TransitRouterAttachmentId
						if _, err := client.ProcessCommonRequest(delRequest); err != nil {
							log.Printf("[WARN] Error deleting VPC attachment %s: %s", att.TransitRouterAttachmentId, err)
						}
					}
				}
			}

			// Wait for VPC attachments to be fully deleted
			time.Sleep(10 * time.Second)

			// Step 3: Delete CEN instance
			log.Printf("[INFO] Deleting CEN instance: %s", cen.CenId)
			delCenRequest := client.NewCommonRequest("POST", "Cbn", "2017-09-12", "DeleteCen", "")
			delCenRequest.QueryParams["CenId"] = cen.CenId
			if _, err := client.ProcessCommonRequest(delCenRequest); err != nil {
				log.Printf("[WARN] Error deleting CEN %s: %s", cen.CenId, err)
			}
		}



	return nil
}

func TestAccAlibabacloudStackCenInstance0(t *testing.T) {
	var v *CenInstance

	resourceId := "alibabacloudstack_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenInstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribecensRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccceninstance%d", rand)
	modify_name := fmt.Sprintf("tf-testaccceninstancemodify%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenInstanceBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name":          name,
					"description":                name,
					"transit_router_name":        name,
					"transit_router_description": name,
					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.0.0/16"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":            CHECKSET,
						"transit_router_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": modify_name,
					"description":       modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": modify_name,
						"description":       modify_name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_name":        modify_name,
					"transit_router_description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_name":        modify_name,
						"transit_router_description": modify_name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.10.2/24"},
						{"cidr": "10.10.11.2/24"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_cidrs.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.10.2/24"},
						{"cidr": "10.10.12.2/24"},
						{"cidr": "10.10.13.2/24"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_cidrs.#": "3",
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

var AlibabacloudTestAccCenInstanceCheckmap = map[string]string{

	"status": CHECKSET,

	"protection_level": CHECKSET,

	"create_time": CHECKSET,

	"cen_id": CHECKSET,
}

func AlibabacloudTestAccCenInstanceBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
