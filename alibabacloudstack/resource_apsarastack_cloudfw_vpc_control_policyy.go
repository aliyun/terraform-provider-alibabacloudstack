package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackCloudfwVpcControlPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"acl_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proto": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_order": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest_port_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest_port_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hit_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_group_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_group_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dest_port_group_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCloudfwVpcControlPolicyCreate, resourceAlibabacloudStackCloudfwVpcControlPolicyRead, resourceAlibabacloudStackCloudfwVpcControlPolicyUpdate, resourceAlibabacloudStackCloudfwVpcControlPolicyDelete)
	return resource
}

func resourceAlibabacloudStackCloudfwVpcControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"SourceCode":      "yundun",
		"AclAction":       d.Get("acl_action").(string),
		"ApplicationId":   d.Get("application_id").(string),
		"ApplicationName": d.Get("application_name").(string),
		"Description":     d.Get("description").(string),
		"Destination":     d.Get("destination").(string),
		"DestinationType": d.Get("destination_type").(string),
		"VpcFirewallId":   d.Get("vpc_firewall_id").(string),
		"Proto":           d.Get("proto").(string),
		"Source":          d.Get("source").(string),
		"SourceType":      d.Get("source_type").(string),
		"NewOrder":        d.Get("new_order").(string),
	}

	// Add optional parameters if they exist
	if v, ok := d.GetOk("dest_port"); ok {
		request["DestPort"] = v.(string)
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		request["DestPortType"] = v.(string)
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		request["DestPortGroup"] = v.(string)
	}
	if v, ok := d.GetOk("release"); ok {
		request["Release"] = v.(string)
	}

	response, err := client.DoTeaRequest("POST", "cloudfw", "2017-12-07", "CreateVpcFirewallVpcControlPolicy", "", nil, request, nil)
	if err != nil {
		return err
	}

	// Get the AclUuid from response
	aclUuid, ok := response["AclUuid"]
	if !ok {
		return fmt.Errorf("AclUuid not found in response")
	}

	// Set the resource ID using the AclUuid
	d.SetId(aclUuid.(string))

	return nil
}

func resourceAlibabacloudStackCloudfwVpcControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudfwVpcControlPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("acl_action", object["AclAction"])
	d.Set("application_id", object["ApplicationId"])
	d.Set("application_name", object["ApplicationName"])
	d.Set("description", object["Description"])
	d.Set("destination", object["Destination"])
	d.Set("destination_type", object["DestinationType"])
	d.Set("vpc_firewall_id", object["VpcFirewallId"])
	d.Set("proto", object["Proto"])
	d.Set("source", object["Source"])
	d.Set("source_type", object["SourceType"])
	d.Set("dest_port", object["DestPort"])
	d.Set("dest_port_type", object["DestPortType"])
	d.Set("dest_port_group", object["DestPortGroup"])
	d.Set("release", object["Release"])
	d.Set("direction", object["Direction"])
	d.Set("acl_uuid", object["AclUuid"])
	d.Set("order", object["Order"])
	d.Set("hit_times", object["HitTimes"])

	if v, ok := object["SourceGroupCidrs"]; ok {
		sourceGroupCidrs := make([]string, 0)
		if cidrs, ok := v.([]interface{}); ok {
			for _, cidr := range cidrs {
				if cidrStr, ok := cidr.(string); ok {
					sourceGroupCidrs = append(sourceGroupCidrs, cidrStr)
				}
			}
		}
		d.Set("source_group_cidrs", sourceGroupCidrs)
	}

	if v, ok := object["DestinationGroupCidrs"]; ok {
		destinationGroupCidrs := make([]string, 0)
		if cidrs, ok := v.([]interface{}); ok {
			for _, cidr := range cidrs {
				if cidrStr, ok := cidr.(string); ok {
					destinationGroupCidrs = append(destinationGroupCidrs, cidrStr)
				}
			}
		}
		d.Set("destination_group_cidrs", destinationGroupCidrs)
	}

	if v, ok := object["DestPortGroupPorts"]; ok {
		destPortGroupPorts := make([]string, 0)
		if ports, ok := v.([]interface{}); ok {
			for _, port := range ports {
				if portStr, ok := port.(string); ok {
					destPortGroupPorts = append(destPortGroupPorts, portStr)
				}
			}
		}
		d.Set("dest_port_group_ports", destPortGroupPorts)
	}

	return nil
}

func resourceAlibabacloudStackCloudfwVpcControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// For new resources, skip update as Create already sets all attributes
	if d.IsNewResource() {
		return nil
	}

	// Prepare request parameters for update
	requestInfo := map[string]interface{}{
		"AclAction":       d.Get("acl_action"),
		"ApplicationId":   d.Get("application_id"),
		"ApplicationName": d.Get("application_name"),
		"Description":     d.Get("description"),
		"Destination":     d.Get("destination"),
		"DestinationType": d.Get("destination_type"),
		"VpcFirewallId":   d.Get("vpc_firewall_id"),
		"Proto":           d.Get("proto"),
		"Source":          d.Get("source"),
		"SourceType":      d.Get("source_type"),
		"AclUuid":         d.Id(), // Use the resource ID as AclUuid
		"Order":           d.Get("new_order"),
		"SourceCode":      "yundun", // Required parameter
	}

	// Add optional parameters if they exist
	if v, ok := d.GetOk("dest_port"); ok {
		requestInfo["DestPort"] = v
	}
	if v, ok := d.GetOk("dest_port_type"); ok {
		requestInfo["DestPortType"] = v
	}
	if v, ok := d.GetOk("dest_port_group"); ok {
		requestInfo["DestPortGroup"] = v
	}
	if v, ok := d.GetOk("release"); ok {
		requestInfo["Release"] = v
	}

	// Call the modify API
	_, err := client.DoTeaRequest("POST", "cloudfw", "2017-12-07", "ModifyVpcFirewallVpcControlPolicy", "", nil, requestInfo, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_cloudfw_control_policy", "ModifyVpcFirewallVpcControlPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackCloudfwVpcControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"SourceCode":    "yundun",
		"AclUuid":       d.Id(),
		"VpcFirewallId": d.Get("vpc_firewall_id"),
	}

	_, err := client.DoTeaRequest("POST", "cloudfw", "2017-12-07", "DeleteVpcFirewallVpcControlPolicy", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
