package alibabacloudstack

import (
	"encoding/json"

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
