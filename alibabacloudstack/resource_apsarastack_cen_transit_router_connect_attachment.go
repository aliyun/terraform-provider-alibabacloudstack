package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCenTransitRouterConnectAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_attachment_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackCenTransitRouterConnectAttachmentCreate,
		resourceAlibabacloudStackCenTransitRouterConnectAttachmentRead,
		nil,
		resourceAlibabacloudStackCenTransitRouterConnectAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackCenTransitRouterConnectAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}

	request := make(map[string]interface{})
	request["CenId"] = d.Get("cen_id")
	request["TransitRouterId"] = d.Get("transit_router_id")
	if v, ok := d.GetOk("transit_router_attachment_name"); ok {
		request["TransitRouterAttachmentName"] = v
	}

	action := "CreateTransitRouterConnectAttachment"
	response, err := client.DoTeaRequest("POST", "Cbn", "2017-09-12", action, "", nil, request, nil)
	if err != nil {
		return err
	}

	attachmentId, ok := response["TransitRouterAttachmentId"].(string)
	if !ok || attachmentId == "" {
		return fmt.Errorf("unable to find TransitRouterAttachmentId in response")
	}

	// Construct the ID using CenId:TransitRouterId:TransitRouterAttachmentId
	cenId := d.Get("cen_id").(string)
	transitRouterId := d.Get("transit_router_id").(string)
	d.SetId(fmt.Sprintf("%s:%s:%s", cenId, transitRouterId, attachmentId))

	// Wait for the attachment to be in the "Attached" state
	stateConf := BuildStateConf([]string{"Creating", "Attaching"}, []string{"Attached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterConnectAttachmentStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for CEN transit router connect attachment to be attached failed: %v", err)
	}

	return nil
}

func resourceAlibabacloudStackCenTransitRouterConnectAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}

	object, err := cbnService.DescribeCenTransitRouterConnectAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_cen_transit_router_connect_attachment cbnService.DescribeCenTransitRouterConnectAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid id format, expected {CenId}:{TransitRouterId}:{TransitRouterAttachmentId}")
	}

	d.Set("cen_id", parts[0])
	d.Set("transit_router_id", parts[1])
	d.Set("transit_router_attachment_id", object["TransitRouterAttachmentId"])
	d.Set("resource_id", object["ResourceId"])
	d.Set("creation_time", object["CreationTime"])
	d.Set("resource_type", object["ResourceType"])
	d.Set("resource_owner_id", object["ResourceOwnerId"])
	d.Set("transit_router_attachment_name", object["TransitRouterAttachmentName"])

	return nil
}

func resourceAlibabacloudStackCenTransitRouterConnectAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	// No update API is supported for this resource
	return nil
}

func resourceAlibabacloudStackCenTransitRouterConnectAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cbnService := CenService{client}

	// Parse the resource ID to get the TransitRouterAttachmentId
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	transitRouterAttachmentId := parts[2]

	// Prepare the request query parameters
	reqQuery := map[string]interface{}{
		"CenId":                     parts[0],
		"TransitRouterAttachmentId": transitRouterAttachmentId,
	}

	// Call the delete API
	_, err = client.DoTeaRequest("POST", "Cbn", "2017-09-12", "DeleteTransitRouterConnectAttachment", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteTransitRouterConnectAttachment", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Wait for the resource to be deleted
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, cbnService.CenTransitRouterConnectAttachmentStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return errmsgs.WrapError(err)
}
