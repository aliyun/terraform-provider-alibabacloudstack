package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RdsService struct {
	client *connectivity.AlibabacloudStackClient
}

//	_______________                      _______________                       _______________
//	|              | ______param______\  |              |  _____request_____\  |              |
//	|   Business   |                     |    Service   |                      |    SDK/API   |
//	|              | __________________  |              |  __________________  |              |
//	|______________| \    (obj, err)     |______________|  \ (status, cont)    |______________|
//	                    |                                    |
//	                    |A. {instance, nil}                  |a. {200, content}
//	                    |B. {nil, error}                     |b. {200, nil}
//	               					  |c. {4xx, nil}
//
// The API return 200 for resource not found.
// When getInstance is empty, then throw InstanceNotfound error.
// That the business layer only need to check error.
var DBInstanceStatusCatcher = Catcher{"OperationDenied.DBInstanceStatus", 60, 5}

func (s *RdsService) DescribeDBInstance(id string) (*rds.DBInstanceAttribute, error) {
	instance := &rds.DBInstanceAttribute{}
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceAttribute(request)
	})
	response, ok := raw.(*rds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.Items.DBInstanceAttribute) < 1 {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &response.Items.DBInstanceAttribute[0], nil
}

func (s *RdsService) DescribeTasks(id string) (task *rds.DescribeTasksResponse, err error) {
	request := rds.CreateDescribeTasksRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeTasks(request)
	})
	response, ok := raw.(*rds.DescribeTasksResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return task, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return task, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

func (s *RdsService) DescribeDBReadonlyInstance(id string) (*rds.DBInstanceAttribute, error) {
	instance := &rds.DBInstanceAttribute{}
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceAttribute(request)
	})
	response, ok := raw.(*rds.DescribeDBInstanceAttributeResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return instance, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return instance, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(response.Items.DBInstanceAttribute) < 1 {
		return instance, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return &response.Items.DBInstanceAttribute[0], nil
}

func (s *RdsService) DoRdsDescribeaccountsRequest(id string) (*rds.DBInstanceAccount, error) {
    return s.DescribeDBAccount(id)
}
func (s *RdsService) DescribeDBAccount(id string) (*rds.DBInstanceAccount, error) {
	ds := &rds.DBInstanceAccount{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return ds, errmsgs.WrapError(err)
	}
	request := rds.CreateDescribeAccountsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var raw interface{}
	err = invoker.Run(func() error {
		raw, err = s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	response, ok := raw.(*rds.DescribeAccountsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ds, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return ds, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if len(response.Accounts.DBInstanceAccount) < 1 {
		return ds, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBAccount", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return &response.Accounts.DBInstanceAccount[0], nil
}

func (s *RdsService) DescribeDBAccountPrivilege(id string) (*rds.DBInstanceAccount, error) {
	ds := &rds.DBInstanceAccount{}
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return ds, errmsgs.WrapError(err)
	}
	request := rds.CreateDescribeAccountsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]
	invoker := NewInvoker()
	invoker.AddCatcher(DBInstanceStatusCatcher)
	var raw interface{}
	err = invoker.Run(func() error {
		raw, err = s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeAccounts(request)
		})
		if err != nil {
			return err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})
	response, ok := raw.(*rds.DescribeAccountsResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ds, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return ds, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	if len(response.Accounts.DBInstanceAccount) < 1 {
		return ds, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBAccountPrivilege", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return &response.Accounts.DBInstanceAccount[0], nil
}

func (s *RdsService) DoRdsDescribedatabasesRequest(id string) (*rds.Database, error) {
    return s.DescribeDBDatabase(id)
}
func (s *RdsService) DescribeDBDatabase(id string) (*rds.Database, error) {
	ds := &rds.Database{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return ds, errmsgs.WrapError(err)
	}
	dbName := parts[1]
	request := rds.CreateDescribeDatabasesRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.DBName = dbName

	err = resource.Retry(30*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDatabases(request)
		})
		response, ok := raw.(*rds.DescribeDatabasesResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR))
			}
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR))
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if len(response.Databases.Database) < 1 {
			return resource.NonRetryableError(errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBDatabase", dbName)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR))
		}
		ds = &response.Databases.Database[0]
		return nil
	})
	return ds, err
}

func (s *RdsService) DescribeParameters(id string) (*rds.DescribeParametersResponse, error) {
	ds := &rds.DescribeParametersResponse{}
	request := rds.CreateDescribeParametersRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeParameters(request)
	})
	response, ok := raw.(*rds.DescribeParametersResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return ds, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return ds, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response, err
}

