package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type MongoDBService struct {
	client *connectivity.AlibabacloudStackClient
}

type DdsDescribeaccountsResponse struct {
	Accounts struct {
		Account []struct {
			Roles struct {
				AccountRole []struct {
					DatabaseName string `json:"DatabaseName"`
					RoleName     string `json:"RoleName"`
				} `json:"AccountRole"`
			} `json:"Roles"`
			DBInstanceId       string `json:"DBInstanceId"`
			AccountName        string `json:"AccountName"`
			DatabaseName       string `json:"DatabaseName"`
			AccountStatus      string `json:"AccountStatus"`
			AccountDescription string `json:"AccountDescription"`
			CharacterType      string `json:"CharacterType"`
			AccountType        string `json:"AccountType"`
			AccountId          int    `json:"AccountId"`
		} `json:"Account"`
	} `json:"Accounts"`
	RequestId string `json:"RequestId"`
}

func (s *MongoDBService) DescribeMongoDBInstance(id string) (instance dds.DBInstance, err error) {
	request := dds.CreateDescribeDBInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	// raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
	// 	return client.DescribeDBInstanceAttribute(request)
	// })
	var raw interface{}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.DescribeDBInstanceAttribute(request)
		})

		if err != nil {
			if sdkErr, ok := err.(*errors.ServerError); ok && sdkErr.ErrorCode() == "InternalError" {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	bresponse, ok := raw.(*dds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse == nil || len(bresponse.DBInstances.DBInstance) == 0 {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("MongoDB Instance", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return bresponse.DBInstances.DBInstance[0], nil
}

// WaitForInstance waits for instance to given statusid
func (s *MongoDBService) WaitForMongoDBInstance(instanceId string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		instance, err := s.DescribeMongoDBInstance(instanceId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}

		if instance.DBInstanceStatus == string(status) {
			return nil
		}

		if status == Updating {
			if instance.DBInstanceStatus == "NodeCreating" ||
				instance.DBInstanceStatus == "NodeDeleting" ||
				instance.DBInstanceStatus == "DBInstanceClassChanging" {
				return nil
			}
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, instanceId, GetFunc(1), timeout, instance.DBInstanceStatus, string(status), errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

func (s *MongoDBService) RdsMongodbDBInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *MongoDBService) DescribeMongoDBSecurityIps(instanceId string) (ips []string, err error) {
	request := dds.CreateDescribeSecurityIpsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityIps(request)
	})
	bresponse, ok := raw.(*dds.DescribeSecurityIpsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return ips, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range bresponse.SecurityIpGroups.SecurityIpGroup {
		if ip.SecurityIpGroupAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIpList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}

	return finalIps, nil
}

func (s *MongoDBService) ModifyMongoDBSecurityIps(instanceId, ips string) error {
	request := dds.CreateModifySecurityIpsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId
	request.SecurityIps = ips

	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	bresponse, ok := raw.(*dds.ModifySecurityIpsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForMongoDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBSecurityGroupId(id string) (*dds.DescribeSecurityGroupConfigurationResponse, error) {
	response := &dds.DescribeSecurityGroupConfigurationResponse{}
	request := dds.CreateDescribeSecurityGroupConfigurationRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	if err := s.WaitForMongoDBInstance(id, Running, DefaultTimeoutMedium); err != nil {
		return response, errmsgs.WrapError(err)
	}
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeSecurityGroupConfiguration(request)
	})
	bresponse, ok := raw.(*dds.DescribeSecurityGroupConfigurationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeSecurityGroupConfigurationResponse)

	return response, nil
}

func (server *MongoDBService) ModifyMongodbShardingInstanceNode(
	instanceID string, nodeType MongoDBShardingNodeType, stateList, diffList []interface{}) error {
	client := server.client

	err := server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	//create node
	if len(stateList) < len(diffList) {
		createList := diffList[len(stateList):]
		diffList = diffList[:len(stateList)]

		for _, item := range createList {
			node := item.(map[string]interface{})

			request := dds.CreateCreateNodeRequest()
			server.client.InitRpcRequest(*request.RpcRequest)
			request.DBInstanceId = instanceID
			request.NodeClass = node["node_class"].(string)
			request.NodeType = string(nodeType)
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(node["node_storage"].(int))
			}

			raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
				return client.CreateNode(request)
			})
			bresponse, ok := raw.(*dds.CreateNodeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceID, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			err = server.WaitForMongoDBInstance(instanceID, Updating, DefaultLongTimeout)
			if err != nil {
				return errmsgs.WrapError(err)
			}

			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	} else if len(stateList) > len(diffList) {
		deleteList := stateList[len(diffList):]
		stateList = stateList[:len(diffList)]

		for _, item := range deleteList {
			node := item.(map[string]interface{})

			request := dds.CreateDeleteNodeRequest()
			server.client.InitRpcRequest(*request.RpcRequest)
			request.DBInstanceId = instanceID
			request.NodeId = node["node_id"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
				return client.DeleteNode(request)
			})
			bresponse, ok := raw.(*dds.DeleteNodeResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceID, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			addDebug(request.GetActionName(), raw, request.RpcRequest, request)

			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}

	//modify node
	for key := 0; key < len(stateList); key++ {
		state := stateList[key].(map[string]interface{})
		diff := diffList[key].(map[string]interface{})

		if state["node_class"] != diff["node_class"] ||
			state["node_storage"] != diff["node_storage"] {
			request := dds.CreateModifyNodeSpecRequest()
			server.client.InitRpcRequest(*request.RpcRequest)
			request.DBInstanceId = instanceID
			request.NodeClass = diff["node_class"].(string)
			request.ClientToken = buildClientToken(request.GetActionName())

			if nodeType == MongoDBShardingNodeShard {
				request.NodeStorage = requests.NewInteger(diff["node_storage"].(int))
			}
			request.NodeId = state["node_id"].(string)

			raw, err := client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
				return client.ModifyNodeSpec(request)
			})
			bresponse, ok := raw.(*dds.ModifyNodeSpecResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceID, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			err = server.WaitForMongoDBInstance(instanceID, Updating, DefaultLongTimeout)
			if err != nil {
				return errmsgs.WrapError(err)
			}
			err = server.WaitForMongoDBInstance(instanceID, Running, DefaultLongTimeout)
			if err != nil {
				return errmsgs.WrapError(err)
			}
		}
	}
	return nil
}

func (s *MongoDBService) DescribeMongoDBBackupPolicy(id string) (*dds.DescribeBackupPolicyResponse, error) {
	response := &dds.DescribeBackupPolicyResponse{}
	request := dds.CreateDescribeBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeBackupPolicy(request)
	})
	bresponse, ok := raw.(*dds.DescribeBackupPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeBackupPolicyResponse)
	return response, nil
}

func (s *MongoDBService) DescribeMongoDBTDEInfo(id string) (*dds.DescribeDBInstanceTDEInfoResponse, error) {

	response := &dds.DescribeDBInstanceTDEInfoResponse{}
	request := dds.CreateDescribeDBInstanceTDEInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	statErr := s.WaitForMongoDBInstance(id, Running, DefaultLongTimeout)
	if statErr != nil {
		return response, errmsgs.WrapError(statErr)
	}
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeDBInstanceTDEInfo(request)
	})
	bresponse, ok := raw.(*dds.DescribeDBInstanceTDEInfoResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceTDEInfoResponse)
	return response, nil
}

