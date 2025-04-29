package alibabacloudstack

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsDomainAttachment() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackDnsDomainAttachmentCreate,
		resourceApasaraStackDnsDomainAttachmentRead, resourceAlibabacloudStackDnsDomainAttachmentUpdate, resourceAlibabacloudStackDnsdomainAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackDnsDomainAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("instance_id").(string))
	return resourceAlibabacloudStackDnsDomainAttachmentUpdate(d, meta)
}

func resourceApasaraStackDnsDomainAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}

	object, err := dnsService.DescribeDnsDomainAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("instance_id", d.Id())
	d.Set("domain_names", flatten(object))
	return nil
}

func resourceAlibabacloudStackDnsDomainAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
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
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.DomainNames = removeNames
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UnbindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*alidns.UnbindInstanceDomainsResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	if len(add) > 0 {
		addNames := strings.Join(add, ",")
		request := alidns.CreateBindInstanceDomainsRequest()
		client.InitRpcRequest(*request.RpcRequest)
		request.InstanceId = d.Id()
		request.DomainNames = addNames
		raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.BindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if response, ok := raw.(*alidns.BindInstanceDomainsResponse); ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	if err := dnsService.WaitForAlidnsDomainAttachment(d.Id(), map[string]interface{}{"Domain": d.Get("domain_names").(*schema.Set).List()}, false, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func resourceAlibabacloudStackDnsdomainAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := DnsService{client}

	domainNames := d.Get("domain_names").(*schema.Set).List()
	deleteSli := make([]string, 0)
	for _, v := range domainNames {
		deleteSli = append(deleteSli, v.(string))
	}

	request := alidns.CreateUnbindInstanceDomainsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = d.Id()
	request.DomainNames = strings.Join(deleteSli, ",")

	raw, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.UnbindInstanceDomains(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		errmsg := ""
		if response, ok := raw.(*alidns.UnbindInstanceDomainsResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return errmsgs.WrapError(dnsService.WaitForAlidnsDomainAttachment(d.Id(), nil, true, DefaultTimeout))
}

func flatten(input alidns.DescribeInstanceDomainsResponse) []string {
	domainNames := make([]string, 0)
	for _, v := range input.InstanceDomains {
		domainNames = append(domainNames, v.DomainName)
	}
	return domainNames
}