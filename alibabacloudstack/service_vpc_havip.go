package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type VpcDescribehavipsResponse struct {
	HaVips struct {
		HaVip []struct {
			AssociatedInstances struct {
				AssociatedInstance []string `json:"associatedInstance"`
			} `json:"AssociatedInstances"`

			AssociatedEipAddresses struct {
				AssociatedEipAddress []string `json:"associatedEipAddresse"`
			} `json:"AssociatedEipAddresses"`
			HaVipId                string `json:"HaVipId"`
			RegionId               string `json:"RegionId"`
			VpcId                  string `json:"VpcId"`
			VSwitchId              string `json:"VSwitchId"`
			IpAddress              string `json:"IpAddress"`
			Status                 string `json:"Status"`
			MasterInstanceId       string `json:"MasterInstanceId"`
			Description            string `json:"Description"`
			Name                   string `json:"Name"`
			ChargeType             string `json:"ChargeType"`
			CreateTime             string `json:"CreateTime"`
			AssociatedInstanceType string `json:"AssociatedInstanceType"`
		} `json:"HaVip"`
	} `json:"HaVips"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
}

func (s *VpcService) DoVpcDescribehavipsRequest(id string) (*VpcDescribehavipsResponse, error) {
	// api: Vpc - 2016-04-28 - DescribeHaVips
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DescribeHaVips", "")
	VpcDescribehavipsResponseObj := &VpcDescribehavipsResponse{}
	request.QueryParams["Filter.1.Key"] = "HaVipId"
	request.QueryParams["Filter.1.Value.1"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeHaVips", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcDescribehavipsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeHaVips", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcDescribehavipsResponseObj, nil
}

func (s *VpcService) VpcHaVipStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := s.DoVpcDescribehavipsRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		havip := response.HaVips.HaVip[0]
		for _, failState := range failStates {
			if havip.Status == failState {
				return havip, havip.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, havip.Status))
			}
		}
		return havip, havip.Status, nil
	}
}

func (s *VpcService) AssociateHaVip(id, instance_type, instance_id string) error {
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "AssociateHaVip", "")
	request.QueryParams["HaVipId"] = id
	request.QueryParams["InstanceId"] = instance_id
	request.QueryParams["InstanceType"] = instance_type
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_vpc_ha_vip", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	stateConf := BuildStateConf([]string{"Associating"}, []string{"Available", "InUse"}, 5*time.Minute, 10*time.Second, s.VpcHaVipStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
	}
	return nil
}

func (s *VpcService) UnassociateHaVip(id, instance_type, instance_id string) error {
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "UnassociateHaVip", "")
	request.QueryParams["HaVipId"] = id
	request.QueryParams["InstanceId"] = instance_id
	request.QueryParams["InstanceType"] = instance_type
	request.QueryParams["Force"] = "true"
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_vpc_ha_vip", "UnassociateHaVip", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	stateConf := BuildStateConf([]string{"Associating", "Unassociating"}, []string{"Available", "InUse"}, 5*time.Minute, 10*time.Second, s.VpcHaVipStateRefreshFunc(id, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
	}
	return nil
}
