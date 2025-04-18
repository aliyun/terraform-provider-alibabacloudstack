package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolardbService struct {
	client *connectivity.AlibabacloudStackClient
}

type PolardbCheckaccountnameavailableResponse struct {
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbCheckaccountnameavailableRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbCheckaccountnameavailableResponse, error) {
	// api: polardb - 2024-01-30 - CheckAccountNameAvailable
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "CheckAccountNameAvailable", "")
	PolardbCheckaccountnameavailableResponse := &PolardbCheckaccountnameavailableResponse{}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "CheckAccountNameAvailable", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbCheckaccountnameavailableResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "CheckAccountNameAvailable", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbCheckaccountnameavailableResponse, nil
}

type PolardbDescribeaccountsResponse struct {
	Accounts struct {
		DBInstanceAccount []struct {
			DatabasePrivileges struct {
				DatabasePrivilege []struct {
					DBName                 string `json:"DBName"`
					AccountPrivilege       string `json:"AccountPrivilege"`
					AccountPrivilegeDetail string `json:"AccountPrivilegeDetail"`
				} `json:"DatabasePrivilege"`
			} `json:"DatabasePrivileges"`
			DBInstanceId       string `json:"DBInstanceId"`
			AccountName        string `json:"AccountName"`
			AccountStatus      string `json:"AccountStatus"`
			AccountType        string `json:"AccountType"`
			AccountDescription string `json:"AccountDescription"`
			PrivExceeded       string `json:"PrivExceeded"`
			ValidUntil         string `json:"ValidUntil"`
			CreateDB           string `json:"CreateDB"`
			Replication        string `json:"Replication"`
			CreateRole         string `json:"CreateRole"`
			BypassRLS          string `json:"BypassRLS"`
		} `json:"DBInstanceAccount"`
	} `json:"Accounts"`
	RequestId                             string `json:"RequestId"`
	SystemAdminAccountStatus              string `json:"SystemAdminAccountStatus"`
	SystemAdminAccountFirstActivationTime string `json:"SystemAdminAccountFirstActivationTime"`
}