func (s *RdsService) RefreshParameters(d *schema.ResourceData, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		d.Set(attribute, param)
		return nil
	}
	object, err := s.DescribeParameters(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}

	var parameters = make(map[string]interface{})
	for _, i := range object.RunningParameters.DBInstanceParameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, i := range object.ConfigParameters.DBInstanceParameter {
		if i.ParameterName != "" {
			parameter := map[string]interface{}{
				"name":  i.ParameterName,
				"value": i.ParameterValue,
			}
			parameters[i.ParameterName] = parameter
		}
	}

	for _, parameter := range documented.(*schema.Set).List() {
		name := parameter.(map[string]interface{})["name"]
		for _, value := range parameters {
			if value.(map[string]interface{})["name"] == name {
				param = append(param, value.(map[string]interface{}))
				break
			}
		}
	}
	if err := d.Set(attribute, param); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *RdsService) ModifyParameters(d *schema.ResourceData, attribute string) error {
	request := rds.CreateModifyParameterRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	request.Forcerestart = requests.NewBoolean(d.Get("force_restart").(bool))
	config := make(map[string]string)
	allConfig := make(map[string]string)
	o, n := d.GetChange(attribute)
	os, ns := o.(*schema.Set), n.(*schema.Set)
	add := ns.Difference(os).List()
	if len(add) > 0 {
		for _, i := range add {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			config[key] = value
		}
		cfg, _ := json.Marshal(config)
		request.Parameters = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForDBInstance(d.Id(), Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
		// Need to check whether some parameter needs restart
		if !d.Get("force_restart").(bool) {
			req := rds.CreateDescribeParameterTemplatesRequest()
			s.client.InitRpcRequest(*req.RpcRequest)
			req.DBInstanceId = d.Id()
			req.Engine = d.Get("engine").(string)
			req.EngineVersion = d.Get("engine_version").(string)
			req.ClientToken = buildClientToken(req.GetActionName())
			forceRestartMap := make(map[string]string)
			raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.DescribeParameterTemplates(req)
			})
			response, ok := raw.(*rds.DescribeParameterTemplatesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
			for _, para := range response.Parameters.TemplateRecord {
				if para.ForceRestart == "true" {
					forceRestartMap[para.ParameterName] = para.ForceRestart
				}
			}
			if len(forceRestartMap) > 0 {
				for key, _ := range config {
					if _, ok := forceRestartMap[key]; ok {
						return errmsgs.WrapError(fmt.Errorf("Modifying RDS instance's parameter '%s' requires setting 'force_restart = true'.", key))
					}
				}
			}
		}
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyParameter(request)
		})
		response, ok := raw.(*rds.ModifyParameterResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		if err := s.WaitForDBParameter(d.Id(), DefaultTimeoutMedium, allConfig); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	//d.SetPartial(attribute)
	return nil
}

