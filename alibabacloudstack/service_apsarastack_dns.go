package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
	"time"
)

type DnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DnsService) DescribeDnsRecord(id string) (response *DnsRecord, err error) {
	var requestInfo *ecs.Client

	if err != nil {
		return response, WrapError(err)
	}
	ZoneId := id
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "CloudDns",
		"Action":          "DescribeGlobalZoneRecords",
		"Version":         "2021-06-24",
		"ZoneId":          ZoneId,
	}
	request.Method = "POST"
	request.Product = "CloudDns"
	request.Version = "2021-06-24"
	request.ServiceCode = "CloudDns"
	request.Domain = s.client.Domain
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(2)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DescribeGlobalZoneRecords"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &DnsRecord{}
	raw, err := s.client.WithEcsClient(func(cmsClient *ecs.Client) (interface{}, error) {
		return cmsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRecordNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "DescribeGlobalZoneRecords", AlibabacloudStackSdkGoERROR)

	}
	addDebug("DescribeGlobalZoneRecords", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.AsapiSuccess == true {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *DnsService) DescribeDnsGroup(id string) (alidns.DomainGroup, error) {
	var group alidns.DomainGroup
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = s.client.Department
	request.QueryParams["ResourceGroup"] = s.client.ResourceGroup
	request.RegionId = s.client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(2)
	for {
		raw, err := s.client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			return group, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if domainGroup.GroupId == id {
				return domainGroup, nil
			}
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return group, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return group, WrapErrorf(Error(GetNotFoundMessage("DnsGroup", id)), NotFoundMsg, ProviderERROR)
}

func (s *DnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = s.client.Department
	request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.ListTagResourcesResponse)
	return *response, nil
}
func (s *DnsService) DescribeDnsDomainAttachment(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
	request.QueryParams["Department"] = s.client.Department
	request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

	request.InstanceId = id

	raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DnsDomainAttachment", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}

func (s *DnsService) WaitForAlidnsDomainAttachment(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDnsDomainAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, ProviderERROR)
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
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = s.client.Department
		request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.Headers = map[string]string{"RegionId": s.client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": s.client.SecretKey, "Product": "alidns"}
		request.QueryParams["Department"] = s.client.Department
		request.QueryParams["ResourceGroup"] = s.client.ResourceGroup

		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
	}
	return nil
}

func (s *DnsService) DescribeDnsDomain(id string) (response *DnsDomains, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()
	request.Method = "POST"          // Set request method
	request.Product = "CloudDns"     // Specify product
	request.Domain = s.client.Domain // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2022-06-24"   // Specify product version
	request.PageNumber = requests.NewInteger(2)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DescribeGlobalZones"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"AccessKeyId":     s.client.AccessKey,
		"Product":         "CloudDns",
		"RegionId":        s.client.RegionId,
		"Action":          "DescribeGlobalZones",
		"Version":         "2022-06-24",
		//"Id":              did[1],
		"Name":              did[0],
		"Forwardedregionid": s.client.RegionId,
		"SignatureVersion":  "2.1",
		"PageNumber":        fmt.Sprint(1),
		"PageSize":          fmt.Sprint(PageSizeLarge),
	}
	resp := &DnsDomains{}
	raw, err := s.client.WithEcsClient(func(cmsClient *ecs.Client) (interface{}, error) {
		return cmsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorDomainNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "DescribeGlobalZones", AlibabacloudStackSdkGoERROR)

	}
	addDebug("DescribeGlobalZones", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.AsapiSuccess == true {
		return resp, WrapError(err)
	}

	return resp, nil
}
