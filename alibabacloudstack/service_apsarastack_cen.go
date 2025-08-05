package alibabacloudstack

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

const ChildInstanceTypeVpc = "VPC"
const ChildInstanceTypeVbr = "VBR"
const ChildInstanceTypeCcn = "CCN"

type CenService struct {
	client *connectivity.AlibabacloudStackClient
}

type CbnDescribecensResponse struct {
	Cens struct {
		Cen []struct {
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
		} `json:"Cen"`
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

func (s *CenService) DoCbnDescribecensRequest(id string) (*CbnDescribecensResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "DescribeCens", "")
	CbnDescribecensResponseObj := &CbnDescribecensResponse{}
	//调用request_params_handler
	request.QueryParams["Filter.1.Key"] = "CenId"
	request.QueryParams["Filter.1.Value.1"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeCens", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &CbnDescribecensResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeCens", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return CbnDescribecensResponseObj, nil
}

func (s *CenService) DoCbnDescribeTransitRoutersRequest(id string) (*TransitRouterResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouters", "")
	TransitRouterResponseObj := &TransitRouterResponse{}
	//调用request_params_handler
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

func (s *CenService) DoCbnDescribeTransitRouterRouteEntriesRequest(id string) (*CbnDescribeTransitRouterRouteEntriesResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("GET", "Cbn", "2017-09-12", "ListTransitRouterRouteEntries", "")
	CbnDescribeRouterRouteEntriesResponseObj := &CbnDescribeTransitRouterRouteEntriesResponse{}
	//调用request_params_handler
	request.QueryParams["TransitRouterRouteTableId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
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
	//调用request_params_handler
	parts := strings.Split(id, ":")
	transit_router_id := parts[0]
	request.QueryParams["TransitRouterId"] = transit_router_id

	bresponse, err := s.client.ProcessCommonRequest(request)
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
