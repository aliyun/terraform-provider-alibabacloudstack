package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const ChildInstanceTypeVpc = "VPC"
const ChildInstanceTypeVbr = "VBR"
const ChildInstanceTypeCcn = "CCN"

type CenService struct {
	client *connectivity.AlibabacloudStackClient
}

type CenInstance struct {
	CenBandwidthPackageIds struct {
		CenBandwidthPackageId []string `json:"CenBandwidthPackageId"`
	} `json:"CenBandwidthPackageIds"`
	CenId                           string `json:"CenId"`
	CreationTime                    string `json:"CreationTime"`
	Ipv6Level                       string `json:"Ipv6Level"`
	ProtectionLevel                 string `json:"ProtectionLevel"`
	SecurityLevelTag                string `json:"SecurityLevelTag"`
	SecurityLevelTagBackgroundColer string `json:"SecurityLevelTagBackgroundColer"`
	SecurityLevelTagTextColer       string `json:"SecurityLevelTagTextColer"`
	Status                          string `json:"Status"`
	Description                     string `json:"Description"`
	Name                            string `json:"Name"`
}

type CbnDescribecensResponse struct {
	Cens struct {
		Cen []CenInstance `json:"Cen"`
	} `json:"Cens"`
	RequestId string `json:"RequestId"`
}

type TransitRouterResponse struct {
	TransitRouters []struct {
		TransitRouterCidrList []struct {
			Cidr                string `json:"Cidr"`
			PublishCidrRoute    bool   `json:"PublishCidrRoute"`
			TransitRouterCidrId string `json:"TransitRouterCidrId"`
		} `json:"TransitRouterCidrList"`
		CenId                    string `json:"CenId"`
		CreationTime             string `json:"CreationTime"`
		ServiceMode              string `json:"ServiceMode"`
		RegionId                 string `json:"RegionId"`
		SupportMulticast         bool   `json:"SupportMulticast"`
		TransitRouterDescription string `json:"TransitRouterDescription"`
		TransitRouterId          string `json:"TransitRouterId"`
		Status                   string `json:"Status"`
		TransitRouterName        string `json:"TransitRouterName"`
		Type                     string `json:"Type"`
	} `json:"TransitRouters"`
	RequestId string `json:"RequestId"`
}

type CbnDescribeTransitRouterRouteEntriesResponse struct {
	TransitRouterRouteEntries []struct {
		TransitRouterRouteEntryDestinationCidrBlock string `json:"TransitRouterRouteEntryDestinationCidrBlock"`
		TransitRouterRouteEntryNextHopId            string `json:"TransitRouterRouteEntryNextHopId"`
		TransitRouterRouteEntryType                 string `json:"TransitRouterRouteEntryType"`
		CreateTime                                  string `json:"CreateTime"`
		TransitRouterRouteEntryNextHopType          string `json:"TransitRouterRouteEntryNextHopType"`
		TransitRouterRouteEntryName                 string `json:"TransitRouterRouteEntryName"`
		OperationalMode                             bool   `json:"OperationalMode,"`
		TransitRouterRouteEntryId                   string `json:"TransitRouterRouteEntryId"`
		TransitRouterRouteEntryStatus               string `json:"TransitRouterRouteEntryStatus"`
		TransitRouterRouteEntryDescription          string `json:"TransitRouterRouteEntryDescription"`
	} `json:"TransitRouterRouteEntries"`
	RequestId string `json:"RequestId"`
}

type CbnDescribeTransitRouterRouteTablesResponse struct {
	TransitRouterRouteTables []struct {
		TransitRouterRouteTableId          string `json:"TransitRouterRouteTableId"`
		TransitRouterRouteTableStatus      string `json:"TransitRouterRouteTableStatus"`
		TransitRouterRouteTableType        string `json:"TransitRouterRouteTableType"`
		TransitRouterRouteTableDescription string `json:"TransitRouterRouteTableDescription"`
		CreateTime                         string `json:"CreateTime"`
		TransitRouterRouteTableName        string `json:"TransitRouterRouteTableName"`
		Tags                               []Tag  `json:"Tags"`
	} `json:"TransitRouterRouteTables"`
	RequestId string `json:"RequestId"`
}

