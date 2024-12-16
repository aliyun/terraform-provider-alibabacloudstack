package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DnsService) DescribeDnsRecord(id string) (response *DnsRecord, err error) {
	var requestInfo *ecs.Client

	ZoneId := id
	request := s.client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalZoneRecords", "")
	request.QueryParams["ZoneId"] = ZoneId
	var resp = &DnsRecord{}
	raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRecordNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "DescribeGlobalZoneRecords", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("DescribeGlobalZoneRecords", response, requestInfo, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.AsapiSuccess == true {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(2)
	for {
		raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		bresponse, ok := raw.(*alidns.DescribeDomainGroupsResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return group, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		groups := bresponse.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return group, errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DnsGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *DnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	bresponse, ok := raw.(*alidns.ListTagResourcesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return *bresponse, nil
}

func (s *DnsService) DescribeDnsDomainAttachment(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	bresponse, ok := raw.(*alidns.DescribeInstanceDomainsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DnsDomainAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(bresponse.InstanceDomains) < 1 {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DnsDomainAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return
	}
	return *bresponse, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDnsDomainAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		domainNames := make(map[string]interface{}, 0)
		for _, v := range object.InstanceDomains {
			domainNames[v.DomainName] = v.DomainName
		}

		exceptDomainNames := make(map[string]interface{}, 0)
		for _, v := range expected {
			for _, vv := range v.([]interface{}) {
				exceptDomainNames[vv.(string)] = vv.(string)
			}
		}

		if reflect.DeepEqual(domainNames, exceptDomainNames) {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *DnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidns.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidns.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidns.CreateUntagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		addDebug(request.GetActionName(), raw)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	return nil
}

func (s *DnsService) DescribeDnsDomain(id string) (response *DnsDomains, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "CloudDns", "2021-06-24", "DescribeGlobalZones", "")
	request.QueryParams["Name"] = did[0]
	request.QueryParams["Forwardedregionid"] = s.client.RegionId
	request.QueryParams["SignatureVersion"] = "2.1"
	request.QueryParams["PageNumber"] = fmt.Sprint(1)
	request.QueryParams["PageSize"] = fmt.Sprint(PageSizeLarge)
	resp := &DnsDomains{}
	raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.ProcessCommonRequest(request)
	})
	addDebug("DescribeGlobalZones", response, nil, request)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ErrorDomainNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "DescribeGlobalZones", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("DescribeGlobalZones", response, nil, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.AsapiSuccess == true {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}
