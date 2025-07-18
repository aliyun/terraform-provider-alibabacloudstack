package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type RdsDescribedbproxyResponse struct {
	DBProxyConnectStringItems struct {
		DBProxyConnectStringItems []struct {
			DBProxyEndpointId               int    `json:"DBProxyEndpointId"`
			DBProxyConnectString            string `json:"DBProxyConnectString"`
			DBProxyConnectStringPort        string `json:"DBProxyConnectStringPort"`
			DBProxyConnectStringNetType     string `json:"DBProxyConnectStringNetType"`
			DBProxyVpcInstanceId            string `json:"DBProxyVpcInstanceId"`
			DBProxyEndpointName             string `json:"DBProxyEndpointName"`
			DBProxyConnectStringNetWorkType int    `json:"DBProxyConnectStringNetWorkType"`
		} `json:"DBProxyConnectStringItems"`
	} `json:"DBProxyConnectStringItems"`

	DBProxyEndpointItems struct {
		DBProxyEndpointItems []struct {
			DBProxyEndpointName    string `json:"DBProxyEndpointName"`
			DBProxyEndpointType    string `json:"DBProxyEndpointType"`
			DBProxyEndpointAliases string `json:"DBProxyEndpointAliases"`
			DBProxyReadWriteMode   string `json:"DBProxyReadWriteMode"`
		} `json:"DBProxyEndpointItems"`
	} `json:"DBProxyEndpointItems"`
	RequestId                          string `json:"RequestId"`
	DBProxyServiceStatus               string `json:"DBProxyServiceStatus"`
	DBProxyInstanceType                string `json:"DBProxyInstanceType"`
	DBProxyInstanceNum                 int    `json:"DBProxyInstanceNum"`
	DBProxyInstanceStatus              string `json:"DBProxyInstanceStatus"`
	DBProxyInstanceCurrentMinorVersion string `json:"DBProxyInstanceCurrentMinorVersion"`
	DBProxyInstanceLatestMinorVersion  string `json:"DBProxyInstanceLatestMinorVersion"`
}

func (s *RdsService) DoRdsDescribedbproxyRequest(id string) (*RdsDescribedbproxyResponse, error) {
	// api: Rds - 2014-08-15 - DescribeDBProxy
	request := s.client.NewCommonRequest("POST", "Rds", "2014-08-15", "DescribeDBProxy", "")
	RdsDescribedbproxyResponseObj := &RdsDescribedbproxyResponse{}

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBProxy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &RdsDescribedbproxyResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBProxy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return RdsDescribedbproxyResponseObj, nil
}

func (s *RdsService) RdsProxyStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoRdsDescribedbproxyRequest(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBProxyInstanceStatus == failState {
				return object, object.DBProxyInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.DBProxyInstanceStatus))
			}
		}
		return object, object.DBProxyInstanceStatus, nil
	}
}