type CbnDescribeTransitRouterMulticastDomainsResponse struct {
	TransitRouterMulticastDomains []struct {
		Status                                  string        `json:"Status"`
		TransitRouterMulticastDomainId          string        `json:"TransitRouterMulticastDomainId"`
		TransitRouterMulticastDomainName        string        `json:"TransitRouterMulticastDomainName,omitempty"`
		TransitRouterMulticastDomainDescription string        `json:"TransitRouterMulticastDomainDescription,omitempty"`
		Tags                                    []interface{} `json:"Tags"`
		TransitRouterId                         string        `json:"TransitRouterId"`
	} `json:"TransitRouterMulticastDomains"`
	RequestId string `json:"RequestId"`
}

type CbnTransitRouterMulticastGroupData struct {
	TransitRouterAttachmentId      string `json:"TransitRouterAttachmentId"`
	Status                         string `json:"Status"`
	GroupMember                    bool   `json:"GroupMember"`
	ResourceId                     string `json:"ResourceId"`
	VSwitchId                      string `json:"VSwitchId"`
	SourceType                     string `json:"SourceType"`
	MemberType                     string `json:"MemberType"`
	TransitRouterMulticastDomainId string `json:"TransitRouterMulticastDomainId"`
	ResourceType                   string `json:"ResourceType"`
	NetworkInterfaceId             string `json:"NetworkInterfaceId"`
	GroupSource                    bool   `json:"GroupSource"`
	ResourceOwnerId                int64  `json:"ResourceOwnerId"`
	GroupIpAddress                 string `json:"GroupIpAddress"`
	ConnectPeerId                  string `json:"ConnectPeerId"`
}

type CbnDescribeTransitRouterMulticastDomainSourceResponse struct {
	TransitRouterMulticastGroups []CbnTransitRouterMulticastGroupData `json:"TransitRouterMulticastGroups"`
	RequestId                    string                               `json:"RequestId"`
}

type CbnDescribeTransitRouterMulticastDomainAssociationsResponse struct {
	TransitRouterMulticastAssociations []struct {
		TransitRouterAttachmentId      string `json:"TransitRouterAttachmentId"`
		Status                         string `json:"Status"`
		ResourceId                     string `json:"ResourceId"`
		TransitRouterMulticastDomainId string `json:"TransitRouterMulticastDomainId"`
		VSwitchId                      string `json:"VSwitchId"`
		ResourceType                   string `json:"ResourceType"`
		ResourceOwnerId                int64  `json:"ResourceOwnerId"`
	} `json:"TransitRouterMulticastAssociations"`
	RequestId string `json:"RequestId"`
}

type ZoneMapping struct {
	ZoneId             string `json:"ZoneId"`
	VSwitchId          string `json:"VSwitchId"`
	NetworkInterfaceId string `json:"NetworkInterfaceId"`
}

type CbnDescribeTransitRouterVpcAttachmentsResponse struct {
	TransitRouterAttachments []struct {
		Status                             string        `json:"Status"`
		TransitRouterAttachmentId          string        `json:"TransitRouterAttachmentId"`
		TransitRouterAttachmentName        string        `json:"TransitRouterAttachmentName"`
		ResourceType                       string        `json:"ResourceType"`
		ZoneMappings                       []ZoneMapping `json:"ZoneMappings"`
		VpcOwnerId                         int64         `json:"VpcOwnerId"`
		AutoPublishRouteEnabled            bool          `json:"AutoPublishRouteEnabled"`
		VpcId                              string        `json:"VpcId"`
		ChargeType                         string        `json:"ChargeType"`
		CreationTime                       string        `json:"CreationTime"`
		VpcRegionId                        string        `json:"VpcRegionId"`
		TransitRouterAttachmentDescription string        `json:"TransitRouterAttachmentDescription"`
		Tags                               []interface{} `json:"Tags"`
		TransitRouterId                    string        `json:"TransitRouterId"`
	} `json:"TransitRouterAttachments"`
	RequestId string `json:"RequestId"`
}