func (s *MongoDBService) DescribeDBInstanceSSL(id string) (*dds.DescribeDBInstanceSSLResponse, error) {
	response := &dds.DescribeDBInstanceSSLResponse{}
	request := dds.CreateDescribeDBInstanceSSLRequest()

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		instance, err := s.DescribeMongoDBInstance(id)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		if instance.DBInstanceStatus == "SSLModifying" {
			return resource.RetryableError(fmt.Errorf("SSLModifying"))
		}
		addDebug(request.GetActionName(), instance, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.DescribeDBInstanceSSL(request)
	})
	bresponse, ok := raw.(*dds.DescribeDBInstanceSSLResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*dds.DescribeDBInstanceSSLResponse)
	return response, nil
}

func (s *MongoDBService) MotifyMongoDBBackupPolicy(d *schema.ResourceData) error {
	if err := s.WaitForMongoDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	periodList := expandStringList(connectivity.GetResourceData(d, "preferred_backup_period", "backup_period").(*schema.Set).List())
	backupPeriod := fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	backupTime := connectivity.GetResourceData(d, "preferred_backup_time", "backup_time").(string)

	request := dds.CreateModifyBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.PreferredBackupPeriod = backupPeriod
	request.PreferredBackupTime = backupTime
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ModifyBackupPolicy(request)
	})
	bresponse, ok := raw.(*dds.ModifyBackupPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err := s.WaitForMongoDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *MongoDBService) ResetAccountPassword(d *schema.ResourceData, password string) error {
	request := dds.CreateResetAccountPasswordRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.AccountName = "root"
	request.AccountPassword = password
	raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
		return client.ResetAccountPassword(request)
	})
	bresponse, ok := raw.(*dds.ResetAccountPasswordResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *MongoDBService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})

	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := dds.CreateUntagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = "INSTANCE"
		request.TagKey = &tagKey
		raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		bresponse, ok := raw.(*dds.UntagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := dds.CreateTagResourcesRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = "INSTANCE"
		raw, err := s.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		bresponse, ok := raw.(*dds.TagResourcesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	//d.SetPartial("tags")
	return nil
}

func (s *MongoDBService) tagsToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTag(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) tagsInAttributeToMap(tags []dds.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTagInAttribute(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *MongoDBService) ignoreTagInAttribute(t dds.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *MongoDBService) diffTags(oldTags, newTags []dds.TagResourcesTag) ([]dds.TagResourcesTag, []dds.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []dds.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *MongoDBService) tagsFromMap(m map[string]interface{}) []dds.TagResourcesTag {
	result := make([]dds.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, dds.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *MongoDBService) DoDdsDescribeaccountsRequest(id string) (*DdsDescribeaccountsResponse, error) {
	// api: Dds - 2022-11-21 - DescribeAccounts
	request := s.client.NewCommonRequest("POST", "Dds", "2022-11-21", "DescribeAccounts", "")
	DdsDescribeaccountsResponseObj := &DdsDescribeaccountsResponse{}
	//调用request_params_handler
	parts := strings.Split(id, COLON_SEPARATED)
	instance_id := parts[1]
	request.QueryParams["DBInstanceId"] = instance_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DdsDescribeaccountsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DdsDescribeaccountsResponseObj, nil
}

type DdsDescribeauditpolicyResponse struct {
	RequestId      string `json:"RequestId"`
	LogAuditStatus string `json:"LogAuditStatus"`
}

func (s *MongoDBService) DoDdsDescribeauditpolicyRequest(id string) (*DdsDescribeauditpolicyResponse, error) {
	// api: Dds - 2015-12-01 - DescribeAuditPolicy
	request := s.client.NewCommonRequest("POST", "Dds", "2015-12-01", "DescribeAuditPolicy", "")
	DdsDescribeauditpolicyResponseObj := &DdsDescribeauditpolicyResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAuditPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DdsDescribeauditpolicyResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAuditPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DdsDescribeauditpolicyResponseObj, nil
}

type DdsDescribeauditlogfilterResponse struct {
	RequestId string `json:"RequestId"`
	Filter    string `json:"Filter"`
	RoleType  string `json:"RoleType"`
}

func (s *MongoDBService) DoDdsDescribeauditlogfilterRequest(id string) (*DdsDescribeauditlogfilterResponse, error) {
	// api: Dds - 2015-12-01 - DescribeAuditLogFilter
	const maxRetries = 24
	var instance dds.DBInstance

	for i := 0; i < maxRetries; i++ {
		var err error
		instance, err = s.DescribeMongoDBInstance(id)

		if err != nil {
			log.Printf("DescribeMongoDBInstance failed: %v", err)
			continue
		}

		if instance.DBInstanceStatus != "CONFIG_SWITCHING" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	// 使用最终获取的 instance 变量
	request := s.client.NewCommonRequest("POST", "Dds", "2015-12-01", "DescribeAuditLogFilter", "")
	DdsDescribeauditlogfilterResponseObj := &DdsDescribeauditlogfilterResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAuditLogFilter", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DdsDescribeauditlogfilterResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAuditLogFilter", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DdsDescribeauditlogfilterResponseObj, nil
}

func (s *MongoDBService) FormatMongodbTime(mongodbtime string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", mongodbtime)
	if err != nil {
		return "", fmt.Errorf("parse mongodb time failed: %w", err)
	}
	outputLayout := "2006-01-02T15:04Z"
	format_time := t.Format(outputLayout)
	return format_time, nil
}

type DdsDescribebackupsResponse struct {
	Backups struct {
		Backup []struct {
			BackupDBNames             string `json:"BackupDBNames"`
			BackupId                  int    `json:"BackupId"`
			BackupStatus              string `json:"BackupStatus"`
			BackupStartTime           string `json:"BackupStartTime"`
			BackupEndTime             string `json:"BackupEndTime"`
			BackupType                string `json:"BackupType"`
			BackupMode                string `json:"BackupMode"`
			BackupMethod              string `json:"BackupMethod"`
			BackupDownloadURL         string `json:"BackupDownloadURL"`
			BackupIntranetDownloadURL string `json:"BackupIntranetDownloadURL"`
			BackupSize                int    `json:"BackupSize"`
		} `json:"Backup"`
	} `json:"Backups"`
	RequestId  string `json:"RequestId"`
	PageNumber int    `json:"PageNumber"`
	PageSize   int    `json:"PageSize"`
	TotalCount int    `json:"TotalCount"`
}

func (s *MongoDBService) DoDdsDescribebackupsRequest(id string) (*DdsDescribebackupsResponse, error) {
	// api: Dds - 2015-12-01 - DescribeBackups
	request := s.client.NewCommonRequest("POST", "Dds", "2015-12-01", "DescribeBackups", "")
	DdsDescribebackupsResponseObj := &DdsDescribebackupsResponse{}
	parts := strings.Split(id, "&")
	instance_id := parts[1]
	start_time := parts[2]
	end_time := time.Now().UTC().Add(time.Hour).Format("2006-01-02T15:04Z")
	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = instance_id
	request.QueryParams["StartTime"] = start_time
	request.QueryParams["EndTime"] = end_time
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DdsDescribebackupsResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackups", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return DdsDescribebackupsResponseObj, nil
}
