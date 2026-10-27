package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_cen_instance", &resource.Sweeper{
		Name: "alibabacloudstack_cen_instance",
		F: func(region string) error {
			return SweeperFunc(
				sweeperEntry{fn: testSweepCenInstances, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource", "tf-testaccceninstance"}, broad: broadPrefixes, fatal: true, waitSec: 10},
				sweeperEntry{fn: testSweepVSwitches, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource_vsw"}},
				sweeperEntry{fn: testSweepVPCs, precise: []string{"tf-testAccMulticastDomainAssociationsDatasource_vpc"}},
				sweeperEntry{fn: testSweepVirtualBorderRouters, precise: []string{"tf-testacc-ecvbr"}},
			)
		},
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

	leftover := []string{}

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

		// Step 2: enumerate the transit routers. Connect attachments, connect peers and
		// transit router CIDR blocks are only reachable through a transit router, and a
		// transit router created by alibabacloudstack_cen_instance carries no name, so it
		// cannot be swept on its own and is cleaned up here as a child of the matched CEN.
		transitRouters, err := cenSweepListTransitRouters(client, cen.CenId)
		if err != nil {
			log.Printf("[WARN] Error listing transit routers for CEN %s: %s", cen.CenId, err)
		}

		// Step 3: delete the connect peers, then the connect attachments that hold the
		// CIDR blocks allocated for them.
		for _, transitRouter := range transitRouters {
			connectListParams := map[string]string{"CenId": cen.CenId, "TransitRouterId": transitRouter.TransitRouterId}
			for _, attachId := range cenSweepIdsOrWarn(client, "ListTransitRouterConnectAttachments", connectListParams, "TransitRouterAttachmentId") {
				peerListParams := map[string]string{"CenId": cen.CenId, "TransitRouterAttachmentId": attachId}
				for _, peerId := range cenSweepIdsOrWarn(client, "ListTransitRouterConnectPeers", peerListParams, "TransitRouterConnectPeerId") {
					log.Printf("[INFO] Deleting connect peer: %s", peerId)
					if err := cenSweepDelete(client, "DeleteTransitRouterConnectPeer", map[string]string{"TransitRouterConnectPeerId": peerId}); err != nil {
						log.Printf("[WARN] Error deleting connect peer %s: %s", peerId, err)
						leftover = append(leftover, peerId)
						continue
					}
					if err := cenSweepWaitGone(client, "ListTransitRouterConnectPeers", peerListParams, "TransitRouterConnectPeerId", peerId); err != nil {
						log.Printf("[WARN] %s", err)
						leftover = append(leftover, peerId)
					}
				}
				if err := cenSweepDeleteAttachment(client, "DeleteTransitRouterConnectAttachment", "ListTransitRouterConnectAttachments", "TransitRouterAttachmentId", connectListParams, attachId); err != nil {
					log.Printf("[WARN] %s", err)
					leftover = append(leftover, attachId)
				}
			}
		}

		// Step 4: detach the VBR and VPC attachments. ListTransitRouterVbrAttachments
		// needs TransitRouterId next to CenId on 17E, while the VPC list accepts CenId
		// alone. Both delete calls take CenId plus TransitRouterAttachmentId, and neither
		// needs its route table entries dissociated first, as observed while clearing the
		// connect_peer chain.
		for _, transitRouter := range transitRouters {
			vbrListParams := map[string]string{"CenId": cen.CenId, "TransitRouterId": transitRouter.TransitRouterId}
			for _, attachId := range cenSweepIdsOrWarn(client, "ListTransitRouterVbrAttachments", vbrListParams, "TransitRouterAttachmentId") {
				if err := cenSweepDeleteAttachment(client, "DeleteTransitRouterVbrAttachment", "ListTransitRouterVbrAttachments", "TransitRouterAttachmentId", vbrListParams, attachId); err != nil {
					log.Printf("[WARN] %s", err)
					leftover = append(leftover, attachId)
				}
			}
		}
		vpcListParams := map[string]string{"CenId": cen.CenId}
		for _, attachId := range cenSweepIdsOrWarn(client, "ListTransitRouterVpcAttachments", vpcListParams, "TransitRouterAttachmentId") {
			if err := cenSweepDeleteAttachment(client, "DeleteTransitRouterVpcAttachment", "ListTransitRouterVpcAttachments", "TransitRouterAttachmentId", vpcListParams, attachId); err != nil {
				log.Printf("[WARN] %s", err)
				leftover = append(leftover, attachId)
			}
		}

		// Step 5: release the CIDR blocks. The server rejects this with
		// OperationFailed.CidrBlockAllocated while a connect peer still holds a block
		// carved out of them, hence it runs after steps 3 and 4.
		for _, transitRouter := range transitRouters {
			for _, cidr := range transitRouter.CidrList {
				if err := cenSweepDelete(client, "DeleteTransitRouterCidr", map[string]string{"TransitRouterId": transitRouter.TransitRouterId, "TransitRouterCidrId": cidr.TransitRouterCidrId}); err != nil {
					log.Printf("[WARN] Error deleting transit router cidr %s: %s", cidr.TransitRouterCidrId, err)
					leftover = append(leftover, cidr.TransitRouterCidrId)
				}
			}
		}

		// Step 6: delete the transit routers themselves.
		trListParams := map[string]string{"CenId": cen.CenId}
		for _, transitRouter := range transitRouters {
			log.Printf("[INFO] Deleting transit router: %s", transitRouter.TransitRouterId)
			if err := cenSweepDelete(client, "DeleteTransitRouter", map[string]string{"CenId": cen.CenId, "TransitRouterId": transitRouter.TransitRouterId}); err != nil {
				log.Printf("[WARN] Error deleting transit router %s: %s", transitRouter.TransitRouterId, err)
				leftover = append(leftover, transitRouter.TransitRouterId)
				continue
			}
			if err := cenSweepWaitGone(client, "ListTransitRouters", trListParams, "TransitRouterId", transitRouter.TransitRouterId); err != nil {
				log.Printf("[WARN] %s", err)
				leftover = append(leftover, transitRouter.TransitRouterId)
			}
		}

		// Step 7: delete CEN instance
		log.Printf("[INFO] Deleting CEN instance: %s", cen.CenId)
		delCenRequest := client.NewCommonRequest("POST", "Cbn", "2017-09-12", "DeleteCen", "")
		delCenRequest.QueryParams["CenId"] = cen.CenId
		if _, err := client.ProcessCommonRequest(delCenRequest); err != nil {
			log.Printf("[WARN] Error deleting CEN %s: %s", cen.CenId, err)
			leftover = append(leftover, cen.CenId)
		}
	}

	if len(leftover) > 0 {
		return fmt.Errorf("CEN sweep left %d resource(s) behind: %s", len(leftover), strings.Join(leftover, ", "))
	}
	return nil
}