func (s *RdsService) DescribeDBInstanceNetInfo(id string) ([]rds.DBInstanceNetInfo, error) {
	request := rds.CreateDescribeDBInstanceNetInfoRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceNetInfo(request)
	})

	response, ok := raw.(*rds.DescribeDBInstanceNetInfoResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(response.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return response.DBInstanceNetInfos.DBInstanceNetInfo, nil
}

func (s *RdsService) DescribeDBConnection(id string) (*rds.DBInstanceNetInfo, error) {
	info := &rds.DBInstanceNetInfo{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return info, errmsgs.WrapError(err)
	}
	object, err := s.DescribeDBInstanceNetInfo(parts[0])

	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound"}) {
			return info, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return info, errmsgs.WrapError(err)
	}

	if object != nil {
		for _, o := range object {
			if strings.HasPrefix(o.ConnectionString, parts[1]) {
				return &o, nil
			}
		}
	}

	return info, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *RdsService) DescribeDBReadWriteSplittingConnection(id string) (*rds.DBInstanceNetInfo, error) {
	ds := &rds.DBInstanceNetInfo{}
	object, err := s.DescribeDBInstanceNetInfo(id)
	if err != nil && !errmsgs.NotFoundError(err) {
		return ds, err
	}

	if object != nil {
		for _, conn := range object {
			if conn.ConnectionStringType != "ReadWriteSplitting" {
				continue
			}
			if conn.MaxDelayTime == "" {
				continue
			}
			if _, err := strconv.Atoi(conn.MaxDelayTime); err != nil {
				return ds, err
			}
			return &conn, nil
		}
	}

	return ds, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ReadWriteSplittingConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *RdsService) GrantAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := rds.CreateGrantAccountPrivilegeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName
	request.AccountPrivilege = parts[2]

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.GrantAccountPrivilege(request)
		})
		response, ok := raw.(*rds.GrantAccountPrivilegeResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if err := s.WaitForAccountPrivilege(id, dbName, Available, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func (s *RdsService) RevokeAccountPrivilege(id, dbName string) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := rds.CreateRevokeAccountPrivilegeRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = parts[0]
	request.AccountName = parts[1]
	request.DBName = dbName

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.RevokeAccountPrivilege(request)
		})
		response, ok := raw.(*rds.RevokeAccountPrivilegeResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if err := s.WaitForAccountPrivilegeRevoked(id, dbName, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

func (s *RdsService) ReleaseDBPublicConnection(instanceId, connection string) error {
	request := rds.CreateReleaseInstancePublicConnectionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId
	request.CurrentConnectionString = connection

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ReleaseInstancePublicConnection(request)
	})
	response, ok := raw.(*rds.ReleaseInstancePublicConnectionResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil
}

func (s *RdsService) ModifyDBBackupPolicy(d *schema.ResourceData, updateForData, updateForLog bool) error {
	enableBackupLog := "1"

	backupPeriod := ""
	if v, ok := d.GetOk("preferred_backup_period"); ok && v.(*schema.Set).Len() > 0 {
		periodList := expandStringList(v.(*schema.Set).List())
		backupPeriod = fmt.Sprintf("%s", strings.Join(periodList[:], COMMA_SEPARATED))
	}

	backupTime := "02:00Z-03:00Z"
	if v, ok := d.GetOk("preferred_backup_time"); ok && v.(string) != "02:00Z-03:00Z" {
		backupTime = v.(string)
	}

	retentionPeriod := "7"
	if v, ok := d.GetOk("backup_retention_period"); ok && v.(int) != 7 {
		retentionPeriod = strconv.Itoa(v.(int))
	}

	logBackupRetentionPeriod := ""
	if v, ok := d.GetOk("log_backup_retention_period"); ok && v.(int) != 0 {
		logBackupRetentionPeriod = strconv.Itoa(v.(int))
	}

	localLogRetentionHours := ""
	if v, ok := d.GetOk("local_log_retention_hours"); ok {
		localLogRetentionHours = strconv.Itoa(v.(int))
	}

	localLogRetentionSpace := ""
	if v, ok := d.GetOk("local_log_retention_space"); ok {
		localLogRetentionSpace = strconv.Itoa(v.(int))
	}

	highSpaceUsageProtection := d.Get("high_space_usage_protection").(string)

	if !d.Get("enable_backup_log").(bool) {
		enableBackupLog = "0"
	}

	if d.HasChange("log_backup_retention_period") {
		if d.Get("log_backup_retention_period").(int) > d.Get("backup_retention_period").(int) {
			logBackupRetentionPeriod = retentionPeriod
		}
	}

	logBackupFrequency := ""
	if v, ok := d.GetOk("log_backup_frequency"); ok {
		logBackupFrequency = v.(string)
	}
	compressType := ""
	if v, ok := d.GetOk("compress_type"); ok {
		compressType = v.(string)
	}

	archiveBackupRetentionPeriod := "0"
	if v, ok := d.GetOk("archive_backup_retention_period"); ok {
		archiveBackupRetentionPeriod = strconv.Itoa(v.(int))
	}

	archiveBackupKeepCount := 1
	if v, ok := d.GetOk("archive_backup_keep_count"); ok {
		archiveBackupKeepCount = v.(int)
	}

	archiveBackupKeepPolicy := "0"
	if v, ok := d.GetOk("archive_backup_keep_policy"); ok {
		archiveBackupKeepPolicy = v.(string)
	}

	instance, err := s.DescribeDBInstance(d.Id())
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if updateForData {
		request := rds.CreateModifyBackupPolicyRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.PreferredBackupPeriod = backupPeriod
		request.PreferredBackupTime = backupTime
		request.BackupRetentionPeriod = retentionPeriod
		request.CompressType = compressType
		request.BackupPolicyMode = "DataBackupPolicy"
		if instance.Engine == "SQLServer" && logBackupFrequency == "LogInterval" {
			request.LogBackupFrequency = logBackupFrequency
		}
		if instance.Engine == "MySQL" && instance.DBInstanceStorageType == "local_ssd" {

			request.ArchiveBackupRetentionPeriod = archiveBackupRetentionPeriod
			request.ArchiveBackupKeepCount = requests.NewInteger(archiveBackupKeepCount)
			request.ArchiveBackupKeepPolicy = archiveBackupKeepPolicy
		}
		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyBackupPolicy(request)
		})

		response, ok := raw.(*rds.ModifyBackupPolicyResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if err := s.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
	}

	// At present, the sql server database does not support setting logBackupRetentionPeriod
	if updateForLog && instance.Engine != "SQLServer" {
		request := rds.CreateModifyBackupPolicyRequest()
		s.client.InitRpcRequest(*request.RpcRequest)
		request.DBInstanceId = d.Id()
		request.EnableBackupLog = enableBackupLog
		request.LocalLogRetentionHours = localLogRetentionHours
		request.LocalLogRetentionSpace = localLogRetentionSpace
		request.HighSpaceUsageProtection = highSpaceUsageProtection
		request.BackupPolicyMode = "LogBackupPolicy"
		request.LogBackupRetentionPeriod = logBackupRetentionPeriod

		raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ModifyBackupPolicy(request)
		})
		response, ok := raw.(*rds.ModifyBackupPolicyResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if err := s.WaitForDBInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	return nil
}