func (s *PolardbService) DescribeDBAccount(id string) (*PolardbDescribeaccountsResponse, error) {
	parts, _ := ParseResourceId(id, 2)
	// api: polardb - 2024-01-30 - DescribeAccounts
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeAccounts", "")
	PolardbDescribeaccountsResponse := &PolardbDescribeaccountsResponse{}

	request.QueryParams["AccountName"] = parts[1]

	// 常规参数填充
	request.QueryParams["DBInstanceId"] = parts[0]

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
			return nil, nil
		}
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeaccountsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribeaccountsResponse.Accounts.DBInstanceAccount) < 1 {
		return PolardbDescribeaccountsResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("accounts", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribeaccountsResponse, nil
}

func (s *PolardbService) DoPolardbDescribeaccountsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeaccountsResponse, error) {
	// api: polardb - 2024-01-30 - DescribeAccounts
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeAccounts", "")
	PolardbDescribeaccountsResponse := &PolardbDescribeaccountsResponse{}

	//调用request_params_handler

	// 常规参数填充
	if v, ok := d.GetOk("account_name"); ok && v != "" {
		//调用requestin_handler
		request.QueryParams["AccountName"] = v.(string)
	}

	// 常规参数填充
	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		//调用requestin_handler
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return nil, fmt.Errorf("DataBaseInstanceId is required")
	}

	// 常规参数填充
	if v, ok := d.GetOk("page_number"); ok {
		//调用requestin_handler
		request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
	}

	// 常规参数填充
	if v, ok := d.GetOk("page_size"); ok {
		//调用requestin_handler
		request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeAccounts", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeaccountsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeAccounts", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribeaccountsResponse.Accounts.DBInstanceAccount) < 1 {
		return PolardbDescribeaccountsResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("accounts", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return PolardbDescribeaccountsResponse, nil
}

type PolardbDescribedatabasesResponse struct {
	Databases struct {
		Database []struct {
			Accounts struct {
				AccountPrivilegeInfo []struct {
					Account          string `json:"Account"`
					AccountPrivilege string `json:"AccountPrivilege"`
				} `json:"AccountPrivilegeInfo"`
			} `json:"Accounts"`
			DBName           string `json:"DBName"`
			DBInstanceId     string `json:"DBInstanceId"`
			Engine           string `json:"Engine"`
			DBStatus         string `json:"DBStatus"`
			CharacterSetName string `json:"CharacterSetName"`
			DBDescription    string `json:"DBDescription"`
			Collate          string `json:"Collate"`
			Ctype            string `json:"Ctype"`
			ConnLimit        int    `json:"ConnLimit"`
			Tablespace       string `json:"Tablespace"`
		} `json:"Database"`
	} `json:"Databases"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DescribeDBDatabase(id string) (*PolardbDescribedatabasesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDatabases
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDatabases", "")
	PolardbDescribedatabasesResponse := &PolardbDescribedatabasesResponse{}

	parts, err := ParseResourceId(id, 2)

	request.QueryParams["DBInstanceId"] = parts[0]

	request.QueryParams["DBName"] = parts[1]
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
			return nil, nil
		}
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, []string{"InvalidDBName.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDatabases", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedatabasesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDatabases", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedatabasesResponse.Databases.Database) < 1 {
		return PolardbDescribedatabasesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Databases", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedatabasesResponse, nil
}
func (s *PolardbService) DoPolardbDescribedatabasesRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedatabasesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDatabases
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDatabases", "")
	PolardbDescribedatabasesResponse := &PolardbDescribedatabasesResponse{}

	//调用request_params_handler

	// 常规参数填充
	if v, ok := d.GetOk("data_base_instance_id"); ok && v != "" {
		//调用requestin_handler
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return nil, fmt.Errorf("DataBaseInstanceId is required")
	}

	// 常规参数填充
	if v, ok := d.GetOk("data_base_name"); ok && v != "" {
		//调用requestin_handler
		request.QueryParams["DBName"] = v.(string)
	}

	// 常规参数填充
	if v, ok := d.GetOk("page_number"); ok {
		//调用requestin_handler
		request.QueryParams["PageNumber"] = strconv.Itoa(v.(int))
	}

	// 常规参数填充
	if v, ok := d.GetOk("page_size"); ok {
		//调用requestin_handler
		request.QueryParams["PageSize"] = strconv.Itoa(v.(int))
	}

	// 常规参数填充
	if v, ok := d.GetOk("status"); ok && v != "" {
		//调用requestin_handler
		request.QueryParams["DBStatus"] = v.(string)
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDatabases", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedatabasesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDatabases", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedatabasesResponse.Databases.Database) < 1 {
		return PolardbDescribedatabasesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Databases", "")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedatabasesResponse, nil
}

type PolardbDescribebackuppolicyResponse struct {
	RequestId                     string `json:"RequestId"`
	BackupRetentionPeriod         int    `json:"BackupRetentionPeriod"`
	PreferredNextBackupTime       string `json:"PreferredNextBackupTime"`
	MaxRetention                  int    `json:"MaxRetention"`
	PreferredBackupTime           string `json:"PreferredBackupTime"`
	PreferredBackupPeriod         string `json:"PreferredBackupPeriod"`
	BackupLog                     string `json:"BackupLog"`
	LogBackupRetentionPeriod      string `json:"LogBackupRetentionPeriod"`
	EnableBackupLog               string `json:"EnableBackupLog"`
	LocalLogRetentionHours        string `json:"LocalLogRetentionHours"`
	LocalLogRetentionSpace        string `json:"LocalLogRetentionSpace"`
	CompressType                  string `json:"CompressType"`
	Duplication                   string `json:"Duplication"`
	DuplicationContent            string `json:"DuplicationContent"`
	HighSpaceUsageProtection      string `json:"HighSpaceUsageProtection"`
	LogBackupFrequency            string `json:"LogBackupFrequency"`
	ArchiveBackupRetentionPeriod  int    `json:"ArchiveBackupRetentionPeriod"`
	ArchiveBackupKeepPolicy       int    `json:"ArchiveBackupKeepPolicy"`
	ArchiveBackupKeepCount        int    `json:"ArchiveBackupKeepCount"`
	ReleasedKeepPolicy            string `json:"ReleasedKeepPolicy"`
	LogBackupLocalRetentionNumber int    `json:"LogBackupLocalRetentionNumber"`
	BackupMethod                  string `json:"BackupMethod"`

	DuplicationLocation struct {
		Sotrage string `json:"Sotrage"`

		Location struct {
			Endpoint string `json:"Endpoint"`
			Bucket   string `json:"Bucket"`
		} `json:"Location"`
	} `json:"DuplicationLocation"`
}

func (s *PolardbService) DoPolardbDescribebackuppolicyRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribebackuppolicyResponse, error) {
	// api: polardb - 2024-01-30 - DescribeBackupPolicy
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeBackupPolicy", "")
	PolardbDescribebackuppolicyResponse := &PolardbDescribebackuppolicyResponse{}

	//调用request_params_handler

	if v, ok := d.GetOk("backup_policy_mode"); ok {
		//调用requestin_handler
		request.QueryParams["BackupPolicyMode"] = v.(string)
	}

	if v, ok := d.GetOk("db_instance_id"); ok {
		//调用requestin_handler
		request.QueryParams["DBInstanceId"] = v.(string)
	} else {
		return nil, fmt.Errorf("DBInstanceId is required")
	}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeBackupPolicy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribebackuppolicyResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeBackupPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribebackuppolicyResponse, nil
}

type PolardbDescriberegionsResponse struct {
	Regions struct {
		RDSRegion []struct {
			RegionId  string `json:"RegionId"`
			ZoneId    string `json:"ZoneId"`
			SubDomain string `json:"SubDomain"`
			SubZoneId string `json:"SubZoneId"`
		} `json:"RDSRegion"`
	} `json:"Regions"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescriberegionsRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescriberegionsResponse, error) {
	// api: polardb - 2024-01-30 - DescribeRegions
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeRegions", "")
	PolardbDescriberegionsResponse := &PolardbDescriberegionsResponse{}

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeRegions", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescriberegionsResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeRegions", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescriberegionsResponse, nil
}

type PolardbDescribedbinstancenetinfoResponse struct {
	DBInstanceNetInfos struct {
		DBInstanceNetInfo []struct {
			SecurityIPGroups struct {
				securityIPGroup []struct {
					SecurityIPGroupName string `json:"SecurityIPGroupName"`
					SecurityIPs         string `json:"SecurityIPs"`
				} `json:"securityIPGroup"`
			} `json:"SecurityIPGroups"`

			DBInstanceWeights struct {
				DBInstanceWeight []struct {
					DBInstanceId   string `json:"DBInstanceId"`
					DBInstanceType string `json:"DBInstanceType"`
					Availability   string `json:"Availability"`
					Weight         string `json:"Weight"`
				} `json:"DBInstanceWeight"`
			} `json:"DBInstanceWeights"`
			Upgradeable          string `json:"Upgradeable"`
			ExpiredTime          string `json:"ExpiredTime"`
			ConnectionString     string `json:"ConnectionString"`
			IPAddress            string `json:"IPAddress"`
			IPType               string `json:"IPType"`
			Port                 string `json:"Port"`
			VPCId                string `json:"VPCId"`
			VSwitchId            string `json:"VSwitchId"`
			ConnectionStringType string `json:"ConnectionStringType"`
			MaxDelayTime         string `json:"MaxDelayTime"`
			DistributionType     string `json:"DistributionType"`
		} `json:"DBInstanceNetInfo"`
	} `json:"DBInstanceNetInfos"`
	RequestId           string `json:"RequestId"`
	InstanceNetworkType string `json:"InstanceNetworkType"`
	SecurityIPMode      string `json:"SecurityIPMode"`
}

func (s *PolardbService) DoPolardbDescribedbinstancenetinfoRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string) (*PolardbDescribedbinstancenetinfoResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceNetInfo
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceNetInfo", "")
	PolardbDescribedbinstancenetinfoResponse := &PolardbDescribedbinstancenetinfoResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceNetInfo", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancenetinfoResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceNetInfo", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstancenetinfoResponse, nil
}