type CbnDescribeTransitRouterVbrAttachmentsResponse struct {
	TransitRouterAttachments []struct {
		Status                             string        `json:"Status"`
		TransitRouterAttachmentId          string        `json:"TransitRouterAttachmentId"`
		AutoPublishRouteEnabled            bool          `json:"AutoPublishRouteEnabled"`
		VbrOwnerId                         int64         `json:"VbrOwnerId"`
		CreationTime                       string        `json:"CreationTime"`
		ResourceType                       string        `json:"ResourceType"`
		TransitRouterAttachmentName        string        `json:"TransitRouterAttachmentName"`
		TransitRouterAttachmentDescription string        `json:"TransitRouterAttachmentDescription"`
		VbrRegionId                        string        `json:"VbrRegionId"`
		VbrId                              string        `json:"VbrId"`
		Tags                               []interface{} `json:"Tags"`
		TransitRouterId                    string        `json:"TransitRouterId"`
	} `json:"TransitRouterAttachments"`
	RequestId string `json:"RequestId"`
}

type TransitRouterRouteTableAssociationsResponse struct {
	TransitRouterAssociations []struct {
		TransitRouterAttachmentId string `json:"TransitRouterAttachmentId"`
		Status                    string `json:"Status"`
		TransitRouterRouteTableId string `json:"TransitRouterRouteTableId"`
		ResourceId                string `json:"ResourceId"`
		ResourceType              string `json:"ResourceType"`
	} `json:"TransitRouterAssociations"`
	RequestId string `json:"RequestId"`
}

type TransitRouterRouteTablePropagationsResponse struct {
	TransitRouterPropagations []struct {
		Status                    string `json:"Status"`
		TransitRouterAttachmentId string `json:"TransitRouterAttachmentId"`
		TransitRouterRouteTableId string `json:"TransitRouterRouteTableId"`
		ResourceId                string `json:"ResourceId"`
		ResourceType              string `json:"ResourceType"`
	} `json:"TransitRouterPropagations"`
	RequestId string `json:"RequestId"`
}

type TransitRouterAttachmentsResponse struct {
	TransitRouterAttachments []struct {
		Status                             string      `json:"Status"`
		TransitRouterAttachmentId          string      `json:"TransitRouterAttachmentId"`
		ResourceRegionId                   string      `json:"ResourceRegionId"`
		Association                        interface{} `json:"Association"`
		ResourceId                         string      `json:"ResourceId"`
		CreationTime                       string      `json:"CreationTime"`
		TransitRouterAttachmentName        string      `json:"TransitRouterAttachmentName"`
		ResourceType                       string      `json:"ResourceType"`
		ResourceOwnerId                    int         `json:"ResourceOwnerId"`
		TransitRouterAttachmentDescription string      `json:"TransitRouterAttachmentDescription"`
	} `json:"TransitRouterAttachments"`
	RequestId string `json:"RequestId"`
}

type SourceRegionIds struct {
	SourceRegionId []string `json:"SourceRegionId"`
}

type DestinationChildInstanceTypes struct {
	DestinationChildInstanceType []string `json:"DestinationChildInstanceType"`
}
type SourceChildInstanceTypes struct {
	SourceChildInstanceType []string `json:"SourceChildInstanceType"`
}
type DestinationRouteTableIds struct {
	DestinationRouteTableId []string `json:"DestinationRouteTableId"`
}
type SourceInstanceIds struct {
	SourceInstanceId []string `json:"SourceInstanceId"`
}

