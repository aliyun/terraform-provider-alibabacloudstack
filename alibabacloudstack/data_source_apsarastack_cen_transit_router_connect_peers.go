package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCenTransitRouterConnectPeers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCenTransitRouterConnectPeersRead,

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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of Transit Router Connect Peer IDs to filter results by.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by the Transit Router Connect Peer name.",
			},
			"peer_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"peers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connect_attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCenTransitRouterConnectPeersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare request parameters
	reqQuery := make(map[string]interface{})
	reqQuery["CenId"] = d.Get("cen_id")
	reqQuery["TransitRouterAttachmentId"] = d.Get("connect_attachment_id")

	// Call the API to list transit router connect peers
	resp, err := client.DoTeaRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterConnectPeers", "", nil, reqQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to list transit router connect peers: %v", err)
	}

	// Parse response
	peersRaw, ok := resp["TransitRouterConnectPeers"]
	if !ok || peersRaw == nil || len(peersRaw.([]interface{})) == 0 {
		// Return empty list if no peers found
		d.SetId(dataResourceIdHash([]string{}))
		if err := d.Set("peers", []interface{}{}); err != nil {
			return fmt.Errorf("error setting 'peers': %v", err)
		}
		return nil
	}
	// Process filtering
	idsMap := getIdsStringFilter(d)
	// Prepare result
	var ids []string
	var result []map[string]interface{}
	for _, peerRaw := range peersRaw.([]interface{}) {
		peer := peerRaw.(map[string]interface{})
		peerId := fmt.Sprintf("%s:%s:%s", peer["CenId"], peer["TransitRouterAttachmentId"], peer["TransitRouterConnectPeerId"])
		if v, ok := d.GetOk("peer_name"); ok && v.(string) != peer["TransitRouterConnectPeerName"].(string) {
			continue
		}
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(peer["TransitRouterConnectPeerName"].(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, exist := idsMap[peerId]; !exist {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":                    peerId,
			"cen_id":                peer["CenId"],
			"connect_attachment_id": peer["TransitRouterAttachmentId"],
			"peer_id":               peer["TransitRouterConnectPeerId"],
			"name":                  peer["TransitRouterConnectPeerName"],
			"local_ip":              peer["TransitRouterConnectPeerLocalIp"],
			"peer_ip":               peer["TransitRouterConnectPeerPeerIp"],
			"region_id":             peer["TransitRouterRegionId"],
			"status":                peer["Status"],
			"creation_time":         peer["CreationTime"],
		}
		ids = append(ids, peerId)
		result = append(result, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("peers", result); err != nil {
		return fmt.Errorf("error setting 'peers': %v", err)
	}

	return nil
}
