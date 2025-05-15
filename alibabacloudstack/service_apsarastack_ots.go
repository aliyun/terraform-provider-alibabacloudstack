package alibabacloudstack

import (
	"encoding/json"
	"strings"
	"time"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type OtsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *OtsService) getPrimaryKeyType(primaryKeyType string) tablestore.PrimaryKeyType {
	var keyType tablestore.PrimaryKeyType
	t := PrimaryKeyTypeString(primaryKeyType)
	switch t {
	case IntegerType:
		keyType = tablestore.PrimaryKeyType_INTEGER
	case StringType:
		keyType = tablestore.PrimaryKeyType_STRING
	case BinaryType:
		keyType = tablestore.PrimaryKeyType_BINARY
	}
	return keyType
}

func (s *OtsService) ListOtsTable(instanceName string) (table *tablestore.ListTableResponse, err error) {
	if _, err := s.DescribeOtsInstance(instanceName); err != nil {
		return nil, errmsgs.WrapError(err)
	}
	var raw interface{}
	var requestInfo *tablestore.TableStoreClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.ListTable()
		})
		if err != nil {
			if strings.HasSuffix(err.Error(), "no such host") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListTable", raw, requestInfo)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return table, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AliyunTablestoreGoSdk)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, instanceName, "ListTable", errmsgs.AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.ListTableResponse)
	if table == nil {
		return table, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OtsTable", instanceName)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return
}

func (s *OtsService) DescribeOtsTable(id string) (*tablestore.DescribeTableResponse, error) {
	table := &tablestore.DescribeTableResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return table, errmsgs.WrapError(err)
	}
	instanceName, tableName := parts[0], parts[1]
	request := new(tablestore.DescribeTableRequest)
	request.TableName = tableName

	if _, err := s.DescribeOtsInstance(instanceName); err != nil {
		return table, errmsgs.WrapError(err)
	}
	var raw interface{}
	var requestInfo *tablestore.TableStoreClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.DescribeTable(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DescribeTable", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return table, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AliyunTablestoreGoSdk)
		}
		return table, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, "DescribeTable", errmsgs.AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.DescribeTableResponse)
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		return table, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OtsTable", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return table, nil
}

func (s *OtsService) WaitForOtsTable(instanceName, tableName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	id := fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName)

	for {
		object, err := s.DescribeOtsTable(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.TableMeta.TableName == tableName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.TableMeta.TableName, tableName, errmsgs.ProviderERROR)
		}
	}
}

func (s *OtsService) convertPrimaryKeyType(t tablestore.PrimaryKeyType) PrimaryKeyTypeString {
	var typeString PrimaryKeyTypeString
	switch t {
	case tablestore.PrimaryKeyType_INTEGER:
		typeString = IntegerType
	case tablestore.PrimaryKeyType_BINARY:
		typeString = BinaryType
	case tablestore.PrimaryKeyType_STRING:
		typeString = StringType
	}
	return typeString
}

func (s *OtsService) ListOtsInstance(pageSize int, pageNum int) ([]string, error) {
	req := ots.CreateListInstanceRequest()
	s.client.InitRpcRequest(*req.RpcRequest)
	req.Method = "GET"
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNum = requests.NewInteger(pageNum)
	req.Domain = s.client.Domain

	var allInstanceNames []string

	for {
		raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		response, ok := raw.(*ots.ListInstanceResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alicloud_ots_instances", req.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(req.GetActionName(), raw, req.RpcRequest, req)

		if response == nil || len(response.InstanceInfos.InstanceInfo) < 1 {
			break
		}

		for _, instance := range response.InstanceInfos.InstanceInfo {
			allInstanceNames = append(allInstanceNames, instance.InstanceName)
		}

		if len(response.InstanceInfos.InstanceInfo) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNum); err != nil {
			return nil, errmsgs.WrapError(err)
		} else {
			req.PageNum = page
		}
	}
	return allInstanceNames, nil
}

