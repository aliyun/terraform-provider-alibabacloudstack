package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDnsGroup() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDnsGroupCreate, resourceAlibabacloudStackDnsGroupRead, resourceAlibabacloudStackDnsGroupUpdate, resourceAlibabacloudStackDnsGroupDelete)
	return resource
}

func resourceAlibabacloudStackDnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	GroupName := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "GenesisDns", "2018-07-20", "AddDomainGroup", "")
	request.QueryParams["GroupName"] = GroupName
	raw, err := client.WithEcsClient(func(dnsClient *ecs.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_group", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request)
	response, _ := raw.(*alidns.AddDomainGroupResponse)
	d.SetId(response.GroupId)
	return nil
}

func resourceAlibabacloudStackDnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateUpdateDomainGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.GroupId = d.Id()

	if d.HasChange("name") {
		request.GroupName = d.Get("name").(string)
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainGroup(request)
		})
		bresponse, ok := raw.(*alidns.UpdateDomainGroupResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func resourceAlibabacloudStackDnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsGroup(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("name", object.GroupName)
	return nil
}

func resourceAlibabacloudStackDnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateDeleteDomainGroupRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.GroupId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(request)
		})
		bresponse, ok := raw.(*alidns.DeleteDomainGroupResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			if errmsgs.IsExpectedErrors(err, []string{"Fobidden.NotEmptyGroup"}) {
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultTimeoutMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
}