type DestinationCidrBlocks struct {
	DestinationCidrBlock []string `json:"DestinationCidrBlock"`
}
type RouteTypes struct {
	RouteType []string `json:"RouteType"`
}

type MatchAsns struct {
	MatchAsn []int `json:"MatchAsn"`
}

type PrependAsPath struct {
	AsPath []int `json:"AsPath"`
}

type OperateCommunitySet struct {
	OperateCommunity []string `json:"OperateCommunity"`
}

type MatchCommunitySet struct {
	MatchCommunity []string `json:"MatchCommunity"`
}

type DestinationInstanceIds struct {
	DestinationInstanceId []string `json:"DestinationInstanceId"`
}

type SourceRouteTableIds struct {
	SourceRouteTableId []string `json:"SourceRouteTableId"`
}

type CbnDescribeCenRouteMapsResponse struct {
	RouteMaps struct {
		RouteMap []struct {
			Status                             string                        `json:"Status"`
			TransitRouterRouteTableId          string                        `json:"TransitRouterRouteTableId"`
			Preference                         int                           `json:"Preference,omitempty"`
			Priority                           int                           `json:"Priority"`
			TransmitDirection                  string                        `json:"TransmitDirection"`
			CenId                              string                        `json:"CenId"`
			NextPriority                       int                           `json:"NextPriority,omitempty"`
			CenRegionId                        string                        `json:"CenRegionId"`
			RouteMapId                         string                        `json:"RouteMapId"`
			Description                        string                        `json:"Description"`
			MapResult                          string                        `json:"MapResult"`
			CommunityOperateMode               string                        `json:"CommunityOperateMode"`
			MatchAsns                          MatchAsns                     `json:"MatchAsns,omitempty"`
			MatchAddressType                   string                        `json:"MatchAddressType"`
			AsPathMatchMode                    string                        `json:"AsPathMatchMode"`
			CidrMatchMode                      string                        `json:"CidrMatchMode"`
			CommunityMatchMode                 string                        `json:"CommunityMatchMode"`
			SourceInstanceIdsReverseMatch      bool                          `json:"SourceInstanceIdsReverseMatch"`
			DestinationInstanceIdsReverseMatch bool                          `json:"DestinationInstanceIdsReverseMatch"`
			RouteTypes                         RouteTypes                    `json:"RouteTypes,omitempty"`
			SourceRegionIds                    SourceRegionIds               `json:"SourceRegionIds,omitempty"`
			PrependAsPath                      PrependAsPath                 `json:"PrependAsPath,omitempty"`
			OperateCommunitySet                OperateCommunitySet           `json:"OperateCommunitySet,omitempty"`
			MatchCommunitySet                  MatchCommunitySet             `json:"MatchCommunitySet,omitempty"`
			DestinationInstanceIds             DestinationInstanceIds        `json:"DestinationInstanceIds,omitempty"`
			DestinationChildInstanceTypes      DestinationChildInstanceTypes `json:"DestinationChildInstanceTypes,omitempty"`
			DestinationRouteTableIds           DestinationRouteTableIds      `json:"DestinationRouteTableIds,omitempty"`
			SourceRouteTableIds                SourceRouteTableIds           `json:"SourceRouteTableIds,omitempty"`
			DestinationCidrBlocks              DestinationCidrBlocks         `json:"DestinationCidrBlocks,omitempty"`
			SourceChildInstanceTypes           SourceChildInstanceTypes      `json:"SourceChildInstanceTypes,omitempty"`
			SourceInstanceIds                  SourceInstanceIds             `json:"SourceInstanceIds,omitempty"`
		} `json:"RouteMap"`
	} `json:"RouteMaps"`
}

func (s *CenService) DoCbnDescribecensRequest(id string) (*CenInstance, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "DescribeCens", "")
	cbnDescribecensResponseObj := &CbnDescribecensResponse{}
	// Call request_params_handler
	request.QueryParams["Filter.1.Key"] = "CenId"
	request.QueryParams["Filter.1.Value.1"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeCens", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &cbnDescribecensResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeCens", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if len(cbnDescribecensResponseObj.Cens.Cen) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString("Not Found Cent Instance " + id)
	}

	return &cbnDescribecensResponseObj.Cens.Cen[0], nil
}