func (s *RdsService) ModifyDBSecurityIps(instanceId, ips string) error {
	request := rds.CreateModifySecurityIpsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId
	request.SecurityIps = ips

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifySecurityIps(request)
	})
	response, ok := raw.(*rds.ModifySecurityIpsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err := s.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func (s *RdsService) DescribeDBSecurityIps(instanceId string) (ips []rds.DBInstanceIPArray, err error) {
	request := rds.CreateDescribeDBInstanceIPArrayListRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = instanceId

	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceIPArrayList(request)
	})
	resp, ok := raw.(*rds.DescribeDBInstanceIPArrayListResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, instanceId, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resp.Items.DBInstanceIPArray, nil
}

func (s *RdsService) GetSecurityIps(instanceId string) ([]string, error) {
	object, err := s.DescribeDBSecurityIps(instanceId)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range object {
		if ip.DBInstanceIPArrayAttribute == "hidden" {
			continue
		}
		ips += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}

	for _, ip := range strings.Split(ips, COMMA_SEPARATED) {
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

func (s *RdsService) DescribeSecurityGroupConfiguration(id string) ([]string, error) {
	request := rds.CreateDescribeSecurityGroupConfigurationRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeSecurityGroupConfiguration(request)
	})

	response, ok := raw.(*rds.DescribeSecurityGroupConfigurationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	groupIds := make([]string, 0)
	for _, v := range response.Items.EcsSecurityGroupRelation {
		groupIds = append(groupIds, v.SecurityGroupId)
	}
	return groupIds, nil
}

func (s *RdsService) DescribeDBInstanceSSL(id string) (*rds.DescribeDBInstanceSSLResponse, error) {
	response := &rds.DescribeDBInstanceSSLResponse{}
	request := rds.CreateDescribeDBInstanceSSLRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceSSL(request)
	})
	response, ok := raw.(*rds.DescribeDBInstanceSSLResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

func (s *RdsService) DescribeRdsTDEInfo(id string) (*rds.DescribeDBInstanceTDEResponse, error) {
	response := &rds.DescribeDBInstanceTDEResponse{}
	request := rds.CreateDescribeDBInstanceTDERequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	statErr := s.WaitForDBInstance(id, Running, DefaultLongTimeout)
	if statErr != nil {
		return response, errmsgs.WrapError(statErr)
	}
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeDBInstanceTDE(request)
	})
	response, ok := raw.(*rds.DescribeDBInstanceTDEResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

func (s *RdsService) ModifySecurityGroupConfiguration(id string, groupid string) error {
	request := rds.CreateModifySecurityGroupConfigurationRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	// openapi required that input "Empty" if groupid is ""
	if len(groupid) == 0 {
		groupid = "Empty"
	}
	request.SecurityGroupId = groupid
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.ModifySecurityGroupConfiguration(request)
	})
	response, ok := raw.(*rds.ModifySecurityGroupConfigurationResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

// return multiIZ list of current region
func (s *RdsService) DescribeMultiIZByRegion() (izs []string, err error) {
	request := rds.CreateDescribeRegionsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RegionId = s.client.RegionId
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeRegions(request)
	})
	response, ok := raw.(*rds.DescribeRegionsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeMultiIZByRegion", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	regions := response.Regions.RDSRegion

	zoneIds := []string{}
	for _, r := range regions {
		if r.RegionId == string(s.client.Region) && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, r.ZoneId)
		}
	}

	return zoneIds, nil
}

