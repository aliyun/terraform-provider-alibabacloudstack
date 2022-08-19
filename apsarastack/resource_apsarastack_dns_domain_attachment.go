package apsarastack

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDnsDomainAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDnsDomainAttachmentCreate,
		Read:   resourceApasaraStackDnsDomainAttachmentRead,
		Update: resourceApsaraStackDnsDomainAttachmentUpdate,
		Delete: resourceApsaraStackDnsdomainAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceApsaraStackDnsDomainAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("instance_id").(string))
	return resourceApsaraStackDnsDomainAttachmentUpdate(d, meta)
}

func resourceApasaraStackDnsDomainAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}

	object, err := dnsService.DescribeDnsDomainAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", d.Id())
	d.Set("domain_names", flatten(object))
	return nil

}
func resourceApsaraStackDnsDomainAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}

	o, n := d.GetChange("domain_names")
	oldmap := make(map[string]string)
	newmap := make(map[string]string)
	add := make([]string, 0)
	remove := make([]string, 0)
	for _, v := range o.(*schema.Set).List() {
		oldmap[v.(string)] = v.(string)
	}
	for _, v := range n.(*schema.Set).List() {
		if _, ok := oldmap[v.(string)]; !ok {
			add = append(add, v.(string))
		}
	}

	for _, v := range n.(*schema.Set).List() {
		newmap[v.(string)] = v.(string)
	}
	for _, v := range o.(*schema.Set).List() {
		if _, ok := newmap[v.(string)]; !ok {
			remove = append(remove, v.(string))
		}
	}
	if len(remove) > 0 {
		removeNames := strings.Join(remove, ",")
		request := alidns.CreateUnbindInstanceDomainsRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = client.Department
		request.QueryParams["ResourceGroup"] = client.ResourceGroup
		request.InstanceId = d.Id()
		request.DomainNames = removeNames
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UnbindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	if len(add) > 0 {
		addNames := strings.Join(add, ",")
		request := alidns.CreateBindInstanceDomainsRequest()
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = client.Department
		request.QueryParams["ResourceGroup"] = client.ResourceGroup
		request.InstanceId = d.Id()
		request.DomainNames = addNames
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.BindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
	}
	if err := dnsService.WaitForAlidnsDomainAttachment(d.Id(), map[string]interface{}{"Domain": d.Get("domain_names").(*schema.Set).List()}, false, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceApasaraStackDnsDomainAttachmentRead(d, meta)
}

func resourceApsaraStackDnsdomainAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dnsService := DnsService{client}

	domainNames := d.Get("domain_names").(*schema.Set).List()
	deleteSli := make([]string, 0)
	for _, v := range domainNames {
		deleteSli = append(deleteSli, v.(string))
	}

	request := alidns.CreateUnbindInstanceDomainsRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.InstanceId = d.Id()
	request.DomainNames = strings.Join(deleteSli, ",")

	raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.UnbindInstanceDomains(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	return WrapError(dnsService.WaitForAlidnsDomainAttachment(d.Id(), nil, true, DefaultTimeout))
}

func flatten(input alidns.DescribeInstanceDomainsResponse) []string {
	domainNames := make([]string, 0)
	for _, v := range input.InstanceDomains {
		domainNames = append(domainNames, v.DomainName)
	}
	return domainNames
}