func (s *CenService) DoCbnDescribeTransitRoutersRequest(id string) (*TransitRouterResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouters", "")
	TransitRouterResponseObj := &TransitRouterResponse{}
	// Call request_params_handler
	request.QueryParams["CenId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &TransitRouterResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouters", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return TransitRouterResponseObj, nil
}

func (s *CenService) WaitForCenInstance(instanceId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DoCbnDescribecensRequest(instanceId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if instance.Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, instanceId, GetFunc(1), timeout, instance.Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) WaitForTransitRouterInstance(instanceId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DoCbnDescribeTransitRoutersRequest(instanceId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if instance.TransitRouters[0].Status == string(status) {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, instanceId, GetFunc(1), timeout, instance.TransitRouters[0].Status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) WaitForTransitRouterMulticastDomain(instanceId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(instanceId, COLON_SEPARATED)
	domain_id := parts[1]
	for {
		instance, err := s.DoCbnDescribeTransitRouterMuliticastDomainsRequest(instanceId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		for _, domain := range instance.TransitRouterMulticastDomains {
			if domain.TransitRouterMulticastDomainId == domain_id && domain.Status == string(status) {
				return nil

			}
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, instanceId, GetFunc(1), timeout, string(status), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) DoCbnDescribeTransitRouterRouteEntriesRequest(id string) (*CbnDescribeTransitRouterRouteEntriesResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts

	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterRouteEntries", "")
	CbnDescribeRouterRouteEntriesResponseObj := &CbnDescribeTransitRouterRouteEntriesResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	route_table_id := parts[0]
	request.QueryParams["TransitRouterRouteTableId"] = route_table_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterRouteEntries", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterRouteEntriesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterRouteEntries", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterRouteEntriesResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterRouteTablesRequest(id string) (*CbnDescribeTransitRouterRouteTablesResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterRouteTables", "")
	CbnDescribeRouterRouteTablesResponseObj := &CbnDescribeTransitRouterRouteTablesResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	transit_router_id := parts[0]
	request.QueryParams["TransitRouterId"] = transit_router_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterRouteTables", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterRouteTablesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterRouteTables", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterRouteTablesResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterMuliticastDomainsRequest(id string) (*CbnDescribeTransitRouterMulticastDomainsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterMulticastDomains", "")
	CbnDescribeRouterMulticastDomainResponseObj := &CbnDescribeTransitRouterMulticastDomainsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	transit_router_id := parts[0]
	request.QueryParams["TransitRouterId"] = transit_router_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterMulticastDomains", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterMulticastDomainResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterMulticastDomains", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterMulticastDomainResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterMuliticastDomainSourceRequest(id string) (*CbnTransitRouterMulticastGroupData, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterMulticastGroups", "")
	CbnDescribeRouterMulticastDomainSourceResponseObj := &CbnDescribeTransitRouterMulticastDomainSourceResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	group_ip_address := parts[0]
	transit_router_multicast_domain_id := parts[1]
	resource_type := parts[2]
	key := parts[3]
	request.QueryParams["IsGroupSource"] = "true"
	request.QueryParams["GroupIpAddress"] = group_ip_address
	request.QueryParams["TransitRouterMulticastDomainId"] = transit_router_multicast_domain_id
	if resource_type == "VPC" {
		request.QueryParams["NetworkInterfaceIds.1"] = key
	} else {
		request.QueryParams["ConnectPeerIds.1"] = key
	}

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterMulticastGroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterMulticastDomainSourceResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterMulticastDomains", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(CbnDescribeRouterMulticastDomainSourceResponseObj.TransitRouterMulticastGroups) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString("Not Found CentTransitRouterMuliticastDomainMember " + id)
	}
	return &CbnDescribeRouterMulticastDomainSourceResponseObj.TransitRouterMulticastGroups[0], nil
}

func (s *CenService) DoCbnDescribeTransitRouterMuliticastDomainMemberRequest(id string) (*CbnTransitRouterMulticastGroupData, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterMulticastGroups", "")
	CbnDescribeRouterMulticastDomainMemberResponseObj := &CbnDescribeTransitRouterMulticastDomainSourceResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	request.QueryParams["IsGroupMember"] = "true"
	request.QueryParams["GroupIpAddress"] = parts[0]
	request.QueryParams["VSwitchIds.1"] = parts[1]
	request.QueryParams["TransitRouterMulticastDomainId"] = parts[2]
	request.QueryParams["NetworkInterfaceIds.1"] = parts[3]

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterMulticastGroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterMulticastDomainMemberResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterMulticastDomains", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(CbnDescribeRouterMulticastDomainMemberResponseObj.TransitRouterMulticastGroups) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString("Not Found CentTransitRouterMuliticastDomainMember " + id)
	}
	return &CbnDescribeRouterMulticastDomainMemberResponseObj.TransitRouterMulticastGroups[0], nil
}

func (s *CenService) CbnTransitRouterMuliticastDomainMemberStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoCbnDescribeTransitRouterMuliticastDomainMemberRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CenService) CbnTransitRouterMuliticastDomainSourceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoCbnDescribeTransitRouterMuliticastDomainSourceRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Status == failState {
				return object, object.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Status))
			}
		}
		return object, object.Status, nil
	}
}

func (s *CenService) DoCbnDescribeTransitRouterMuliticastDomainAssociationsRequest(id string) (*CbnDescribeTransitRouterMulticastDomainAssociationsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterMulticastDomainAssociations", "")
	CbnDescribeRouterMulticastDomainAssociationResponseObj := &CbnDescribeTransitRouterMulticastDomainAssociationsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	domain_id := parts[1]
	request.QueryParams["TransitRouterMulticastDomainId"] = domain_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterMulticastDomainAssociations", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterMulticastDomainAssociationResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterMulticastDomainAssociations", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterMulticastDomainAssociationResponseObj, nil
}

func (s *CenService) WaitForTransitRouterMulticastDomainAssociation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(id, ":")
	attchment_id := parts[0]
	domain_id := parts[1]
	vswithch_id := parts[2]
	for {
		domain_associations, err := s.DoCbnDescribeTransitRouterMuliticastDomainAssociationsRequest(":" + domain_id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		delete := true
		for _, domain_association := range domain_associations.TransitRouterMulticastAssociations {
			if domain_association.TransitRouterAttachmentId == attchment_id && domain_association.VSwitchId == vswithch_id && domain_association.TransitRouterMulticastDomainId == domain_id {
				delete = false
				if domain_association.Status == string(status) {
					return nil
				}
			}
		}
		if status == Deleted && delete == true {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, string(status), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) WaitForTransitRouterTable(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(id, ":")
	table_id := parts[1]
	for {
		tables, err := s.DoCbnDescribeTransitRouterRouteTablesRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		table_status := ""
		for _, table := range tables.TransitRouterRouteTables {
			table_status = table.TransitRouterRouteTableStatus
			if table.TransitRouterRouteTableId == table_id && table_status == string(status) {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, table_id, GetFunc(1), timeout, table_status, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) DoCbnDescribeTransitRouterAttachmentsRequest(id string) (*TransitRouterAttachmentsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterAttachments", "")
	DescribeRouterattachmentsResponseObj := &TransitRouterAttachmentsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	cen_id := parts[0]
	transit_router_id := parts[1]

	request.QueryParams["TransitRouterId"] = transit_router_id
	request.QueryParams["CenId"] = cen_id
	request.QueryParams["ResourceTypes.1"] = "VPC"
	request.QueryParams["ResourceTypes.2"] = "VBR"
	request.QueryParams["ResourceTypes.3"] = "Connect"

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterAttachments", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DescribeRouterattachmentsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterAttachments", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DescribeRouterattachmentsResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterVbrAttachmentsRequest(id string) (*CbnDescribeTransitRouterVbrAttachmentsResponse, error) {

	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterVbrAttachments", "")
	CbnDescribeRouterVbrattachmentsResponseObj := &CbnDescribeTransitRouterVbrAttachmentsResponse{}

	parts := strings.Split(id, ":")
	cen_id := parts[0]
	transit_router_id := parts[1]
	attachment_id := parts[2]
	request.QueryParams["TransitRouterId"] = transit_router_id
	request.QueryParams["CenId"] = cen_id
	request.QueryParams["TransitRouterAttachmentId"] = attachment_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterVbrAttachments", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterVbrattachmentsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterVbrAttachments", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterVbrattachmentsResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterVpcAttachmentsRequest(id string) (*CbnDescribeTransitRouterVpcAttachmentsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterVpcAttachments", "")
	CbnDescribeRouterVpcattachmentsResponseObj := &CbnDescribeTransitRouterVpcAttachmentsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	cen_id := parts[0]
	transit_router_id := parts[1]
	attachment_id := parts[2]
	request.QueryParams["TransitRouterId"] = transit_router_id
	request.QueryParams["CenId"] = cen_id
	request.QueryParams["TransitRouterAttachmentId"] = attachment_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterVpcAttachments", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouterVpcattachmentsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterVpcAttachments", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouterVpcattachmentsResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterRouteTableAssociationsRequest(id string) (*TransitRouterRouteTableAssociationsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterRouteTableAssociations", "")
	TransitRouterRouteTableAssociationsResponseObj := &TransitRouterRouteTableAssociationsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	transit_router_table := parts[0]
	if len(parts) > 1 {
		attachment_id := parts[1]
		request.QueryParams["TransitRouterAttachmentId"] = attachment_id
	}
	request.QueryParams["TransitRouterRouteTableId"] = transit_router_table

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterRouteTableAssociations", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &TransitRouterRouteTableAssociationsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterRouteTableAssociations", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return TransitRouterRouteTableAssociationsResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRouterRouteTablePropagationsRequest(id string) (*TransitRouterRouteTablePropagationsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterRouteTablePropagations", "")
	TransitRouterRouteTablePropagationsResponseObj := &TransitRouterRouteTablePropagationsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	transit_router_table := parts[0]
	if len(parts) > 1 {
		attachment_id := parts[1]
		request.QueryParams["TransitRouterAttachmentId"] = attachment_id
	}
	request.QueryParams["TransitRouterRouteTableId"] = transit_router_table
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListTransitRouterRouteTablePropagations", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &TransitRouterRouteTablePropagationsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "ListTransitRouterRouteTablePropagations", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return TransitRouterRouteTablePropagationsResponseObj, nil
}

func (s *CenService) WaitForAttachmentInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(id, COLON_SEPARATED)
	cen_id := parts[0]
	transit_router_id := parts[1]
	attach_ment_id := parts[2]
	for {
		instance, err := s.DoCbnDescribeTransitRouterAttachmentsRequest(cen_id + ":" + transit_router_id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		for _, v := range instance.TransitRouterAttachments {
			if v.TransitRouterAttachmentId == attach_ment_id && v.Status == string(status) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, attach_ment_id, GetFunc(1), timeout, string(status), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) WaitForTransitRouterTableAssociation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(id, COLON_SEPARATED)
	transit_router_table_id := parts[0]
	attach_ment_id := parts[1]
	for {
		instance, err := s.DoCbnDescribeTransitRouterRouteTableAssociationsRequest(transit_router_table_id + ":" + attach_ment_id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if len(instance.TransitRouterAssociations) == 0 && status == Deleted {
			return nil
		}
		for _, v := range instance.TransitRouterAssociations {
			if v.TransitRouterAttachmentId == attach_ment_id && v.Status == string(status) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, attach_ment_id, GetFunc(1), timeout, string(status), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) WaitForTransitRouterTablePropagation(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts := strings.Split(id, COLON_SEPARATED)
	transit_router_table_id := parts[0]
	attach_ment_id := parts[1]
	for {
		instance, err := s.DoCbnDescribeTransitRouterRouteTablePropagationsRequest(transit_router_table_id + ":" + attach_ment_id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if len(instance.TransitRouterPropagations) == 0 && status == Deleted {
			return nil
		}
		for _, v := range instance.TransitRouterPropagations {
			if v.TransitRouterAttachmentId == attach_ment_id && v.Status == string(status) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, attach_ment_id, GetFunc(1), timeout, string(status), string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *CenService) DoCbnDescribeCenRouteMapsRequest(id string) (*CbnDescribeCenRouteMapsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "DescribeCenRouteMaps", "")
	CbnDescribeRouteMapsResponseObj := &CbnDescribeCenRouteMapsResponse{}
	// Call request_params_handler
	parts := strings.Split(id, ":")
	cen_id := parts[0]
	route_table_id := parts[2]
	request.QueryParams["CenId"] = cen_id
	request.QueryParams["TransitRouterRouteTableId"] = route_table_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeCenRouteMaps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribeRouteMapsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeCenRouteMaps", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribeRouteMapsResponseObj, nil
}

func (s *CenService) DescribeCenTransitRouterConnectAttachment(id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, expected {CenId}:{TransitRouterId}:{TransitRouterAttachmentId}")
	}
	cenId := parts[0]
	transitRouterId := parts[1]
	transitRouterAttachmentId := parts[2]

	reqQuery := map[string]interface{}{
		"CenId":                     cenId,
		"TransitRouterId":           transitRouterId,
		"TransitRouterAttachmentId": transitRouterAttachmentId,
	}

	response, err := s.client.DoTeaRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterConnectAttachments", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := response["TransitRouterAttachments"]; !ok || fmt.Sprint(response["TotalCount"]) == "0" {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Cen TransitRouterConnectAttachment %s not found", id))
	}

	attachments := response["TransitRouterAttachments"].([]interface{})
	for _, v := range attachments {
		attachment := v.(map[string]interface{})
		if attachment["TransitRouterAttachmentId"].(string) == transitRouterAttachmentId {
			return attachment, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Cen TransitRouterConnectAttachment %s not found", id))
}

func (s *CenService) CenTransitRouterConnectAttachmentStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCenTransitRouterConnectAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		status := object["Status"].(string)
		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *CenService) DescribeCenVbrHealthCheck(id string) (map[string]interface{}, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected CenId:VbrInstanceId, got %s", id)
	}
	cenId := parts[0]
	vbrInstanceId := parts[1]

	reqQuery := map[string]interface{}{
		"CenId":               cenId,
		"VbrInstanceId":       vbrInstanceId,
		"VbrInstanceRegionId": s.client.RegionId,
	}

	response, err := s.client.DoTeaRequest("GET", "Cbn", "2017-09-12", "DescribeCenVbrHealthCheck", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := response["VbrHealthChecks"]; !ok || response["VbrHealthChecks"] == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("CenVbrHealthCheck not found for CenId: %s, VbrInstanceId: %s", cenId, vbrInstanceId))
	}

	vbrHealthChecks := response["VbrHealthChecks"].(map[string]interface{})
	if vbrHealthCheckList, ok := vbrHealthChecks["VbrHealthCheck"].([]interface{}); ok && len(vbrHealthCheckList) > 0 {
		// Return the first matched health check item
		return vbrHealthCheckList[0].(map[string]interface{}), nil
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("CenVbrHealthCheck not found for CenId: %s, VbrInstanceId: %s", cenId, vbrInstanceId))
}