// Polling budget for the CEN sweep. The CEN team's agreed threshold is that an
// asynchronous operation still unfinished after roughly 300s needs manual
// intervention, so the sweeper reports such a resource as leftover instead of
// stretching its wait to cover an environment-side stall.
const (
	cenSweepWaitTimeout  = 5 * time.Minute
	cenSweepPollInterval = 10 * time.Second
)

// cenSweepTransitRouter carries only the fields the sweep needs: the transit router id
// and the ids of the CIDR blocks configured on it.
type cenSweepTransitRouter struct {
	TransitRouterId string `json:"TransitRouterId"`
	CidrList        []struct {
		TransitRouterCidrId string `json:"TransitRouterCidrId"`
	} `json:"TransitRouterCidrList"`
}

func cenSweepListTransitRouters(client *connectivity.AlibabacloudStackClient, cenId string) ([]cenSweepTransitRouter, error) {
	request := client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouters", "")
	request.QueryParams["CenId"] = cenId
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	var response struct {
		RequestId      string                  `json:"RequestId"`
		TransitRouters []cenSweepTransitRouter `json:"TransitRouters"`
	}
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &response); err != nil {
		return nil, err
	}
	if response.RequestId == "" {
		return nil, fmt.Errorf("ListTransitRouters for %s came back without a RequestId", cenId)
	}
	return response.TransitRouters, nil
}

