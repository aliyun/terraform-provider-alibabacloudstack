package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCenTransitRouterConnectPeer() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connect_attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peer_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource,
		resourceAlibabacloudStackCenTransitRouterConnectPeerCreate,
		resourceAlibabacloudStackCenTransitRouterConnectPeerRead,
		nil, resourceAlibabacloudStackCenTransitRouterConnectPeerDelete)
	return resource
}

func resourceAlibabacloudStackCenTransitRouterConnectPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}

	// Prepare the request parameters
	reqQuery := make(map[string]interface{})
	reqQuery["TransitRouterAttachmentId"] = d.Get("connect_attachment_id")
	reqQuery["TransitRouterConnectPeerName"] = d.Get("name")
	if v, ok := d.GetOk("local_ip"); ok {
		reqQuery["TransitRouterConnectPeerLocalIp"] = v
	}
	reqQuery["TransitRouterConnectPeerPeerIp"] = d.Get("peer_ip")

	// Call the API to create the resource
	resp, err := client.DoTeaRequest("POST", "cbn", "2017-09-12", "CreateTransitRouterConnectPeer", "", nil, reqQuery, nil)
	if err != nil {
		return err
	}

	// Extract the TransitRouterConnectPeerId from response
	transitRouterConnectPeerId, ok := resp["TransitRouterConnectPeerId"].(string)
	if !ok {
		return fmt.Errorf("unable to find TransitRouterConnectPeerId in response")
	}

	tempId := fmt.Sprintf("%s:%s:%s", d.Get("cen_id"), reqQuery["TransitRouterAttachmentId"], transitRouterConnectPeerId)
	d.SetId(tempId)

	// Wait for the resource to be active
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 3*time.Second,
		cbnService.CenTransitRouterConnectPeerStateRefreshFunc(tempId, []string{"Failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for CEN transit router connect peer to be active failed: %v", err)
	}

	return nil
}

func resourceAlibabacloudStackCenTransitRouterConnectPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}
	object, err := cbnService.DescribeCenTransitRouterConnectPeer(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cen_id", object["CenId"])
	d.Set("connect_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("peer_id", object["TransitRouterConnectPeerId"])
	d.Set("name", object["TransitRouterConnectPeerName"])
	d.Set("local_ip", object["TransitRouterConnectPeerLocalIp"])
	d.Set("peer_ip", object["TransitRouterConnectPeerPeerIp"])
	d.Set("region_id", object["TransitRouterRegionId"])

	return nil
}

func resourceAlibabacloudStackCenTransitRouterConnectPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}

	// Parse the resource ID to get the TransitRouterConnectPeerId
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	transitRouterConnectPeerId := parts[2]

	// Prepare the request query for deleting the transit router connect peer
	reqQuery := map[string]interface{}{
		"TransitRouterConnectPeerId": transitRouterConnectPeerId,
	}

	// Call the delete API
	_, err = client.DoTeaRequest("POST", "Cbn", "2017-09-12", "DeleteTransitRouterConnectPeer", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteTransitRouterConnectPeer", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Wait for the resource to be deleted
	stateConf := BuildStateConf([]string{"Available", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second,
		cbnService.CenTransitRouterConnectPeerStateRefreshFunc(d.Id(), []string{"Failed"}))

	_, err = stateConf.WaitForState()
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func (s *CenService) DescribeCenTransitRouterConnectPeer(id string) (map[string]interface{}, error) {
	// The resource ID is composed of CenId:TransitRouterAttachmentId:TransitRouterConnectPeerId
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, expected CenId:TransitRouterAttachmentId:TransitRouterConnectPeerId")
	}
	cenId := parts[0]
	transitRouterAttachmentId := parts[1]
	transitRouterConnectPeerId := parts[2]

	reqQuery := map[string]interface{}{
		"CenId":                      cenId,
		"TransitRouterAttachmentId":  transitRouterAttachmentId,
		"TransitRouterConnectPeerId": transitRouterConnectPeerId,
	}

	response, err := s.client.DoTeaRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterConnectPeers", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := response["TransitRouterConnectPeers"]; !ok || response["TransitRouterConnectPeers"] == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("cen transit router connect peer %s not found", transitRouterConnectPeerId))
	}

	peers := response["TransitRouterConnectPeers"].([]interface{})
	for _, peer := range peers {
		peerMap := peer.(map[string]interface{})
		if peerMap["TransitRouterConnectPeerId"].(string) == transitRouterConnectPeerId {
			return peerMap, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("cen transit router connect peer %s not found", transitRouterConnectPeerId))
}

func (s *CenService) CenTransitRouterConnectPeerStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterConnectPeer(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"] == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"]))
			}
		}
		return object, object["Status"].(string), nil
	}
}
