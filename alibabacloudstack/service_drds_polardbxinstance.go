package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolardbxDescribedbinstanceattributeResponse struct {
	RequestId string `json:"RequestId"`

	DBInstance struct {
		ReadDBInstances struct {
			ReadDBInstance []string `json:"ReadDBInstance"`
		} `json:"ReadDBInstances"`

		LTSVersions struct {
			LTSVersion []string `json:"LTSVersion"`
		} `json:"LTSVersions"`

		DBNodes struct {
			DBNode []struct {
				Id            string `json:"Id"`
				NodeClass     string `json:"NodeClass"`
				RegionId      string `json:"RegionId"`
				ZoneId        string `json:"ZoneId"`
				ComputeNodeId string `json:"ComputeNodeId"`
				DataNodeId    string `json:"DataNodeId"`
			} `json:"DBNode"`
		} `json:"DBNodes"`

		ConnAddrs struct {
			ConnAddr []struct {
				ConnectionString string `json:"ConnectionString"`
				Port             string `json:"Port"`
				Type             string `json:"Type"`
				VPCId            string `json:"VPCId"`
				VSwitchId        string `json:"VSwitchId"`
				VpcInstanceId    string `json:"VpcInstanceId"`
			} `json:"ConnAddr"`
		} `json:"ConnAddrs"`

		TagSet struct {
			TagSet []struct {
				Key   string `json:"Key"`
				Value string `json:"Value"`
			} `json:"TagSet"`
		} `json:"TagSet"`
		Status                  string `json:"Status"`
		Description             string `json:"Description"`
		ZoneId                  string `json:"ZoneId"`
		VPCId                   string `json:"VPCId"`
		CreateTime              string `json:"CreateTime"`
		Expired                 string `json:"Expired"`
		PayType                 string `json:"PayType"`
		DBType                  string `json:"DBType"`
		LockMode                string `json:"LockMode"`
		StorageUsed             int    `json:"StorageUsed"`
		Storage                 int    `json:"Storage"`
		DBVersion               string `json:"DBVersion"`
		Network                 string `json:"Network"`
		RegionId                string `json:"RegionId"`
		Engine                  string `json:"Engine"`
		Id                      string `json:"Id"`
		ConnectionString        string `json:"ConnectionString"`
		Port                    string `json:"Port"`
		MinorVersion            string `json:"MinorVersion"`
		LatestMinorVersion      string `json:"LatestMinorVersion"`
		DBNodeCount             int    `json:"DBNodeCount"`
		DBInstanceType          string `json:"DBInstanceType"`
		SpecSeries              string `json:"SpecSeries"`
		DNNodeCount             int    `json:"DNNodeCount"`
		DNNodeClass             string `json:"DNNodeClass"`
		CNNodeCount             int    `json:"CNNodeCount"`
		CNNodeClass             string `json:"CNNodeClass"`
		MaintainStartTime       string `json:"MaintainStartTime"`
		MaintainEndTime         string `json:"MaintainEndTime"`
		VSwitchId               string `json:"VSwitchId"`
		CommodityCode           string `json:"CommodityCode"`
		ExpireDate              string `json:"ExpireDate"`
		Type                    string `json:"Type"`
		DBNodeClass             string `json:"DBNodeClass"`
		RightsSeparationStatus  bool   `json:"RightsSeparationStatus"`
		RightsSeparationEnabled bool   `json:"RightsSeparationEnabled"`
		KindCode                int    `json:"KindCode"`
		ResourceGroupId         string `json:"ResourceGroupId"`
		Series                  string `json:"Series"`
		CpuType                 string `json:"CpuType"`
	} `json:"DBInstance"`
}

func (s *DrdsService) DoPolardbxDescribedbinstanceattributeRequest(id string) (*PolardbxDescribedbinstanceattributeResponse, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstanceAttribute
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstanceAttribute", "")
	PolardbxDescribedbinstanceattributeResponseObj := &PolardbxDescribedbinstanceattributeResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceName"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbxDescribedbinstanceattributeResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbxDescribedbinstanceattributeResponseObj, nil
}

type PolardbxDescribedbinstancesResponse struct {
	DBInstances struct {
		DBInstance []struct {
			ReadDBInstances struct {
				ReadDBInstance []string `json:"ReadDBInstance"`
			} `json:"ReadDBInstances"`

			Nodes struct {
				PolarDBXNode []struct {
					Id        string `json:"Id"`
					ClassCode string `json:"ClassCode"`
					RegionId  string `json:"RegionId"`
					ZoneId    string `json:"ZoneId"`
				} `json:"PolarDBXNode"`
			} `json:"Nodes"`

			TagSet struct {
				TagSet []struct {
					Key   string `json:"Key"`
					Value string `json:"Value"`
				} `json:"TagSet"`
			} `json:"TagSet"`
			Id              string `json:"Id"`
			Description     string `json:"Description"`
			PayType         string `json:"PayType"`
			CreateTime      string `json:"CreateTime"`
			ExpireTime      string `json:"ExpireTime"`
			Expired         bool   `json:"Expired"`
			RegionId        string `json:"RegionId"`
			ZoneId          string `json:"ZoneId"`
			Network         string `json:"Network"`
			VPCId           string `json:"VPCId"`
			Engine          string `json:"Engine"`
			DBType          string `json:"DBType"`
			DBVersion       string `json:"DBVersion"`
			Status          string `json:"Status"`
			LockMode        string `json:"LockMode"`
			LockReason      string `json:"LockReason"`
			NodeCount       int    `json:"NodeCount"`
			NodeClass       string `json:"NodeClass"`
			SpecSeries      string `json:"SpecSeries"`
			DNNodeCount     int    `json:"DNNodeCount"`
			DNNodeClass     string `json:"DNNodeClass"`
			CNNodeCount     int    `json:"CNNodeCount"`
			CNNodeClass     string `json:"CNNodeClass"`
			StorageUsed     int    `json:"StorageUsed"`
			Storage         int    `json:"Storage"`
			CommodityCode   string `json:"CommodityCode"`
			Type            string `json:"Type"`
			MinorVersion    string `json:"MinorVersion"`
			ResourceGroupId string `json:"ResourceGroupId"`
			DBInstanceName  string `json:"DBInstanceName"`
			Series          string `json:"Series"`
			CpuType         string `json:"CpuType"`
		} `json:"DBInstance"`
	} `json:"DBInstances"`
	RequestId   string `json:"RequestId"`
	PageNumber  int    `json:"PageNumber"`
	PageSize    int    `json:"PageSize"`
	TotalNumber int    `json:"TotalNumber"`
}

func (s *DrdsService) DoPolardbxDescribedbinstancesRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbxDescribedbinstancesResponse, error) {
	// api: polardbx - 2020-02-02 - DescribeDBInstances
	request := s.client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DescribeDBInstances", "")
	PolardbxDescribedbinstancesResponseObj := &PolardbxDescribedbinstancesResponse{}

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbxDescribedbinstancesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbxDescribedbinstancesResponseObj, nil
}
