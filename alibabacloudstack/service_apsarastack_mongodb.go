package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/alibabacloud-go/tea/tea"
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
			if _, ok := err.(*errors.ServerError); ok {
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
		if errmsgs.IsExpectedErrors(err, "InvalidDBInstanceId.NotFound") {
			return instance, err
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if bresponse == nil || len(bresponse.DBInstances.DBInstance) == 0 {
		return instance, errmsgs.WrapErrorf(errmsgs.GetNotFoundErrorFromString(errmsgs.GetNotFoundMessage("MongoDB Instance", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return bresponse.DBInstances.DBInstance[0], nil
}

func (s *MongoDBService) MongoDbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
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

func (s *MongoDBService) MongoDbInstanceNodeAddressStateRefreshFunc(nodeid, netType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeShardingInstanceNode(nodeid)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		if object["enable_"+netType+"_connection"].(bool) {
			return object, "Enable", nil
		} else {
			return object, "Disable", nil
		}

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

	return nil
}

func (s *MongoDBService) DescribeMongoDBSecurityGroupId(id string) (*dds.DescribeSecurityGroupConfigurationResponse, error) {
	response := &dds.DescribeSecurityGroupConfigurationResponse{}
	request := dds.CreateDescribeSecurityGroupConfigurationRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
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

func (server *MongoDBService) ModifyMongodbShardingInstanceNode(d *schema.ResourceData, param string) error {

	instanceID := d.Id()
	nodeType := map[string]string{
		"mongo_list":        string(MongoDBShardingNodeMongos),
		"shard_list":        string(MongoDBShardingNodeShard),
		"configserver_list": string(MongoDBShardingNodeCs),
	}[param]

	stateConf := BuildStateConf(MongoDBChangingStatus, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, server.MongoDbInstanceStateRefreshFunc(d.Id(), []string{"failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	oldMap := map[string]interface{}{}
	newMap := map[string]interface{}{}

	old, new := d.GetChange(param)
	for _, n := range new.(*schema.Set).List() {
		item := n.(map[string]interface{})
		newMap[item["description"].(string)] = item
	}
	if !d.IsNewResource() {
		for _, o := range old.(*schema.Set).List() {
			item := o.(map[string]interface{})
			oldMap[item["description"].(string)] = item
		}

		// create new node
		for key, value := range newMap {
			if _, exist := oldMap[key]; !exist {
				node := value.(map[string]interface{})
				request := dds.CreateCreateNodeRequest()
				server.client.InitRpcRequest(*request.RpcRequest)
				request.DBInstanceId = instanceID
				request.NodeClass = node["node_class"].(string)
				request.NodeType = nodeType
				request.ClientToken = buildClientToken(request.GetActionName())

				if param == "shard_list" {
					request.NodeStorage = requests.NewInteger(node["node_storage"].(int))
				}

				raw, err := server.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
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

				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
				}

				nodeId := bresponse.NodeId
				reqQuery := map[string]interface{}{
					"DBInstanceId":          instanceID,
					"NodeId":                nodeId,
					"DBInstanceDescription": key,
				}
				if _, err := server.client.DoTeaRequest("POST", "Dds", "2015-12-01", "ModifyDBInstanceDescription", "", nil, reqQuery, nil); err != nil {
					return err
				}
			}
		}

		// remove old node
		for key, value := range oldMap {
			if _, exist := newMap[key]; !exist {
				node := value.(map[string]interface{})

				request := dds.CreateDeleteNodeRequest()
				server.client.InitRpcRequest(*request.RpcRequest)
				request.DBInstanceId = instanceID
				request.NodeId = node["node_id"].(string)
				request.ClientToken = buildClientToken(request.GetActionName())

				raw, err := server.client.WithDdsClient(func(client *dds.Client) (interface{}, error) {
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

				if _, err := stateConf.WaitForState(); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
				}
			}
		}
	}

	if response, err := server.DdsDescribeShardingInstanceNodes(instanceID); err != nil {
		return err
	} else {
		oldMap = response[param]
	}

	//modify node
	for key, value := range newMap {
		newNode := value.(map[string]interface{})
		var oldNode map[string]interface{}
		if v, exist := oldMap[key]; !exist {
			return fmt.Errorf("Lost Node %s", key)
		} else {
			oldNode = v.(map[string]interface{})
		}
		if newNode["node_class"].(string) != oldNode["node_class"].(string) ||
			newNode["node_storage"] != oldNode["node_storage"] {
			// node specification change
			if param == "configserver_list" {
				nodesInfo := map[string]interface{}{
					"ConfigSvrs": []map[string]interface{}{
						{
							"DBInstanceClass": newNode["node_class"],
							"Storage":         newNode["node_storage"],
							"DBInstanceName":  oldNode["node_id"],
						},
					},
				}
				jsonByte, err := json.Marshal(nodesInfo)
				if err != nil {
					return err
				}
				reqQuery := map[string]interface{}{
					"DBInstanceId": d.Id(),
					"NodesInfo":    string(jsonByte),
				}
				if _, err := server.client.DoTeaRequest("POST", "Dds", "2015-12-01", "ModifyNodeSpecBatch", "", nil, reqQuery, nil); err != nil {
					return err
				}
			} else {
				reqQuery := map[string]interface{}{
					"DBInstanceId": d.Id(),
					"NodeClass":    newNode["node_class"],
					"NodeId":       oldNode["node_id"],
				}
				if param == "shard_list" {
					reqQuery["NodeStorage"] = newNode["node_storage"]
				}
				if _, err := server.client.DoTeaRequest("POST", "Dds", "2015-12-01", "ModifyNodeSpec", "", nil, reqQuery, nil); err != nil {
					return err
				}

			}
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
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
	stateConf := BuildStateConf(MongoDBChangingStatus, []string{"Running"}, 15*time.Minute, 10*time.Second, s.MongoDbInstanceStateRefreshFunc(id, []string{"failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
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
			if errmsgs.IsExpectedErrors(err, "InvalidDBInstanceId.NotFound") {
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
	stateConf := BuildStateConf(MongoDBChangingStatus, []string{"Running"}, 15*time.Minute, 10*time.Second, s.MongoDbInstanceStateRefreshFunc(d.Id(), []string{"failed"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
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

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}

func (s *MongoDBService) ResetAccountPassword(dbInstanceId, account, password string) error {
	request := dds.CreateResetAccountPasswordRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = dbInstanceId
	request.AccountName = account
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
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, dbInstanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return err
}

func (s *MongoDBService) ResetDbAccountPassword(dbInstanceId, nodeId, account, password string) error {
	reqQuery := map[string]interface{}{
		"DBInstanceId":    dbInstanceId,
		"AccountName":     account,
		"AccountPassword": password,
		"CharacterType":   "db",
		"NodeId":          nodeId,
	}
	if _, err := s.client.DoTeaRequest("POST", "Dds", "2015-12-01", "ResetAccountPassword", "", nil, reqQuery, nil); err != nil {
		return err
	}
	return nil
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
	request := s.client.NewCommonRequest("GET", "Dds", "2022-11-21", "DescribeAccounts", "")
	DdsDescribeaccountsResponseObj := &DdsDescribeaccountsResponse{}
	// Call request_params_handler
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
	request := s.client.NewCommonRequest("GET", "Dds", "2015-12-01", "DescribeAuditPolicy", "")
	DdsDescribeauditpolicyResponseObj := &DdsDescribeauditpolicyResponse{}

	// Call request_params_handler

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

func (s *MongoDBService) doDdsDescribeauditlogfilterRequest(id string) (*DdsDescribeauditlogfilterResponse, error) {
	// api: Dds - 2015-12-01 - DescribeAuditLogFilter
	stateConf := BuildStateConf([]string{"CONFIG_SWITCHING"}, []string{"Running"}, 10*time.Minute, 10*time.Second, s.MongoDbInstanceStateRefreshFunc(id, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return nil, errmsgs.WrapError(err)
	}

	// Use the final obtained instance variable
	request := s.client.NewCommonRequest("GET", "Dds", "2015-12-01", "DescribeAuditLogFilter", "")
	DdsDescribeauditlogfilterResponseObj := &DdsDescribeauditlogfilterResponse{}

	// Call request_params_handler

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
	request := s.client.NewCommonRequest("GET", "Dds", "2015-12-01", "DescribeBackups", "")
	DdsDescribebackupsResponseObj := &DdsDescribebackupsResponse{}
	parts := strings.Split(id, "&")
	instance_id := parts[1]
	start_time := parts[2]
	end_time := time.Now().UTC().Add(time.Hour).Format("2006-01-02T15:04Z")
	// Call request_params_handler

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

type DdsDescribeshardingnetworkaddressResponse struct {
	NetworkAddresses struct {
		NetworkAddress []struct {
			NetworkAddress string `json:"NetworkAddress"`
			IPAddress      string `json:"IPAddress"`
			NetworkType    string `json:"NetworkType"`
			Port           string `json:"Port"`
			VPCId          string `json:"VPCId"`
			VswitchId      string `json:"VswitchId"`
			NodeId         string `json:"NodeId"`
			ExpiredTime    string `json:"ExpiredTime"`
			NodeType       string `json:"NodeType"`
			Role           string `json:"Role"`
		} `json:"NetworkAddress"`
	} `json:"NetworkAddresses"`

	CompatibleConnections struct {
		CompatibleConnection []struct {
			NetworkAddress string `json:"NetworkAddress"`
			IPAddress      string `json:"IPAddress"`
			NetworkType    string `json:"NetworkType"`
			Port           string `json:"Port"`
			VPCId          string `json:"VPCId"`
			VswitchId      string `json:"VswitchId"`
			ExpiredTime    string `json:"ExpiredTime"`
		} `json:"CompatibleConnection"`
	} `json:"CompatibleConnections"`
	RequestId string `json:"RequestId"`
}

type DdsDescribeDBInstancesResponse struct {
	TotalCount  int    `json:"TotalCount" xml:"TotalCount"`
	RequestId   string `json:"RequestId" xml:"RequestId"`
	PageSize    int    `json:"PageSize" xml:"PageSize"`
	PageNumber  int    `json:"PageNumber" xml:"PageNumber"`
	DBInstances struct {
		DBInstance []struct {
			ReplicaSetName              string                                            `json:"ReplicaSetName" xml:"ReplicaSetName"`
			Engine                      string                                            `json:"Engine" xml:"Engine"`
			ProvisionedIops             int64                                             `json:"ProvisionedIops" xml:"ProvisionedIops"`
			DBInstanceOrderStatus       string                                            `json:"DBInstanceOrderStatus" xml:"DBInstanceOrderStatus"`
			DBInstanceClass             string                                            `json:"DBInstanceClass" xml:"DBInstanceClass"`
			VpcAuthMode                 string                                            `json:"VpcAuthMode" xml:"VpcAuthMode"`
			LastDowngradeTime           string                                            `json:"LastDowngradeTime" xml:"LastDowngradeTime"`
			MaxConnections              int                                               `json:"MaxConnections" xml:"MaxConnections"`
			HiddenZoneId                string                                            `json:"HiddenZoneId" xml:"HiddenZoneId"`
			DBInstanceType              string                                            `json:"DBInstanceType" xml:"DBInstanceType"`
			UseClusterBackup            bool                                              `json:"UseClusterBackup" xml:"UseClusterBackup"`
			DBInstanceId                string                                            `json:"DBInstanceId" xml:"DBInstanceId"`
			NetworkType                 string                                            `json:"NetworkType" xml:"NetworkType"`
			ReplicationFactor           string                                            `json:"ReplicationFactor" xml:"ReplicationFactor"`
			EncryptionKey               string                                            `json:"EncryptionKey" xml:"EncryptionKey"`
			MaxIOPS                     int                                               `json:"MaxIOPS" xml:"MaxIOPS"`
			DBInstanceReleaseProtection bool                                              `json:"DBInstanceReleaseProtection" xml:"DBInstanceReleaseProtection"`
			ReplacateId                 string                                            `json:"ReplacateId" xml:"ReplacateId"`
			EngineVersion               string                                            `json:"EngineVersion" xml:"EngineVersion"`
			VPCId                       string                                            `json:"VPCId" xml:"VPCId"`
			BurstingEnabled             bool                                              `json:"BurstingEnabled" xml:"BurstingEnabled"`
			VPCCloudInstanceIds         string                                            `json:"VPCCloudInstanceIds" xml:"VPCCloudInstanceIds"`
			MaintainStartTime           string                                            `json:"MaintainStartTime" xml:"MaintainStartTime"`
			DBInstanceStorage           int                                               `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
			SecondaryZoneId             string                                            `json:"SecondaryZoneId" xml:"SecondaryZoneId"`
			Encrypted                   bool                                              `json:"Encrypted" xml:"Encrypted"`
			CurrentKernelVersion        string                                            `json:"CurrentKernelVersion" xml:"CurrentKernelVersion"`
			StorageType                 string                                            `json:"StorageType" xml:"StorageType"`
			ZoneId                      string                                            `json:"ZoneId" xml:"ZoneId"`
			PaymentType                 string                                            `json:"PaymentType" xml:"PaymentType"`
			LockMode                    string                                            `json:"LockMode" xml:"LockMode"`
			DBInstanceDescription       string                                            `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
			ChargeType                  string                                            `json:"ChargeType" xml:"ChargeType"`
			ReadonlyReplicas            string                                            `json:"ReadonlyReplicas" xml:"ReadonlyReplicas"`
			CapacityUnit                string                                            `json:"CapacityUnit" xml:"CapacityUnit"`
			DestroyTime                 string                                            `json:"DestroyTime" xml:"DestroyTime"`
			RegionId                    string                                            `json:"RegionId" xml:"RegionId"`
			ResourceGroupId             string                                            `json:"ResourceGroupId" xml:"ResourceGroupId"`
			CloudType                   string                                            `json:"CloudType" xml:"CloudType"`
			MaintainEndTime             string                                            `json:"MaintainEndTime" xml:"MaintainEndTime"`
			ExpireTime                  string                                            `json:"ExpireTime" xml:"ExpireTime"`
			SyncPercent                 string                                            `json:"SyncPercent" xml:"SyncPercent"`
			VSwitchId                   string                                            `json:"VSwitchId" xml:"VSwitchId"`
			CreationTime                string                                            `json:"CreationTime" xml:"CreationTime"`
			StorageEngine               string                                            `json:"StorageEngine" xml:"StorageEngine"`
			DBInstanceStatus            string                                            `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
			ProtocolType                string                                            `json:"ProtocolType" xml:"ProtocolType"`
			KindCode                    int                                               `json:"KindCode" xml:"KindCode"`
			ReplicaSets                 dds.ReplicaSetsInDescribeDBInstanceAttribute      `json:"ReplicaSets" xml:"ReplicaSets"`
			Tags                        dds.TagsInDescribeDBInstanceAttribute             `json:"Tags" xml:"Tags"`
			ConfigserverList            dds.ConfigserverList                              `json:"ConfigserverList" xml:"ConfigserverList"`
			ShardList                   dds.ShardListInDescribeDBInstances                `json:"ShardList" xml:"ShardList"`
			NetworkAddresses            dds.NetworkAddressesInDescribeDBInstanceAttribute `json:"NetworkAddresses" xml:"NetworkAddresses"`
			MongosList                  dds.MongosListInDescribeDBInstanceAttribute       `json:"MongosList" xml:"MongosList"`
		}
	}
}

func (s *MongoDBService) DoDdsDescribeshardingnetworkaddressRequest(id string) (*DdsDescribeshardingnetworkaddressResponse, error) {
	// api: Dds - 2015-12-01 - DescribeShardingNetworkAddress
	request := s.client.NewCommonRequest("GET", "Dds", "2015-12-01", "DescribeShardingNetworkAddress", "")
	DdsDescribeshardingnetworkaddressResponseObj := &DdsDescribeshardingnetworkaddressResponse{}

	// Call request_params_handler
	parts := strings.Split(id, COLON_SEPARATED)
	db_instance_id := parts[0]
	request.QueryParams["DBInstanceId"] = db_instance_id

	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("DescribeShardingNetworkAddress", bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeShardingNetworkAddress", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DdsDescribeshardingnetworkaddressResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeShardingNetworkAddress", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return DdsDescribeshardingnetworkaddressResponseObj, nil
}

func (s *MongoDBService) ModifyAuditLogFilter(d *schema.ResourceData) error {
	stateConf := BuildStateConf([]string{"CONFIG_SWITCHING"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, s.MongoDbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapError(err)
	}

	request := s.client.NewCommonRequest("POST", "Dds", "2015-12-01", "ModifyAuditLogFilter", "")
	// DdsModifyauditlogfilterResponseObj := DdsModifyauditlogfilterResponse{}
	request.QueryParams["DBInstanceId"] = d.Id()

	old, new := d.GetChange("audit_filter")
	oldValue := map[string][]string{}
	newValue := map[string][]string{}
	for _, v := range old.(*schema.Set).List() {
		i := v.(map[string]interface{})
		roleType := i["role_type"].(string)
		filters := []string{}
		for _, f := range i["filters"].(*schema.Set).List() {
			filters = append(filters, f.(string))
		}
		sort.Strings(filters)
		oldValue[roleType] = filters
	}
	for _, v := range new.(*schema.Set).List() {
		i := v.(map[string]interface{})
		roleType := i["role_type"].(string)
		filters := []string{}
		for _, f := range i["filters"].(*schema.Set).List() {
			filters = append(filters, f.(string))
		}
		sort.Strings(filters)
		newValue[roleType] = filters
	}

	for roleType := range newValue {
		if _, exist := oldValue[roleType]; exist && strings.Join(newValue[roleType], ",") == strings.Join(oldValue[roleType], ",") {
			continue
		}
		request.QueryParams["RoleType"] = roleType
		request.QueryParams["Filter"] = strings.Join(newValue[roleType], ",")

		bresponse, err := s.client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_mongo_db_audit_log_filter", "ModifyAuditLogFilter", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		stateConf := BuildStateConf([]string{"CONFIG_SWITCHING"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, s.MongoDbInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func (s *MongoDBService) GetAuditLogFilter(id string) ([]map[string]interface{}, error) {
	if response, err := s.doDdsDescribeauditlogfilterRequest(id); err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mongodb_auditlogfilter", errmsgs.AlibabacloudStackSdkGoERROR)
	} else {
		auditFilters := []map[string]interface{}{}
		for _, item := range strings.Split(response.Filter, "-") {
			var filters []string
			var filterType string
			parts := strings.Split(item, "@")
			if len(parts) == 2 {
				filterType = parts[0]
				filters = strings.Split(parts[1], ",")
			} else {
				filterType = "db"
				filters = strings.Split(parts[0], ",")
			}
			auditFilters = append(auditFilters, map[string]interface{}{
				"role_type": filterType,
				"filters":   filters,
			})
		}
		return auditFilters, nil
	}
}

func (s *MongoDBService) DdsDescribeShardingInstanceNodes(id string) (map[string]map[string]interface{}, error) {
	result := map[string]map[string]interface{}{
		"mongo_list":        {},
		"shard_list":        {},
		"configserver_list": {},
	}

	if response, err := s.DescribeMongoDBInstance(id); err != nil {
		return result, err
	} else {
		for _, node := range response.MongosList.MongosAttribute {
			result["mongo_list"][node.NodeDescription] = map[string]interface{}{
				"node_class":  node.NodeClass,
				"node_id":     node.NodeId,
				"description": node.NodeDescription,
			}
		}
		for _, node := range response.ShardList.ShardAttribute {
			result["shard_list"][node.NodeDescription] = map[string]interface{}{
				"node_class":   node.NodeClass,
				"node_id":      node.NodeId,
				"description":  node.NodeDescription,
				"node_storage": node.NodeStorage,
			}
		}
		for _, node := range response.ConfigserverList.ConfigserverAttribute {
			result["configserver_list"][node.NodeDescription] = map[string]interface{}{
				"node_class":   node.NodeClass,
				"node_id":      node.NodeId,
				"description":  node.NodeDescription,
				"node_storage": node.NodeStorage,
			}
		}
	}

	return result, nil
}

func (s *MongoDBService) DescribeShardingInstanceNode(id string) (map[string]interface{}, error) {

	result := map[string]interface{}{
		"enable_public_connection":  false,
		"enable_private_connection": false,
	}
	parts := strings.SplitN(id, ":", 2)
	instanceId := parts[0]
	nodeId := parts[1]

	reqQuery := map[string]interface{}{"DBInstanceId": instanceId}
	if response, err := s.client.DoTeaRequest("GET", "Dds", "2015-12-01", "DescribeShardingNetworkAddress", "", nil, reqQuery, nil); err != nil {
		return nil, err
	} else {
		for _, v := range response["NetworkAddresses"].(map[string]interface{})["NetworkAddress"].([]interface{}) {
			address := v.(map[string]interface{})
			if address["NodeId"].(string) != nodeId {
				continue
			}
			port, err := strconv.Atoi(address["Port"].(string))
			if err != nil {
				return nil, err
			}

			if address["NetworkType"].(string) == "Public" {
				result["enable_public_connection"] = true
				result["public_connect_string"] = address["NetworkAddress"].(string)
				result["public_connect_port"] = port
			} else {
				result["enable_private_connection"] = true
				result["private_connect_string"] = address["NetworkAddress"].(string)
				result["private_connect_port"] = port
			}

		}
	}
	return result, nil
}

func (s *MongoDBService) UpdateInstanceConnection(id string, existedConnections, targetConnections map[string]map[string]interface{}) error {
	updatedList1 := []map[string]interface{}{}
	updatedList2 := []map[string]interface{}{}
	for key, v := range targetConnections {
		if _, exist := existedConnections[key]; !exist {
			updatedList1 = append(updatedList1, v)
		}
	}

	for key, v := range existedConnections {
		if _, exist := targetConnections[key]; !exist {
			updatedList2 = append(updatedList2, v)
		}
	}

	if len(updatedList1) != len(updatedList2) {
		return fmt.Errorf("The items to be updated are inconsistent")
	}

	for index := range updatedList1 {
		targetConnection := updatedList1[index]
		existedConnection := updatedList2[index]
		reqQuery := map[string]interface{}{
			"DBInstanceId":            id,
			"NodeId":                  nil,
			"NewConnectionString":     targetConnection["connect_string_prefix"],
			"CurrentConnectionString": existedConnection["connect_string"],
			"NewPort":                 targetConnection["connect_port"],
			"OldPort":                 existedConnection["connect_port"],
		}
		s.ModifyDBInstanceConnectionString(reqQuery)
	}
	return nil
}

func (s *MongoDBService) ModifyDBInstanceConnectionString(reqQuery map[string]interface{}) error {
	id := reqQuery["DBInstanceId"].(string)
	stateConf := BuildStateConf(MongoDBChangingStatus, []string{"Running"}, 10*time.Minute, 10*time.Second, s.MongoDbInstanceStateRefreshFunc(id, []string{"Deleting"}))

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := s.client.DoTeaRequest("POST", "Dds", "2015-12-01", "ModifyDBInstanceConnectionString", "", nil, reqQuery, nil)

		if err == nil {
			return nil
		}

		if sdkError, ok := err.(*tea.SDKError); ok && *sdkError.Code == "OperationDenied.DBInstanceStatus" {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, id)
	}
	return nil
}

func (s *MongoDBService) SwtichNodeConnection(action, dbInstanceId, nodeId string) error {
	reqQuery := map[string]interface{}{
		"DBInstanceId": dbInstanceId,
		"NodeId":       nodeId,
	}
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := s.client.DoTeaRequest("POST", "Dds", "2015-12-01", action, "", nil, reqQuery, nil)

		if err == nil {
			return nil
		}

		if sdkError, ok := err.(*tea.SDKError); ok && *sdkError.Code == "OperationDenied.DBInstanceStatus" {
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}