func (s *RdsService) DescribeBackupPolicy(id string) (*rds.DescribeBackupPolicyResponse, error) {
	policy := &rds.DescribeBackupPolicyResponse{}
	request := rds.CreateDescribeBackupPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeBackupPolicy(request)
	})

	response, ok := raw.(*rds.DescribeBackupPolicyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return policy, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "DescribeMultiIZByRegion", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return raw.(*rds.DescribeBackupPolicyResponse), nil
}

func (s *RdsService) DescribeDbInstanceMonitor(id string) (monitoringPeriod int, err error) {
	request := rds.CreateDescribeDBInstanceMonitorRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		return client.DescribeDBInstanceMonitor(request)
	})
	response, ok := raw.(*rds.DescribeDBInstanceMonitorResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return 0, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	monPeriod, err := strconv.Atoi(response.Period)
	if err != nil {
		return 0, errmsgs.WrapError(err)
	}
	return monPeriod, nil
}

func (s *RdsService) DescribeSQLCollectorPolicy(id string) (collectorPolicy *rds.DescribeSQLCollectorPolicyResponse, err error) {
	request := rds.CreateDescribeSQLCollectorPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeSQLCollectorPolicy(request)
	})
	response, ok := raw.(*rds.DescribeSQLCollectorPolicyResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return collectorPolicy, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

func (s *RdsService) DescribeSQLCollectorRetention(id string) (collectorRetention *rds.DescribeSQLCollectorRetentionResponse, err error) {
	request := rds.CreateDescribeSQLCollectorRetentionRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = id
	raw, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeSQLCollectorRetention(request)
	})
	response, ok := raw.(*rds.DescribeSQLCollectorRetentionResponse)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return collectorRetention, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

// WaitForInstance waits for instance to given status
func (s *RdsService) WaitForDBInstance(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && strings.ToLower(object.DBInstanceStatus) == strings.ToLower(string(status)) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBInstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) RdsDBInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstance(id)
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

func (s *RdsService) RdsTaskStateRefreshFunc(id string, taskAction string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTasks(id)
		if err != nil {

			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, t := range object.Items.TaskProgressInfo {
			if t.TaskAction == taskAction {
				return object, t.Status, nil
			}
		}

		return object, "Pending", nil
	}
}

// WaitForDBParameter waits for instance parameter to given value.
// Status of DB instance is Running after ModifyParameters API was
// call, so we can not just wait for instance status become
// Running, we should wait until parameters have expected values.
func (s *RdsService) WaitForDBParameter(instanceId string, timeout int, expects map[string]string) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeParameters(instanceId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		var actuals = make(map[string]string)
		for _, i := range object.RunningParameters.DBInstanceParameter {
			actuals[i.ParameterName] = i.ParameterValue
		}
		for _, i := range object.ConfigParameters.DBInstanceParameter {
			actuals[i.ParameterName] = i.ParameterValue
		}

		match := true

		got_value := ""
		expected_value := ""

		for name, expect := range expects {
			if actual, ok := actuals[name]; ok {
				if expect != actual {
					match = false
					got_value = actual
					expected_value = expect
					break
				}
			} else {
				match = false
			}
		}

		if match {
			break
		}

		time.Sleep(DefaultIntervalShort * time.Second)

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, instanceId, GetFunc(1), timeout, got_value, expected_value, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForDBConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && object.ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *RdsService) WaitForDBReadWriteSplitting(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBReadWriteSplittingConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if err == nil {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForAccount(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBAccount(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil {
			if object.AccountStatus == string(status) {
				break
			} else if object.AccountStatus == "Lock" {
				request := rds.CreateDeleteAccountRequest()
				s.client.InitRpcRequest(*request.RpcRequest)
				request.DBInstanceId = object.DBInstanceId
				request.AccountName = object.AccountName

				_, err := s.client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
					return rdsClient.DeleteAccount(request)
				})
				if err != nil && !errmsgs.IsExpectedErrors(err, []string{"InvalidAccountName.NotFound"}) {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
				}
			}
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.AccountStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilege(id, dbName string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(parts[0] + ":" + dbName)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		ready := false
		if object != nil {
			for _, account := range object.Accounts.AccountPrivilegeInfo {
				// At present, postgresql response has a bug, DBOwner will be changed to ALL
				if account.Account == parts[1] && (account.AccountPrivilege == parts[2] || (parts[2] == "DBOwner" && account.AccountPrivilege == "ALL")) {
					ready = true
					break
				}
			}
		}
		if status == Deleted && !ready {
			break
		}
		if ready {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *RdsService) WaitForAccountPrivilegeRevoked(id, dbName string, timeout int) error {
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(parts[0] + ":" + dbName)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return errmsgs.WrapError(err)
		}

		exist := false
		if object != nil {
			for _, account := range object.Accounts.AccountPrivilegeInfo {
				if account.Account == parts[1] && (account.AccountPrivilege == parts[2] || (parts[2] == "DBOwner" && account.AccountPrivilege == "ALL")) {
					exist = true
					break
				}
			}
		}

		if !exist {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", dbName, errmsgs.ProviderERROR)
		}

	}
	return nil
}

func (s *RdsService) WaitForDBDatabase(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeDBDatabase(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errmsgs.WrapError(err)
		}
		if object != nil && object.DBName == parts[1] {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.DBName, parts[1], errmsgs.ProviderERROR)
		}
	}
	return nil
}

