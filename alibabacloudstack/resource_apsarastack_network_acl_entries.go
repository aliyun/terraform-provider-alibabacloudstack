package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNetworkAclEntries() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackNetworkAclEntriesCreate,
		Read:   resourceAlibabacloudStackNetworkAclEntriesRead,
		Update: resourceAlibabacloudStackNetworkAclEntriesUpdate,
		Delete: resourceAlibabacloudStackNetworkAclEntriesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ingress": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"entry_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"egress": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"destination_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"entry_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlibabacloudStackNetworkAclEntriesCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("network_acl_id").(string) + COLON_SEPARATED + resource.UniqueId())
	return resourceAlibabacloudStackNetworkAclEntriesUpdate(d, meta)
}

func resourceAlibabacloudStackNetworkAclEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	object, err := vpcService.DescribeNetworkAcl(parts[0])
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	var ingress []map[string]interface{}
	if ingressAclEntryList, ok := object["IngressAclEntries"].(map[string]interface{})["IngressAclEntry"].([]interface{}); ok {
		for _, ob := range ingressAclEntryList {
			if v, ok := ob.(map[string]interface{}); ok {
				mapping := map[string]interface{}{
					"description":     v["Description"],
					"source_cidr_ip":  v["SourceCidrIp"],
					"entry_type":      "custom",
					"name":            v["NetworkAclEntryName"],
					"policy":          v["Policy"],
					"port":            v["Port"],
					"protocol":        v["Protocol"],
				}
				ingress = append(ingress, mapping)
			}
		}
	}

	var egress []map[string]interface{}
	if egressAclEntryList, ok := object["EgressAclEntries"].(map[string]interface{})["EgressAclEntry"].([]interface{}); ok {
		for _, ob := range egressAclEntryList {
			if v, ok := ob.(map[string]interface{}); ok {
				mapping := map[string]interface{}{
					"description":         v["Description"],
					"destination_cidr_ip": v["DestinationCidrIp"],
					"entry_type":          "custom",
					"name":                v["NetworkAclEntryName"],
					"policy":              v["Policy"],
					"port":                v["Port"],
					"protocol":            v["Protocol"],
				}
				egress = append(egress, mapping)
			}
		}
	}

	d.Set("network_acl_id", object["NetworkAclId"])
	d.Set("egress", egress)
	d.Set("ingress", ingress)

	return nil
}

func resourceAlibabacloudStackNetworkAclEntriesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclId := parts[0]

	request := vpc.CreateUpdateNetworkAclEntriesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NetworkAclId = networkAclId

	if d.HasChange("ingress") {
		ingress := []vpc.UpdateNetworkAclEntriesIngressAclEntries{}
		for _, e := range d.Get("ingress").([]interface{}) {
			ingress = append(ingress, vpc.UpdateNetworkAclEntriesIngressAclEntries{
				Protocol:             e.(map[string]interface{})["protocol"].(string),
				Port:                 e.(map[string]interface{})["port"].(string),
				SourceCidrIp:         e.(map[string]interface{})["source_cidr_ip"].(string),
				NetworkAclEntryName:  e.(map[string]interface{})["name"].(string),
				EntryType:            e.(map[string]interface{})["entry_type"].(string),
				Policy:               e.(map[string]interface{})["policy"].(string),
				Description:          e.(map[string]interface{})["description"].(string),
			})
		}
		request.IngressAclEntries = &ingress
		request.UpdateIngressAclEntries = requests.NewBoolean(true)
	}

	if d.HasChange("egress") {
		egress := []vpc.UpdateNetworkAclEntriesEgressAclEntries{}
		for _, e := range d.Get("egress").([]interface{}) {
			egress = append(egress, vpc.UpdateNetworkAclEntriesEgressAclEntries{
				Protocol:             e.(map[string]interface{})["protocol"].(string),
				Port:                 e.(map[string]interface{})["port"].(string),
				DestinationCidrIp:    e.(map[string]interface{})["destination_cidr_ip"].(string),
				NetworkAclEntryName:  e.(map[string]interface{})["name"].(string),
				EntryType:            e.(map[string]interface{})["entry_type"].(string),
				Policy:               e.(map[string]interface{})["policy"].(string),
				Description:          e.(map[string]interface{})["description"].(string),
			})
		}
		request.EgressAclEntries = &egress
		request.UpdateEgressAclEntries = requests.NewBoolean(true)
	}

	if err := vpcService.WaitForNetworkAcl(networkAclId, Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UpdateNetworkAclEntries(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*vpc.UpdateNetworkAclEntriesResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return vpcService.WaitForNetworkAcl(networkAclId, Available, DefaultTimeout)
}

func resourceAlibabacloudStackNetworkAclEntriesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	networkAclId := parts[0]

	request := vpc.CreateUpdateNetworkAclEntriesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.NetworkAclId = networkAclId

	ingress := []vpc.UpdateNetworkAclEntriesIngressAclEntries{}
	egress := []vpc.UpdateNetworkAclEntriesEgressAclEntries{}
	request.IngressAclEntries = &ingress
	request.EgressAclEntries = &egress
	request.UpdateIngressAclEntries = requests.NewBoolean(true)
	request.UpdateEgressAclEntries = requests.NewBoolean(true)

	if err := vpcService.WaitForNetworkAcl(networkAclId, Available, DefaultTimeout); err != nil {
		return errmsgs.WrapError(err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UpdateNetworkAclEntries(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*vpc.UpdateNetworkAclEntriesResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return vpcService.WaitForNetworkAcl(networkAclId, Available, DefaultTimeout)
}