type PolardbDescribedbinstanceattributeResponse struct {
	Items struct {
		DBInstanceAttribute []struct {
			SlaveZones struct {
				SlaveZone []struct {
					ZoneId string `json:"ZoneId"`
				} `json:"SlaveZone"`
			} `json:"SlaveZones"`

			ReadOnlyDBInstanceIds struct {
				ReadOnlyDBInstanceId []struct {
					DBInstanceId string `json:"DBInstanceId"`
				} `json:"ReadOnlyDBInstanceId"`
			} `json:"ReadOnlyDBInstanceIds"`
			IPType                            string      `json:"IPType"`
			DBInstanceDiskUsed                string      `json:"DBInstanceDiskUsed"`
			GuardDBInstanceName               string      `json:"GuardDBInstanceName"`
			CanTempUpgrade                    bool        `json:"CanTempUpgrade"`
			TempUpgradeTimeStart              string      `json:"TempUpgradeTimeStart"`
			TempUpgradeTimeEnd                string      `json:"TempUpgradeTimeEnd"`
			TempUpgradeRecoveryTime           string      `json:"TempUpgradeRecoveryTime"`
			TempUpgradeRecoveryClass          string      `json:"TempUpgradeRecoveryClass"`
			TempUpgradeRecoveryCpu            int         `json:"TempUpgradeRecoveryCpu"`
			TempUpgradeRecoveryMemory         int         `json:"TempUpgradeRecoveryMemory"`
			TempUpgradeRecoveryMaxIOPS        string      `json:"TempUpgradeRecoveryMaxIOPS"`
			SqlGrammarCompatibility           string      `json:"SqlGrammarCompatibility"`
			TempUpgradeRecoveryMaxConnections string      `json:"TempUpgradeRecoveryMaxConnections"`
			InsId                             int         `json:"InsId"`
			DBInstanceId                      string      `json:"DBInstanceId"`
			PayType                           string      `json:"PayType"`
			DBInstanceClassType               string      `json:"DBInstanceClassType"`
			DBInstanceType                    string      `json:"DBInstanceType"`
			RegionId                          string      `json:"RegionId"`
			ConnectionString                  string      `json:"ConnectionString"`
			SlaveConnectionString             string      `json:"SlaveConnectionString"`
			Port                              string      `json:"Port"`
			Engine                            string      `json:"Engine"`
			EngineVersion                     string      `json:"EngineVersion"`
			DBInstanceClass                   string      `json:"DBInstanceClass"`
			DBInstanceMemory                  int         `json:"DBInstanceMemory"`
			DBInstanceStorage                 int         `json:"DBInstanceStorage"`
			VpcCloudInstanceId                string      `json:"VpcCloudInstanceId"`
			DBInstanceNetType                 interface{} `json:"DBInstanceNetType"`
			DBInstanceStatus                  string      `json:"DBInstanceStatus"`
			DBInstanceDescription             string      `json:"DBInstanceDescription"`
			LockMode                          string      `json:"LockMode"`
			LockReason                        string      `json:"LockReason"`
			ReadDelayTime                     string      `json:"ReadDelayTime"`
			DBMaxQuantity                     int         `json:"DBMaxQuantity"`
			AccountMaxQuantity                int         `json:"AccountMaxQuantity"`
			CreationTime                      string      `json:"CreationTime"`
			ExpireTime                        string      `json:"ExpireTime"`
			MaintainTime                      string      `json:"MaintainTime"`
			AvailabilityValue                 string      `json:"AvailabilityValue"`
			MaxIOPS                           int         `json:"MaxIOPS"`
			MaxConnections                    int         `json:"MaxConnections"`
			MasterInstanceId                  string      `json:"MasterInstanceId"`
			DBInstanceCPU                     string      `json:"DBInstanceCPU"`
			IncrementSourceDBInstanceId       string      `json:"IncrementSourceDBInstanceId"`
			GuardDBInstanceId                 string      `json:"GuardDBInstanceId"`
			ReplicateId                       string      `json:"ReplicateId"`
			TempDBInstanceId                  string      `json:"TempDBInstanceId"`
			SecurityIPList                    string      `json:"SecurityIPList"`
			ZoneId                            string      `json:"ZoneId"`
			InstanceNetworkType               string      `json:"InstanceNetworkType"`
			DBInstanceStorageType             string      `json:"DBInstanceStorageType"`
			EncryptionKey                     string      `json:"EncryptionKey"`
			AdvancedFeatures                  string      `json:"AdvancedFeatures"`
			Category                          string      `json:"Category"`
			AccountType                       string      `json:"AccountType"`
			SupportUpgradeAccountType         string      `json:"SupportUpgradeAccountType"`
			SupportCreateSuperAccount         string      `json:"SupportCreateSuperAccount"`
			VpcId                             string      `json:"VpcId"`
			VSwitchId                         string      `json:"VSwitchId"`
			ConnectionMode                    string      `json:"ConnectionMode"`
			CurrentKernelVersion              string      `json:"CurrentKernelVersion"`
			LatestKernelVersion               string      `json:"LatestKernelVersion"`
			CurrentKernelShowVersion          string      `json:"CurrentKernelShowVersion"`
			ResourceGroupId                   string      `json:"ResourceGroupId"`
			ReadonlyInstanceSQLDelayedTime    string      `json:"ReadonlyInstanceSQLDelayedTime"`
			SecurityIPMode                    string      `json:"SecurityIPMode"`
			TimeZone                          string      `json:"TimeZone"`
			Collation                         string      `json:"Collation"`
			DispenseMode                      string      `json:"DispenseMode"`
			MasterZone                        string      `json:"MasterZone"`
			AutoUpgradeMinorVersion           string      `json:"AutoUpgradeMinorVersion"`
			ProxyType                         int         `json:"ProxyType"`
			ConsoleVersion                    string      `json:"ConsoleVersion"`
			CpuType                           string      `json:"CpuType"`
			Vip                               string      `json:"Vip"`
			Vip_v6                            string      `json:"Vip_v6"`
			Vport                             string      `json:"Vport"`

			Extra struct {
				DBInstanceId struct {
					DBInstanceId []string `json:"DBInstanceId"`
				} `json:"DBInstanceId"`
				ReplicaGroupID            string `json:"ReplicaGroupID"`
				ReplicaGroupStatus        string `json:"ReplicaGroupStatus"`
				ActiveReplicaDBInstanceID string `json:"ActiveReplicaDBInstanceID"`
			} `json:"Extra"`
		} `json:"DBInstanceAttribute"`
	} `json:"Items"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescribedbinstanceattributeRequest(id string, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstanceattributeResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceAttribute
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceAttribute", "")
	PolardbDescribedbinstanceattributeResponse := &PolardbDescribedbinstanceattributeResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = id

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceattributeResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstanceattributeResponse.Items.DBInstanceAttribute) < 1 {
		return PolardbDescribedbinstanceattributeResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstanceattributeResponse, nil
}

type PolardbDescribedbinstancemonitorResponse struct {
	RequestId string `json:"RequestId"`
	Period    string `json:"Period"`
}

func (s *PolardbService) DoPolardbDescribedbinstancemonitorRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstancemonitorResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceMonitor
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceMonitor", "")
	PolardbDescribedbinstancemonitorResponse := &PolardbDescribedbinstancemonitorResponse{}

	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceMonitor", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancemonitorResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceMonitor", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribedbinstancemonitorResponse, nil
}

type PolardbDescribeparametersResponse struct {
	ConfigParameters struct {
		DBInstanceParameter []struct {
			ParameterName        string `json:"ParameterName"`
			ParameterValue       string `json:"ParameterValue"`
			ParameterDescription string `json:"ParameterDescription"`
		} `json:"DBInstanceParameter"`
	} `json:"ConfigParameters"`

	RunningParameters struct {
		DBInstanceParameter []struct {
			ParameterName        string `json:"ParameterName"`
			ParameterValue       string `json:"ParameterValue"`
			ParameterDescription string `json:"ParameterDescription"`
		} `json:"DBInstanceParameter"`
	} `json:"RunningParameters"`
	RequestId     string `json:"RequestId"`
	Engine        string `json:"Engine"`
	EngineVersion string `json:"EngineVersion"`
}

func (s *PolardbService) DoPolardbDescribeparametersRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeparametersResponse, error) {
	// api: polardb - 2024-01-30 - DescribeParameters
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeParameters", "")
	PolardbDescribeparametersResponse := &PolardbDescribeparametersResponse{}

	//调用request_params_handler
	request.QueryParams["DBInstanceId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeParameters", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeparametersResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeParameters", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribeparametersResponse, nil
}

type PolardbDescribedbinstanceiparraylistResponse struct {
	Items struct {
		DBInstanceIPArray []struct {
			DBInstanceIPArrayName      string `json:"DBInstanceIPArrayName"`
			DBInstanceIPArrayAttribute string `json:"DBInstanceIPArrayAttribute"`
			SecurityIPType             string `json:"SecurityIPType"`
			SecurityIPList             string `json:"SecurityIPList"`
			WhitelistNetworkType       string `json:"WhitelistNetworkType"`
		} `json:"DBInstanceIPArray"`
	} `json:"Items"`
	RequestId string `json:"RequestId"`
}

func (s *PolardbService) DoPolardbDescribedbinstanceiparraylistRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstanceiparraylistResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstanceIPArrayList
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceIPArrayList", "")
	PolardbDescribedbinstanceiparraylistResponse := &PolardbDescribedbinstanceiparraylistResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = d.Id()

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceIPArrayList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceiparraylistResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceIPArrayList", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribedbinstanceiparraylistResponse, nil
}

// not finish
func (s *PolardbService) WaitForPolardbConnection(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		parts, err := ParseResourceId(d.Id(), 2)
		object, err := s.DoPolardbDescribedbinstancenetinfoRequest(d, client, parts[0])
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && object.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, object.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString, d.Id(), errmsgs.ProviderERROR)
		}
	}
}
func (s *PolardbService) WaitForDBInstance(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DoPolardbDescribedbinstancesRequest(d.Id(), client)
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
			if status == Deleted && len(object.Items.DBInstance) == 0 {
				break
			} else if strings.ToLower(object.Items.DBInstance[0].DBInstanceStatus) == strings.ToLower(string(status)) {
				break
			}
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, d.Id(), GetFunc(1), timeout, object.Items.DBInstance[0].DBInstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *PolardbService) GetSecurityIps(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) ([]string, error) {
	object, err := s.DoPolardbDescribedbinstanceiparraylistRequest(d, client)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	var ips, separator string
	ipsMap := make(map[string]string)
	for _, ip := range object.Items.DBInstanceIPArray {
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

func (s *PolardbService) RefreshParameters(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, attribute string) error {
	var param []map[string]interface{}
	documented, ok := d.GetOk(attribute)
	if !ok {
		d.Set(attribute, param)
		return nil
	}
	object, err := s.DoPolardbDescribeparametersRequest(d, client)
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

func (s *PolardbService) ModifyDBSecurityIps(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, ips string) error {
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifySecurityIps", "")
	PolardbModifysecurityipsResponse := PolardbModifysecurityipsResponse{}
	request.QueryParams["DBInstanceId"] = d.Id()
	request.QueryParams["SecurityIps"] = ips

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
			"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbModifysecurityipsResponse)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_polardb_db_instance", "ModifySecurityIps", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}

type PolardbDescribeinstanceautorenewalattributeResponse struct {
	Items struct {
		Item []struct {
			DBInstanceId string `json:"DBInstanceId"`
			RegionId     string `json:"RegionId"`
			Duration     int    `json:"Duration"`
			Status       string `json:"Status"`
			AutoRenew    string `json:"AutoRenew"`
		} `json:"Item"`
	} `json:"Items"`
	RequestId        string `json:"RequestId"`
	PageNumber       int    `json:"PageNumber"`
	TotalRecordCount int    `json:"TotalRecordCount"`
	PageRecordCount  int    `json:"PageRecordCount"`
}

func (s *PolardbService) DoPolardbDescribeinstanceautorenewalattributeRequest(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient) (*PolardbDescribeinstanceautorenewalattributeResponse, error) {
	// api: polardb - 2024-01-30 - DescribeInstanceAutoRenewalAttribute
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeInstanceAutoRenewalAttribute", "")
	PolardbDescribeinstanceautorenewalattributeResponse := &PolardbDescribeinstanceautorenewalattributeResponse{}

	request.QueryParams["DBInstanceId"] = d.Id()
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeInstanceAutoRenewalAttribute", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribeinstanceautorenewalattributeResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeInstanceAutoRenewalAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return PolardbDescribeinstanceautorenewalattributeResponse, nil
}

func (s *PolardbService) ModifyParameters(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, attribute string) error {
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "ModifyParameters", "")

	request.QueryParams["DBInstanceId"] = d.Id()
	request.QueryParams["Forcerestart"] = d.Get("force_restart").(string)
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
		request.QueryParams["Parameters"] = string(cfg)
		// wait instance status is Normal before modifying
		if err := s.WaitForDBInstance(d, client, Running, DefaultLongTimeout); err != nil {
			return errmsgs.WrapError(err)
		}
		// Need to check whether some parameter needs restart
		if !d.Get("force_restart").(bool) {
			req := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeParameterTemplates", "")
			req.QueryParams["DBInstanceId"] = d.Id()
			req.QueryParams["Engine"] = d.Get("engine").(string)
			req.QueryParams["EngineVersion"] = d.Get("engine_version").(string)
			req.QueryParams["ClientToken"] = buildClientToken(req.GetActionName())
			forceRestartMap := make(map[string]string)
			bresponse, err := client.ProcessCommonRequest(req)
			DescribeParameterTemplatesResponse := DescribeParameterTemplatesResponse{}
			if err != nil {
				if bresponse == nil {
					return errmsgs.WrapErrorf(err, "Process Common Request Failed")
				}
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			err = json.Unmarshal(bresponse.GetHttpContentBytes(), &DescribeParameterTemplatesResponse)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
					"alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			for _, para := range DescribeParameterTemplatesResponse.Parameters.TemplateRecord {
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
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_polardb_db_instance", "CreateReadOnlyDBInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		// wait instance parameter expect after modifying
		for _, i := range ns.List() {
			key := i.(map[string]interface{})["name"].(string)
			value := i.(map[string]interface{})["value"].(string)
			allConfig[key] = value
		}
		//待实现
		// if err := s.WaitForDBParameter(d.Id(), DefaultTimeoutMedium, allConfig); err != nil {
		// 	return errmsgs.WrapError(err)
		// }
	}
	//d.SetPartial(attribute)
	return nil
}

type DescribeParameterTemplatesResponse struct {
	*responses.BaseResponse
	RequestId      string     `json:"RequestId" xml:"RequestId"`
	Engine         string     `json:"Engine" xml:"Engine"`
	ParameterCount string     `json:"ParameterCount" xml:"ParameterCount"`
	EngineVersion  string     `json:"EngineVersion" xml:"EngineVersion"`
	Parameters     Parameters `json:"Parameters" xml:"Parameters"`
}
type Parameters struct {
	TemplateRecord []TemplateRecord `json:"TemplateRecord" xml:"TemplateRecord"`
}
type TemplateRecord struct {
	CheckingCode         string `json:"CheckingCode" xml:"CheckingCode"`
	ParameterName        string `json:"ParameterName" xml:"ParameterName"`
	ParameterValue       string `json:"ParameterValue" xml:"ParameterValue"`
	ForceModify          string `json:"ForceModify" xml:"ForceModify"`
	ForceRestart         string `json:"ForceRestart" xml:"ForceRestart"`
	ParameterDescription string `json:"ParameterDescription" xml:"ParameterDescription"`
}

func (s *PolardbService) PolardbDBInstanceStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DoPolardbDescribedbinstancesRequest(d.Id(), client)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Items.DBInstance[0].DBInstanceStatus == failState {
				return object, object.Items.DBInstance[0].DBInstanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Items.DBInstance[0].DBInstanceStatus))
			}
		}
		return object, object.Items.DBInstance[0].DBInstanceStatus, nil
	}
}

func (s *PolardbService) PolardbDBInstanceTdeStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstanceTDE(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["TDEStatus"].(string) == failState {
				return object, object["TDEStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["TDEStatus"].(string)))
			}
		}
		return object, object["TDEStatus"].(string), nil
	}
}

func (s *PolardbService) DescribeDBInstanceTDE(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceTDE", "")
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceTDE", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

func (s *PolardbService) PolardbDBInstanceSslStateRefreshFunc(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstanceSSL(d.Id())
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["SSLEnabled"].(string) == failState {
				return object, object["SSLEnabled"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["SSLEnabled"].(string)))
			}
		}
		return object, object["SSLEnabled"].(string), nil
	}
}

func (s *PolardbService) DescribeDBInstanceSSL(id string) (map[string]interface{}, error) {
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceSSL", "")
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceSSL", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &result)
	return result, nil
}

func (s *PolardbService) WaitForDBConnection(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		response, err := s.DescribeDBConnection(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		data := response
		if data != nil && data.DBInstanceNetInfos.DBInstanceNetInfo[0].ConnectionString != "" {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, errmsgs.ProviderERROR)
		}
	}
}

func (s *PolardbService) WaitForConnectionDBInstance(d *schema.ResourceData, client *connectivity.AlibabacloudStackClient, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DoPolardbDescribedbinstancesRequest(id, client)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object != nil && strings.ToLower(object.Items.DBInstance[0].DBInstanceStatus) == strings.ToLower(string(status)) {
			break
		}
		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Items.DBInstance[0].DBInstanceStatus, status, errmsgs.ProviderERROR)
		}
	}
	return nil
}

func (s *PolardbService) DoPolardbDescribedbinstancesRequest(id string, client *connectivity.AlibabacloudStackClient) (*PolardbDescribedbinstancesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstances
	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstances", "")
	PolardbDescribedbinstancesResponse := &PolardbDescribedbinstancesResponse{}
	request.QueryParams["DBInstanceId"] = id
	request.QueryParams["InstanceLevel"] = "1"
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "1"
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancesResponse.Items.DBInstance) < 1 {
		return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}
	return PolardbDescribedbinstancesResponse, nil
}

func (s *PolardbService) Describedbinstances(id string) (*PolardbDescribedbinstancesResponse, error) {
	// api: polardb - 2024-01-30 - DescribeDBInstances
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstances", "")
	PolardbDescribedbinstancesResponse := &PolardbDescribedbinstancesResponse{}
	request.QueryParams["DBInstanceId"] = id
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancesResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancesResponse.Items.DBInstance) < 1 {
		return PolardbDescribedbinstancesResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("PolardbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return PolardbDescribedbinstancesResponse, nil
}

func (s *PolardbService) DescribeDBSecurityIps(instanceId string) (*PolardbDescribedbinstanceiparraylistResponse, error) {
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceIPArrayList", "")
	PolardbDescribedbinstanceiparraylistResponse := &PolardbDescribedbinstanceiparraylistResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = instanceId

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeDBInstanceIPArrayList", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstanceiparraylistResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceIPArrayList", errmsgs.AlibabacloudStackSdkGoERROR)

	}
	return PolardbDescribedbinstanceiparraylistResponse, nil

}

func (s *PolardbService) flattenDBSecurityIPs(resp *PolardbDescribedbinstanceiparraylistResponse) []map[string]interface{} {
	list := resp.Items.DBInstanceIPArray
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"security_ips": i.SecurityIPList,
		}
		result = append(result, l)
	}
	return result
}

func (s *PolardbService) DescribeDBConnection(id string) (*PolardbDescribedbinstancenetinfoResponse, error) {
	parts, _ := ParseResourceId(id, 2)
	request := s.client.NewCommonRequest("POST", "polardb", "2024-01-30", "DescribeDBInstanceNetInfo", "")
	PolardbDescribedbinstancenetinfoResponse := &PolardbDescribedbinstancenetinfoResponse{}

	//调用request_params_handler

	request.QueryParams["DBInstanceId"] = parts[0]

	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound"}) {
			return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapError(err)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &PolardbDescribedbinstancenetinfoResponse)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeDBInstanceNetInfo", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo) < 1 {
		return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBInstanceNetInfo", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	object := PolardbDescribedbinstancenetinfoResponse.DBInstanceNetInfos.DBInstanceNetInfo
	if object != nil {
		for _, o := range object {
			if strings.HasPrefix(o.ConnectionString, parts[1]) {
				return PolardbDescribedbinstancenetinfoResponse, nil
			}
		}
	}

	return PolardbDescribedbinstancenetinfoResponse, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DBConnection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}
