package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type VpngatewayService struct {
	client *connectivity.AlibabacloudStackClient
}

type VpnGatewayVpnPbrRouteEntry struct {
	VpnInstanceId string `json:"VpnInstanceId"`
	RouteSource   string `json:"RouteSource"`
	RouteDest     string `json:"RouteDest"`
	NextHop       string `json:"NextHop"`
	Weight        int    `json:"Weight"`
	CreateTime    int    `json:"CreateTime"`
	State         string `json:"State"`
}

type VpcDescribevpnpbrrouteentriesResponse struct {
	VpnPbrRouteEntries struct {
		VpnPbrRouteEntry []VpnGatewayVpnPbrRouteEntry `json:"VpnPbrRouteEntry"`
	} `json:"VpnPbrRouteEntries"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *VpngatewayService) DoVpcDescribevpnpbrrouteentriesRequest(id string) (*VpnGatewayVpnPbrRouteEntry, error) {
	// api: Vpc - 2016-04-28 - DescribeVpnPbrRouteEntries
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeVpnPbrRouteEntries", "")
	VpcDescribevpnpbrrouteentriesResponseObj := &VpcDescribevpnpbrrouteentriesResponse{}
	result := &VpnGatewayVpnPbrRouteEntry{}
	param := strings.Split(id, "_")
	request.QueryParams["VpnGatewayId"] = param[3]
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeVpnPbrRouteEntries", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribevpnpbrrouteentriesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeVpnPbrRouteEntries", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	for _, data := range VpcDescribevpnpbrrouteentriesResponseObj.VpnPbrRouteEntries.VpnPbrRouteEntry {
		if data.NextHop == param[0] && data.RouteDest == param[1] && data.RouteSource == param[2] && data.VpnInstanceId == param[3] && fmt.Sprint(data.Weight) == param[4] {
			result = &data
		}
	}
	if result == nil {
		return nil, errmsgs.Error(fmt.Sprintf(errmsgs.NotFoundMsg, "VpnPbrRouteEntry"))
	}

	return result, nil
}
