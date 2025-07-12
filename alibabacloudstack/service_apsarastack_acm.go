package alibabacloudstack

import (
	"encoding/json"
	"strings"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
)

type AcmService struct {
	client *connectivity.AlibabacloudStackClient
}

type AcmDescribeconfigurationResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`

	Configuration struct {
		AppName   string `json:"AppName"`
		Content   string `json:"Content"`
		DataId    string `json:"DataId"`
		Desc      string `json:"Desc"`
		Group     string `json:"Group"`
		Md5       string `json:"Md5"`
		Tags      string `json:"Tags"`
		Type      string `json:"Type"`
		UdVersion string `json:"UdVersion"`
	} `json:"Configuration"`
}

func (s *AcmService) DoAcmDescribeconfigurationRequest(id string) (*AcmDescribeconfigurationResponse, error) {
	// api: acm - 2020-02-06 - DescribeConfiguration
	request := s.client.NewCommonRequest("GET", "acm", "2020-02-06", "DescribeConfiguration", "/diamond-ops/pop/configuration")
	AcmDescribeconfigurationResponseObj := &AcmDescribeconfigurationResponse{}

	params := strings.Split(id, ":")

	request.QueryParams["DataId"] = params[0]

	request.QueryParams["Group"] = params[1]

	request.QueryParams["NamespaceId"] = params[2]

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if sdkErr, ok := err.(*errors.ServerError); ok && sdkErr.ErrorCode() == "ConfigurationNotExists" {
			return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Configuration %s Not Exists", id))
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeConfiguration", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &AcmDescribeconfigurationResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeConfiguration", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return AcmDescribeconfigurationResponseObj, nil
}

type AcmDescribeconfigurationsResponse struct {
	Configurations struct {
		Configuration []struct {
			AppName          string `json:"AppName"`
			DataId           string `json:"DataId"`
			EncryptedDataKey string `json:"EncryptedDataKey"`
			Group            string `json:"Group"`
			Md5              string `json:"Md5"`
			NamespaceId      string `json:"NamespaceId"`
			UdVersion        string `json:"UdVersion"`
		} `json:"Configuration"`
	} `json:"Configurations"`
	Code       string `json:"Code"`
	Message    string `json:"Message"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageSize   int    `json:"PageSize"`
	PageNumber int    `json:"PageNumber"`
}

func (s *AcmService) DoAcmDescribeconfigurationsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*AcmDescribeconfigurationsResponse, error) {
	// api: acm - 2020-02-06 - DescribeConfigurations
	request := s.client.NewCommonRequest("GET", "acm", "2020-02-06", "DescribeConfigurations", "")
	AcmDescribeconfigurationsResponseObj := &AcmDescribeconfigurationsResponse{}

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeConfigurations", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &AcmDescribeconfigurationsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeConfigurations", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return AcmDescribeconfigurationsResponseObj, nil
}

func (s *AcmService) DoAcmDescribeBetaConfigurationRequest(id string) (*AcmDescribeconfigurationResponse, error) {
	// api: acm - 2020-02-06 - DescribeBetaConfiguration
	request := s.client.NewCommonRequest("GET", "acm", "2020-02-06", "DescribeBetaConfiguration", "/diamond-ops/pop/configuration/beta")
	AcmDescribeconfigurationResponseObj := &AcmDescribeconfigurationResponse{}

	params := strings.Split(id, ":")

	request.QueryParams["DataId"] = params[0]

	request.QueryParams["Group"] = params[1]

	request.QueryParams["NamespaceId"] = params[2]

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBetaConfiguration", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &AcmDescribeconfigurationResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBetaConfiguration", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return AcmDescribeconfigurationResponseObj, nil
}
