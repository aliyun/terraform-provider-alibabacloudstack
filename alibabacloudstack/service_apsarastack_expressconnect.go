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
	PeerAsn     string `json:"PeerAsn"`
	AuthKey     string `json:"AuthKey"`
	RouterId    string `json:"RouterId"`
	Status      string `json:"Status"`
	Keepalive   int    `json:"Keepalive"`
	LocalAsn    string `json:"LocalAsn"`
	Hold        string `json:"Hold"`
	IsFake      string `json:"IsFake"`
	RouteLimit  string `json:"RouteLimit"`
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
		return nil, errmsgs.Error(fmt.Sprintf("NotFound ExpressConnect BgpGroup:%s", id))
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