// turn period to TimeType
func (s *RdsService) TransformPeriod2Time(period int, chargeType string) (ut int, tt common.TimeType) {
	if chargeType == string(Postpaid) {
		return 1, common.Day
	}

	if period >= 1 && period <= 9 {
		return period, common.Month
	}

	if period == 12 {
		return 1, common.Year
	}

	if period == 24 {
		return 2, common.Year
	}
	return 0, common.Day

}

// turn TimeType to Period
func (s *RdsService) TransformTime2Period(ut int, tt common.TimeType) (period int) {
	if tt == common.Year {
		return 12 * ut
	}

	return ut

}

func (s *RdsService) flattenDBSecurityIPs(list []rds.DBInstanceIPArray) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIPList,
		}
		result = append(result, l)
	}
	return result
}

func (s *RdsService) setInstanceTags(d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		remove, add := diffRdsTags(o, n)

		if len(remove) > 0 {
			request := rds.CreateUntagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.ResourceType = "INSTANCE"
			request.TagKey = &remove
			wait := incrementalWait(1*time.Second, 2*time.Second)
			var raw interface{}
			var err error
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err = s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
					return client.UntagResources(request)
				})
				if err != nil {
					if errmsgs.IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			response, ok := raw.(*rds.UntagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}

		if len(add) > 0 {
			request := rds.CreateTagResourcesRequest()
			s.client.InitRpcRequest(*request.RpcRequest)
			request.ResourceId = &[]string{d.Id()}
			request.Tag = &add
			request.ResourceType = "INSTANCE"
			wait := incrementalWait(1*time.Second, 2*time.Second)
			var raw interface{}
			var err error
			err = resource.Retry(10*time.Minute, func() *resource.RetryError {
				raw, err = s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
					return client.TagResources(request)
				})
				if err != nil {
					if errmsgs.IsThrottling(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			response, ok := raw.(*rds.TagResourcesResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
		}

		//d.SetPartial("tags")
	}

	return nil
}

func (s *RdsService) describeTags(d *schema.ResourceData) (tags []Tag, err error) {
	request := rds.CreateDescribeTagsRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DBInstanceId = d.Id()
	raw, err := s.client.WithRdsClient(func(client *rds.Client) (interface{}, error) {
		return client.DescribeTags(request)
	})
	response, ok := raw.(*rds.DescribeTagsResponse)
	if err != nil {
		tmp := make([]Tag, 0)
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return tmp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return s.respToTags(response.Items.TagInfos), nil
}

func (s *RdsService) respToTags(tagSet []rds.TagInfos) (tags []Tag) {
	result := make([]Tag, 0, len(tagSet))
	for _, t := range tagSet {
		tag := Tag{
			Key:   t.TagKey,
			Value: t.TagValue,
		}
		result = append(result, tag)
	}

	return result
}

func (s *RdsService) tagsToMap(tags []Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func (s *RdsService) ignoreTag(t Tag) bool {
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

func (s *RdsService) tagsToString(tags []Tag) string {
	v, _ := json.Marshal(s.tagsToMap(tags))

	return string(v)
}