// cenSweepQuery issues a Cbn list request and picks `idKey` out of every object of the
// array the response carries. The second return value reports whether the response is a
// usable read at all: 17E intermittently answers with an empty body, and that must not be
// mistaken for "the resource is gone". An empty body fails to unmarshal and a successful
// response always carries a RequestId, so either one marks the read as unusable.
func cenSweepQuery(client *connectivity.AlibabacloudStackClient, action string, params map[string]string, idKey string) ([]string, bool, error) {
	request := client.NewCommonRequest("GET", "Cbn", "2017-09-12", action, "")
	for key, value := range params {
		request.QueryParams[key] = value
	}
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, false, err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &raw); err != nil {
		return nil, false, err
	}
	if _, ok := raw["RequestId"]; !ok {
		return nil, false, nil
	}
	ids := []string{}
	for _, value := range raw {
		objects, ok := value.([]interface{})
		if !ok {
			continue
		}
		for _, item := range objects {
			object, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if id, ok := object[idKey].(string); ok && id != "" {
				ids = append(ids, id)
			}
		}
	}
	return ids, true, nil
}

func cenSweepIdsOrWarn(client *connectivity.AlibabacloudStackClient, action string, params map[string]string, idKey string) []string {
	ids, valid, err := cenSweepQuery(client, action, params, idKey)
	if err != nil {
		log.Printf("[WARN] Error listing %s: %s", action, err)
		return nil
	}
	if !valid {
		log.Printf("[WARN] %s returned a response without a RequestId, skipping", action)
		return nil
	}
	if len(ids) > 0 {
		log.Printf("[INFO] Found %d %s for %v", len(ids), idKey, params)
	}
	return ids
}

func cenSweepDelete(client *connectivity.AlibabacloudStackClient, action string, params map[string]string) error {
	request := client.NewCommonRequest("POST", "Cbn", "2017-09-12", action, "")
	for key, value := range params {
		request.QueryParams[key] = value
	}
	_, err := client.ProcessCommonRequest(request)
	return err
}

// cenSweepDeleteAttachment deletes one attachment and waits until it disappears from the
// list it was found in. Attachments go through Detaching asynchronously, and a CEN cannot
// be deleted while any of them survives, so the wait is part of the delete.
func cenSweepDeleteAttachment(client *connectivity.AlibabacloudStackClient, delAction, listAction, idKey string, listParams map[string]string, attachId string) error {
	log.Printf("[INFO] Deleting %s: %s", idKey, attachId)
	params := map[string]string{"CenId": listParams["CenId"], "TransitRouterAttachmentId": attachId}
	if err := cenSweepDelete(client, delAction, params); err != nil {
		return err
	}
	return cenSweepWaitGone(client, listAction, listParams, idKey, attachId)
}

// cenSweepWaitGone reports the resource as gone only after two consecutive usable reads
// leave it out, and gives up with an error once the polling budget runs out.
func cenSweepWaitGone(client *connectivity.AlibabacloudStackClient, listAction string, listParams map[string]string, idKey, id string) error {
	deadline := time.Now().Add(cenSweepWaitTimeout)
	absent := 0
	for time.Now().Before(deadline) {
		ids, valid, err := cenSweepQuery(client, listAction, listParams, idKey)
		switch {
		case err != nil:
			log.Printf("[WARN] Poll of %s for %s failed: %s", listAction, id, err)
		case !valid:
			log.Printf("[WARN] Poll of %s for %s returned a response without a RequestId", listAction, id)
		default:
			found := false
			for _, got := range ids {
				if got == id {
					found = true
					break
				}
			}
			if found {
				absent = 0
			} else {
				absent++
				if absent >= 2 {
					return nil
				}
			}
		}
		time.Sleep(cenSweepPollInterval)
	}
	return fmt.Errorf("%s %s is still present after %s, manual intervention required", idKey, id, cenSweepWaitTimeout)
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