func (s *OtsService) DescribeOtsInstance(id string) (inst InstanceInfo, err error) {
	request := s.client.NewCommonRequest("GET", "Ots", "2016-06-20", "GetInstance", "")
	request.QueryParams["InstanceName"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request.QueryParams, errmsgs.AlibabacloudStackSdkGoERROR)
	// OTS instance not found error code is "NotFound"
	if err != nil {
		if bresponse == nil {
			return inst, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		if errmsgs.NotFoundError(err) {
			return inst, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return inst, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var instmap GetInstanceResponse

	e := json.Unmarshal(bresponse.GetHttpContentBytes(), &instmap)
	if e != nil {
		return inst, errmsgs.WrapErrorf(e, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if instmap.AsapiSuccess {
		inst = instmap.InstanceInfo
		return inst, nil
	} else {
		return inst, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OtsInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
}

func (s *OtsService) DescribeOtsInstanceAttachment(id string) (inst ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.Method = "GET"
	request.InstanceName = id

	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	resp, ok := raw.(*ots.ListVpcInfoByInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		if errmsgs.NotFoundError(err) {
			return inst, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return inst, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if resp.TotalCount < 1 {
		return inst, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OtsInstanceAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func (s *OtsService) WaitForOtsInstanceVpc(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeOtsInstanceAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.InstanceName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceName, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *OtsService) ListOtsInstanceVpc(id string) (inst []ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.Method = "GET"
	request.InstanceName = id

	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	resp, ok := raw.(*ots.ListVpcInfoByInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return inst, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alicloud_ots_instance_attachments", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	if resp.TotalCount < 1 {
		return inst, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OtsInstanceAttachment", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	var retInfos []ots.VpcInfo
	for _, vpcInfo := range resp.VpcInfos.VpcInfo {
		vpcInfo.InstanceName = id
		retInfos = append(retInfos, vpcInfo)
	}
	return retInfos, nil
}

func (s *OtsService) WaitForOtsInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeOtsInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Status == convertOtsInstanceStatus(status) {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object.Status), status, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *OtsService) DescribeOtsInstanceTypes() (types []string, err error) {
	request := ots.CreateListClusterTypeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.Method = requests.GET

	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListClusterType(request)
	})
	resp, ok := raw.(*ots.ListClusterTypeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alicloud_ots_instance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if resp != nil {
		return resp.ClusterTypeInfos.ClusterType, nil
	}
	return
}

type TagInfos struct {
	TagInfo []map[string]string `json:"TagInfo" xml:"TagInfo"`
}

type InstanceInfo struct {
	InstanceName string            `json:"InstanceName" xml:"InstanceName"`
	Status       int               `json:"Status" xml:"Status"`
	TagInfos     TagInfos          `json:"TagInfos" xml:"TagInfos"`
	Description  string            `json:"Description" xml:"Description"`
	Quota        map[string]interface{} `json:"Quota" xml:"Quota"`
	UserId       string            `json:"UserId" xml:"UserId"`
	Network      string            `json:"Network" xml:"Network"`
	CreateTime   string            `json:"CreateTime" xml:"CreateTime"`
	ClusterType  string            `json:"ClusterType" xml:"ClusterType"`
	WriteCapacity int              `json:"WriteCapacity" xml:"WriteCapacity"`
	ReadCapacity  int              `json:"ReadCapacity" xml:"ReadCapacity"`
}

type GetInstanceResponse struct {
	EagleEyeTraceId string       `json:"eagleEyeTraceId" xml:"eagleEyeTraceId"`
	AsapiSuccess    bool         `json:"asapiSuccess" xml:"asapiSuccess"`
	RequestId       string       `json:"RequestId" xml:"RequestId"`
	AsapiRequestId  string       `json:"asapiRequestId" xml:"asapiRequestId"`
	InstanceInfo    InstanceInfo `json:"InstanceInfo" xml:"InstanceInfo"`
}