package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ExpressconnectService struct {
	client *connectivity.AlibabacloudStackClient
}

type ExpressconnectBgpGroup struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	BgpGroupId  string `json:"BgpGroupId"`
	PeerAsn     int    `json:"PeerAsn"`
	AuthKey     string `json:"AuthKey"`
	RouterId    string `json:"RouterId"`
	Status      string `json:"Status"`
	Keepalive   int    `json:"Keepalive"`
	LocalAsn    int    `json:"LocalAsn"`
	Hold        int    `json:"Hold"`
	IsFake      bool   `json:"IsFake"`
	RouteLimit  int    `json:"RouteLimit"`
	RegionId    string `json:"RegionId"`
	IpVersion   string `json:"IpVersion"`
}

type VpcDescribebgpgroupsResponse struct {
	BgpGroups struct {
		BgpGroup []ExpressconnectBgpGroup `json:"BgpGroup"`
	} `json:"BgpGroups"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *ExpressconnectService) DoVpcDescribebgpgroupsRequest(id string) (*ExpressconnectBgpGroup, error) {
	// api: Vpc - 2016-04-28 - DescribeBgpGroups
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeBgpGroups", "")
	VpcDescribebgpgroupsResponseObj := &VpcDescribebgpgroupsResponse{}
	request.QueryParams["BgpGroupId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBgpGroups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribebgpgroupsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBgpGroups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(VpcDescribebgpgroupsResponseObj.BgpGroups.BgpGroup) > 0 {
		return &VpcDescribebgpgroupsResponseObj.BgpGroups.BgpGroup[0], nil
	} else {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("BgpGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
}

func (s *ExpressconnectService) ExpressconnectBgpGroupsStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoVpcDescribebgpgroupsRequest(id)
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

type ExpressconnectBgpPeer struct {
	Name          string `json:"Name"`
	Description   string `json:"Description"`
	BgpPeerId     string `json:"BgpPeerId"`
	BgpGroupId    string `json:"BgpGroupId"`
	PeerIpAddress string `json:"PeerIpAddress"`
	PeerAsn       int    `json:"PeerAsn"`
	AuthKey       string `json:"AuthKey"`
	RouterId      string `json:"RouterId"`
	BgpStatus     string `json:"BgpStatus"`
	Status        string `json:"Status"`
	Keepalive     int    `json:"Keepalive"`
	LocalAsn      int    `json:"LocalAsn"`
	Hold          int    `json:"Hold"`
	IsFake        bool   `json:"IsFake"`
	RouteLimit    int    `json:"RouteLimit"`
	RegionId      string `json:"RegionId"`
	EnableBfd     bool   `json:"EnableBfd"`
	IpVersion     string `json:"IpVersion"`
	BfdMultiHop   int    `json:"BfdMultiHop"`
}

type VpcDescribebgppeersResponse struct {
	BgpPeers struct {
		BgpPeer []ExpressconnectBgpPeer `json:"BgpPeer"`
	} `json:"BgpPeers"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *ExpressconnectService) DoVpcDescribebgppeersRequest(id string) (*ExpressconnectBgpPeer, error) {
	// api: Vpc - 2016-04-28 - DescribeBgpPeers
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeBgpPeers", "")
	VpcDescribebgppeersResponseObj := &VpcDescribebgppeersResponse{}
	request.QueryParams["BgpPeerId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBgpPeers", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribebgppeersResponseObj)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBgpPeers", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(VpcDescribebgppeersResponseObj.BgpPeers.BgpPeer) > 0 {
		return &VpcDescribebgppeersResponseObj.BgpPeers.BgpPeer[0], nil
	} else {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("BgpPeer", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
}

func (s *ExpressconnectService) DoVpcDescribeVbrHaRequest(id string) (map[string]interface{}, error) {
	
	reqQuery := map[string]interface{}{
		"VbrHaId" : id,
	}
	// api: Vpc - 2016-04-28 - DescribeBgpPeers
	response, err := s.client.DoTeaRequest("GET", "Vpc", "2016-04-28", "DescribeVbrHa", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	
	if _, exist := response["VbrHaId"]; !exist {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Vbr Ha %s Not Found", id))
	}
	
	return response, nil
}

func (s *ExpressconnectService) ExpressconnectBgpPeersStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoVpcDescribebgppeersRequest(id)
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


func (s *ExpressconnectService) ExpressconnectVbrHaStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoVpcDescribeVbrHaRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		var objectStatus string
		if v, exist := object["Status"]; !exist {
			return nil, "", nil
		} else {
			objectStatus = v.(string)
		}
		for _, failState := range failStates {
			if objectStatus == failState {
				return object, objectStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, objectStatus))
			}
		}

		return object, objectStatus, nil
	}
}



type VpcDescribebgpnetworksResponse struct {
	BgpNetworks struct {
		BgpNetwork []struct {
			VpcId        string `json:"VpcId"`
			DstCidrBlock string `json:"DstCidrBlock"`
			RouterId     string `json:"RouterId"`
			Status       string `json:"Status"`
		} `json:"BgpNetwork"`
	} `json:"BgpNetworks"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *ExpressconnectService) DoVpcDescribebgpnetworksRequest(id string) (*VpcDescribebgpnetworksResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeBgpNetworks
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeBgpNetworks", "")
	VpcDescribebgpnetworksResponseObj := &VpcDescribebgpnetworksResponse{}

	//调用request_params_handler

	request.QueryParams["RouterId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBgpNetworks", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribebgpnetworksResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBgpNetworks", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribebgpnetworksResponseObj, nil
}
