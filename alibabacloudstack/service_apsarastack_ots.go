package alibabacloudstack

import (
	"encoding/json"
	"strings"

	"time"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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
		return nil, WrapError(err)
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
			return table, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return nil, WrapErrorf(err, DataDefaultErrorMsg, instanceName, "ListTable", AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.ListTableResponse)
	if table == nil {
		return table, WrapErrorf(Error(GetNotFoundMessage("OtsTable", instanceName)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *OtsService) DescribeOtsTable(id string) (*tablestore.DescribeTableResponse, error) {
	table := &tablestore.DescribeTableResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return table, WrapError(err)
	}
	instanceName, tableName := parts[0], parts[1]
	request := new(tablestore.DescribeTableRequest)
	request.TableName = tableName

	if _, err := s.DescribeOtsInstance(instanceName); err != nil {
		return table, WrapError(err)
	}
	var raw interface{}
	var requestInfo *tablestore.TableStoreClient
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
			requestInfo = tableStoreClient
			return tableStoreClient.DescribeTable(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTableIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DescribeTable", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "OTSObjectNotExist") {
			return table, WrapErrorf(err, NotFoundMsg, AliyunTablestoreGoSdk)
		}
		return table, WrapErrorf(err, DefaultErrorMsg, id, "DescribeTable", AliyunTablestoreGoSdk)
	}
	table, _ = raw.(*tablestore.DescribeTableResponse)
	if table == nil || table.TableMeta == nil || table.TableMeta.TableName != tableName {
		return table, WrapErrorf(Error(GetNotFoundMessage("OtsTable", id)), NotFoundMsg, ProviderERROR)
	}
	return table, nil
}

func (s *OtsService) WaitForOtsTable(instanceName, tableName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	id := fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName)

	for {
		object, err := s.DescribeOtsTable(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.TableMeta.TableName == tableName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.TableMeta.TableName, tableName, ProviderERROR)
		}

	}
}

// Convert tablestore.PrimaryKeyType to PrimaryKeyTypeString
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
	req.RegionId = s.client.RegionId
	req.Method = "GET"
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNum = requests.NewInteger(pageNum)
	req.Headers = map[string]string{"RegionId": s.client.RegionId}
	req.Domain = s.client.Domain

	req.QueryParams = map[string]string{ "Product": "Ots", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	var allInstanceNames []string

	for {
		raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.ListInstance(req)
		})
		if err != nil {
			return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instances", req.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(req.GetActionName(), raw, req.RpcRequest, req)
		response, _ := raw.(*ots.ListInstanceResponse)

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
			return nil, WrapError(err)
		} else {
			req.PageNum = page
		}
	}
	return allInstanceNames, nil
}

func (s *OtsService) DescribeOtsInstance(id string) (inst InstanceInfo, err error) {
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "Ots"
	request.Version = "2016-06-20"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetInstance"
	request.RegionId = s.client.RegionId
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		
		"InstanceName":    id,
		"RegionId":        s.client.RegionId,
		"apiName":         "GetInstance",
		"Action":          "GetInstance",
		"version":         "2016-06-20",
		
		"Product":         "Ots",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"ResourceOwnerId": s.client.ResourceGroup,
	}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	addDebug(request.GetActionName(), raw, request.QueryParams, AlibabacloudStackSdkGoERROR)
	// OTS instance not found error code is "NotFound"
	if err != nil {
		if NotFoundError(err) {
			return inst, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return inst, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	var instmap GetInstanceResponse
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse == nil {
		return inst, WrapErrorf(Error(GetNotFoundMessage("OtsInstance", id)), NotFoundMsg, ProviderERROR)
	}
	e := json.Unmarshal(bresponse.GetHttpContentBytes(), &instmap)
	if e != nil {
		return inst, WrapErrorf(e, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	if instmap.AsapiSuccess {
		inst = instmap.InstanceInfo
		return inst, nil
	} else {
		return inst, WrapErrorf(Error(GetNotFoundMessage("OtsInstance", id)), NotFoundMsg, ProviderERROR)
	}
}

func (s *OtsService) DescribeOtsInstanceAttachment(id string) (inst ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	request.RegionId = s.client.RegionId
	request.Method = "GET"
	request.InstanceName = id
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	request.QueryParams = map[string]string{ "Product": "Ots", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return inst, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return inst, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp.TotalCount < 1 {
		return inst, WrapErrorf(Error(GetNotFoundMessage("OtsInstanceAttachment", id)), NotFoundMsg, ProviderERROR)
	}
	return resp.VpcInfos.VpcInfo[0], nil
}

func (s *OtsService) WaitForOtsInstanceVpc(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeOtsInstanceAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.InstanceName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.InstanceName, id, ProviderERROR)
		}

	}
}

func (s *OtsService) ListOtsInstanceVpc(id string) (inst []ots.VpcInfo, err error) {
	request := ots.CreateListVpcInfoByInstanceRequest()
	request.RegionId = s.client.RegionId
	request.Method = "GET"
	request.InstanceName = id
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	request.QueryParams = map[string]string{ "Product": "Ots", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListVpcInfoByInstance(request)
	})
	if err != nil {
		return inst, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ots_instance_attachments", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListVpcInfoByInstanceResponse)
	if resp.TotalCount < 1 {
		return inst, WrapErrorf(Error(GetNotFoundMessage("OtsInstanceAttachment", id)), NotFoundMsg, ProviderERROR)
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
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.Status == convertOtsInstanceStatus(status) {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, fmt.Sprint(object.Status), status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *OtsService) DescribeOtsInstanceTypes() (types []string, err error) {
	request := ots.CreateListClusterTypeRequest()
	request.Method = requests.GET
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	request.QueryParams = map[string]string{ "Product": "Ots", "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	raw, err := s.client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.ListClusterType(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instance", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ots.ListClusterTypeResponse)
	if resp != nil {
		return resp.ClusterTypeInfos.ClusterType, nil
	}
	return
}

type TagInfos struct {
	TagInfo []map[string]string `json:"TagInfo" xml:"TagInfo"`
}

type InstanceInfo struct {
	InstanceName  string                 `json:"InstanceName" xml:"InstanceName"`
	Status        int                    `json:"Status" xml:"Status"`
	TagInfos      TagInfos               `json:"TagInfos" xml:"TagInfos"`
	Description   string                 `json:"Description" xml:"Description"`
	Quota         map[string]interface{} `json:"Quota" xml:"Quota"`
	UserId        string                 `json:"UserId" xml:"UserId"`
	Network       string                 `json:"Network" xml:"Network"`
	CreateTime    string                 `json:"CreateTime" xml:"CreateTime"`
	ClusterType   string                 `json:"ClusterType" xml:"ClusterType"`
	WriteCapacity int                    `json:"WriteCapacity" xml:"WriteCapacity"`
	ReadCapacity  int                    `json:"ReadCapacity" xml:"ReadCapacity"`
}

type GetInstanceResponse struct {
	EagleEyeTraceId string       `json:"eagleEyeTraceId" xml:"eagleEyeTraceId"`
	AsapiSuccess    bool         `json:"asapiSuccess" xml:"asapiSuccess"`
	RequestId       string       `json:"RequestId" xml:"RequestId"`
	AsapiRequestId  string       `json:"asapiRequestId" xml:"asapiRequestId"`
	InstanceInfo    InstanceInfo `json:"InstanceInfo" xml:"InstanceInfo"`
}
