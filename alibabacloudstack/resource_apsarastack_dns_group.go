package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDnsGroupCreate,
		Read:   resourceAlibabacloudStackDnsGroupRead,
		Update: resourceAlibabacloudStackDnsGroupUpdate,
		Delete: resourceAlibabacloudStackDnsGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlibabacloudStackDnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	//request := alidns.CreateAddDomainGroupRequest()
	GroupName := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"        // Set request method
	request.Product = "GenesisDns" // Specify product
	request.Domain = client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2018-07-20" // Specify product version
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "AddDomainGroup"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "GenesisDns",
		"RegionId":        client.RegionId,
		"Action":          "AddDomainGroup",
		"Version":         "2018-07-20",
		"GroupName":       GroupName,
	}
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_dns_group", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request)
	response, _ := raw.(*alidns.AddDomainGroupResponse)
	d.SetId(response.GroupId)
	return resourceAlibabacloudStackDnsGroupRead(d, meta)
}

func resourceAlibabacloudStackDnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateUpdateDomainGroupRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "GenesisDns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	if d.HasChange("name") {
		request.GroupName = d.Get("name").(string)
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlibabacloudStackDnsGroupRead(d, meta)
}

func resourceAlibabacloudStackDnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.GroupName)
	return nil
}

func resourceAlibabacloudStackDnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateDeleteDomainGroupRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Fobidden.NotEmptyGroup"}) {
				return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
}
