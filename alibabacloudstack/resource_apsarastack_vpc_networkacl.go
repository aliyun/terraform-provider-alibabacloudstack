package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNetworkAcl() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"egress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
					},
				},
			},
			"ingress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"network_acl_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.122.0. New field 'network_acl_name' instead",
				ConflictsWith: []string{"network_acl_name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNetworkAclCreate, resourceAlibabacloudStackNetworkAclRead, resourceAlibabacloudStackNetworkAclUpdate, resourceAlibabacloudStackNetworkAclDelete)
	return resource
}

func resourceAlibabacloudStackNetworkAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	action := "CreateNetworkAcl"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := connectivity.GetResourceDataOk(d, "network_acl_name", "name"); ok {
		request["NetworkAclName"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateNetworkAcl")
	response, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	responseNetworkAclAttribute := response["NetworkAclAttribute"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseNetworkAclAttribute["NetworkAclId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeNetworkAcl(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_network_acl vpcService.DescribeNetworkAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("description", object["Description"])

	egressAclEntry := make([]map[string]interface{}, 0)
	if egressAclEntryList, ok := object["EgressAclEntries"].(map[string]interface{})["EgressAclEntry"].([]interface{}); ok {
		for _, v := range egressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"destination_cidr_ip":    m1["DestinationCidrIp"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
				}
				egressAclEntry = append(egressAclEntry, temp1)
			}
		}
	}
	if err := d.Set("egress_acl_entries", egressAclEntry); err != nil {
		return errmsgs.WrapError(err)
	}

	ingressAclEntry := make([]map[string]interface{}, 0)
	if ingressAclEntryList, ok := object["IngressAclEntries"].(map[string]interface{})["IngressAclEntry"].([]interface{}); ok {
		for _, v := range ingressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
					"source_cidr_ip":         m1["SourceCidrIp"],
				}
				ingressAclEntry = append(ingressAclEntry, temp1)
			}
		}
	}
	if err := d.Set("ingress_acl_entries", ingressAclEntry); err != nil {
		return errmsgs.WrapError(err)
	}
	connectivity.SetResourceData(d, object["NetworkAclName"], "network_acl_name", "name")

	resourceMap := make([]map[string]interface{}, 0)
	if resourceMapList, ok := object["Resources"].(map[string]interface{})["Resource"].([]interface{}); ok {
		for _, v := range resourceMapList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"resource_id":   m1["ResourceId"],
					"resource_type": m1["ResourceType"],
				}
				resourceMap = append(resourceMap, temp1)
			}
		}
	}
	if err := d.Set("resources", resourceMap); err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}

func resourceAlibabacloudStackNetworkAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChanges("name", "network_acl_name") {
		update = true
		request["NetworkAclName"] = connectivity.GetResourceData(d, "network_acl_name", "name").(string)
	}
	if update {
		action := "ModifyNetworkAclAttributes"
		request["ClientToken"] = buildClientToken("ModifyNetworkAclAttributes")
		_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	update = false
	updateNetworkAclEntriesReq := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	if d.HasChange("egress_acl_entries") {
		updateNetworkAclEntriesReq["UpdateEgressAclEntries"] = true
		update = true
		EgressAclEntries := make([]map[string]interface{}, len(d.Get("egress_acl_entries").([]interface{})))
		for i, EgressAclEntriesValue := range d.Get("egress_acl_entries").([]interface{}) {
			EgressAclEntriesMap := EgressAclEntriesValue.(map[string]interface{})
			EgressAclEntries[i] = make(map[string]interface{})
			EgressAclEntries[i]["Description"] = EgressAclEntriesMap["description"]
			EgressAclEntries[i]["DestinationCidrIp"] = EgressAclEntriesMap["destination_cidr_ip"]
			EgressAclEntries[i]["NetworkAclEntryName"] = EgressAclEntriesMap["network_acl_entry_name"]
			EgressAclEntries[i]["Policy"] = EgressAclEntriesMap["policy"]
			EgressAclEntries[i]["Port"] = EgressAclEntriesMap["port"]
			EgressAclEntries[i]["Protocol"] = EgressAclEntriesMap["protocol"]
		}
		updateNetworkAclEntriesReq["EgressAclEntries"] = EgressAclEntries
	}
	if d.HasChange("ingress_acl_entries") {
		updateNetworkAclEntriesReq["UpdateIngressAclEntries"] = true
		update = true
		IngressAclEntries := make([]map[string]interface{}, len(d.Get("ingress_acl_entries").([]interface{})))
		for i, IngressAclEntriesValue := range d.Get("ingress_acl_entries").([]interface{}) {
			IngressAclEntriesMap := IngressAclEntriesValue.(map[string]interface{})
			IngressAclEntries[i] = make(map[string]interface{})
			IngressAclEntries[i]["Description"] = IngressAclEntriesMap["description"]
			IngressAclEntries[i]["NetworkAclEntryName"] = IngressAclEntriesMap["network_acl_entry_name"]
			IngressAclEntries[i]["Policy"] = IngressAclEntriesMap["policy"]
			IngressAclEntries[i]["Port"] = IngressAclEntriesMap["port"]
			IngressAclEntries[i]["Protocol"] = IngressAclEntriesMap["protocol"]
			IngressAclEntries[i]["SourceCidrIp"] = IngressAclEntriesMap["source_cidr_ip"]
		}
		updateNetworkAclEntriesReq["IngressAclEntries"] = IngressAclEntries
	}
	if update {
		action := "UpdateNetworkAclEntries"
		updateNetworkAclEntriesReq["ClientToken"] = buildClientToken("UpdateNetworkAclEntries")
		_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, updateNetworkAclEntriesReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	if d.HasChange("resources") {
		oldResources, newResources := d.GetChange("resources")
		oldResourcesSet := oldResources.(*schema.Set)
		newResourcesSet := newResources.(*schema.Set)

		removed := oldResourcesSet.Difference(newResourcesSet)
		added := newResourcesSet.Difference(oldResourcesSet)
		if added.Len() > 0 {
			associatenetworkaclrequest := map[string]interface{}{
				"NetworkAclId": d.Id(),
			}

			resourcesMaps := make([]map[string]interface{}, 0)
			for _, resources := range added.List() {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := map[string]interface{}{
					"ResourceId":   resourcesArg["resource_id"],
					"ResourceType": resourcesArg["resource_type"],
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}
			associatenetworkaclrequest["Resource"] = resourcesMaps
			action := "AssociateNetworkAcl"
			associatenetworkaclrequest["ClientToken"] = buildClientToken("AssociateNetworkAcl")
			_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, associatenetworkaclrequest)
			if err != nil {
				return err
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
		if removed.Len() > 0 {
			unassociatenetworkaclrequest := map[string]interface{}{
				"NetworkAclId": d.Id(),
			}
			resourcesMaps := make([]map[string]interface{}, 0)
			for _, resources := range removed.List() {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := map[string]interface{}{
					"ResourceId":   resourcesArg["resource_id"],
					"ResourceType": resourcesArg["resource_type"],
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}
			unassociatenetworkaclrequest["Resource"] = resourcesMaps
			action := "UnassociateNetworkAcl"
			unassociatenetworkaclrequest["ClientToken"] = buildClientToken("UnassociateNetworkAcl")
			_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, unassociatenetworkaclrequest)
			if err != nil {
				return err
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	// Delete binging resources before delete the ACL
	_, err := vpcService.DeleteAclResources(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	action := "DeleteNetworkAcl"
	request := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	request["ClientToken"] = buildClientToken("DeleteNetworkAcl")
	_, err = client.DoTeaRequest("POST", "Vpc", "2016-04-28", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}