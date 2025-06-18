package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type VpcGetdhcpoptionssetResponse struct {
	AssociateVpcs struct {
		AssociateVpc []struct {
			VpcId           string `json:"VpcId"`
			AssociateStatus string `json:"AssociateStatus"`
		} `json:"AssociateVpc"`
	} `json:"AssociateVpcs"`
	RequestId                 string `json:"RequestId"`
	DhcpOptionsSetName        string `json:"DhcpOptionsSetName"`
	DhcpOptionsSetDescription string `json:"DhcpOptionsSetDescription"`
	DhcpOptionsSetId          string `json:"DhcpOptionsSetId"`
	OwnerId                   int    `json:"OwnerId"`
	Status                    string `json:"Status"`

	DhcpOptions struct {
		DomainNameServers string `json:"DomainNameServers"`
		DomainName        string `json:"DomainName"`
	} `json:"DhcpOptions"`
}

func (s *VpcService) DoVpcGetdhcpoptionssetRequest(id string) (*VpcGetdhcpoptionssetResponse, error) {
	// api: Vpc - 2016-04-28 - GetDhcpOptionsSet
	request := s.client.NewCommonRequest("POST", "Vpc", "2016-04-28", "GetDhcpOptionsSet", "")
	VpcGetdhcpoptionssetResponseObj := &VpcGetdhcpoptionssetResponse{}

	request.QueryParams["DhcpOptionsSetId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "GetDhcpOptionsSet", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcGetdhcpoptionssetResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "GetDhcpOptionsSet", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return VpcGetdhcpoptionssetResponseObj, nil
}
